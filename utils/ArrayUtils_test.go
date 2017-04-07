package utils_test

import (
	"testing"

	. "github.com/guzzlerio/corcel/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestArrayUtils(t *testing.T) {
	Convey("ArrayUtils", t, func() {
		var input = []string{"A", "B", "C"}

		Convey("ContainsString suceeds", func() {
			So(ContainsString(input, "B"), ShouldBeTrue)
		})

		Convey("ContainsString fails", func() {
			So(ContainsString(input, "D"), ShouldBeFalse)
		})
	})
}
