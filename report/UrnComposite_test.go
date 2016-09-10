package report

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type node struct {
	Name     string
	Children []*node
	Parent   *node
	Value    interface{}
}

func (instance node) Child(index int) node {
	return *instance.Children[index]
}

func (instance node) Root() node {
	if instance.Parent == nil {
		return instance
	}

	root := instance.Parent

	for root.Parent != nil {
		root = root.Parent
	}

	return *root
}

func (instance *node) AddValue(urn string, value interface{}) error {
	var next = instance
	split := strings.Split(urn, ":")
	for index, item := range split {
		if index == 0 && item != instance.Root().Name {
			return fmt.Errorf("Multiple root elements")
		}

		if item == next.Name {
			continue
		}

		found := false
		for _, nodeElement := range next.Children {
			if nodeElement.Name == item {
				found = true
				next = nodeElement
				break
			}
		}

		if !found {
			childNode := createNode(item, next)
			next.Children = append(next.Children, childNode)
			next = childNode
		}
	}
	next.Value = value
	return nil
}

func createNode(name string, parent *node) *node {
	return &node{
		Name:     name,
		Children: []*node{},
		Parent:   parent,
	}
}

func createUrnAggregate(urns ...string) (*node, error) {

	var root *node
	var next *node

	for i := 0; i < len(urns); i++ {
		split := strings.Split(urns[i], ":")
		if root != nil {
			next = root
		}
		for index, item := range split {
			if root == nil {
				root = createNode(item, nil)
				next = root
				continue
			}

			if root != nil && index == 0 && item != root.Name {
				return nil, fmt.Errorf("Multiple root elements")
			}

			if item == next.Name {
				continue
			}

			found := false
			for _, nodeElement := range next.Children {
				if nodeElement.Name == item {
					found = true
					next = nodeElement
					break
				}
			}

			if !found {
				childNode := createNode(item, next)
				next.Children = append(next.Children, childNode)
				next = childNode
			}

		}
	}

	return root, nil
}

var _ = Describe("UrnComposite", func() {

	It("Combines two urns into a composite", func() {
		urn1 := "urn:a:b:c"
		urn2 := "urn:a:b:d"

		aggregate, err := createUrnAggregate(urn1, urn2)

		Expect(err).To(BeNil())
		Expect(aggregate.Name).To(Equal("urn"))
		Expect(len(aggregate.Children)).To(Equal(1))
		Expect(len(aggregate.Child(0).Child(0).Children)).To(Equal(2))
	})

	It("Detects multiple root elements", func() {

		urn1 := "urn:a:b:c"
		urn2 := "uri:a:b:d"

		_, err := createUrnAggregate(urn1, urn2)

		Expect(err).ToNot(BeNil())
		Expect(err).To(MatchError("Multiple root elements"))
	})

	It("Combines more than two urns into a composite", func() {
		urn1 := "urn:a:b:c"
		urn2 := "urn:a:b:d"
		urn3 := "urn:a:c:a"
		urn4 := "urn:a:b:e"

		aggregate, err := createUrnAggregate(urn1, urn2, urn3, urn4)

		Expect(err).To(BeNil())
		Expect(aggregate.Name).To(Equal("urn"))
		Expect(len(aggregate.Children)).To(Equal(1))
		Expect(len(aggregate.Child(0).Children)).To(Equal(2))
		Expect(len(aggregate.Child(0).Child(0).Children)).To(Equal(3))
	})

	It("Can locate the root", func() {
		urn1 := "urn:a:b:c"
		aggregate, _ := createUrnAggregate(urn1)
		Expect(aggregate.Child(0).Child(0).Child(0).Root().Name).To(Equal("urn"))

	})

	It("Can add a value to the node", func() {
		urn1 := "urn:a:b"
		urn2 := "urn:a:b:d"
		urn3 := "urn:a:b:e"
		urn4 := "urn:a:b:d:f"
		value := []int{1, 2, 3, 4, 5, 6}

		aggregate, _ := createUrnAggregate(urn1)
		aggregate.AddValue(urn2, value)
		aggregate.AddValue(urn3, value)
		aggregate.AddValue(urn4, value)
		Expect(aggregate.Child(0).Child(0).Child(0).Value).To(Equal(value))
		Expect(aggregate.Child(0).Child(0).Child(1).Value).To(Equal(value))
		Expect(aggregate.Child(0).Child(0).Child(0).Child(0).Value).To(Equal(value))
	})

})
