package v012

import (
	"github.com/avinor/tau/pkg/config"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hcl2 "github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Backend struct {
	ctx *hcl.EvalContext
}

func (b *Backend) ProcessBackendConfig(source *config.Source) (map[string]cty.Value, error) {
	return b.processBackendBody(source.Config.Backend.Config)
}

func (b *Backend) processBackendBody(body hcl2.Body) (map[string]cty.Value, error) {
	values := map[string]cty.Value{}
	diags := gohcl2.DecodeBody(body, b.ctx, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	return values, nil
}
