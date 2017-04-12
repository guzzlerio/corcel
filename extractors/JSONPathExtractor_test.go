package extractors

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONPathExtractor(t *testing.T) {
	Convey("JSONPathExtractor", t, func() {

		Convey("extracts a single key", func() {
			var extractor = JSONPathExtractor{
				Name:     "something",
				Key:      "targetKey",
				JSONPath: "$.aKey",
			}

			var result = core.ExecutionResult{
				"targetKey": `{"aKey":32}`,
			}

			var extractionResult = extractor.Extract(result)

			So(extractionResult["something"], ShouldEqual, float64(32))
		})

		Convey("invalid json path", func() {
			var extractor = JSONPathExtractor{
				Name:     "something",
				Key:      "targetKey",
				JSONPath: "talula",
			}

			var result = core.ExecutionResult{
				"targetKey": `{"aKey":32}`,
			}

			var extractionResult = extractor.Extract(result)

			So(extractionResult["something"], ShouldEqual, ErrInvalidJsonPath)
		})
	})
}
