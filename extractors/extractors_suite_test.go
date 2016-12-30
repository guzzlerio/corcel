package extractors_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestExtractors(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extractors Suite")
}
