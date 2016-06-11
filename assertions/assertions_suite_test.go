package assertions

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAssertions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Assertions Suite")
}
