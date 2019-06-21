package lang

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/lang"
	"github.com/zclconf/go-cty/cty"
)

// EvalContext creates an evaluation context or hcl
func EvalContext() *hcl.EvalContext {
	s := lang.Scope{}
	funcs := s.Functions()

	funcs["env"] = EnvFunc

	return &hcl.EvalContext{
		Functions: funcs,
	}
}

// ChildEvalContext returns a child context with variables set
func ChildEvalContext(context *hcl.EvalContext, vars map[string]cty.Value) *hcl.EvalContext {
	child := context.NewChild()
	child.Variables = vars

	return child
}
