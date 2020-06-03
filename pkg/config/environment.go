package config

import (
	"regexp"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"

	"github.com/avinor/tau/pkg/config/comp"
)

var (
	// envRegexp validates environment variable names against the POSIX standard
	envRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

	// envVariableNotMatch is error if regexp for env variables do not match
	envVariableNotMatch = errors.Errorf("environment variable contains invalid character")

	// envCannotContainList is returned if some attributes in env variables is a list of items
	envCannotContainList = errors.Errorf("environment variables cannot contain a list of items")

	// envCannotContainMap is returned if some attributes in env variables is a map of items
	envCannotContainMap = errors.Errorf("environment variables cannot contain a map of items")
)

// Environment variables that should be added to the shell commands run (specfifically terraform).
// Define variables using attributes, blocks not supported
type Environment struct {
	Config hcl.Body `hcl:",remain"`

	comp.Remainer
}

// Merge current environment with config from source
func (e *Environment) Merge(src *Environment) error {
	if src == nil {
		return nil
	}

	e.Config = hcl.MergeBodies([]hcl.Body{e.Config, src.Config})

	return nil
}

// Validate checks that all variables contain valid names
func (e Environment) Validate() (bool, error) {
	attrs, diags := e.Config.JustAttributes()
	if diags.HasErrors() {
		return false, diags
	}

	for _, attr := range attrs {
		if !envRegexp.MatchString(attr.Name) {
			return false, envVariableNotMatch
		}

		if _, diags := hcl.ExprMap(attr.Expr); diags == nil {
			return false, envCannotContainMap
		}

		if _, diags := hcl.ExprList(attr.Expr); diags == nil {
			return false, envCannotContainList
		}
	}

	return true, nil
}

// Parse parses the config and returns all the environment variables defined in config.
func (e *Environment) Parse(context *hcl.EvalContext) (map[string]string, error) {
	env := map[string]string{}

	if e == nil {
		return env, nil
	}

	values := map[string]cty.Value{}
	diags := gohcl.DecodeBody(e.Config, context, &values)

	if diags.HasErrors() {
		return nil, diags
	}

	for key, value := range values {
		env[key] = value.AsString()
	}

	return env, nil
}

// mergeEnvironments merges only the environment from all configurations in srcs into dest
func mergeEnvironments(dest *Config, srcs []*Config) error {
	for _, src := range srcs {
		if src.Environment == nil {
			continue
		}

		if dest.Environment == nil {
			dest.Environment = src.Environment
			continue
		}

		if err := dest.Environment.Merge(src.Environment); err != nil {
			return err
		}
	}

	return nil
}
