package terraform

import (
	"regexp"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	v012 "github.com/avinor/tau/pkg/terraform/v012"
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

const (
	versionPattern = "Terraform v(\\d+.\\d+)"
)

var (
	versionRegex = regexp.MustCompile(versionPattern)
)

// NewEngine creates a terraform engine for the currently installed terraform version
func NewEngine() *Engine {
	version := version()

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

	}

	return &Engine{
		Version:       version,
		Compatibility: compatibility,
		Generator:     generator,
		Processor:     processor,
		Backend:       backend,
	}
}

func version() string {
	buffer := &processors.Buffer{}

	options := &shell.Options{
		Stdout: shell.Processors(buffer),
		Stderr: shell.Processors(buffer),
	}

	if err := shell.Execute(options, "terraform", "version"); err != nil {
		return ""
	}

	matches := versionRegex.FindAllStringSubmatch(buffer.Stdout(), -1)

	if len(matches) < 1 && len(matches[0]) < 2 {
		return ""
	}

	return matches[0][1]
}
