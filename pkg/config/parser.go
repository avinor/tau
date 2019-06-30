package config

import (
	"github.com/avinor/tau/pkg/hclcontext"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
)

var (
	parser *hcl2parse.Parser
)

func init() {
	parser = hcl2parse.NewParser()
}

// ParseFile parses file content into hcl.File. It will not read the filename from disk
// so filename argument can be anything, but it will not try to load same file twice.
//
// To convert the file into a struct use ParseBody
func ParseFile(source *SourceFile) (*hcl.File, error) {
	f, diags := parser.ParseHCL(source.Content, source.File)
	if diags.HasErrors() {
		return nil, diags
	}

	return f, nil
}

// ParseBody parses the hcl.Body into a value map.
func ParseBody(body hcl.Body, val interface{}) error {
	diags := gohcl2.DecodeBody(body, hclcontext.Default, val)

	if diags.HasErrors() {
		return diags
	}

	return nil
}
