package report

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUrnComposite(t *testing.T) {
	Convey("UrnComposite", t, func() {

		Convey("Combines two urns into a composite", func() {
			urn1 := "urn:a:b:c"
			urn2 := "urn:a:b:d"

			composite, err := createUrnComposite(urn1, urn2)

			So(err, ShouldBeNil)
			So(composite.Name, ShouldEqual, "urn")
			So(len(composite.Children), ShouldEqual, 1)
			So(len(composite.Child(0).Child(0).Children), ShouldEqual, 2)
		})

		Convey("Detects multiple root elements", func() {

			urn1 := "urn:a:b:c"
			urn2 := "uri:a:b:d"

			_, err := createUrnComposite(urn1, urn2)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldContainSubstring, "Multiple root elements")
		})

		Convey("Combines more than two urns into a composite", func() {
			urn1 := "urn:a:b:c"
			urn2 := "urn:a:b:d"
			urn3 := "urn:a:c:a"
			urn4 := "urn:a:b:e"

			composite, err := createUrnComposite(urn1, urn2, urn3, urn4)

			So(err, ShouldBeNil)
			So(composite.Name, ShouldEqual, "urn")
			So(len(composite.Children), ShouldEqual, 1)
			So(len(composite.Child(0).Children), ShouldEqual, 2)
			So(len(composite.Child(0).Child(0).Children), ShouldEqual, 3)
		})

		Convey("Can locate the root", func() {
			urn1 := "urn:a:b:c"
			composite, _ := createUrnComposite(urn1)
			So(composite.Child(0).Child(0).Child(0).Root().Name, ShouldEqual, "urn")

		})

		Convey("Can add a value to the node", func() {
			urn1 := "urn:a:b"
			urn2 := "urn:a:b:d"
			urn3 := "urn:a:b:e"
			urn4 := "urn:a:b:d:f"
			value := []int{1, 2, 3, 4, 5, 6}

			composite, _ := createUrnComposite(urn1)
			composite.AddValue(urn2, value)
			composite.AddValue(urn3, value)
			composite.AddValue(urn4, value)
			So(composite.Child(0).Child(0).Child(0).Value, ShouldResemble, value)
			So(composite.Child(0).Child(0).Child(1).Value, ShouldResemble, value)
			So(composite.Child(0).Child(0).Child(0).Child(0).Value, ShouldResemble, value)
		})

		Convey("Can report its depth", func() {
			urn1 := "urn:a:b"
			composite, _ := createUrnComposite(urn1)
			So(composite.Child(0).Depth(), ShouldEqual, 1)
			So(composite.Child(0).Child(0).Depth(), ShouldEqual, 2)
		})

		Convey("Can report its metric type", func() {
			urn1 := "urn:action:counter:a"
			urn2 := "urn:action:meter:b"
			composite, _ := createUrnComposite(urn1, urn2)
			firstChild, _ := composite.Child(0).Child(0).Child(0).MetricType()
			secondChild, _ := composite.Child(0).Child(1).Child(0).MetricType()
			So(firstChild, ShouldEqual, "counter")
			So(secondChild, ShouldEqual, "meter")
		})

		Convey("Cannot report its metric type if Depth < 2 and multiple metric types exist", func() {
			urn1 := "urn:action:counter:a"
			urn2 := "urn:action:meter:b"
			composite, _ := createUrnComposite(urn1, urn2)
			_, err := composite.MetricType()
			So(err.Error(), ShouldContainSubstring, "Possible multiple metric types")
		})

		Convey("Can report its connector", func() {
			urn1 := "urn:http:counter:a"
			urn2 := "urn:http:meter:b"
			composite, _ := createUrnComposite(urn1, urn2)
			firstChild, _ := composite.Child(0).Child(0).Child(0).Connector()
			secondChild, _ := composite.Child(0).Child(1).Child(0).Connector()
			So(firstChild, ShouldEqual, "http")
			So(secondChild, ShouldEqual, "http")
		})

		Convey("Cannot report its connector", func() {
			urn1 := "urn:action:counter:a"
			urn2 := "urn:action:meter:b"
			composite, _ := createUrnComposite(urn1, urn2)
			_, err := composite.Connector()
			So(err.Error(), ShouldContainSubstring, "Possible multiple connector types")
		})

		Convey("Can render", func() {
			urn1 := "urn"
			urn2 := "urn:http:counter:d"
			urn3 := "urn:http:counter:e"
			urn4 := "urn:http:meter:e"
			urn5 := "urn:http:counter:d:a"
			value := []int64{1, 2, 3, 4, 5, 6}

			composite, _ := createUrnComposite(urn1)
			composite.AddValue(urn2, value)
			composite.AddValue(urn3, value)
			composite.AddValue(urn4, value)
			composite.AddValue(urn5, value)

			registry := NewRendererRegistry()
			registry.Add("counter", RenderCounter)

			result := composite.Render(registry, []int64{1, 2, 3, 4, 5, 6})

			So(result, ShouldNotBeNil)
		})

		Convey("Can report all connectors", func() {
			urn1 := "urn:http:counter:d"
			urn2 := "urn:action:counter:e"

			composite, _ := createUrnComposite(urn1, urn2)

			So(composite.Connectors(), ShouldResemble, []string{"http", "action"})
		})

	})
}
