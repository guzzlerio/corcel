package processor_test

import (
	. "ci.guzzler.io/guzzler/corcel/processor"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListRingIterator", func() {

	It("Loops round", func() {
		peopleKey := "People"
		data := map[string][]map[string]interface{}{}
		data[peopleKey] = []map[string]interface{}{}

		bob := map[string]interface{}{
			"name": "bob",
			"age":  30,
		}
		carol := map[string]interface{}{
			"name": "carol",
			"age":  31,
		}
		alice := map[string]interface{}{
			"name": "alice",
			"age":  32,
		}

		data[peopleKey] = append(data[peopleKey], bob)
		data[peopleKey] = append(data[peopleKey], carol)
		data[peopleKey] = append(data[peopleKey], alice)

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

})
