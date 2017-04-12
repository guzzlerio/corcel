package core_test

import (
	"testing"

	. "github.com/guzzlerio/corcel/core"
	"github.com/guzzlerio/corcel/infrastructure/http"
	"github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegistry(t *testing.T) {
	Convey("Registry", t, func() {

		Convey("Creates a Registry", func() {
			var registry = CreateRegistry()

			So(len(registry.AssertionParsers), ShouldEqual, 0)
			So(len(registry.ActionParsers), ShouldEqual, 0)
			So(len(registry.ResultProcessors), ShouldEqual, 0)
			So(len(registry.ExtractorParsers), ShouldEqual, 0)
		})

		Convey("AddExtractorParser", func() {
			var registry = CreateRegistry().AddExtractorParser(yaml.RegexExtractorParser{})
			So(len(registry.AssertionParsers), ShouldEqual, 0)
			So(len(registry.ActionParsers), ShouldEqual, 0)
			So(len(registry.ResultProcessors), ShouldEqual, 0)
			So(len(registry.ExtractorParsers), ShouldEqual, 1)
		})

		Convey("AddAssertionParser", func() {
			var registry = CreateRegistry().AddAssertionParser(yaml.ExactAssertionParser{})
			So(len(registry.AssertionParsers), ShouldEqual, 1)
			So(len(registry.ActionParsers), ShouldEqual, 0)
			So(len(registry.ResultProcessors), ShouldEqual, 0)
			So(len(registry.ExtractorParsers), ShouldEqual, 0)
		})

		Convey("AddActionParser", func() {
			var registry = CreateRegistry().AddActionParser(http.YamlHTTPRequestParser{})
			So(len(registry.AssertionParsers), ShouldEqual, 0)
			So(len(registry.ActionParsers), ShouldEqual, 1)
			So(len(registry.ResultProcessors), ShouldEqual, 0)
			So(len(registry.ExtractorParsers), ShouldEqual, 0)
		})
	})
}
