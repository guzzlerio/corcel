package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotEmptyAssertion", func() {

	key := "some:key"

	It("Fails when empty string", func() {
		executionResult := core.ExecutionResult{
			key: "",
		}

		assertion := NotEmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: value is empty"))
	})

	It("Fails when empty string of whitespace", func() {
		executionResult := core.ExecutionResult{
			key: "    ",
		}

		assertion := NotEmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: value is empty"))
	})

	It("Fails when nil", func() {
		executionResult := core.ExecutionResult{}

		assertion := NotEmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: value is empty"))
	})

	It("Succeeds when value is not nil", func() {

		executionResult := core.ExecutionResult{
			key: 8,
		}

		assertion := NotEmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})
})
