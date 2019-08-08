package v012

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/hclcontext"
	"github.com/avinor/tau/pkg/terraform/def"
	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// Generator implements the def.Generator interface and can generate files for terraform 0.12 version
type Generator struct {
	executor *Executor
}

// GenerateOverrides generates overrides file bytes
func (g *Generator) GenerateOverrides(file *loader.ParsedFile) ([]byte, bool, error) {
	if file.Config.Backend == nil {
		return nil, false, nil
	}

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	tfBlock := rootBody.AppendNewBlock("terraform", nil)
	tfBody := tfBlock.Body()
	backendBlock := tfBody.AppendNewBlock("backend", []string{file.Config.Backend.Type})
	backendBody := backendBlock.Body()

	values, err := processBackendBody(file.Config.Backend.Config, file.EvalContext())
	if err != nil {
		return nil, false, err
	}

	for k, v := range values {
		backendBody.SetAttributeValue(k, v)
	}

	return f.Bytes(), true, nil
}

// GenerateDependencies returns a list of all dependency processors that will generate dependencies.
func (g *Generator) GenerateDependencies(file *loader.ParsedFile) ([]def.DependencyProcessor, bool, error) {
	trav, err := file.Config.Inputs.ResolveVariables(file.Config.Inputs.Config)
	if err != nil {
		return nil, false, err
	}

	if len(trav) == 0 {
		return nil, false, nil
	}

	if len(file.Config.Datas) == 0 && len(file.Config.Dependencies) == 0 {
		return nil, false, nil
	}

	processors := []def.DependencyProcessor{}

	if len(file.Config.Datas) != 0 {
		dataProcessor, err := g.generateDataProcessor(file, trav)
		if err != nil {
			return nil, false, err
		}

		processors = append(processors, dataProcessor)
	}

	for _, dep := range file.Config.Dependencies {
		depProcessor, err := g.generateDepProcessor(file, dep, trav)
		if err != nil {
			return nil, false, err
		}

		processors = append(processors, depProcessor)
	}

	return processors, true, nil
}

