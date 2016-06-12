package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExactAssertion", func() {
	key := "some:key"

	It("Exact Assertion Succeeds", func() {
		expectedValue := 7

		executionResult := core.ExecutionResult{
			key: expectedValue,
		}

		assertion := ExactAssertion{
			Key:   key,
			Value: expectedValue,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(true))
		Expect(result["message"]).To(BeNil())
	})

	It("Exact Assertion Fails", func() {
		expectedValue := 7

		executionResult := core.ExecutionResult{
			key: 8,
		}

		assertion := ExactAssertion{
			Key:   key,
			Value: expectedValue,
		}

		result := assertion.Assert(executionResult)
		Expect(result["result"]).To(Equal(false))
		Expect(result["message"]).To(Equal("FAIL: 8 does not match 7"))
	})

	//NOTHING is currently using the message when an assertion fails but we will need
	//it for when we put the errors into the report.  One of the edge cases with the message
	//is that say the actual value was a string "7" and the expected is an int 7.  The message
	//will not include the quotes so the message would read 7 does not equal 7 as opposed
	//to "7" does not equal 7.  Notice this is a type mismatch
	/*
		PIt("Exact Assertion Fails when actual and expected are different types", func() {
			key := "some:key"
			expectedValue := 7

			executionResult := core.ExecutionResult{
				key: "7",
			}

			assertion := ExactAssertion{
				Key:   key,
				Value: expectedValue,
			}

			result := assertion.Assert(executionResult)
			Expect(result["result"]).To(Equal(false))
			Expect(result["message"]).To(Equal("FAIL: \"7\" does not match 7"))
		})
	*/
})
