package inproc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInproc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inproc Suite")
}
