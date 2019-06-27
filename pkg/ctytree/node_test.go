package ctytree

import (
	"fmt"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/zclconf/go-cty/cty"
)

func TestCreateTest(t *testing.T) {
	tests := []struct {
		values map[string]cty.Value
		tree   *Node
	}{
		{
			map[string]cty.Value{
				"dependency.test.outputs.id": cty.StringVal("value"),
			},
			&Node{
				Children: map[string]*Node{
					"dependency": {
						Children: map[string]*Node{
							"test": {
								Children: map[string]*Node{
									"outputs": {
										Children: map[string]*Node{
											"id": {
												Value: cty.StringVal("value"),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	prettyConfig := &pretty.Config{
		Compact:           true,
		IncludeUnexported: true,
		PrintStringers:    true,
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			tree := CreateTree(test.values)

			if same, err := compareTree(tree, test.tree); !same {
				t.Errorf("created and expected trees do not match: %s\ngot:\n%s\nexpected:\n%s",
					err,
					prettyConfig.Sprint(tree),
					prettyConfig.Sprint(test.tree))
			}
		})
	}
}

func compareTree(source, dest *Node) (bool, error) {
	if source.Value.Equals(dest.Value) == cty.False {
		return false, fmt.Errorf("values do not match")
	}

	matches := 0

	if len(source.Children) != len(dest.Children) {
		return false, fmt.Errorf("different number of children")
	}

	for sk, sv := range source.Children {
		for dk, dv := range dest.Children {
			if sk == dk {
				matches++
				if same, err := compareTree(sv, dv); !same {
					return false, err
				}
			}
		}
	}

	if len(source.Children) != matches {
		return false, fmt.Errorf("map contains different elements")
	}

	return true, nil
}