// GenerateVariables generates the input variables
func (g *Generator) GenerateVariables(file *loader.ParsedFile) ([]byte, error) {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()

	values := map[string]cty.Value{}
	diags := gohcl.DecodeBody(file.Config.Inputs.Config, file.EvalContext(), &values)

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
		diags := gohcl.DecodeExpression(attr.Expr, hclcontext.NewContext(), &value)

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

func (g *Generator) generateRemoteBackendBlock(file *loader.ParsedFile, name, backend string, bodies ...hcl.Body) (*hclwrite.Block, error) {
	block := hclwrite.NewBlock("data", []string{"terraform_remote_state", name})
	blockBody := block.Body()

	values := map[string]cty.Value{}
	for _, body := range bodies {
		if body == nil {
			continue
		}

		vals, err := processBackendBody(body, file.EvalContext())
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

func (g *Generator) generateDataProcessor(file *loader.ParsedFile, trav []hcl.Traversal) (*DependencyProcessor, error) {
	dataProcessor := NewDependencyProcessor(file, file, g.executor)

	// TODO Make sure we use azurerm data provider < 2.0
	azblock := hclwrite.NewBlock("required_providers", []string{})
	azblock.Body().SetAttributeValue("azurerm", cty.StringVal("< 2.0.0"))
	tblock := hclwrite.NewBlock("terraform", []string{})
	tblock.Body().AppendBlock(azblock)
	dataProcessor.File.Body().AppendBlock(tblock)

	for _, data := range file.Config.Datas {
		block, err := g.generateHclWriterBlock("data", []string{data.Type, data.Name}, data.Config.(*hclsyntax.Body))
		if err != nil {
			return nil, err
		}

		dataProcessor.File.Body().AppendBlock(block)
	}

	// Find variables with data source
	for _, block := range generateOutputBlocks(trav, "data", "") {
		dataProcessor.File.Body().AppendBlock(block)
	}

	return dataProcessor, nil
}

func (g *Generator) generateDepProcessor(file *loader.ParsedFile, dep *config.Dependency, trav []hcl.Traversal) (*DependencyProcessor, error) {
	depFile, ok := file.Dependencies[dep.Name]
	if !ok {
		return nil, errors.Errorf("Could not find dependency %s", dep.Name)
	}

	if depFile.Config.Backend == nil {
		return nil, errors.Errorf("Dependencies must have a backend")
	}

	if dep.Backend != nil && depFile.Config.Backend.Type != dep.Backend.Type {
		return nil, errors.Errorf("Dependency backend type and override backend type must match")
	}

	var depBackend hcl.Body
	if dep.Backend != nil {
		depBackend = dep.Backend.Config
	}

	block, err := g.generateRemoteBackendBlock(depFile, dep.Name, depFile.Config.Backend.Type, depFile.Config.Backend.Config, depBackend)
	if err != nil {
		return nil, err
	}

	depProcessor := NewDependencyProcessor(file, depFile, g.executor)
	depProcessor.File.Body().AppendBlock(block)

	// Find variables using this dependency
	for _, block := range generateOutputBlocks(trav, "dependency", dep.Name) {
		depProcessor.File.Body().AppendBlock(block)
	}

	return depProcessor, nil
}

func generateOutputBlocks(trav []hcl.Traversal, rootName, name string) []*hclwrite.Block {
	blocks := map[string]*hclwrite.Block{}

	for _, t := range trav {
		// For some reason this does not work.. using workaround under instead to convert
		// to a hclwrite.Expression and then to token
		// tokens := hclwrite.TokensForTraversal(t)

		if t.RootName() != rootName {
			continue
		}

		expr := hclwrite.NewExpressionAbsTraversal(t)
		tokens := expr.BuildTokens(nil)
		fullname := tokens.Bytes()

		if _, ok := blocks[string(fullname)]; ok {
			continue
		}

		if name != "" {
			if len(tokens) < 3 || string(tokens[2].Bytes) != name {
				continue
			}
		}

		// Need to "rewrite" root for dependencies
		if t.RootName() == "dependency" {
			split := t.SimpleSplit()
			root := hcl.TraverseRoot{
				Name: "data.terraform_remote_state",
			}
			t = hcl.TraversalJoin([]hcl.Traverser{root}, split.Rel)
		}

		block := hclwrite.NewBlock("output", []string{encodeName(fullname)})
		blockBody := block.Body()

		blockBody.SetAttributeTraversal("value", t)

		blocks[string(fullname)] = block
	}

	ret := []*hclwrite.Block{}
	for _, block := range blocks {
		ret = append(ret, block)
	}

	return ret
}

func generateOutputTraversalBlock(t hcl.Traversal, rootname string, name string) *hclwrite.Block {
	// For some reason this does not work.. using workaround under instead to convert
	// to a hclwrite.Expression and then to token
	// tokens := hclwrite.TokensForTraversal(t)

	if t.RootName() != rootname {
		return nil
	}

	expr := hclwrite.NewExpressionAbsTraversal(t)
	tokens := expr.BuildTokens(nil)
	outputName := encodeName(tokens.Bytes())

	if name != "" {
		if len(tokens) < 3 || string(tokens[2].Bytes) != name {
			return nil
		}
	}

	// Need to "rewrite" root for dependencies
	if t.RootName() == "dependency" {
		split := t.SimpleSplit()
		root := hcl.TraverseRoot{
			Name: "data.terraform_remote_state",
		}
		t = hcl.TraversalJoin([]hcl.Traverser{root}, split.Rel)
	}

	block := hclwrite.NewBlock("output", []string{outputName})
	blockBody := block.Body()

	blockBody.SetAttributeTraversal("value", t)

	return block
}

// processBackendBody returns a map of backend data processed in context of `context`
func processBackendBody(body hcl.Body, context *hcl.EvalContext) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl.DecodeBody(body, context, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}
