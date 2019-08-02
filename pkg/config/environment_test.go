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
)

var (
	envFile1 = fileFromString("env1", envTest1)
	envFile2 = fileFromString("env2", envTest2)
	envFile3 = fileFromString("env3", envTest3)
	envFile4 = fileFromString("env4", envTest4)
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
