package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	dataTest1 = `
		data "type" "name" {
			attr = "test"
		}
	`

	dataTest2 = `
		data "type" "name" {
			another = "test"
		}
	`

	dataTest3 = `
		data "type" "name" {
			attr = "crash"
		}
	`

	dataTest4 = `
		data "type" "two" {
			attr = "test"
		}
	`
)

var (
	dataFile1 = fileFromString("data1", dataTest1)
	dataFile2 = fileFromString("data2", dataTest2)
	dataFile3 = fileFromString("data3", dataTest3)
	dataFile4 = fileFromString("data4", dataTest4)
)

type ExpectedData struct {
	Type       string
	Name       string
	Attributes map[string]string
}

func TestDataMerge(t *testing.T) {
	tests := []struct {
		Files     []*File
		Expected  []ExpectedData
		Error     error
		AttrError bool
	}{
		{
			[]*File{dataFile1, dataFile3},
			nil,
			nil,
			true,
		},
		{
			[]*File{dataFile1, dataFile2},
			[]ExpectedData{
				{
					Type: "type",
					Name: "name",
					Attributes: map[string]string{
						"attr":    "test",
						"another": "test",
					},
				},
			},
			nil,
			false,
		},
		{
			[]*File{dataFile1, dataFile4},
			[]ExpectedData{
				{
					Type: "type",
					Name: "name",
					Attributes: map[string]string{
						"attr": "test",
					},
				},
				{
					Type: "type",
					Name: "two",
					Attributes: map[string]string{
						"attr": "test",
					},
				},
			},
			nil,
			false,
		},
		{
			[]*File{dataFile1, dataFile2, dataFile4},
			[]ExpectedData{
				{
					Type: "type",
					Name: "name",
					Attributes: map[string]string{
						"attr":    "test",
						"another": "test",
					},
				},
				{
					Type: "type",
					Name: "two",
					Attributes: map[string]string{
						"attr": "test",
					},
				},
			},
			nil,
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeDatas(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual := []ExpectedData{}
			for _, data := range config.Datas {
				attr, err := getBodyAttributes(data.Config)
				if err != nil {
					if test.AttrError {
						return
					}

					t.Fatalf("failed parsing attributes on %s.%s", data.Type, data.Name)
				}

				actual = append(actual, ExpectedData{Type: data.Type, Name: data.Name, Attributes: attr})
			}

			if test.AttrError {
				t.Fatal("expected attribute failure")
			}

			assert.Equal(t, test.Expected, actual)
		})
	}
}
