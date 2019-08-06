package v012

import (
	"encoding/json"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

// Resolver is used to resolve output or config
type Resolver struct{}

// ResolveVariables decodes the config and returns a list of all variables found in hcl.Body.
// This is used later to be able to determine what variables it needs to resolve before
// finally resolving the body.
func (r *Resolver) ResolveVariables(body hcl.Body) ([]hcl.Traversal, error) {
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

// ResolveStateOutput takes the output from terraform command and parses the output into
// a map of string -> cty.Value
func (r *Resolver) ResolveStateOutput(output []byte) (map[string]cty.Value, error) {
	type OutputMeta struct {
		Sensitive bool            `json:"sensitive"`
		Type      json.RawMessage `json:"type"`
		Value     json.RawMessage `json:"value"`
	}
	outputs := map[string]OutputMeta{}
	values := map[string]cty.Value{}

	if err := json.Unmarshal(output, &outputs); err != nil {
		return nil, err
	}

	for name, meta := range outputs {
		ctyType, err := ctyjson.UnmarshalType(meta.Type)
		if err != nil {
			return nil, err
		}
		ctyValue, err := ctyjson.Unmarshal(meta.Value, ctyType)
		if err != nil {
			return nil, err
		}
		name = decodeName(name)
		values[name] = ctyValue
	}

	return values, nil
}
