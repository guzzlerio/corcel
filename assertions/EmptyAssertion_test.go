package assertions

import (
	"testing"

	"github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEmptyAssertion(t *testing.T) {
	Convey("EmptyAssertion", t, func() {

		key := "some:key"

		Convey("Succeeds when empty string", func() {
			executionResult := core.ExecutionResult{
				key: "",
			}

			assertion := EmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})

		Convey("Succeeds when empty string of whitespace", func() {
			executionResult := core.ExecutionResult{
				key: "    ",
			}

			assertion := EmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})

		Convey("Succeeds when nil", func() {
			executionResult := core.ExecutionResult{}

			assertion := EmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})

		Convey("Fails when value is not nil", func() {

			executionResult := core.ExecutionResult{
				key: 8,
			}

			assertion := EmptyAssertion{
				Key: key,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: value is not empty")
		})
	})
}
