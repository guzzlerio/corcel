package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNotEmptyAssertionParser(t *testing.T) {
	Convey("NotEmptyAssertionParser", t, func() {

		Convey("Parses", func() {

			var expected = "talula"
			var input = map[string]interface{}{
				"key": expected,
			}

			var parser = NotEmptyAssertionParser{}
			assertion, err := parser.Parse(input)
			emptyAssertion := assertion.(*assertions.NotEmptyAssertion)
			So(err, ShouldBeNil)

			So(emptyAssertion.Key, ShouldEqual, expected)
		})

		Convey("Fails to parse", func() {
			var expected = "talula"
			var input = map[string]interface{}{
				"bang": expected,
			}

			var parser = NotEmptyAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})
	})
}
