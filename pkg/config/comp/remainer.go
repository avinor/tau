package comp

import (
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
)

// Remainer includes a hcl remain tag to parse all remaining items in block
type Remainer struct {
	// Config hcl.Body `hcl:",remain"`
}

// ResolveVariables decodes the config and returns a list of all variables found in hcl.Body.
// This is used later to be able to determine what variables it needs to resolve before
// finally resolving the body.
func (r *Remainer) ResolveVariables(body hcl.Body) ([]hcl.Traversal, error) {
	exprs := map[string]hcl.Expression{}
	diags := gohcl.DecodeBody(body, nil, &exprs)

	if diags.HasErrors() {
		return nil, diags
	}

	trav := []hcl.Traversal{}
	for _, expr := range exprs {
		vars := expr.Variables()
		if len(vars) == 0 {
			continue
		}

		trav = append(trav, vars...)
	}

	return trav, nil
}
