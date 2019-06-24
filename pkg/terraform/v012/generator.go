package v012

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/gohcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
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

func (g *Generator) GenerateDependencies(source *config.Source) ([]byte, bool, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	if len(source.Config.Datas) == 0 && len(source.Config.Dependencies) == 0 {
		return nil, false, nil
	}

	for _, data := range source.Config.Datas {
		block, err := g.generateHclWriterBlock("data", []string{data.Type, data.Name}, data.Config.(*hclsyntax.Body))
		if err != nil {
			return nil, false, err
		}

		rootBody.AppendBlock(block)
	}

	for _, dep := range source.Config.Dependencies {
		depsource, ok := source.Dependencies[dep.Name]
		if !ok {
			return nil, false, errors.Errorf("Could not find dependency %s", dep.Name)
		}

		if depsource.Config.Backend == nil {
			return nil, false, errors.Errorf("Dependencies must have a backend")
		}

		if dep.Backend != nil && depsource.Config.Backend.Type != dep.Backend.Type {
			return nil, false, errors.Errorf("Dependency backend type and override backend type must match")
		}

		bodies := []hcl2.Body{depsource.Config.Backend.Config}

		if dep.Backend != nil {
			bodies = append(bodies, dep.Backend.Config)
		}

		mergedBody := hcl2.MergeBodies(bodies)

		block, err := g.generateRemoteBackendBlock(dep.Name, depsource.Config.Backend.Type, mergedBody)
		if err != nil {
			return nil, false, err
		}

		rootBody.AppendBlock(block)
	}

	return f.Bytes(), true, nil
}

func (g *Generator) GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error) {
	return nil, nil
}

func (g *Generator) generateHclWriterBlock(typeName string, labels []string, body *hclsyntax.Body) (*hclwrite.Block, error) {
	block := hclwrite.NewBlock(typeName, labels)
	blockBody := block.Body()

	for _, attr := range body.Attributes {
		value := cty.Value{}
		diags := gohcl.DecodeExpression(attr.Expr, g.backend.ctx, &value)

		if diags.HasErrors() {
			return nil, diags
		}

		blockBody.SetAttributeValue(attr.Name, value)
	}

	for _, block := range body.Blocks {
		subBlock, err := g.generateHclWriterBlock(block.Type, block.Labels, block.Body)
		if err != nil {
			return nil, err
		}

		blockBody.AppendBlock(subBlock)
	}

	return block, nil
}

func (g *Generator) generateRemoteBackendBlock(name, backend string, body hcl2.Body) (*hclwrite.Block, error) {
	block := hclwrite.NewBlock("data", []string{"terraform_remote_state", name})
	blockBody := block.Body()

	vals, err := g.backend.processBackendBody(body)
	if err != nil {
		return nil, err
	}

	blockBody.SetAttributeValue("backend", cty.StringVal(backend))
	blockBody.SetAttributeValue("config", cty.MapVal(vals))

	return block, nil
}
