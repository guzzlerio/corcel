package yaml

import (
	"testing"

	"github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/extractors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestXPathExtractorParser(t *testing.T) {
	Convey("XPathExtractorParser", t, func() {
		Convey("Parses", func() {
			var expectedName = "A"
			var expectedKey = "B"
			var expectedXpath = "C"

			var input = map[string]interface{}{
				"name":  expectedName,
				"key":   expectedKey,
				"xpath": expectedXpath,
				"scope": core.PlanScope,
			}

			var parser = XPathExtractorParser{}

			var result, err = parser.Parse(input)

			So(err, ShouldBeNil)

			So(result.(extractors.XPathExtractor), ShouldNotBeNil)

			var xpathExtractor = result.(extractors.XPathExtractor)
			So(xpathExtractor.Key, ShouldEqual, expectedKey)
			So(xpathExtractor.Name, ShouldEqual, expectedName)
			So(xpathExtractor.XPath, ShouldEqual, expectedXpath)
			So(xpathExtractor.Scope, ShouldEqual, core.PlanScope)
		})

		Convey("Fails to parse with empty name", func() {
			var input = map[string]interface{}{
				"key":   "key",
				"xpath": "path",
				"scope": core.PlanScope,
			}
			var parser = XPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse with empty key", func() {
			var input = map[string]interface{}{
				"name":  "name",
				"xpath": "path",
				"scope": core.PlanScope,
			}
			var parser = XPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse with empty xpath", func() {
			var input = map[string]interface{}{
				"key":   "key",
				"name":  "name",
				"scope": core.PlanScope,
			}
			var parser = XPathExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})
	})
}
