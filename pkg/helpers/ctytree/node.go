package ctytree

import (
	"strings"

	"github.com/zclconf/go-cty/cty"
)

// Node is a node in tree structure that has children or value
// If value is set it will assume there are no children, only used for last node in tree
type Node struct {
	Children map[string]*Node
	Value    cty.Value
}

// CreateTree creates a new node tree based on the value map sent in. It will parse the
// keys as dot separated paths, ie. dependency.test.outputs.id and create a node tree
// where instead of key being dot separated it will have a root node called
// "dependency" that has child "test" and so on.
//
// Returns the root node with all top level nodes
func CreateTree(values map[string]cty.Value) *Node {
	rootNode := &Node{
		Children: map[string]*Node{},
	}

	for name, value := range values {
		node := rootNode.getNodePath(name)
		node.Value = value
	}

	return rootNode
}

// ToCtyMap converts the node tree back to a map of cty.Value's where cty.Values
// are type cty.ObjectVal to represent the hierarchy of nodes
func (n Node) ToCtyMap() map[string]cty.Value {
	objValues := map[string]cty.Value{}

	for name, subnode := range n.Children {
		if len(subnode.Children) == 0 {
			objValues[name] = subnode.Value
			continue
		}

		objValues[name] = cty.ObjectVal(subnode.ToCtyMap())
	}

	return objValues
}

// getNodePath creates the path defined and returns the node for path
func (n *Node) getNodePath(path string) *Node {
	split := strings.Split(path, ".")
	nextPathPart := split[0]

	if _, exists := n.Children[nextPathPart]; !exists {
		n.Children[nextPathPart] = &Node{
			Children: map[string]*Node{},
		}
	}

	if !strings.Contains(path, ".") {
		return n.Children[nextPathPart]
	}

	return n.Children[nextPathPart].getNodePath(strings.Join(split[1:], "."))
}
