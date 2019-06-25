package v012

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

type Processor struct {
	ctx *hcl.EvalContext
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
	debugLog := &processors.Log{
		Debug: true,
	}

	options := &shell.Options{
		Stdout: shell.Processors(debugLog),
		Stderr: shell.Processors(debugLog),
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

	// vals, err := e.Processor.ProcessOutput([]byte(buffer.Stdout()))
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}
