package yaml

import (
	"fmt"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("EmptyAssertionParser", func() {

	It("Parses", func() {

		var expected = "talula"
		var input = map[string]interface{}{
			"key": expected,
		}

		var parser = EmptyAssertionParser{}
		assertion, err := parser.Parse(input)
		emptyAssertion := assertion.(*assertions.EmptyAssertion)
		Expect(err).To(BeNil())

		Expect(emptyAssertion.Key).To(Equal(expected))
	})

	It("Fails to parse", func() {
		var expected = "talula"
		var input = map[string]interface{}{
			"bang": expected,
		}

		var parser = EmptyAssertionParser{}
		_, err := parser.Parse(input)

		Expect(err).ToNot(BeNil())
		Expect(fmt.Sprintf("%v", err)).To(ContainSubstring(""))
	})
})
