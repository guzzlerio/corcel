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
}

func (instance node) Child(index int) node {
	return *instance.Children[index]
}

func createNode(name string) *node {
	return &node{
		Name:     name,
		Children: []*node{},
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
				root = createNode(item)
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
				childNode := createNode(item)
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

})
