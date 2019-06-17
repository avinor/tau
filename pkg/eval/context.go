package eval

import (
	"github.com/hashicorp/terraform/lang"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// CreateContext creates an evaluation context or hcl
func CreateContext(vals map[string]cty.Value) (*hcl.EvalContext, error) {
	s := lang.Scope{}
	funcs := s.Functions()

	funcs["env"] = EnvFunc

	return &hcl.EvalContext{
		Variables: vals,
		Functions: funcs,
	}, nil
}