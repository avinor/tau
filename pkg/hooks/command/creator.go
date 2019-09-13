package command

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/hooks/def"
)

type Creator struct{}

func (c *Creator) CanCreate(hook *config.Hook) bool {
	if !hook.HasCommand() {
		return false
	}

	return true
}

func (c *Creator) Create(hook *config.Hook) def.Executor {
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
	}
}
