package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLessThanAssertionParser(t *testing.T) {
	Convey("LessThanAssertionParser", t, func() {
		Convey("Parses", func() {

			var expectedKey = "talula"
			var input = map[string]interface{}{
				"key":      expectedKey,
				"expected": 7,
			}

			var parser = LessThanAssertionParser{}
			assertion, err := parser.Parse(input)
			var ltAssertion = assertion.(*assertions.LessThanAssertion)

			So(err, ShouldBeNil)
			So(ltAssertion.Key, ShouldEqual, expectedKey)
			So(ltAssertion.Value, ShouldEqual, 7)
		})

		Convey("Returns Key", func() {
			So(LessThanAssertionParser{}.Key(), ShouldEqual, "LessThanAssertion")
		})

		Convey("Fails to parse without key", func() {

			var input = map[string]interface{}{
				"bang":     "talula",
				"expected": "boomboom",
			}

			var parser = LessThanAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})

		Convey("Fails to parse without expected", func() {

			var input = map[string]interface{}{
				"key":  "talula",
				"bang": "boomboom",
			}

			var parser = LessThanAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "expected is not present")
		})

	})
}
