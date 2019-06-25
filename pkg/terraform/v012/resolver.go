package v012

import (
	stdstr "strings"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/strings"
	gohcl2 "github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	hcljson "github.com/hashicorp/hcl2/hcl/json"
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
	values := map[string]cty.Value{}

	file, diags := hcljson.Parse(output, strings.SecureRandomAlphaString(16))
	if diags.HasErrors() {
		return nil, diags
	}

	attrs, diag := file.Body.JustAttributes()
	if diag.HasErrors() {
		return nil, diag
	}

	for name, attr := range attrs {
		value, vdiag := attr.Expr.Value(nil)
		if vdiag.HasErrors() {
			return nil, vdiag
		}

		name = decodeName(name)
		values[name] = value
	}

	return convertFlatMapToObjects(values)
}

func convertFlatMapToObjects(values map[string]cty.Value) (map[string]cty.Value, error) {
	type node struct {
		name     string
		children []node
	}

	objMap := map[string][]string{}

	for name := range values {
		split := stdstr.Split(name, ".")

		for ii := len(split) - 1; ii > 0; ii-- {
			partname := stdstr.Join(split[0:ii], ".")
			nextname := stdstr.Join(split[0:ii+1], ".")

			if _, ok := objMap[partname]; !ok {
				objMap[partname] = []string{}
			}

			objMap[partname] = append(objMap[partname], nextname)
		}
	}

	for k, v := range objMap {
		log.Warnf("%s => %s", k, v)
	}

	return values, nil
}

func getRootElements(values map[string]cty.Value) (root map[string]cty.Value) {
	for name, value := range values {
		if stdstr.Index(name, ".") >= 0 {
			continue
		}

		root[name] = value
	}

	return root
}
