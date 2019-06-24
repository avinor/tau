package terraform

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	v012 "github.com/avinor/tau/pkg/terraform/v012"
	"github.com/fatih/color"
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
	GenerateDependencies(source *config.Source) ([]byte, error)
	GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error)
}

// Processor for processing sources or terraform config
type Processor interface {
	ProcessOutput(output []byte) (map[string]cty.Value, error)
}

// Backend processes backend config
type Backend interface {
	ProcessBackendConfig(source *config.Source) (map[string]cty.Value, error)
}

// Engine to process
type Engine struct {
	Version string

	Compatibility VersionCompatibility
	Generator     Generator
	Processor     Processor
	Backend       Backend
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
	var backend Backend

	switch version {
	case "0.12":
		v012Engine := v012.NewEngine()
		compatibility = v012Engine
		generator = v012Engine
		processor = v012Engine
		backend = v012Engine
	default:
		log.Fatal(color.RedString("Unsupported terraform version!"))
	}

	return &Engine{
		Version:       version,
		Compatibility: compatibility,
		Generator:     generator,
		Processor:     processor,
		Backend:       backend,
	}
}

func (e *Engine) CreateOverrides(source *config.Source) error {
	return nil
}

func (e *Engine) CreateValues(source *config.Source) error {
	return nil
}
