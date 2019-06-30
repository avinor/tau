package hclcontext

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/terraform/lang"
	"github.com/zclconf/go-cty/cty"
)

var (
	// Default context with no variables but all functions defined
	Default = NewContext()
)

// NewContext creates a new evaluation context that supports all terraform functions and
// custom functions defined in tau
func NewContext() *hcl.EvalContext {
	s := lang.Scope{}
	funcs := s.Functions()

	funcs["env"] = EnvFunc

	return &hcl.EvalContext{
		Functions: funcs,
	}
}

// WithVariables returns a new child context with variables added.
// If requiring a new scope that defines additional variables this should be used and not
// add variables to Default or parent scope as they would be available to all.
func WithVariables(context *hcl.EvalContext, vars map[string]cty.Value) *hcl.EvalContext {
	child := context.NewChild()
	child.Variables = vars

	return child
}
