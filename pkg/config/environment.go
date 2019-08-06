package config

import (
	"regexp"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

var (
	// envRegexp validates environment variable names against the POSIX standard
	envRegexp = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")

	// envVariableNotMatch is error if regexp for env variables do not match
	envVariableNotMatch = errors.Errorf("environment variable contains invalid character")
)

// Environment variables that should be added to the shell commands run (specfifically terraform).
// Define variables using attributes, blocks not supported
type Environment struct {
	Config hcl.Body `hcl:",remain"`
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
