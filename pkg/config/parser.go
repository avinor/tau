package config

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
)

var (
	Parser *parser
)

func init() {
	Parser = newParser()
}

// Parser can load config files
type parser struct {
	parser *hcl2parse.Parser
}

// NewParser returns a new parser instance
func newParser() *parser {
	return &parser{
		parser: hcl2parse.NewParser(),
	}
}

// Parse file and return the complete Config
func (p *parser) Parse(content []byte, filename string) (*Config, error) {
	f, diags := p.parser.ParseHCL(content, filename)
	if diags.HasErrors() {
		return nil, diags
	}

	config := &Config{}
	
	if diags := gohcl2.DecodeBody(f.Body, nil, config); diags.HasErrors() {
		return nil, diags
	}

	return config, nil
}
