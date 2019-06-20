package terraform

// import (
// 	"github.com/avinor/tau/pkg/eval"
// 	"github.com/hashicorp/hcl2/hcl"
// 	"github.com/zclconf/go-cty/cty"
// 	"github.com/hashicorp/hcl2/hclwrite"
// 	gohcl2 "github.com/hashicorp/hcl2/gohcl"
// 	hcl2parse "github.com/hashicorp/hcl2/hclparse"
// )

// var (
// 	Parser *parser
// )

// func init() {
// 	Parser = newParser()
// }

// // Parser can load config files
// type parser struct {
// 	parser *hcl2parse.Parser
// }

// // NewParser returns a new parser instance
// func newParser() *parser {
// 	return &parser{
// 		parser: hcl2parse.NewParser(),
// 	}
// }

// // Parse file and return the complete Config
// func (p *parser) Parse(content []byte, filename string) (*Config, error) {
// 	f, diags := p.parser.ParseHCL(content, filename)
// 	if diags.HasErrors() {
// 		return nil, diags
// 	}

// 	config := &Config{}

// 	if diags := gohcl2.DecodeBody(f.Body, nil, config); diags.HasErrors() {
// 		return nil, diags
// 	}

// 	return config, nil
// }

// func GetTerraformOverride(config *Config) ([]byte, error) {
// 	ctx, err := eval.CreateContext(nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	f := hclwrite.NewEmptyFile()
// 	rootBody := f.Body()
// 	tfBlock := rootBody.AppendNewBlock("terraform", nil)
// 	tfBody := tfBlock.Body()
// 	backendBlock := tfBody.AppendNewBlock("backend", []string{config.Backend.Type})
// 	backendBody := backendBlock.Body()

// 	bc, err := GetBackendConfig(config, ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for k, v := range bc {
// 		backendBody.SetAttributeValue(k, v)
// 	}

// 	return f.Bytes(), nil
// }

// func GetBackendConfig(config *Config, ctx *hcl.EvalContext) (map[string]cty.Value, error) {
// 	values := map[string]cty.Value{}
// 	diags := gohcl2.DecodeBody(config.Backend.Config, ctx, &values)
	
// 	if diags.HasErrors() {
// 		return nil, diags
// 	}

// 	return values, nil
// }
