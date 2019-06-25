package v012

import (
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hcljson "github.com/hashicorp/hcl2/hcl/json"
	"github.com/zclconf/go-cty/cty"
)

type Processor struct {
	ctx      *hcl.EvalContext
	executor *Executor
}

func (p *Processor) ProcessBackendBody(body hcl.Body) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl2.DecodeBody(body, p.ctx, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}

func (p *Processor) ProcessDependencies(dest string) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}

	debugLog := &processors.Log{
		Debug: true,
	}

	options := &shell.Options{
		Stdout:           shell.Processors(debugLog),
		Stderr:           shell.Processors(debugLog),
		WorkingDirectory: dest,
	}

	if err := p.executor.Execute(options, "init"); err != nil {
		return nil, err
	}

	if err := p.executor.Execute(options, "apply"); err != nil {
		return nil, err
	}

	buffer := &processors.Buffer{}
	options.Stdout = shell.Processors(buffer)

	if err := p.executor.Execute(options, "output", "-json"); err != nil {
		return nil, err
	}

	file, diags := hcljson.Parse([]byte(buffer.Stdout()), "test")
	if diags.HasErrors() {
		return nil, diags
	}

	attrs, diag := file.Body.JustAttributes()
	if diag.HasErrors() {
		return nil, diag
	}

	for name, attr := range attrs {
		value, vdiag := attr.Expr.Value(nil)
		if vdiag.HasErrors() {
			return nil, vdiag
		}

		values[name] = value
	}

	return values, nil
}
