package hcl

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/lang"
	"github.com/zclconf/go-cty/cty"
)

// NewContext creates a new evaluation context that supports all terraform functions and
// custom functions defined in tau
func NewContext() *hcl.EvalContext {
	s := lang.Scope{}
	funcs := s.Functions()

	funcs["env"] = EnvFunc

	return &hcl.EvalContext{
		Variables: map[string]cty.Value{},
		Functions: funcs,
	}
}
