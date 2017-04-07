package assertions

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExactAssertion(t *testing.T) {
	Convey("ExactAssertion", t, func() {
		key := "some:key"

		Convey("Exact Assertion Succeeds", func() {
			expectedValue := 7

			executionResult := core.ExecutionResult{
				key: expectedValue,
			}

			assertion := ExactAssertion{
				Key:   key,
				Value: expectedValue,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})

		Convey("Exact Assertion Fails", func() {
			expectedValue := 7

			executionResult := core.ExecutionResult{
				key: 8,
			}

			assertion := ExactAssertion{
				Key:   key,
				Value: expectedValue,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: 8 int does not match 7 int")
		})

		//NOTHING is currently using the message when an assertion fails but we will need
		//it for when we put the errors into the report.  One of the edge cases with the message
		//is that say the actual value was a string "7" and the expected is an int 7.  The message
		//will not include the quotes so the message would read 7 does not equal 7 as opposed
		//to "7" does not equal 7.  Notice this is a type mismatch
		/*
			PConvey("Exact Assertion Fails when actual and expected are different types", func() {
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
				So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
				So(result[core.AssertionMessageUrn.String()]).To(Equal("FAIL: \"7\" does not match 7"))
			})
		*/
	})
}
