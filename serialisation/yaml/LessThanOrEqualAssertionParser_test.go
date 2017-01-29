package yaml

import (
	"fmt"

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
		assertion, err := parser.Parse(input)
		var ltAssertion = assertion.(*assertions.LessThanOrEqualAssertion)

		Expect(err).To(BeNil())
		Expect(ltAssertion.Key).To(Equal(expectedKey))
		Expect(ltAssertion.Value).To(Equal(7))
	})

	It("Returns Key", func() {
		Expect(LessThanOrEqualAssertionParser{}.Key()).To(Equal("LessThanOrEqualAssertion"))
	})

	It("Fails to parse without key", func() {

		var input = map[string]interface{}{
			"bang":     "talula",
			"expected": "boomboom",
		}

		var parser = LessThanOrEqualAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("key is not present"))
	})

	It("Fails to parse without expected", func() {

		var input = map[string]interface{}{
			"key":  "talula",
			"bang": "boomboom",
		}

		var parser = LessThanOrEqualAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("expected is not present"))
	})

})
