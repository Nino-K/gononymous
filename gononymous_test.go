package main_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"

	"github.com/gorilla/websocket"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Gononymous", func() {
	var (
		gononymousPath string
		err            error
	)

	BeforeEach(func() {
		gononymousPath, err = gexec.Build("main.go", "-race")
		Expect(err).NotTo(HaveOccurred())
	})

	Context("Argument handling", func() {
		It("listens on a default 9797 port if no port is provided", func() {
			command := exec.Command(gononymousPath)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say("gononymous is listening on 127.0.0.1:9797"))
			session.Kill()
		})
		It("listens on 127.0.0.1 if no IP is provided", func() {
			command := exec.Command(gononymousPath, "-port", "8888")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say("gononymous is listening on 127.0.0.1:8888"))
			session.Kill()
		})
		It("returns an error if an invalid port is given", func() {
			command := exec.Command(gononymousPath, "-port", "6500034")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Err).Should(gbytes.Say("-port must be within range 1024-65535"))
			Eventually(session).Should(gexec.Exit(1))
			session.Kill()
		})
		It("listens on the provided IP and port", func() {
			command := exec.Command(gononymousPath, "-port", "9494", "-addr", "localhost")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say("gononymous is listening on localhost:9494"))
			session.Kill()
		})
	})

	Context("http errors", func() {
		It("returns http 400 error if CLIENT_ID header is not provided", func() {
			testPort := testPort()
			command := exec.Command(gononymousPath, "-port", strconv.Itoa(testPort), "-addr", "localhost")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf("gononymous is listening on localhost:%d", testPort)))

			_, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d", testPort), nil)
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			Expect(responseBody(resp.Body)).To(ContainSubstring("CLIENT_ID header must be provided"))

			session.Kill()
		})
		It("returns http 400 error if sessionId is not provided", func() {
			testPort := testPort()
			command := exec.Command(gononymousPath, "-port", strconv.Itoa(testPort), "-addr", "localhost")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf("gononymous is listening on localhost:%d", testPort)))

			header := make(http.Header)
			header.Add("CLIENT_ID", "testClient")
			_, resp, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost:%d", testPort), header)
			Expect(err).To(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

			Expect(responseBody(resp.Body)).To(ContainSubstring("sessionId must be provided"))

			session.Kill()
		})
	})

	Context("Succesfull connection", func() {
		var session *gexec.Session
		var randPort int

		BeforeEach(func() {
			randPort = testPort()
			command := exec.Command(gononymousPath, "-port", strconv.Itoa(randPort), "-addr", "localhost")
			var err error
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session.Out).Should(gbytes.Say(fmt.Sprintf("gononymous is listening on localhost:%d", randPort)))
		})

		AfterEach(func() {
			session.Kill()
		})

		It("connects multiple clients with the same sessionId", func() {
			wsClientOne := wsClient(fmt.Sprintf("ws://localhost:%d/testSession", randPort), "client1")
			wsClientTwo := wsClient(fmt.Sprintf("ws://localhost:%d/testSession", randPort), "client2")

			err = wsClientOne.WriteMessage(websocket.TextMessage, []byte("yolo"))
			Expect(err).ToNot(HaveOccurred())

			msgType, msg, err := wsClientTwo.ReadMessage()
			Expect(err).ToNot(HaveOccurred())
			Expect(msgType).To(Equal(websocket.TextMessage))
			Expect(msg).To(Equal([]byte("yolo")))
		})
	})
})

func wsClient(url, clientId string) *websocket.Conn {
	header := make(http.Header)
	header.Add("CLIENT_ID", clientId)
	wsClient, _, err := websocket.DefaultDialer.Dial(url, header)
	Expect(err).ToNot(HaveOccurred())
	return wsClient
}

func responseBody(b io.ReadCloser) string {
	body, err := ioutil.ReadAll(b)
	Expect(err).NotTo(HaveOccurred())
	return string(body)
}

func testPort() int {
	add, _ := net.ResolveTCPAddr("tcp", ":0")
	l, _ := net.ListenTCP("tcp", add)
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	return port
}
