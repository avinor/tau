package v012

import (
	"encoding/json"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/ctytree"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
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

	return ctytree.CreateTree(values).ToCtyMap(), nil
}
