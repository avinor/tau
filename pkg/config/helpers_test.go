package config

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// ValidationResult is used when creating test struct to check validation results
type ValidationResult struct {
	Result bool
	Error  error
}

// getConfigFromFiles parses the files and returns the config structures
// If it fails to parse it will call t.Fatal to stop test
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

// getBodyAttributes returns map of string -> string with all attributes
// on hcl.Body. It will use nil evalContext so no functions can be used
func getBodyAttributes(body hcl.Body) (map[string]string, error) {
	attrs, diags := body.JustAttributes()
	if diags.HasErrors() {
		return nil, diags
	}

	actual := map[string]string{}
	for _, attr := range attrs {
		value, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, diags
		}

		actual[attr.Name] = value.AsString()
	}

	return actual, nil
}

// getCtyBodyAttributes returns map of string -> cty.Value with all attributes
// on hcl.Body. It will use nil evalContext so no functions can be used
func getCtyBodyAttributes(body hcl.Body) (map[string]cty.Value, error) {
	attrs, diags := body.JustAttributes()
	if diags.HasErrors() {
		return nil, diags
	}

	actual := map[string]cty.Value{}
	for _, attr := range attrs {
		value, diags := attr.Expr.Value(nil)
		if diags.HasErrors() {
			return nil, diags
		}

		actual[attr.Name] = value
	}

	return actual, nil
}
