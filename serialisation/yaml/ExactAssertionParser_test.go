package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EactAssertionParser", func() {

	It("Parses", func() {

		var expectedKey = "talula"
		var expectedValue = "boomboom"
		var input = map[string]interface{}{
			"key":      expectedKey,
			"expected": expectedValue,
		}

		var parser = ExactAssertionParser{}
		var assertion = parser.Parse(input).(*assertions.ExactAssertion)

		Expect(assertion.Key).To(Equal(expectedKey))
		Expect(assertion.Value).To(Equal(expectedValue))
	})

})
