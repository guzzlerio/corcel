package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUrn(t *testing.T) {
	Convey("Urn", t, func() {
		Convey("can create a counter", func() {
			urn := NewUrn("http").Counter().Name("status", "all").Name("200").String()
			So(urn, ShouldEqual, "urn:http:counter:status:all:200")
		})

		Convey("can create a gauge", func() {
			urn := NewUrn("http").Gauge().Name("status", "all").Name("200").String()
			So(urn, ShouldEqual, "urn:http:gauge:status:all:200")
		})

		Convey("can create a meter", func() {
			urn := NewUrn("http").Meter().Name("status", "all").Name("200").String()
			So(urn, ShouldEqual, "urn:http:meter:status:all:200")
		})

		Convey("can create a timer", func() {
			urn := NewUrn("http").Timer().Name("status", "all").Name("200").String()
			So(urn, ShouldEqual, "urn:http:timer:status:all:200")
		})

		Convey("can create a histogram", func() {
			urn := NewUrn("http").Timer().Name("status", "all").Name("200").String()
			So(urn, ShouldEqual, "urn:http:timer:status:all:200")
		})

		Convey("ignores the metric is not specified", func() {
			urn := NewUrn("http").Name("status").Name("200").String()
			So(urn, ShouldEqual, "urn:http:status:200")
		})

		Convey("builds the urn in lowercase", func() {

		})
	})
}
