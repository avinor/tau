package def

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

// DependencyProcesser can process a dependency and return the values from output
type DependencyProcesser interface {
	Name() string
	Content() []byte
	Process(dest string) (map[string]cty.Value, bool, error)
}

// VersionCompatibility checks terraform executor for capabilities
type VersionCompatibility interface {
	GetValidCommands() []string
	GetInvalidArgs(command string) []string
}

// Generator for generating terraform assets
type Generator interface {
	GenerateOverrides(source *config.Source) ([]byte, bool, error)
	GenerateDependencies(source *config.Source) ([]DependencyProcesser, bool, error)
	GenerateVariables(source *config.Source, data map[string]cty.Value) ([]byte, error)
}

// Processor for processing terraform config or output
type Processor interface {
	ProcessBackendBody(body hcl.Body, context *hcl.EvalContext) (map[string]cty.Value, error)
}

// Executor executes terraform commands
type Executor interface {
	Execute(options *shell.Options, command string, args ...string) error
}
