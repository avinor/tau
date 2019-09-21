package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	envTest1 = `
		environment_variables {
			test 	       = "test"
			test_var       = "value"
		}
	`

	envTest2 = `
		environment_variables {
			test = "overwrite"
		}
	`

	envTest3 = `
		environment_variables {
			invalid-bind = "invalid"
		}
	`

	envTest4 = `
		environment_variables {
			test123 = "number"
		}
	`

	envTest5 = `
		environment_variables {
			test = {
				mapitem = "list"
			}
		}
	`

	envTest6 = `
		environment_variables {
			list = ["test"]
		}
	`
)

var (
	envFile1, _ = NewFile("/env1", []byte(envTest1))
	envFile2, _ = NewFile("/env2", []byte(envTest2))
	envFile3, _ = NewFile("/env3", []byte(envTest3))
	envFile4, _ = NewFile("/env4", []byte(envTest4))
	envFile5, _ = NewFile("/env5", []byte(envTest5))
	envFile6, _ = NewFile("/env6", []byte(envTest6))
)

func TestEnvironmentRegexp(t *testing.T) {
	tests := []struct {
		Name  string
		Match bool
	}{
		{"test", true},
		{"TEST", true},
		{"testTEST", true},
		{"test_test", true},
		{"test123", true},
		{"test_underscore", true},
		{"_test_underscore", true},
		{"test.dot", false},
		{"/tmp/test.hcl", false},
		{"123test", false},
		{"test-bind", false},
		{"/tmp/test.tau", false},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			assert.Equal(t, test.Match, envRegexp.MatchString(test.Name), test.Name)
		})
	}
}

func TestEnvironmentMerge(t *testing.T) {
	tests := []struct {
		Files     []*File
		Expected  map[string]string
		Error     error
		AttrError bool
	}{
		{
			[]*File{envFile1},
			map[string]string{
				"test":     "test",
				"test_var": "value",
			},
			nil,
			false,
		},
		{
			[]*File{envFile1, envFile2},
			nil,
			nil,
			true,
		},
		{
			[]*File{envFile1, envFile3},
			map[string]string{
				"test":         "test",
				"test_var":     "value",
				"invalid-bind": "invalid", // this is valid merge, just not pass validation
			},
			nil,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeEnvironments(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual, err := getBodyAttributes(config.Environment.Config)
			if err != nil {
				if test.AttrError {
					return
				}

				t.Fatal("failed getting attribute body", err)
			}

			if test.AttrError {
				t.Fatal("expected attribute failure")
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}

func TestEnvironmentValidation(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected ValidationResult
	}{
		{
			[]*File{envFile1},
			ValidationResult{Result: true, Error: nil},
		},
		{
			[]*File{envFile1, envFile4},
			ValidationResult{Result: true, Error: nil},
		},
		{
			[]*File{envFile3},
			ValidationResult{Result: false, Error: envVariableNotMatch},
		},
		{
			[]*File{envFile5},
			ValidationResult{Result: false, Error: envCannotContainMap},
		},
		{
			[]*File{envFile6},
			ValidationResult{Result: false, Error: envCannotContainList},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeEnvironments(config, getConfigFromFiles(t, test.Files))
			assert.NoError(t, err)

			result, err := config.Environment.Validate()

			assert.Equal(t, test.Expected.Result, result)
			if test.Expected.Error != nil {
				assert.EqualError(t, err, test.Expected.Error.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
