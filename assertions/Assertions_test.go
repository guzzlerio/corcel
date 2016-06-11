package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Assertions", func() {

	It("Exact Assertion Succeeds", func() {
		key := "some:key"
		expectedValue := 7

		executionResult := core.ExecutionResult{
			key: expectedValue,
		}

		assertion := ExactAssertion{
			Key:      key,
			Expected: expectedValue,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})

	It("Exact Assertion Fails", func() {
		key := "some:key"
		expectedValue := 7

		executionResult := core.ExecutionResult{
			key: 8,
		}

		assertion := ExactAssertion{
			Key:      key,
			Expected: expectedValue,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: 8 does not match 7"))
	})

})
