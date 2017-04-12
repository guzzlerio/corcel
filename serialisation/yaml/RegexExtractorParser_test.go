package yaml_test

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegexExtractorParser(t *testing.T) {
	Convey("RegexExtractorParser", t, func() {
		Convey("Parses", func() {
			var expectedName = "A"
			var expectedKey = "B"
			var expectedMatch = "C"

			var input = map[string]interface{}{
				"name":  expectedName,
				"key":   expectedKey,
				"match": expectedMatch,
				"scope": core.PlanScope,
			}

			var parser = RegexExtractorParser{}

			var result, err = parser.Parse(input)

			So(err, ShouldBeNil)

			So(result.(extractors.RegexExtractor), ShouldNotBeNil)

			var jsonPathExtractor = result.(extractors.RegexExtractor)

			So(jsonPathExtractor.Name, ShouldEqual, expectedName)
			So(jsonPathExtractor.Key, ShouldEqual, expectedKey)
			So(jsonPathExtractor.Match, ShouldEqual, expectedMatch)
		})

		Convey("Fails to parse without name", func() {

			var input = map[string]interface{}{
				"key":   "key",
				"match": "match",
				"scope": core.PlanScope,
			}

			var parser = RegexExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse without key", func() {

			var input = map[string]interface{}{
				"name":  "name",
				"match": "match",
				"scope": core.PlanScope,
			}

			var parser = RegexExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse without match", func() {

			var input = map[string]interface{}{
				"name":  "name",
				"key":   "key",
				"scope": core.PlanScope,
			}

			var parser = RegexExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})
	})
}
