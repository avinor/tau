package config

import (
	"fmt"
	"testing"

	"github.com/avinor/tau/pkg/helpers/hclcontext"
	"github.com/stretchr/testify/assert"
)

var (
	testFile1 = `
		data "test" "test1" {
			var = "test1"
		}

		data "test" "test2" {
			var = "test1"
		}

		module {
			source = "file1"
			version = "1.0.0"
		}

		inputs {
			var1 = "file1"
			var2 = "file1"
		}
	`
	testFile2 = `
		module {
			source = "file2"
			version = "1.0.0"
		}

		environment_variables {
			var1 = "file2"
			var2 = "file2"
		}

		inputs {
			var1 = "file2"
			var2 = "file2"
		}
	`
	testFile3 = `
		environment_variables {
			var1 = "file3"
			var2 = "file3"
		}
	`
	testFile4 = `
		data "test" "test1" {
			var = "test4"
		}

		environment_variables {
			var3 = "file4"
			var4 = "file4"
		}
	`
	testFile5 = `
		data "test" "test1" {
			var = "test5"
		}

		inputs {
			var3 = "file5"
		}
	`
)

func TestMergeModule(t *testing.T) {
	tests := []struct {
		Files   []string
		Expects *Module
		Err     bool
	}{
		{
			[]string{testFile1},
			&Module{
				Source:  "file1",
				Version: stringToPointer("1.0.0"),
			},
			false,
		},
		{
			[]string{testFile1, testFile2},
			&Module{
				Source:  "file2",
				Version: stringToPointer("1.0.0"),
			},
			false,
		},
		{
			[]string{testFile1, testFile3},
			&Module{
				Source:  "file1",
				Version: stringToPointer("1.0.0"),
			},
			false,
		},
		{
			[]string{testFile3},
			nil,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config, err := getTestConfig(test.Files)

			if test.Err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.Expects, config.Module)
		})
	}
}

type expectsData struct {
	Type   string
	Name   string
	Config map[string]string
}

func TestMergeData(t *testing.T) {
	tests := []struct {
		Files   []string
		Expects []expectsData
		Err     bool
	}{
		{
			[]string{testFile1},
			[]expectsData{
				{
					Type: "test",
					Name: "test1",
					Config: map[string]string{
						"var": "test1",
					},
				},
				{
					Type: "test",
					Name: "test2",
					Config: map[string]string{
						"var": "test1",
					},
				},
			},
			false,
		},
		{
			[]string{testFile2, testFile4},
			[]expectsData{
				{
					Type: "test",
					Name: "test1",
					Config: map[string]string{
						"var": "test4",
					},
				},
			},
			false,
		},
		{
			[]string{testFile1, testFile4},
			[]expectsData{
				{
					Type: "test",
					Name: "test1",
					Config: map[string]string{
						"var": "test4",
					},
				},
				{
					Type: "test",
					Name: "test2",
					Config: map[string]string{
						"var": "test1",
					},
				},
			},
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config, err := getTestConfig(test.Files)

			if test.Err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			expects := []expectsData{}
			for _, data := range config.Datas {
				expect := expectsData{
					Type:   data.Type,
					Name:   data.Name,
					Config: map[string]string{},
				}
				err := ParseBody(data.Config, hclcontext.Default, &expect.Config)
				assert.NoError(t, err)
				expects = append(expects, expect)
			}

			assert.ElementsMatch(t, test.Expects, expects)
		})
	}
}

func TestMergeInputs(t *testing.T) {
	tests := []struct {
		Files   []string
		Expects map[string]string
		Err     bool
	}{
		{
			[]string{testFile1},
			map[string]string{
				"var1": "file1",
				"var2": "file1",
			},
			false,
		},
		{
			[]string{testFile1, testFile2},
			nil,
			true,
		},
		{
			[]string{testFile1, testFile5},
			map[string]string{
				"var1": "file1",
				"var2": "file1",
				"var3": "file5",
			},
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config, err := getTestConfig(test.Files)
			assert.NoError(t, err)

			expects := map[string]string{}
			err = ParseBody(config.Inputs.Config, hclcontext.Default, &expects)
			if test.Err {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, test.Expects, expects)
		})
	}
}

func TestValidateConfig(t *testing.T) {

}

func getTestConfig(files []string) (*Config, error) {
	sf := []*SourceFile{}
	for _, file := range files {
		sf = append(sf, &SourceFile{File: file, Content: []byte(file)})
	}

	return sf[0].ConfigMergedWith(sf[1:])
}

func stringToPointer(str string) *string {
	return &str
}
