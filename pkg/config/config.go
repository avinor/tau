package config

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	moduleRequired = errors.Errorf("module block is required in config")
)

// Config structure for file describing deployment. This includes the module source, inputs
// dependencies, backend etc. One config element is connected to a single deployment
type Config struct {
	Datas        []*Data       `hcl:"data,block"`
	Dependencies []*Dependency `hcl:"dependency,block"`
	Hooks        []*Hook       `hcl:"hook,block"`
	Environment  *Environment  `hcl:"environment_variables,block"`
	Backend      *Backend      `hcl:"backend,block"`
	Module       *Module       `hcl:"module,block"`
	Inputs       *Inputs       `hcl:"inputs,block"`
}

// Merge all sources into current configuration struct.
// Should just call merge on all blocks / attributes of config struct.
func (c *Config) Merge(srcs []*Config) error {
	if err := mergeDatas(c, srcs); err != nil {
		return err
	}

	if err := mergeDependencies(c, srcs); err != nil {
		return err
	}

	if err := mergeHooks(c, srcs); err != nil {
		return err
	}

	if err := mergeEnvironments(c, srcs); err != nil {
		return err
	}

	if err := mergeBackends(c, srcs); err != nil {
		return err
	}

	if err := mergeModules(c, srcs); err != nil {
		return err
	}

	if err := mergeInputs(c, srcs); err != nil {
		return err
	}

	return nil
}

// PostProcess is called after merging all configurations together to perform additional
// processing after config is read. Can modify config elements
func (c *Config) PostProcess(file *File) {
	for _, hook := range c.Hooks {
		if hook.Command != nil && strings.HasPrefix(*hook.Command, ".") {
			fileDir := filepath.Dir(file.FullPath)
			absCommand := filepath.Join(fileDir, *hook.Command)
			hook.Command = &absCommand
		}
	}
}

// Validate that the configuration is correct. Calls validation on all parts of the struct.
// This assumes merge is already done and this is a complete configuration. If it is just a
// partial configuration from a child config it can fail as required blocks might not have
// been set.
func (c Config) Validate() (bool, error) {
	if c.Module == nil {
		return false, moduleRequired
	}

	for _, dep := range c.Dependencies {
		if valid, err := dep.Validate(); !valid {
			return false, err
		}
	}

	if c.Environment != nil {
		if valid, err := c.Environment.Validate(); !valid {
			return false, err
		}
	}

	for _, hook := range c.Hooks {
		if valid, err := hook.Validate(); !valid {
			return false, err
		}
	}

	return true, nil
}
