package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Corcel Suite")
}
