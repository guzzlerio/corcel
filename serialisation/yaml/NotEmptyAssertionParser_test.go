package yaml

import (
	"fmt"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotEmptyAssertionParser", func() {

	It("Parses", func() {

		var expected = "talula"
		var input = map[string]interface{}{
			"key": expected,
		}

		var parser = NotEmptyAssertionParser{}
		assertion, err := parser.Parse(input)
		emptyAssertion := assertion.(*assertions.NotEmptyAssertion)
		Expect(err).To(BeNil())

		Expect(emptyAssertion.Key).To(Equal(expected))
	})

	It("Fails to parse", func() {
		var expected = "talula"
		var input = map[string]interface{}{
			"bang": expected,
		}

		var parser = NotEmptyAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring("key is not present"))
	})
})
