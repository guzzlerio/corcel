package assertions

import (
	"fmt"
	"strconv"

	"ci.guzzler.io/guzzler/corcel/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type AssertionTestCase struct {
	Actual               interface{}
	Instance             interface{}
	ActualStringNumber   bool
	InstanceStringNumber bool
}

func NewAsssertionTestCase(actual interface{}, instance interface{}) (newInstance AssertionTestCase) {
	newInstance.Actual = actual
	newInstance.Instance = instance
	switch actualType := actual.(type) {
	case string:
		_, err := strconv.ParseFloat(actualType, 64)
		if err == nil {
			newInstance.ActualStringNumber = true
		}
	}
	switch instanceType := instance.(type) {
	case string:
		_, err := strconv.ParseFloat(instanceType, 64)
		if err == nil {
			newInstance.InstanceStringNumber = true
		}
	}
	return
}

var _ = FDescribe("Assertions", func() {

	key := "some:key"

	assert := func(testCases []AssertionTestCase, test func(actual interface{}, instance interface{})) {

		for _, testCase := range testCases {
			actualValue := testCase.Actual
			instanceValue := testCase.Instance
			testName := fmt.Sprintf("When Actual is of type %T and Instance is of type %T", actualValue, instanceValue)
			if testCase.ActualStringNumber {
				testName = fmt.Sprintf("%s. Actual value is a STRING NUMBER in this case", testName)
			}
			if testCase.InstanceStringNumber {
				testName = fmt.Sprintf("%s. Instance value is a STRING NUMBER in this case", testName)
			}
			It(testName, func() {
				test(actualValue, instanceValue)
			})
		}
	}

	/*
	   Using this test function and the test cases we get a nice readable output.  If you run ginkgo -v here is an example of the output

	   Assertions Greater Than Succeeds
	     When Actual is of type int and Instance is of type string. Instance value is a STRING NUMBER
	*/

	var nilValue interface{}

	FContext("Greater Than", func() {

		Context("Succeeds", func() {
			var testCases = []AssertionTestCase{
				NewAsssertionTestCase(float64(1.1), nilValue),
				NewAsssertionTestCase(int(1), nilValue),
				NewAsssertionTestCase("1", nilValue),
				NewAsssertionTestCase("a", nilValue),
				NewAsssertionTestCase(float64(5), float64(1)),
				NewAsssertionTestCase(int(5), float64(1)),
				NewAsssertionTestCase("2.2", float64(1)),
				NewAsssertionTestCase(int(5), int(1)),
				NewAsssertionTestCase("5", int(1)),
				NewAsssertionTestCase("abc", "a"),
				NewAsssertionTestCase(float64(1.3), "1.2"),
				NewAsssertionTestCase(int(3), "1"),
				NewAsssertionTestCase("3.1", "2"),
			}

			assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
				key := "some:key"
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
		})

		FContext("Fails", func() {
			var testCases = []AssertionTestCase{
				NewAsssertionTestCase(nilValue, nilValue),
				NewAsssertionTestCase(nilValue, int(5)),
				NewAsssertionTestCase(nilValue, "5.1"),
				NewAsssertionTestCase(nilValue, "fubar"),
				NewAsssertionTestCase(nilValue, float64(6.1)),
				NewAsssertionTestCase(float64(5.1), float64(6.1)),
				NewAsssertionTestCase(int(5), float64(6.1)),
				NewAsssertionTestCase("5", float64(6.1)),
				NewAsssertionTestCase(float64(3.1), int(6)),
				NewAsssertionTestCase(int(3), int(6)),
				NewAsssertionTestCase("5.1", int(6)),
				NewAsssertionTestCase("fubar", int(6)),
				NewAsssertionTestCase("fubar", float64(1.1)),
				NewAsssertionTestCase("fubar", "1.1"),
			}

			assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
				executionResult := core.ExecutionResult{
					key: actualValue,
				}

				assertion := GreaterThanAssertion{
					Key:   key,
					Value: instanceValue,
				}

				result := assertion.Assert(executionResult)
				Expect(result["result"]).To(Equal(false))
				Expect(result["message"]).To(Equal(fmt.Sprintf("FAIL: %v is not greater %v", actualValue, instanceValue)))
			})
			/*
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
			*/
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
