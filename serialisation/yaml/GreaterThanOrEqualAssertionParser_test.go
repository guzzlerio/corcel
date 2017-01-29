package yaml

import (
	"fmt"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GreaterThanOrEqualAssertionParser", func() {
	It("Parses", func() {

		var expectedKey = "talula"
		var input = map[string]interface{}{
			"key":      expectedKey,
			"expected": 7,
		}

		var parser = GreaterThanOrEqualAssertionParser{}
		assertion, err := parser.Parse(input)
		var gteAssertion = assertion.(*assertions.GreaterThanOrEqualAssertion)

		Expect(err).To(BeNil())
		Expect(gteAssertion.Key).To(Equal(expectedKey))
		Expect(gteAssertion.Value).To(Equal(7))
	})

	It("Returns Key", func() {
		Expect(GreaterThanOrEqualAssertionParser{}.Key()).To(Equal("GreaterThanOrEqualAssertion"))
	})

	It("Fails to parse without key", func() {
		var input = map[string]interface{}{
			"boom":     "talula",
			"expected": 7,
		}

		var parser = GreaterThanOrEqualAssertionParser{}
		_, err := parser.Parse(input)
		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("key is not present"))
	})

	It("Fails to parse without expected", func() {
		var input = map[string]interface{}{
			"key":  "talula",
			"boom": 7,
		}

		var parser = GreaterThanOrEqualAssertionParser{}
		_, err := parser.Parse(input)
		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("expected is not present"))
	})
})
