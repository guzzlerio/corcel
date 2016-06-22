package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"
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
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
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
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: 7 does match 7"))
	})
})
