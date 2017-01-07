package main_test

import (
	"os/exec"

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

})
