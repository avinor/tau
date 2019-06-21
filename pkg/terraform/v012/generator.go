package v012

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Generator struct {
	backend *Backend
}

func (g *Generator) GenerateOverrides(source *config.Source) ([]byte, bool, error) {
	if source.Config.Backend == nil {
		return nil, false, nil
	}

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	tfBlock := rootBody.AppendNewBlock("terraform", nil)
	tfBody := tfBlock.Body()
	backendBlock := tfBody.AppendNewBlock("backend", []string{source.Config.Backend.Type})
	backendBody := backendBlock.Body()

	values, err := g.backend.ProcessBackendConfig(source)
	if err != nil {
		return nil, false, err
	}

	for k, v := range values {
		backendBody.SetAttributeValue(k, v)
	}

	return f.Bytes(), true, nil
}

func (g *Generator) GenerateDependencies(source *config.Source) ([]byte, error) {
	return nil, nil
}

func (g *Generator) GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error) {
	return nil, nil
}
