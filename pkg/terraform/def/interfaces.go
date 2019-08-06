package def

import (
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/shell"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// DependencyProcesser can process a dependency and return the values from output.
// Each dependency processor will run in its own context, with separate environment variables.
// All dependency resolving that can be done in same context can be run in one processor, but
// use multiple processors to separate the context they run in
type DependencyProcesser interface {
	Process() (map[string]cty.Value, bool, error)
}

// VersionCompatibility checks terraform executor for capabilities
type VersionCompatibility interface {
	GetValidCommands() []string
	GetInvalidArgs(command string) []string
}

// Generator for generating terraform assets
type Generator interface {
	GenerateOverrides(file *loader.ParsedFile) ([]byte, bool, error)
	GenerateDependencies(file *loader.ParsedFile) ([]DependencyProcesser, bool, error)
	GenerateVariables(file *loader.ParsedFile) ([]byte, error)
}

// Processor for processing terraform config or output
type Processor interface {
	ProcessBackendBody(body hcl.Body, context *hcl.EvalContext) (map[string]cty.Value, error)
}

// Executor executes terraform commands
type Executor interface {
	Execute(options *shell.Options, command string, args ...string) error
}
