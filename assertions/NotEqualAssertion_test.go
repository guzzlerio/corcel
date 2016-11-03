package assertions

import (
	"github.com/guzzlerio/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotEqualAssertion", func() {
	key := "some:key"

	It("Succeeds", func() {
		executionResult := core.ExecutionResult{
			key: 8,
		}

		assertion := NotEqualAssertion{
			Key:   key,
			Value: 7,
		}

		result := assertion.Assert(executionResult)
		Expect(result[core.AssertionResultUrn.String()]).To(Equal(true))
		Expect(result[core.AssertionMessageUrn.String()]).To(BeNil())
	})

	It("Fails", func() {
		executionResult := core.ExecutionResult{
			key: 7,
		}

		assertion := NotEqualAssertion{
			Key:   key,
			Value: 7,
		}

		result := assertion.Assert(executionResult)
		Expect(result[core.AssertionResultUrn.String()]).To(Equal(false))
		Expect(result[core.AssertionMessageUrn.String()]).To(Equal("FAIL: 7 does match 7"))
	})
})
