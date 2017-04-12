package extractors

import (
	"testing"

	"github.com/guzzlerio/corcel/core"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeyValueExtractor(t *testing.T) {
	Convey("KeyValueExtractor", t, func() {

		Convey("Succeeds when key is present", func() {
			var extractor = KeyValueExtractor{
				Key:   "key",
				Name:  "target",
				Scope: core.StepScope,
			}

			var executionResult = core.ExecutionResult{
				"key": "talula",
			}

			var extractionResult = extractor.Extract(executionResult)

			So(extractionResult["target"], ShouldEqual, "talula")
			So(extractionResult["scope"], ShouldEqual, core.StepScope)
		})

		Convey("Extraction result does not contain the name when the key does not exist inside the execution result", func() {
			var extractor = KeyValueExtractor{
				Key:   "key",
				Name:  "target",
				Scope: core.StepScope,
			}

			var executionResult = core.ExecutionResult{
				"lock": "talula",
			}

			var extractionResult = extractor.Extract(executionResult)

			So(extractionResult["target"], ShouldBeNil)
		})

	})
}
