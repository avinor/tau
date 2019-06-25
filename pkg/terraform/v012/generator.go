package v012

import (
	"encoding/base64"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/terraform/lang"
	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/gohcl"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type Generator struct {
	ctx       *hcl.EvalContext
	processor *Processor
	resolver  *Resolver
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

	values, err := g.processor.ProcessBackendBody(source.Config.Backend.Config)
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

	trav, err := g.resolver.ResolveInputExpressions(source)
	if err != nil {
		return nil, false, err
	}

	if len(trav) == 0 {
		return nil, false, nil
	}

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

		var depBackend hcl.Body
		if dep.Backend != nil {
			depBackend = dep.Backend.Config
		}

		block, err := g.generateRemoteBackendBlock(dep.Name, depsource.Config.Backend.Type, depsource.Config.Backend.Config, depBackend)
		if err != nil {
			return nil, false, err
		}

		rootBody.AppendBlock(block)
	}

	for _, t := range trav {
		// For some reason this does not work.. using workaround under instead to convert
		// to a hclwrite.Expression and then to token
		// tokens := hclwrite.TokensForTraversal(t)

		expr := hclwrite.NewExpressionAbsTraversal(t)
		tokens := expr.BuildTokens(nil)
		outputName := base64.RawStdEncoding.EncodeToString(tokens.Bytes())

		// Need to "rewrite" root for dependencies
		if t.RootName() == "dependency" {
			split := t.SimpleSplit()
			root := hcl.TraverseRoot{
				Name: "data.terraform_remote_state",
			}
			t = hcl.TraversalJoin([]hcl.Traverser{root}, split.Rel)
		}

		expr.RenameVariablePrefix([]string{"dependency"}, []string{"remote.state"})

		block := hclwrite.NewBlock("output", []string{outputName})
		blockBody := block.Body()

		blockBody.SetAttributeTraversal("value", t)

		rootBody.AppendBlock(block)
	}

	formatted := hclwrite.Format(f.Bytes())

	return formatted, true, nil
}

func (g *Generator) GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	ctx := lang.ChildEvalContext(g.ctx, data)

	values := map[string]cty.Value{}
	diags := gohcl2.DecodeBody(source.Config.Inputs.Config, ctx, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	for name, value := range values {
		rootBody.SetAttributeValue(name, value)
	}

	return f.Bytes(), nil
}

func (g *Generator) generateHclWriterBlock(typeName string, labels []string, body *hclsyntax.Body) (*hclwrite.Block, error) {
	block := hclwrite.NewBlock(typeName, labels)
	blockBody := block.Body()

	for _, attr := range body.Attributes {
		value := cty.Value{}
		diags := gohcl.DecodeExpression(attr.Expr, g.ctx, &value)

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

func (g *Generator) generateRemoteBackendBlock(name, backend string, bodies ...hcl.Body) (*hclwrite.Block, error) {
	block := hclwrite.NewBlock("data", []string{"terraform_remote_state", name})
	blockBody := block.Body()

	values := map[string]cty.Value{}
	for _, body := range bodies {
		if body == nil {
			continue
		}

		vals, err := g.processor.ProcessBackendBody(body)
		if err != nil {
			return nil, err
		}

		for k, v := range vals {
			values[k] = v
		}
	}

	blockBody.SetAttributeValue("backend", cty.StringVal(backend))
	blockBody.SetAttributeValue("config", cty.MapVal(values))

	return block, nil
}
