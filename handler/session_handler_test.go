package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Nino-K/gononymous/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SessionHandler", func() {
	FDescribe("Join", func() {
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
		})
		It("returns 500 when upgrader fails", func() {
			sm := server.NewSessionManager()

			upgrader := newMockUpgrader()
			upgrader.UpgradeOutput.Ret0 <- nil
			upgrader.UpgradeOutput.Ret1 <- errors.New("bad stuff")

			sessionHandler := NewSessionHandler(sm, upgrader)

			testServer := httptest.NewServer(http.HandlerFunc(sessionHandler.Join))
			testURL, err := url.Parse(testServer.URL)
			Expect(err).NotTo(HaveOccurred())

			client := http.DefaultClient
			header := make(http.Header)
			header.Add("CLIENT_ID", "testId")

			req := &http.Request{
				URL:    testURL,
				Method: "GET",
				Header: header,
			}
			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
		})
		// TODO: Test if sessionId is not provided
	})
})
