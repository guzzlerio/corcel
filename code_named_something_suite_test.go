package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCodeNamedSomething(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CodeNamedSomething Suite")
}
