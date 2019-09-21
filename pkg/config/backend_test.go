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
	backendFile1, _ = NewFile("/backend1", []byte(backendTest1))
	backendFile2, _ = NewFile("/backend2", []byte(backendTest2))
	backendFile3, _ = NewFile("/backend3", []byte(backendTest3))
	backendFile4, _ = NewFile("/backend4", []byte(backendTest4))
)

func TestBackendMerge(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]string
		Error    error
	}{
		{
			[]*File{backendFile1, backendFile2},
			map[string]string{
				"storage_account_name": "overwrite",
				"container_name":       "state",
			},
			nil,
		},
		{
			[]*File{backendFile1, backendFile3},
			map[string]string{
				"storage_account_name": "test",
				"container_name":       "state",
				"key":                  "addition",
			},
			nil,
		},
		{
			[]*File{backendFile4},
			map[string]string{
				"region": "oregon",
			},
			nil,
		},
		{
			[]*File{backendFile1, backendFile4},
			nil,
			differentBackendTypes,
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
				t.Fatal("failed getting attribute body", err)
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}

func TestBackendMergeWithStruct(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]string
		Error    error
	}{
		{
			[]*File{backendFile1, backendFile2},
			map[string]string{
				"storage_account_name": "overwrite",
				"container_name":       "state",
			},
			nil,
		},
		{
			[]*File{backendFile1, backendFile3},
			map[string]string{
				"storage_account_name": "test",
				"container_name":       "state",
				"key":                  "addition",
			},
			nil,
		},
		{
			[]*File{backendFile4},
			map[string]string{
				"region": "oregon",
			},
			nil,
		},
		{
			[]*File{backendFile1, backendFile4},
			nil,
			differentBackendTypes,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			configs := getConfigFromFiles(t, test.Files)

			backend := &Backend{}

			for _, config := range configs {
				if err := backend.Merge(config.Backend); err != nil {
					assert.EqualError(t, err, test.Error.Error())
					return
				}
			}

			actual, err := getBodyAttributes(backend.Config)
			if err != nil {
				t.Fatal("failed getting attribute body", err)
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}
