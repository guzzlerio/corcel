package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNotEqualAssertionParser(t *testing.T) {
	Convey("NotEqualAssertionParser", t, func() {

		Convey("Parses", func() {

			var expectedKey = "talula"
			var expectedValue = "boomboom"
			var input = map[string]interface{}{
				"key":      expectedKey,
				"expected": expectedValue,
			}

			var parser = NotEqualAssertionParser{}
			assertion, err := parser.Parse(input)
			var exactAssertion = assertion.(*assertions.NotEqualAssertion)

			So(err, ShouldBeNil)
			So(exactAssertion.Key, ShouldEqual, expectedKey)
			So(exactAssertion.Value, ShouldEqual, expectedValue)
		})

		Convey("Fails to parse without key", func() {

			var input = map[string]interface{}{
				"bang":     "talula",
				"expected": "boomboom",
			}

			var parser = NotEqualAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})

		Convey("Fails to parse without expected", func() {

			var input = map[string]interface{}{
				"key":  "talula",
				"bang": "boomboom",
			}

			var parser = NotEqualAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "expected is not present")
		})

	})
}
