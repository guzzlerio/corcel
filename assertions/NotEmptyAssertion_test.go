package assertions

import (
	"testing"

	"github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNotEmptyAssertion(t *testing.T) {
	Convey("NotEmptyAssertion", t, func() {

		key := "some:key"

		Convey("Fails when empty string", func() {
			executionResult := core.ExecutionResult{
				key: "",
			}

			assertion := NotEmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: value is empty")
		})

		Convey("Fails when empty string of whitespace", func() {
			executionResult := core.ExecutionResult{
				key: "    ",
			}

			assertion := NotEmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: value is empty")
		})

		Convey("Fails when nil", func() {
			executionResult := core.ExecutionResult{}

			assertion := NotEmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: value is empty")
		})

		Convey("Succeeds when value is not nil", func() {

			executionResult := core.ExecutionResult{
				key: 8,
			}

			assertion := NotEmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})
	})
}
