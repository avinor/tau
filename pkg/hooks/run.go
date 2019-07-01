package hooks

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
)

// Run all hooks in source for a specfic event. Command input can filter hooks that should only be run
// got specific terraform commands.
func Run(source *config.Source, event string, command string) error {
	for _, hook := range source.Config.Hooks {
		cmd := GetCommand(source, &hook)

		if !cmd.ShouldRun(event, command) {
			log.Debugf("hook %s should not run for command %s", hook.Type, command)
			continue
		}

		if err := cmd.Run(); err != nil {
			return err
		}

		for key, value := range cmd.Output {
			log.Debugf("setting env %s", key)
			source.Env[key] = value
		}
	}

	return nil
}
