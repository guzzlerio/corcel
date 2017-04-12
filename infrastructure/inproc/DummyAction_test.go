package inproc_test

import (
	"context"
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/guzzlerio/corcel/infrastructure/inproc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDummyAction(t *testing.T) {
	Convey("DummyAction", t, func() {

		SkipConvey("Execute with string property", func() {
			var action = DummyAction{
				Results: map[string]interface{}{
					"Data": "$value",
				},
			}

			var executionContext = core.ExecutionContext{
				"vars": core.ExecutionContext{
					"$value": "talula",
				},
			}

			var result = action.Execute(context.TODO(), executionContext)

			So(result["Data"], ShouldEqual, "talula")
		})

		SkipConvey("Execute with float64 poerty", func() {
			var action = DummyAction{
				Results: map[string]interface{}{
					"Data": "$value",
				},
			}

			var executionContext = core.ExecutionContext{
				"vars": core.ExecutionContext{
					"$value": float64(100),
				},
			}

			var result = action.Execute(context.TODO(), executionContext)

			So(result["Data"], ShouldEqual, "100")

		})
	})
}
