package hooks

import (
	"fmt"

	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/ui"
)

// Run all hooks in source for a specific event. Command input can filter hooks that should only be run
// got specific terraform commands.
func Run(file *loader.ParsedFile, event string, command string) error {
	for _, hook := range file.Config.Hooks {
		cmd := GetCommand(file, hook)

		if !cmd.ShouldRun(event, command) {
			ui.Debug("hook %s should not run for command %s", hook.Type, command)
			continue
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		for key, value := range cmd.Output {
			ui.Debug("setting env %s", key)
			file.Env[key] = value
		}
	}

	return nil
}

// RunAll executes the hook `event` for all files in collection
func RunAll(files loader.ParsedFileCollection, event string, command string) error {
	ui.Header(fmt.Sprintf("Executing %s hook...", event))
	for _, file := range files {
		if err := Run(file, event, command); err != nil {
			return err
		}
	}

	return nil
}
