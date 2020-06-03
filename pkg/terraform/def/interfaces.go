package def

import (
	"github.com/zclconf/go-cty/cty"

	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/shell"
)

// DependencyProcessor can process a dependency and return the values from output.
// Each dependency processor will run in its own context, with separate environment variables.
// All dependency resolving that can be done in same context can be run in one processor, but
// use multiple processors to separate the context they run in
type DependencyProcessor interface {
	Process() (map[string]cty.Value, bool, error)
}

// OutputProcessor can parse the output from terraform and parse it into a map of values.
// It implements the shell.OutputProcessor interface so it can be sent into shell executor
// and read the values directly. Calling GetOutput after executing shell command should
// return the output values it found
type OutputProcessor interface {
	shell.OutputProcessor

	GetOutput() (map[string]cty.Value, error)
}

// VersionCompatibility checks terraform executor for capabilities
type VersionCompatibility interface {
	GetValidCommands() []string
	GetInvalidArgs(command string) []string
}

// Generator for generating terraform assets
type Generator interface {
	GenerateOverrides(file *loader.ParsedFile) ([]byte, bool, error)
	GenerateDependencies(file *loader.ParsedFile) ([]DependencyProcessor, bool, error)
	GenerateVariables(file *loader.ParsedFile) ([]byte, error)
}

// Executor executes terraform commands
type Executor interface {
	Execute(options *shell.Options, command string, args ...string) error
	NewOutputProcessor() OutputProcessor
}
