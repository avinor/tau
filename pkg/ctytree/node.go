package ctytree

import (
	"strings"
	"github.com/zclconf/go-cty/cty"
)

type Node struct {
	Children map[string]*Node
	Value    *cty.Value
}

func CreateTree(values map[string]cty.Value) *Node {
	rootNode := &Node{
		Children: map[string]*Node{},
	}

	for name, value := range values {
		node := rootNode.getNodePath(name)
		node.Value = &value
	}

	return rootNode
}

func (n Node) ToCtyMap() map[string]cty.Value {
	objValues := map[string]cty.Value{}

	for name, subnode := range n.Children {
		if subnode.Value != nil {
			objValues[name] = *subnode.Value
			continue
		}

		objValues[name] = cty.ObjectVal(subnode.ToCtyMap())
	}

	return objValues
}

func (n *Node) getNodePath(path string) *Node {
	if strings.Index(path, ".") < 0 {
		return n
	}

	split := strings.Split(path, ".")
	nextPathPart := split[0]

	if _, exists := n.Children[nextPathPart]; !exists {
		n.Children[nextPathPart] = &Node{
			Children: map[string]*Node{},
		}
	}

	return n.Children[nextPathPart].getNodePath(strings.Join(split[1:], "."))
}
