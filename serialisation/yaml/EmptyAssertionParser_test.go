package yaml

import (
	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmptyAssertionParser", func() {

	It("Parses", func() {

		var expected = "talula"
		var input = map[string]interface{}{
			"key": expected,
		}

		var parser = EmptyAssertionParser{}
		var assertion = parser.Parse(input).(*assertions.EmptyAssertion)

		Expect(assertion.Key).To(Equal(expected))
	})
})
