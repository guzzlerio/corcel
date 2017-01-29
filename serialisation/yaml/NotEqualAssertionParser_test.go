package yaml

import (
	"fmt"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotEqualAssertionParser", func() {

	It("Parses", func() {

		var expectedKey = "talula"
		var expectedValue = "boomboom"
		var input = map[string]interface{}{
			"key":      expectedKey,
			"expected": expectedValue,
		}

		var parser = NotEqualAssertionParser{}
		assertion, err := parser.Parse(input)
		var exactAssertion = assertion.(*assertions.NotEqualAssertion)

		Expect(err).To(BeNil())
		Expect(exactAssertion.Key).To(Equal(expectedKey))
		Expect(exactAssertion.Value).To(Equal(expectedValue))
	})

	It("Fails to parse without key", func() {

		var input = map[string]interface{}{
			"bang":     "talula",
			"expected": "boomboom",
		}

		var parser = NotEqualAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("key is not present"))
	})

	It("Fails to parse without expected", func() {

		var input = map[string]interface{}{
			"key":  "talula",
			"bang": "boomboom",
		}

		var parser = NotEqualAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("expected is not present"))
	})

})
