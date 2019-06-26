package config

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
)

var (
	parser *hcl2parse.Parser
)

func init() {
	parser = hcl2parse.NewParser()
}

// Parse file and return the complete Config
func Parse(content []byte, filename string, val interface{}) error {
	f, diags := parser.ParseHCL(content, filename)
	if diags.HasErrors() {
		return diags
	}

	if diags := gohcl2.DecodeBody(f.Body, nil, val); diags.HasErrors() {
		return diags
	}

	return nil
}
