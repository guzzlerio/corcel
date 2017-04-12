package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExactAssertionParser(t *testing.T) {
	Convey("ExactAssertionParser", t, func() {

		Convey("Parses", func() {

			var expectedKey = "talula"
			var expectedValue = "boomboom"
			var input = map[string]interface{}{
				"key":      expectedKey,
				"expected": expectedValue,
			}

			var parser = ExactAssertionParser{}
			assertion, err := parser.Parse(input)
			var exactAssertion = assertion.(*assertions.ExactAssertion)

			So(err, ShouldBeNil)
			So(exactAssertion.Key, ShouldEqual, expectedKey)
			So(exactAssertion.Value, ShouldEqual, expectedValue)
		})

		Convey("Fails to parse without key", func() {

			var input = map[string]interface{}{
				"bang":     "talula",
				"expected": "boomboom",
			}

			var parser = ExactAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})

		Convey("Fails to parse without expected", func() {

			var input = map[string]interface{}{
				"key":  "talula",
				"bang": "boomboom",
			}

			var parser = ExactAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "expected is not present")
		})

	})
}
