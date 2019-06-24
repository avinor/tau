package v012

import (
	"github.com/avinor/tau/pkg/config"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

type Resolver struct {
}

func (r *Resolver) ResolveInputExpressions(source *config.Source) ([]hcl.Traversal, error) {
	exprs := map[string]hcl.Expression{}
	diags := gohcl2.DecodeBody(source.Config.Inputs.Config, nil, &exprs)

	if diags.HasErrors() {
		return nil, diags
	}

	trav := []hcl.Traversal{}
	for _, expr := range exprs {
		vars := expr.Variables()
		if len(vars) == 0 {
			continue
		}

		for _, t := range vars {
			trav = append(trav, t)
		}
	}

	return trav, nil
}

func (r *Resolver) ResolveStateOutput(output []byte) (map[string]cty.Value, error) {
	return nil, nil
}
