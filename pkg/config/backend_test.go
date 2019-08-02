package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	backendTest1 = `
		backend "azurerm" {
			storage_account_name = "test"
			container_name       = "state"
		}
	`

	backendTest2 = `
		backend "azurerm" {
			storage_account_name = "overwrite"
		}
	`

	backendTest3 = `
		backend "azurerm" {
			key = "addition"
		}
	`

	backendTest4 = `
		backend "aws" {
			region = "oregon"
		}
	`
)

var (
	backendFile1 = fileFromString("backend1", backendTest1)
	backendFile2 = fileFromString("backend2", backendTest2)
	backendFile3 = fileFromString("backend3", backendTest3)
	backendFile4 = fileFromString("backend4", backendTest4)
)

func TestBackendMerge(t *testing.T) {
	tests := []struct {
		Files     []*File
		Expected  map[string]string
		Error     error
		AttrError bool
	}{
		{
			[]*File{backendFile1, backendFile2},
			nil,
			nil,
			true,
		},
		{
			[]*File{backendFile1, backendFile3},
			map[string]string{
				"storage_account_name": "test",
				"container_name":       "state",
				"key":                  "addition",
			},
			nil,
			false,
		},
		{
			[]*File{backendFile4},
			map[string]string{
				"region": "oregon",
			},
			nil,
			false,
		},
		{
			[]*File{backendFile1, backendFile4},
			nil,
			differentBackendTypes,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeBackends(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual, err := getBodyAttributes(config.Backend.Config)
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
