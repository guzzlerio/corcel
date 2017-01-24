package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotEmptyAssertionParser", func() {
	It("Parses", func() {

		var expectedKey = "talula"
		var input = map[string]interface{}{
			"key": expectedKey,
		}

		var parser = NotEmptyAssertionParser{}
		var assertion = parser.Parse(input).(*assertions.NotEmptyAssertion)

		Expect(assertion.Key).To(Equal(expectedKey))
	})

	It("Returns Key", func() {
		Expect(NotEmptyAssertionParser{}.Key()).To(Equal("NotEmptyAssertion"))
	})
})
