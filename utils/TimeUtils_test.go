package utils_test

import (
	"testing"
	"time"

	. "github.com/guzzlerio/corcel/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTimeUtils(t *testing.T) {
	Convey("TimeUtils", t, func() {
		Convey("Time a function", func() {
			var result = Time(func() {})

			So(result, ShouldBeGreaterThan, time.Duration(1))
		})

		Convey("DurationIsBetween succeeds", func() {
			var a = time.Duration(1)
			var b = time.Duration(2)
			var c = time.Duration(3)

			So(DurationIsBetween(b, a, c), ShouldBeTrue)
		})

		Convey("DurationIsBetween fails", func() {
			var a = time.Duration(1)
			var b = time.Duration(2)
			var c = time.Duration(3)

			So(DurationIsBetween(a, b, c), ShouldBeFalse)
		})
	})
}
