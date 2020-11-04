package terraform

import (
	"io/ioutil"
	"os"

	"github.com/zclconf/go-cty/cty"

	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/ctytree"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/terraform/def"
	v012 "github.com/avinor/tau/pkg/terraform/v012"
	"golang.org/x/mod/semver"
)

// Engine that can process version specific terraform commands
type Engine struct {
	Version string

	Compatibility def.VersionCompatibility
	Generator     def.Generator
	Executor      def.Executor
}

// NewEngine creates a terraform engine for the currently installed terraform version
func NewEngine(options *def.Options) *Engine {

	version := Version()

	if version == "" {
		ui.Fatal("Could not identify terraform version. Make sure terraform is in PATH.")
	}

	ui.Debug("Terraform version: %s", version)
	ui.NewLine()

	var compatibility def.VersionCompatibility
	var generator def.Generator
	var executor def.Executor

	switch {
	case semver.Compare(version, "0.12") >= 0:
		v012Engine := v012.NewEngine(options)
		compatibility = v012Engine
		generator = v012Engine
		executor = v012Engine
	default:
		ui.Fatal("Unsupported terraform version!")
	}

	return &Engine{
		Version:       version,
		Compatibility: compatibility,
		Generator:     generator,
		Executor:      executor,
	}
}

// CreateOverrides create the tau_override file in module folder. This file will overide
// backend settings
func (e *Engine) CreateOverrides(file *loader.ParsedFile) error {
	content, create, err := e.Generator.GenerateOverrides(file)

	if err != nil {
		return err
	}

	if !create {
		return nil
	}

	return ioutil.WriteFile(file.OverrideFile(), content, os.ModePerm)
}

// ResolveDependencies processes the source file and generates terraform modules for each unique
// source. For each source it will generate output arguments and return the merged values
//
// Bool return value indicates if it successfully resolved the dependency and should proceed to create
// source. If it failed to resolve dependencies but error is nil, it should not proceed to create this
// source, but should also not fail application. That generally means that it was a problem resolving
// dependencies for this source only. Other sources can still be generated.
func (e *Engine) ResolveDependencies(file *loader.ParsedFile) (bool, error) {
	processors, create, err := e.Generator.GenerateDependencies(file)

	if err != nil {
		return false, err
	}

	if !create {
		return true, nil
	}

	values := map[string]cty.Value{}

	for _, proc := range processors {
		vals, create, err := proc.Process()
		if err != nil {
			return false, err
		}

		// if not create then resolving dependency failed, but it should not result in an error.
		// it should just skip this source
		if !create {
			return false, nil
		}

		for key, value := range vals {
			values[key] = value
		}
	}

	for k, v := range ctytree.CreateTree(values).ToCtyMap() {
		file.AddToContext(k, v)
	}

	return true, nil
}

// WriteInputVariables write the terraform.tfvars file into module folder. This file is the parsed and
// processed variables where all dependencies and data source have been resolved and replaced with real
// values
func (e *Engine) WriteInputVariables(file *loader.ParsedFile) error {
	content, err := e.Generator.GenerateVariables(file)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(file.VariableFile(), content, os.ModePerm)
}
