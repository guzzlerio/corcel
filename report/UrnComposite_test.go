package report

import (
	"encoding/json"
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

func createUrnAggregate(urns ...string) *node {

	var root *node
	var next *node

	for i := 0; i < len(urns); i++ {
		split := strings.Split(urns[i], ":")
		if root != nil {
			next = root
		}
		for _, item := range split {
			if root == nil {
				root = createNode(item)
				next = root
				continue
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

	return root
}

var _ = Describe("UrnComposite", func() {

	It("does something", func() {
		urn1 := "urn:a:b:c"
		urn2 := "urn:a:b:d"

		aggregate := createUrnAggregate(urn1, urn2)

		data, _ := json.MarshalIndent(aggregate, "", "  ")

		fmt.Println(string(data))

		Expect(aggregate.Name).To(Equal("urn"))
		Expect(len(aggregate.Children)).To(Equal(1))
		Expect(len(aggregate.Child(0).Child(0).Children)).To(Equal(2))

	})

})
