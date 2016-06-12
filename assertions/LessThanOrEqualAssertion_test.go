package assertions

import (
	"fmt"

	"ci.guzzler.io/guzzler/corcel/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LessThanOrEqualAssertion", func() {
	key := "some:key"

	Context("Succeeds", func() {

		var testCases = []AssertionTestCase{
			NewAsssertionTestCase(float64(5.1), float64(6.1)),
			NewAsssertionTestCase(int(5), float64(6.1)),
			NewAsssertionTestCase("5", float64(6.1)),
			NewAsssertionTestCase(float64(3.1), int(6)),
			NewAsssertionTestCase(float64(6.1), float64(6.1)),
			NewAsssertionTestCase("fubar", "fubar"),
			NewAsssertionTestCase("6.1", float64(6.1)),
			NewAsssertionTestCase(int(6), int(6)),
			NewAsssertionTestCase("6", int(6)),
			NewAsssertionTestCase(int(3), int(6)),
			NewAsssertionTestCase("5.1", int(6)),
			NewAsssertionTestCase("f", "fubar"),
			NewAsssertionTestCase("fubar", "1.1"),
		}

		assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
			key := "some:key"
			executionResult := core.ExecutionResult{
				key: actualValue,
			}

			assertion := LessThanOrEqualAssertion{
				Key:   key,
				Value: instanceValue,
			}

			result := assertion.Assert(executionResult)
			Expect(result["result"]).To(Equal(true))
			Expect(result["message"]).To(BeNil())
		})
	})

	Context("Fails", func() {

		var testCases = []AssertionTestCase{
			NewAsssertionTestCase(NilValue, NilValue),
			NewAsssertionTestCase(NilValue, int(5)),
			NewAsssertionTestCase(NilValue, "5.1"),
			NewAsssertionTestCase(NilValue, "fubar"),
			NewAsssertionTestCase(NilValue, float64(6.1)),
			NewAsssertionTestCase("fubar", int(6)),
			NewAsssertionTestCase("fubar", float64(1.1)),
			NewAsssertionTestCase(float64(1.1), NilValue),
			NewAsssertionTestCase(int(1), NilValue),
			NewAsssertionTestCase("1", NilValue),
			NewAsssertionTestCase("a", NilValue),
			NewAsssertionTestCase(float64(5), float64(1)),
			NewAsssertionTestCase(int(5), float64(1)),
			NewAsssertionTestCase("2.2", float64(1)),
			NewAsssertionTestCase(int(5), int(1)),
			NewAsssertionTestCase("5", int(1)),
			NewAsssertionTestCase("abc", "a"),
			NewAsssertionTestCase(float64(1.3), "1.2"),
			NewAsssertionTestCase(int(3), "1"),
			NewAsssertionTestCase("3.1", "2"),
			NewAsssertionTestCase(float64(6.1), "fubar"),
			NewAsssertionTestCase(int(6), "fubar"),
		}

		assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
			executionResult := core.ExecutionResult{
				key: actualValue,
			}

			assertion := LessThanOrEqualAssertion{
				Key:   key,
				Value: instanceValue,
			}

			result := assertion.Assert(executionResult)
			Expect(result["result"]).To(Equal(false))
			Expect(result["message"]).To(Equal(fmt.Sprintf("FAIL: %v is not greater %v", actualValue, instanceValue)))
		})
	})
})
