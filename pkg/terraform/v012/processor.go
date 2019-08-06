package v012

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// Processor implements the def.Processor interface
type Processor struct {
	executor *Executor
}

// ProcessBackendBody returns a map of backend data processed in context of `context`
func (p *Processor) ProcessBackendBody(body hcl.Body, context *hcl.EvalContext) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl.DecodeBody(body, context, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}
