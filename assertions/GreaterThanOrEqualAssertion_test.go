package assertions

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGreaterThanOrEqualAssertion(t *testing.T) {
	Convey("GreaterThanOrEqualAssertion", t, func() {

		key := "some:key"

		Convey("Succeeds", func() {
			var testCases = []AssertionTestCase{
				NewAsssertionTestCase(float64(1.1), NilValue),
				NewAsssertionTestCase(int(1), NilValue),
				NewAsssertionTestCase("1", NilValue),
				NewAsssertionTestCase("a", NilValue),
				NewAsssertionTestCase(float64(5), float64(1)),
				NewAsssertionTestCase(float64(1), float64(1)),
				NewAsssertionTestCase(int(5), float64(1)),
				NewAsssertionTestCase(int(5), float64(5)),
				NewAsssertionTestCase("2.2", float64(1)),
				NewAsssertionTestCase("2.2", float64(2.2)),
				NewAsssertionTestCase(int(5), int(1)),
				NewAsssertionTestCase(int(5), int(5)),
				NewAsssertionTestCase("5", int(1)),
				NewAsssertionTestCase("1", int(1)),
				NewAsssertionTestCase("abc", "a"),
				NewAsssertionTestCase("abc", "abc"),
				NewAsssertionTestCase(float64(1.3), "1.2"),
				NewAsssertionTestCase(float64(1.3), "1.3"),
				NewAsssertionTestCase(int(3), "1"),
				NewAsssertionTestCase(int(1), "1"),
				NewAsssertionTestCase("3.1", "2"),
				NewAsssertionTestCase("2", "2"),
			}

			assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
				key := "some:key"
				executionResult := core.ExecutionResult{
					key: actualValue,
				}

				assertion := GreaterThanOrEqualAssertion{
					Key:   key,
					Value: instanceValue,
				}

				result := assertion.Assert(executionResult)
				So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
				So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
			})
		})

		Convey("Fails", func() {
			var testCases = []AssertionTestCase{
				NewAsssertionTestCase(NilValue, NilValue),
				NewAsssertionTestCase(NilValue, int(5)),
				NewAsssertionTestCase(NilValue, "5.1"),
				NewAsssertionTestCase(NilValue, "fubar"),
				NewAsssertionTestCase(NilValue, float64(6.1)),
				NewAsssertionTestCase(float64(5.1), float64(6.1)),
				NewAsssertionTestCase(int(5), float64(6.1)),
				NewAsssertionTestCase("5", float64(6.1)),
				NewAsssertionTestCase(float64(3.1), int(6)),
				NewAsssertionTestCase(int(3), int(6)),
				NewAsssertionTestCase("5.1", int(6)),
				NewAsssertionTestCase("fubar", int(6)),
				NewAsssertionTestCase("fubar", float64(1.1)),
				NewAsssertionTestCase("fubar", "1.1"),
				NewAsssertionTestCase("f", "fubar"),
				NewAsssertionTestCase(float64(6.1), "fubar"),
				NewAsssertionTestCase(int(6), "fubar"),
			}

			assert(testCases, func(actualValue interface{}, instanceValue interface{}) {
				executionResult := core.ExecutionResult{
					key: actualValue,
				}

				assertion := GreaterThanOrEqualAssertion{
					Key:   key,
					Value: instanceValue,
				}

				result := assertion.Assert(executionResult)
				So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
				So(result[core.AssertionMessageUrn.String()], ShouldEqual, fmt.Sprintf("FAIL: %v is not greater %v", actualValue, instanceValue))
			})
		})
	})
}
