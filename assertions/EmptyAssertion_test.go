package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("EmptyAssertion", func() {

	key := "some:key"

	It("Succeeds when empty string", func() {
		executionResult := core.ExecutionResult{
			key: "",
		}

		assertion := EmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})

	It("Succeeds when empty string of whitespace", func() {
		executionResult := core.ExecutionResult{
			key: "    ",
		}

		assertion := EmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})

	It("Succeeds when nil", func() {
		executionResult := core.ExecutionResult{}

		assertion := EmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})

	It("Fails when value is not nil", func() {

		executionResult := core.ExecutionResult{
			key: 8,
		}

		assertion := EmptyAssertion{
			Key: key,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: value is not empty"))
	})
})
