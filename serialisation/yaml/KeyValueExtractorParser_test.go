package yaml_test

import (
	"testing"

	"github.com/guzzlerio/corcel/extractors"
	. "github.com/guzzlerio/corcel/serialisation/yaml"

	. "github.com/smartystreets/goconvey/convey"
)

func TestKeyValueExtractorParser(t *testing.T) {
	Convey("KeyValueExtractorParser", t, func() {
		Convey("Parses", func() {
			var expectedName = "A"
			var expectedKey = "B"

			var input = map[string]interface{}{
				"name": expectedName,
				"key":  expectedKey,
			}

			var parser = KeyValueExtractorParser{}

			var result, err = parser.Parse(input)

			So(err, ShouldBeNil)

			So(result.(extractors.KeyValueExtractor), ShouldNotBeNil)

			var keyValueExtractor = result.(extractors.KeyValueExtractor)

			So(keyValueExtractor.Name, ShouldEqual, expectedName)
			So(keyValueExtractor.Key, ShouldEqual, expectedKey)
		})

		Convey("Fails to parse when name not set", func() {

			var input = map[string]interface{}{
				"key": "key",
			}

			var parser = KeyValueExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})

		Convey("Fails to parse when key not set", func() {

			var input = map[string]interface{}{
				"name": "name",
			}

			var parser = KeyValueExtractorParser{}

			var _, err = parser.Parse(input)

			So(err, ShouldNotBeNil)
		})
	})
}
