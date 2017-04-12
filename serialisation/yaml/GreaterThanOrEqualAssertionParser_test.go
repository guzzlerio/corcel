package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGreaterThanOrEqualAssertionParser(t *testing.T) {
	Convey("GreaterThanOrEqualAssertionParser", t, func() {
		Convey("Parses", func() {

			var expectedKey = "talula"
			var input = map[string]interface{}{
				"key":      expectedKey,
				"expected": 7,
			}

			var parser = GreaterThanOrEqualAssertionParser{}
			assertion, err := parser.Parse(input)
			var gteAssertion = assertion.(*assertions.GreaterThanOrEqualAssertion)

			So(err, ShouldBeNil)
			So(gteAssertion.Key, ShouldEqual, expectedKey)
			So(gteAssertion.Value, ShouldEqual, 7)
		})

		Convey("Returns Key", func() {
			So(GreaterThanOrEqualAssertionParser{}.Key(), ShouldEqual, "GreaterThanOrEqualAssertion")
		})

		Convey("Fails to parse without key", func() {
			var input = map[string]interface{}{
				"boom":     "talula",
				"expected": 7,
			}

			var parser = GreaterThanOrEqualAssertionParser{}
			_, err := parser.Parse(input)
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})

		Convey("Fails to parse without expected", func() {
			var input = map[string]interface{}{
				"key":  "talula",
				"boom": 7,
			}

			var parser = GreaterThanOrEqualAssertionParser{}
			_, err := parser.Parse(input)
			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "expected is not present")
		})
	})
}
