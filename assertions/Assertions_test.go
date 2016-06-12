package assertions

import (
	"ci.guzzler.io/guzzler/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("Assertions", func() {

	key := "some:key"

	Context("Greater Than Assertion", func() {

		/*
					   Test Required
					                            INSTANCE

			                                 nil      float64    int    string-number   string
					             ACTUAL

					             float64         x          x          x          x           √

					             int             x          x          x          x           √

					             string-number   x          x          x          x           √

					             string          x          x          x          x           x

					             nil             x          x          x          x           x

		*/

		//To set further context I am making the following assumption
		//Something is greater than nil
		//nil is NOT greater than nil
		//nil is NOT greater than Something
		//string which is not a number is NOT greater than any number
		//number is NOT greater than a string which is not a number
		//Attempts will first be made to parse strings into a float64
		Context("Succeeds", func() {

			It("When Actual is float64 and Instance is nil", func() {
				key := "some:key"
				actualValue := 1.1
				var instanceValue struct{}

				executionResult := core.ExecutionResult{
					key: actualValue,
				}

				assertion := GreaterThanAssertion{
					Key:   key,
					Value: instanceValue,
				}

				result := assertion.Assert(executionResult)
				Expect(result["result"]).To(Equal(true))
				Expect(result["message"]).To(BeNil())
			})

			PIt("When Actual is int and Instance is nil", func() {

			})

			PIt("When Actual is string-number and Instance is nil", func() {

			})

			PIt("When Actual is string and Instance is nil", func() {

			})

			PIt("When Actual is float64 and Instance is float64", func() {

			})

			PIt("When Actual is int and Instance is float64", func() {

			})

			PIt("When Actual is string-number and Instance is float64", func() {

			})

			PIt("When Actual is float64 and Instance is int", func() {

			})

			PIt("When Actual is int and Instance is int", func() {

			})

			PIt("When Actual is string-number and Instance is int", func() {

			})

			PIt("When Actual is string and Instance is string", func() {

			})

			PIt("When Actual is float64 and Instance is string-number", func() {

			})

			PIt("When Actual is int and Instance is string-number", func() {

			})

			PIt("When Actual is string-number and Instance is string-number", func() {

			})
		})

		Context("Fails", func() {
			PIt("When Actual is nil and Instance is nil", func() {

			})

			PIt("When Actual is nil and Instance is int", func() {

			})

			PIt("When Actual is nil and Instance is string-number", func() {

			})

			PIt("When Actual is nil and Instance is string", func() {

			})

			PIt("When Actual is nil and Instance is float64", func() {

			})

			PIt("When Actual is float64 and Instance is float64", func() {

			})

			PIt("When Actual is int and Instance is float64", func() {

			})

			PIt("When Actual is string-number and Instance is float64", func() {

			})

			PIt("When Actual is float64 and Instance is int", func() {

			})

			PIt("When Actual is float64 and Instance is int", func() {

			})

			PIt("When Actual is int and Instance is int", func() {

			})

			PIt("When Actual is string-number and Instance is int", func() {

			})

			PIt("When Actual is string and Instance is int", func() {

			})

			PIt("When Actual is string and Instance is float64", func() {

			})

			PIt("When Actual is string and Instance is string-number", func() {

			})

			PIt("When Actual is string and Instance is string", func() {

			})

			PIt("When Actual is float64 and Instance is string", func() {

			})

			PIt("When Actual is int and Instance is string", func() {

			})

			PIt("When Actual string-number int and Instance is string", func() {

			})

		})

	})

	Context("Not Equal Assertion", func() {

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

	Context("Exact Assertion", func() {
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
	})

})
