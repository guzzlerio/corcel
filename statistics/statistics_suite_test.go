package statistics

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCorcel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Request Suite")
}
