package config

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	hcl2parse "github.com/hashicorp/hcl2/hclparse"
)

// Parser can load config files
type Parser struct {
	parser *hcl2parse.Parser
}

// NewParser returns a new parser instance
func NewParser() (*Parser) {
	return &Parser{
		parser: hcl2parse.NewParser(),
	}
}

// Parse file and return the complete Config
func (p *Parser) Parse(filename string) (*Config, error) {
	f, diags := p.parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		return nil, diags
	}

	config := &Config{}

	if diags := gohcl2.DecodeBody(f.Body, nil, &config); diags.HasErrors() {
		return nil, diags
	}

	return config, nil
}