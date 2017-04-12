package assertions

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNotEqualAssertion(t *testing.T) {
	Convey("NotEqualAssertion", t, func() {
		key := "some:key"

		Convey("Succeeds", func() {
			executionResult := core.ExecutionResult{
				key: 8,
			}

			assertion := NotEqualAssertion{
				Key:   key,
				Value: 7,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, true)
			So(result[core.AssertionMessageUrn.String()], ShouldBeNil)
		})

		Convey("Fails", func() {
			executionResult := core.ExecutionResult{
				key: 7,
			}

			assertion := NotEqualAssertion{
				Key:   key,
				Value: 7,
			}

			result := assertion.Assert(executionResult)
			So(result[core.AssertionResultUrn.String()], ShouldEqual, false)
			So(result[core.AssertionMessageUrn.String()], ShouldEqual, "FAIL: 7 does match 7")
		})
	})
}
