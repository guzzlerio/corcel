package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GreaterThanAssertionParser", func() {

	It("Parses", func() {

		var expectedKey = "talula"
		var input = map[string]interface{}{
			"key":      expectedKey,
			"expected": 7,
		}

		var parser = GreaterThanAssertionParser{}
		var assertion = parser.Parse(input).(*assertions.GreaterThanAssertion)

		Expect(assertion.Key).To(Equal(expectedKey))
		Expect(assertion.Value).To(Equal(7))
	})

	It("Returns Key", func() {
		Expect(GreaterThanAssertionParser{}.Key()).To(Equal("GreaterThanAssertion"))
	})
})
