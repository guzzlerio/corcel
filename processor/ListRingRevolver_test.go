package processor_test

import (
	"testing"

	. "github.com/guzzlerio/corcel/processor"

	. "github.com/smartystreets/goconvey/convey"
)

func getPeopleCollection() []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"name": "bob",
			"age":  30,
		},
		map[string]interface{}{
			"name": "carol",
			"age":  31,
		},
		map[string]interface{}{
			"name": "alice",
			"age":  32,
		},
	}
}

func getProductsCollection() []map[string]interface{} {
	return []map[string]interface{}{
		map[string]interface{}{
			"name": "toaster",
			"sku":  "1234",
		},
		map[string]interface{}{
			"name": "grinder",
			"sku":  "5678",
		},
	}
}

func TestListRingIterator(t *testing.T) {
	BeforeTest()
	defer AfterTest()
	Convey("ListRingIterator", t, func() {

		Convey("Loops round", func() {
			data := map[string][]map[string]interface{}{}
			data["People"] = getPeopleCollection()

			iterator := NewListRingRevolver(data)

			values1 := iterator.Values()
			So(values1["People.name"], ShouldEqual, "bob")
			So(values1["People.age"], ShouldEqual, 30)

			values2 := iterator.Values()
			So(values2["People.name"], ShouldEqual, "carol")
			So(values2["People.age"], ShouldEqual, 31)

			values3 := iterator.Values()
			So(values3["People.name"], ShouldEqual, "alice")
			So(values3["People.age"], ShouldEqual, 32)
		})

		Convey("Loops around uneven lists", func() {
			data := map[string][]map[string]interface{}{}
			data["People"] = getPeopleCollection()
			data["Products"] = getProductsCollection()

			iterator := NewListRingRevolver(data)

			values1 := iterator.Values()
			So(values1["People.name"], ShouldEqual, "bob")
			So(values1["People.age"], ShouldEqual, 30)
			So(values1["Products.name"], ShouldEqual, "toaster")
			So(values1["Products.sku"], ShouldEqual, "1234")

			values2 := iterator.Values()
			So(values2["People.name"], ShouldEqual, "carol")
			So(values2["People.age"], ShouldEqual, 31)
			So(values2["Products.name"], ShouldEqual, "grinder")
			So(values2["Products.sku"], ShouldEqual, "5678")

			values3 := iterator.Values()
			So(values3["People.name"], ShouldEqual, "alice")
			So(values3["People.age"], ShouldEqual, 32)
			So(values3["Products.name"], ShouldEqual, "toaster")
			So(values3["Products.sku"], ShouldEqual, "1234")
		})

	})
}
