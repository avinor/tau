package config

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/stretchr/testify/assert"
)

func fileFromString(name string, content string) *File {
	return &File{
		Name:     name,
		FullPath: name,
		Content:  []byte(content),
		Children: []*File{},
	}
}

func getConfigFromFiles(t *testing.T, files []*File) []*Config {
	configs := []*Config{}

	for _, file := range files {
		config, err := file.parse(nil)
		if err != nil {
			t.Fatal("test failed parsing config", err)
		}

		configs = append(configs, config)
	}

	return configs
}

func getBodyAttributes(body hcl.Body) (map[string]string, error) {
	attrs, diags := body.JustAttributes()
	if diags != nil {
		return nil, diags
	}

	actual := map[string]string{}
	for _, attr := range attrs {
		value, diags := attr.Expr.Value(nil)
		if diags != nil {
			return nil, diags
		}

		actual[attr.Name] = value.AsString()
	}

	return actual, nil
}

func testBodyAttributes(t *testing.T, body hcl.Body, expected map[string]string, expectError bool) {
	attrs, diags := body.JustAttributes()
	if diags != nil {
		if expectError {
			return
		}

		t.Fatal("test attributes failed", diags)
	}

	actual := map[string]string{}
	for _, attr := range attrs {
		value, diags := attr.Expr.Value(nil)
		if diags != nil {
			t.Fatal("failed parsing value expression", diags)
		}

		actual[attr.Name] = value.AsString()
	}

	assert.Equal(t, expected, actual)
}
