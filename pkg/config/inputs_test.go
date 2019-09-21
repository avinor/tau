package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

const (
	inputsTest1 = `
		inputs {
			test 	       = "test"
			test_var       = "value"
		}
	`

	inputsTest2 = `
		inputs {
			test = "overwrite"

		}
	`

	inputsTest3 = `
		inputs {
			invalid-bind = "invalid"
		}
	`

	inputsTest4 = `
		inputs {
			test123 = "5"
			list = "string"
		}
	`

	inputsTest5 = `
		inputs {
			list = ["number", "item"]
			tags = {
				app = "tau"
				version = "1.0"
			}
		}
	`

	inputsTest6 = `
		inputs {
			list = ["overwrite"]
			tags = {
				release = "stable"
			}
		}
	`

	inputsTest7 = `
		inputs {
			tags = {
				app = "terraform"
				month = "april"
			}
		}
	`
)

var (
	inputsFile1, _ = NewFile("/inputs1", []byte(inputsTest1))
	inputsFile2, _ = NewFile("/inputs2", []byte(inputsTest2))
	inputsFile3, _ = NewFile("/inputs3", []byte(inputsTest3))
	inputsFile4, _ = NewFile("/inputs4", []byte(inputsTest4))
	inputsFile5, _ = NewFile("/inputs5", []byte(inputsTest5))
	inputsFile6, _ = NewFile("/inputs6", []byte(inputsTest6))
	inputsFile7, _ = NewFile("/inputs7", []byte(inputsTest7))
)

// TestInputsMerge tests the inputs block. It does not test number values
// as they are stored as pointers in cty. This causes checks to fail as
// pointers are not the same in actual and expected result
func TestInputsMerge(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]cty.Value
		Error    error
	}{
		{
			[]*File{inputsFile1},
			map[string]cty.Value{
				"test":     cty.StringVal("test"),
				"test_var": cty.StringVal("value"),
			},
			nil,
		},
		{
			[]*File{inputsFile1, inputsFile2},
			map[string]cty.Value{
				"test":     cty.StringVal("overwrite"),
				"test_var": cty.StringVal("value"),
			},
			nil,
		},
		{
			[]*File{inputsFile1, inputsFile3},
			map[string]cty.Value{
				"test":         cty.StringVal("test"),
				"test_var":     cty.StringVal("value"),
				"invalid-bind": cty.StringVal("invalid"),
			},
			nil,
		},
		{
			[]*File{inputsFile1, inputsFile4},
			map[string]cty.Value{
				"test":     cty.StringVal("test"),
				"test_var": cty.StringVal("value"),
				"test123":  cty.StringVal("5"),
				"list":     cty.StringVal("string"),
			},
			nil,
		},
		{
			[]*File{inputsFile4, inputsFile5},
			map[string]cty.Value{
				"test123": cty.StringVal("5"),
				"list":    cty.TupleVal([]cty.Value{cty.StringVal("number"), cty.StringVal("item")}),
				"tags": cty.ObjectVal(map[string]cty.Value{
					"app":     cty.StringVal("tau"),
					"version": cty.StringVal("1.0"),
				}),
			},
			nil,
		},
		{
			[]*File{inputsFile5, inputsFile6},
			map[string]cty.Value{
				"list": cty.TupleVal([]cty.Value{cty.StringVal("overwrite")}),
				"tags": cty.ObjectVal(map[string]cty.Value{
					"app":     cty.StringVal("tau"),
					"version": cty.StringVal("1.0"),
					"release": cty.StringVal("stable"),
				}),
			},
			nil,
		},
		{
			[]*File{inputsFile5, inputsFile7},
			map[string]cty.Value{
				"list": cty.TupleVal([]cty.Value{cty.StringVal("number"), cty.StringVal("item")}),
				"tags": cty.ObjectVal(map[string]cty.Value{
					"app":     cty.StringVal("terraform"),
					"version": cty.StringVal("1.0"),
					"month":   cty.StringVal("april"),
				}),
			},
			nil,
		},
		{
			[]*File{inputsFile5, inputsFile6, inputsFile7},
			map[string]cty.Value{
				"list": cty.TupleVal([]cty.Value{cty.StringVal("overwrite")}),
				"tags": cty.ObjectVal(map[string]cty.Value{
					"app":     cty.StringVal("terraform"),
					"version": cty.StringVal("1.0"),
					"release": cty.StringVal("stable"),
					"month":   cty.StringVal("april"),
				}),
			},
			nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeInputs(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual, err := getCtyBodyAttributes(config.Inputs.Config)
			if err != nil {
				t.Fatal("failed getting attribute body", err)
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}
