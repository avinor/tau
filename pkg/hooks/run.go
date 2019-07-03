package hooks

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/ui"
)

// Run all hooks in source for a specfic event. Command input can filter hooks that should only be run
// got specific terraform commands.
func Run(source *config.Source, event string, command string) error {
	for _, hook := range source.Config.Hooks {
		cmd := GetCommand(source, &hook)

		if !cmd.ShouldRun(event, command) {
			ui.Debug("hook %s should not run for command %s", hook.Type, command)
			continue
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		for key, value := range cmd.Output {
			ui.Debug("setting env %s", key)
			source.Env[key] = value
		}
	}

	return nil
}
