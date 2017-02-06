package main_test

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"time"

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
		It("returns http 400 error", func() {
			testPort := testPort()
			command := exec.Command(gononymousPath, "-port", strconv.Itoa(testPort), "-addr", "127.0.0.1")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			// make sure server is up
			Eventually(connected(testPort, time.Second*2)).Should(Equal(http.StatusSwitchingProtocols))

			By("not providing CLIENT_ID")

			Eventually(func() string {
				_, resp, err := dialer().Dial(fmt.Sprintf("wss://127.0.0.1:%d/test", testPort), nil)
				Expect(err).To(HaveOccurred())
				// this is due to not reaching to upgrade, and returing prior to that
				Expect(err.Error()).To(ContainSubstring("bad handshake"))
				return responseBody(resp.Body)
			}).Should(ContainSubstring("CLIENT_ID header must be provided"))

			By("not providing sessionId")

			header := make(http.Header)
			header.Add("CLIENT_ID", "testClient")
			Eventually(func() string {
				_, resp, err := dialer().Dial(fmt.Sprintf("wss://127.0.0.1:%d/", testPort), header)
				Expect(err).To(HaveOccurred())
				// this is due to not reaching to upgrade, and returing prior to that
				Expect(err.Error()).To(ContainSubstring("bad handshake"))
				return responseBody(resp.Body)
			}).Should(ContainSubstring("sessionId must be provided"))
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

			// make sure server is up
			Eventually(connected(randPort, time.Second*2)).Should(Equal(http.StatusSwitchingProtocols))
		})

		AfterEach(func() {
			session.Kill()
		})

		It("connects multiple clients with the same sessionId", func() {
			wsClientOne := wsClient(fmt.Sprintf("wss://localhost:%d/testSession", randPort), "client1")
			wsClientTwo := wsClient(fmt.Sprintf("wss://localhost:%d/testSession", randPort), "client2")

			err = wsClientOne.WriteMessage(websocket.TextMessage, []byte("yolo"))
			Expect(err).ToNot(HaveOccurred())

			msgType, msg, err := wsClientTwo.ReadMessage()
			Expect(err).ToNot(HaveOccurred())
			Expect(msgType).To(Equal(websocket.TextMessage))
			Expect(msg).To(Equal([]byte("yolo")))
		})
	})
})

func connected(testPort int, wait time.Duration) int {
	timer := time.NewTimer(wait)
	<-timer.C

	header := make(http.Header)
	header.Add("CLIENT_ID", "testClient")
	_, resp, err := dialer().Dial(fmt.Sprintf("wss://127.0.0.1:%d/testSession", testPort), header)
	Expect(err).ToNot(HaveOccurred())
	return resp.StatusCode
}

func dialer() *websocket.Dialer {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	dialer := websocket.Dialer{
		TLSClientConfig: tlsConfig,
	}
	return &dialer
}

func wsClient(url, clientId string) *websocket.Conn {
	header := make(http.Header)
	header.Add("CLIENT_ID", clientId)
	dialer := dialer()
	wsClient, _, err := dialer.Dial(url, header)
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
