package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	depTest1 = `
		dependency "name" {
			source = "test"
		}
	`

	depTest2 = `
		dependency "name" {
			source = "overwrite"
		}
	`

	depTest3 = `
		dependency "two" {
			source = "test"
		}
	`

	depTest4 = `
		dependency "two" {
			source = ""
		}
	`

	depTest5 = `
		dependency "two" {
			source = "test"
			backend "azurerm" {}
		}
	`

	depTest6 = `
		dependency "two" {
			source = "test"
			backend "aws" {}
		}
	`
)

var (
	depFile1 = fileFromString("dep1", depTest1)
	depFile2 = fileFromString("dep2", depTest2)
	depFile3 = fileFromString("dep3", depTest3)
	depFile4 = fileFromString("dep4", depTest4)
	depFile5 = fileFromString("dep5", depTest5)
	depFile6 = fileFromString("dep6", depTest6)
)

func TestDependencyMerge(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]string
		Error    error
	}{
		{
			[]*File{depFile1},
			map[string]string{
				"name": "test",
			},
			nil,
		},
		{
			[]*File{depFile1, depFile2},
			map[string]string{
				"name": "overwrite",
			},
			nil,
		},
		{
			[]*File{depFile1, depFile3},
			map[string]string{
				"name": "test",
				"two":  "test",
			},
			nil,
		},
		{
			[]*File{depFile3, depFile5},
			map[string]string{
				"two": "test",
			},
			nil,
		},
		{
			[]*File{depFile5, depFile6},
			nil,
			differentBackendTypes,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeDependencies(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual := map[string]string{}
			for _, dep := range config.Dependencies {
				actual[dep.Name] = dep.Source
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}

func TestDependencyValidation(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]ValidationResult
	}{
		{
			[]*File{depFile1},
			map[string]ValidationResult{
				"name": ValidationResult{Result: true, Error: nil},
			},
		},
		{
			[]*File{depFile4},
			map[string]ValidationResult{
				"two": ValidationResult{Result: false, Error: dependencySourceMustBeSet},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeDependencies(config, getConfigFromFiles(t, test.Files))
			assert.NoError(t, err)

			for _, dep := range config.Dependencies {
				result, err := dep.Validate()
				expect := test.Expected[dep.Name]

				assert.Equal(t, expect.Result, result)
				if expect.Error != nil {
					assert.EqualError(t, err, expect.Error.Error())
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
