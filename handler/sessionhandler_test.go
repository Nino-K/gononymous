package handler

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Nino-K/gononymous/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SessionHandler", func() {
	Describe("Join", func() {
		It("returns 400 when CLIENT_ID header is not provided", func() {
			sm := server.NewSessionManager()

			upgrader := newMockUpgrader()
			upgrader.UpgradeOutput.Ret0 <- nil
			upgrader.UpgradeOutput.Ret1 <- nil

			sessionHandler := NewSessionHandler(sm, upgrader)

			testServer := httptest.NewServer(http.HandlerFunc(sessionHandler.Join))
			resp, err := http.Get(testServer.URL)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			Expect(responseString(resp.Body)).To(ContainSubstring(peerIdErr.Error()))
		})
		It("returns 500 when upgrader fails", func() {
			sm := server.NewSessionManager()

			upgrader := newMockUpgrader()
			upgrader.UpgradeOutput.Ret0 <- nil
			upgrader.UpgradeOutput.Ret1 <- errors.New("bad stuff")

			sessionHandler := NewSessionHandler(sm, upgrader)

			testServer := httptest.NewServer(http.HandlerFunc(sessionHandler.Join))
			resp, err := get(testServer.URL)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
		It("returns 400 when sessionId is not provided", func() {
			sm := server.NewSessionManager()

			upgrader := newMockUpgrader()
			upgrader.UpgradeOutput.Ret0 <- nil
			upgrader.UpgradeOutput.Ret1 <- nil

			sessionHandler := NewSessionHandler(sm, upgrader)

			testServer := httptest.NewServer(http.HandlerFunc(sessionHandler.Join))
			resp, err := get(testServer.URL)
			Expect(err).ToNot(HaveOccurred())
			Expect(responseString(resp.Body)).To(ContainSubstring(sessionIDErr.Error()))
		})
	})

	Describe("sessionId", func() {
		It("returns a second segment of the provided url", func() {
			testURL, err := url.Parse("http://127.0.0.1:9090/sessionId")
			Expect(err).NotTo(HaveOccurred())
			sessionId := sessionId(testURL)
			Expect(sessionId).To(Equal("sessionId"))
		})
	})
})

func get(u string) (*http.Response, error) {
	testURL, err := url.Parse(u)
	Expect(err).NotTo(HaveOccurred())

	client := http.DefaultClient
	header := make(http.Header)
	header.Add("CLIENT_ID", "testId")

	req := &http.Request{
		URL:    testURL,
		Method: "GET",
		Header: header,
	}
	return client.Do(req)
}

func responseString(body io.ReadCloser) string {
	content, err := ioutil.ReadAll(body)
	Expect(err).NotTo(HaveOccurred())
	return string(content)
}
