package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGreaterThanAssertionParser(t *testing.T) {
	Convey("GreaterThanAssertionParser", t, func() {

		Convey("Parses", func() {

			var expectedKey = "talula"
			var input = map[string]interface{}{
				"key":      expectedKey,
				"expected": 7,
			}

			var parser = GreaterThanAssertionParser{}
			assertion, err := parser.Parse(input)
			var gtAssertion = assertion.(*assertions.GreaterThanAssertion)

			So(err, ShouldBeNil)
			So(gtAssertion.Key, ShouldEqual, expectedKey)
			So(gtAssertion.Value, ShouldEqual, 7)
		})

		Convey("Returns Key", func() {
			So(GreaterThanAssertionParser{}.Key(), ShouldEqual, "GreaterThanAssertion")
		})

		Convey("Fails to parse without key", func() {
			var input = map[string]interface{}{
				"boom":     "talula",
				"expected": 7,
			}

			var parser = GreaterThanAssertionParser{}
			_, err := parser.Parse(input)
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})

		Convey("Fails to parse without expected", func() {
			var input = map[string]interface{}{
				"key":  "talula",
				"boom": 7,
			}

			var parser = GreaterThanAssertionParser{}
			_, err := parser.Parse(input)
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "expected is not present")
		})
	})
}
