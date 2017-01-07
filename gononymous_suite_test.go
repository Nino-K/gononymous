package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGononymous(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gononymous Suite")
}
