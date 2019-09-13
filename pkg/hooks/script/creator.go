package script

import (
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/hooks/command"
	"github.com/avinor/tau/pkg/hooks/def"
)

// Creator that can handle script download and execution
type Creator struct {
	Options *def.Options
}

// CanCreate checks if hook has a script reference to see if this execution should be run.
func (c *Creator) CanCreate(hook *config.Hook) bool {
	if !hook.HasScript() {
		return false
	}

	return true
}

// Create a new executor from hook that first downloads script then returns a command executor
func (c *Creator) Create(hook *config.Hook) (def.Executor, error) {
	var script string
	var arguments []string
	var workingDir string

	script = *hook.Script

	if hook.Arguments != nil {
		arguments = *hook.Arguments
	}

	if hook.WorkingDir != nil {
		workingDir = *hook.WorkingDir
	}

	dst := filepath.Join(c.Options.CacheDir, hook.Type)
	cmd := filepath.Join(dst, filepath.Base(script))

	if err := c.Options.Getter.GetFile(script, cmd); err != nil {
		return nil, err
	}

	if err := os.Chmod(cmd, 0755); err != nil {
		return nil, err
	}

	return &command.Executor{
		Command:    cmd,
		Arguments:  arguments,
		WorkingDir: workingDir,
	}, nil
}
