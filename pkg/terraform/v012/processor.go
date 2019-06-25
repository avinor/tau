package v012

import (
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Processor struct {
	ctx *hcl.EvalContext
}

func (p *Processor) ProcessBackendBody(body hcl2.Body) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl2.DecodeBody(body, p.ctx, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}

func (p *Processor) ProcessOutput(output []byte) (map[string]cty.Value, error) {
	return nil, nil
}
