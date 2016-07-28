package processor_test

import (
	. "ci.guzzler.io/guzzler/corcel/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("ListRingIterator", func() {

	It("Loops round", func() {
		data := map[string][]map[string]interface{}{}
		data["People"] = getPeopleCollection()

		iterator := NewListRingIterator(data)

		values1 := iterator.Values()
		Expect(values1["$People.name"]).To(Equal("bob"))
		Expect(values1["$People.age"]).To(Equal(30))

		values2 := iterator.Values()
		Expect(values2["$People.name"]).To(Equal("carol"))
		Expect(values2["$People.age"]).To(Equal(31))

		values3 := iterator.Values()
		Expect(values3["$People.name"]).To(Equal("alice"))
		Expect(values3["$People.age"]).To(Equal(32))
	})

	It("Loops around uneven lists", func() {
		data := map[string][]map[string]interface{}{}
		data["People"] = getPeopleCollection()
		data["Products"] = getProductsCollection()

		iterator := NewListRingIterator(data)

		values1 := iterator.Values()
		Expect(values1["$People.name"]).To(Equal("bob"))
		Expect(values1["$People.age"]).To(Equal(30))
		Expect(values1["$Products.name"]).To(Equal("toaster"))
		Expect(values1["$Products.sku"]).To(Equal("1234"))

		values2 := iterator.Values()
		Expect(values2["$People.name"]).To(Equal("carol"))
		Expect(values2["$People.age"]).To(Equal(31))
		Expect(values2["$Products.name"]).To(Equal("grinder"))
		Expect(values2["$Products.sku"]).To(Equal("5678"))

		values3 := iterator.Values()
		Expect(values3["$People.name"]).To(Equal("alice"))
		Expect(values3["$People.age"]).To(Equal(32))
		Expect(values3["$Products.name"]).To(Equal("toaster"))
		Expect(values3["$Products.sku"]).To(Equal("1234"))
	})

})
