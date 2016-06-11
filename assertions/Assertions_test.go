package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Assertions", func() {

	It("Exact Assertion Succeeds", func() {
		expectedKey := "some:key"
		expectedValue := 7

		executionResult := core.ExecutionResult{
			expectedKey: expectedValue,
		}

		assertion := ExactAssertion{
			Key:      expectedKey,
			Expected: expectedValue,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
	})

})
