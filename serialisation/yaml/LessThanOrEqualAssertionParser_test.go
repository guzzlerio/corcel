package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LessThanOrEqualAssertionParser", func() {
	It("Parses", func() {

		var expectedKey = "talula"
		var input = map[string]interface{}{
			"key":      expectedKey,
			"expected": 7,
		}

		var parser = LessThanOrEqualAssertionParser{}
		var assertion = parser.Parse(input).(*assertions.LessThanOrEqualAssertion)

		Expect(assertion.Key).To(Equal(expectedKey))
		Expect(assertion.Value).To(Equal(7))
	})

	It("Returns Key", func() {
		Expect(LessThanOrEqualAssertionParser{}.Key()).To(Equal("LessThanOrEqualAssertion"))
	})
})