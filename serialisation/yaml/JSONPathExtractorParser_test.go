package yaml

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONPathExtractorParser(t *testing.T) {
	Convey("JSONPathExtractorParser", t, func() {
		Convey("Parses", func() {
			var expectedName = "A"
			var expectedKey = "B"
			var expectedJsonPath = "C"

			var input = map[string]interface{}{
				"name":     expectedName,
				"key":      expectedKey,
				"jsonpath": expectedJsonPath,
				"scope":    core.PlanScope,
			}

			var parser = JSONPathExtractorParser{}

			var result, err = parser.Parse(input)

			So(err, ShouldBeNil)

			So(result.(extractors.JSONPathExtractor), ShouldNotBeNil)

			var jsonPathExtractor = result.(extractors.JSONPathExtractor)

			So(jsonPathExtractor.Name, ShouldEqual, expectedName)
			So(jsonPathExtractor.Key, ShouldEqual, expectedKey)
			So(jsonPathExtractor.JSONPath, ShouldEqual, expectedJsonPath)
		})

		Convey("Fails to parse when name not set", func() {

			var input = map[string]interface{}{
				"key":      "key",
				"jsonpath": "path",
				"scope":    core.PlanScope,
			}

			var parser = JSONPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse when key not set", func() {

			var input = map[string]interface{}{
				"name":     "name",
				"jsonpath": "path",
				"scope":    core.PlanScope,
			}

			var parser = JSONPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse when jsonpath not set", func() {

			var input = map[string]interface{}{
				"name":  "name",
				"key":   "key",
				"scope": core.PlanScope,
			}

			var parser = JSONPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})
	})
}
