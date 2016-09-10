package report

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UrnComposite", func() {

	It("Combines two urns into a composite", func() {
		urn1 := "urn:a:b:c"
		urn2 := "urn:a:b:d"

		composite, err := createUrnComposite(urn1, urn2)

		Expect(err).To(BeNil())
		Expect(composite.Name).To(Equal("urn"))
		Expect(len(composite.Children)).To(Equal(1))
		Expect(len(composite.Child(0).Child(0).Children)).To(Equal(2))
	})

	It("Detects multiple root elements", func() {

		urn1 := "urn:a:b:c"
		urn2 := "uri:a:b:d"

		_, err := createUrnComposite(urn1, urn2)

		Expect(err).ToNot(BeNil())
		Expect(err).To(MatchError("Multiple root elements"))
	})

	It("Combines more than two urns into a composite", func() {
		urn1 := "urn:a:b:c"
		urn2 := "urn:a:b:d"
		urn3 := "urn:a:c:a"
		urn4 := "urn:a:b:e"

		composite, err := createUrnComposite(urn1, urn2, urn3, urn4)

		Expect(err).To(BeNil())
		Expect(composite.Name).To(Equal("urn"))
		Expect(len(composite.Children)).To(Equal(1))
		Expect(len(composite.Child(0).Children)).To(Equal(2))
		Expect(len(composite.Child(0).Child(0).Children)).To(Equal(3))
	})

	It("Can locate the root", func() {
		urn1 := "urn:a:b:c"
		composite, _ := createUrnComposite(urn1)
		Expect(composite.Child(0).Child(0).Child(0).Root().Name).To(Equal("urn"))

	})

	It("Can add a value to the node", func() {
		urn1 := "urn:a:b"
		urn2 := "urn:a:b:d"
		urn3 := "urn:a:b:e"
		urn4 := "urn:a:b:d:f"
		value := []int{1, 2, 3, 4, 5, 6}

		composite, _ := createUrnComposite(urn1)
		composite.AddValue(urn2, value)
		composite.AddValue(urn3, value)
		composite.AddValue(urn4, value)
		Expect(composite.Child(0).Child(0).Child(0).Value).To(Equal(value))
		Expect(composite.Child(0).Child(0).Child(1).Value).To(Equal(value))
		Expect(composite.Child(0).Child(0).Child(0).Child(0).Value).To(Equal(value))
	})

	It("Can report its depth", func() {
		urn1 := "urn:a:b"
		composite, _ := createUrnComposite(urn1)
		Expect(composite.Child(0).Depth()).To(Equal(1))
		Expect(composite.Child(0).Child(0).Depth()).To(Equal(2))
	})

})
