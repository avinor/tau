package command

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/hooks/def"
)

// Creator that can create a command executor
type Creator struct{}

// CanCreate checks if this creator can create executor for hook. Will check if hook
// has a command definition.
func (c *Creator) CanCreate(hook *config.Hook) bool {
	if !hook.HasCommand() {
		return false
	}

	return true
}

// Create a new executor for hook.
func (c *Creator) Create(hook *config.Hook) (def.Executor, error) {
	var command string
	var arguments []string
	var workingDir string

	if hook.Command != nil {
		command = *hook.Command
	}

	if hook.Arguments != nil {
		arguments = *hook.Arguments
	}

	if hook.WorkingDir != nil {
		workingDir = *hook.WorkingDir
	}

	return &Executor{
		Command:    command,
		Arguments:  arguments,
		WorkingDir: workingDir,
	}, nil
}
