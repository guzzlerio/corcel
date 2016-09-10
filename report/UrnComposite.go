package report

import (
	"fmt"
	"strings"
)

//UrnComposite ...
type UrnComposite struct {
	Name     string          `json:"name"`
	Children []*UrnComposite `json:"children"`
	parent   *UrnComposite
	Value    interface{} `json:"value"`
}

//Depth ...
func (instance UrnComposite) Depth() int {
	depth := 0
	if instance.parent == nil {
		return depth
	}

	root := instance.parent

	for root.parent != nil {
		depth = depth + 1
		root = root.parent
	}

	return depth + 1
}

//Child ...
func (instance UrnComposite) Child(index int) UrnComposite {
	return *instance.Children[index]
}

//Root ...
func (instance UrnComposite) Root() UrnComposite {
	if instance.parent == nil {
		return instance
	}

	root := instance.parent

	for root.parent != nil {
		root = root.parent
	}

	return *root
}

//AddValue ...
func (instance *UrnComposite) AddValue(urn string, value interface{}) error {
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

func createNode(name string, parent *UrnComposite) *UrnComposite {
	return &UrnComposite{
		Name:     name,
		Children: []*UrnComposite{},
		parent:   parent,
	}
}

func createUrnComposite(urns ...string) (*UrnComposite, error) {

	var root *UrnComposite
	var next *UrnComposite

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
