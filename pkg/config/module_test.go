package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	moduleTest1 = `
		module {
			source = "./"
		}
	`

	moduleTest2 = `
		module {
			source = "./test"
			version = "1.0.0"
		}
	`

	moduleTest3 = `
		module {
			source = "./test"
			version = "1.1.0"
		}
	`
)

var (
	moduleFile1, _ = NewFile("/module1", []byte(moduleTest1))
	moduleFile2, _ = NewFile("/module2", []byte(moduleTest2))
	moduleFile3, _ = NewFile("/module3", []byte(moduleTest3))
)

func TestModuleMerge(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected *Module
	}{
		{
			[]*File{moduleFile1},
			&Module{
				Source: "./",
			},
		},
		{
			[]*File{moduleFile2},
			&Module{
				Source:  "./test",
				Version: "1.0.0",
			},
		},
		{
			[]*File{moduleFile2, moduleFile3},
			&Module{
				Source:  "./test",
				Version: "1.1.0",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeModules(config, getConfigFromFiles(t, test.Files))
			assert.NoError(t, err)

			assert.Equal(t, test.Expected, config.Module)
		})
	}
}
