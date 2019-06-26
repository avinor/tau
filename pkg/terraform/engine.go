package terraform

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
	v012 "github.com/avinor/tau/pkg/terraform/v012"
	"github.com/fatih/color"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// VersionCompatibility checks terraform executor for capabilities
type VersionCompatibility interface {
	GetValidCommands() []string
	GetInvalidArgs(command string) []string
}

// Generator for generating terraform assets
type Generator interface {
	GenerateOverrides(source *config.Source) ([]byte, bool, error)
	GenerateDependencies(source *config.Source) ([]byte, bool, error)
	GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error)
}

// Processor for processing terraform config or output
type Processor interface {
	ProcessBackendBody(body hcl.Body) (map[string]cty.Value, error)
	ProcessDependencies(dest string) (map[string]cty.Value, error)
}

// Executor executes terraform commands
type Executor interface {
	Execute(options *shell.Options, command string, args ...string) error
}

// Engine to process
type Engine struct {
	Version string

	Compatibility VersionCompatibility
	Generator     Generator
	Processor     Processor
	Executor      Executor
}

// NewEngine creates a terraform engine for the currently installed terraform version
func NewEngine() *Engine {

	version := version()

	if version == "" {
		log.Fatal(color.RedString("Could not identify terraform version. Make sure terraform is in PATH."))
	}

	log.Info(color.New(color.Bold).Sprintf("Terraform version: %s", version))
	log.Info("")

	var compatibility VersionCompatibility
	var generator Generator
	var processor Processor
	var executor Executor

	switch version {
	case "0.12":
		v012Engine := v012.NewEngine()
		compatibility = v012Engine
		generator = v012Engine
		processor = v012Engine
		executor = v012Engine
	default:
		log.Fatal(color.RedString("Unsupported terraform version!"))
	}

	return &Engine{
		Version:       version,
		Compatibility: compatibility,
		Generator:     generator,
		Processor:     processor,
		Executor:      executor,
	}
}

func (e *Engine) CreateOverrides(source *config.Source, dest string) error {
	log.Info(color.New(color.Bold).Sprint("Creating overrides..."))
	log.Info("")

	content, create, err := e.Generator.GenerateOverrides(source)

	if err != nil {
		return err
	}

	if !create {
		return nil
	}

	file := filepath.Join(dest, "tau_override.tf")

	return ioutil.WriteFile(file, content, os.ModePerm)
}

func (e *Engine) ResolveDependencies(source *config.Source, dest string) (map[string]cty.Value, error) {
	content, create, err := e.Generator.GenerateDependencies(source)

	if err != nil {
		return nil, err
	}

	if !create {
		return nil, nil
	}

	file := filepath.Join(dest, "main.tf")
	if err := ioutil.WriteFile(file, content, os.ModePerm); err != nil {
		return nil, err
	}

	test, err := e.Processor.ProcessDependencies(dest)
	if err != nil {
		return nil, err
	}

	log.Warnf("%v", test)

	return test, nil
}

func (e *Engine) WriteInputVariables(source *config.Source, dest string, variables map[string]cty.Value) error {
	content, err := e.Generator.GenerateVariables(source, variables)

	if err != nil {
		return err
	}

	file := filepath.Join(dest, "terraform.tfvars")

	return ioutil.WriteFile(file, content, os.ModePerm)
}
