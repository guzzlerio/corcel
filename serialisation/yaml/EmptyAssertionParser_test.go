package yaml

import (
	"fmt"
	"testing"

	"github.com/guzzlerio/corcel/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEmptyAssertionParser(t *testing.T) {
	Convey("EmptyAssertionParser", t, func() {

		Convey("Parses", func() {

			var expected = "talula"
			var input = map[string]interface{}{
				"key": expected,
			}

			var parser = EmptyAssertionParser{}
			assertion, err := parser.Parse(input)
			emptyAssertion := assertion.(*assertions.EmptyAssertion)
			So(err, ShouldBeNil)

			So(emptyAssertion.Key, ShouldEqual, expected)
		})

		Convey("Fails to parse", func() {
			var expected = "talula"
			var input = map[string]interface{}{
				"bang": expected,
			}

			var parser = EmptyAssertionParser{}
			_, err := parser.Parse(input)

			So(err, ShouldNotBeNil)
			So(fmt.Sprintf("%v", err), ShouldContainSubstring, "key is not present")
		})
	})
}
