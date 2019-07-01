package v012

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Processor struct {
	executor *Executor
}

func (p *Processor) ProcessBackendBody(body hcl.Body, context *hcl.EvalContext) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl2.DecodeBody(body, context, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}
