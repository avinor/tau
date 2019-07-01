package hooks

import (
	"strings"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
)

func Run(source *config.Source, event string, command string) error {
	command = strings.ToLower(command)
	event = strings.ToLower(event)

	for _, hook := range source.Config.Hooks {
		split := strings.Split(hook.TriggerOn, ":")
		hookEvent := strings.ToLower(split[0])

		// Check if this is correct event
		if hookEvent != event {
			continue
		}

		// Check if hook should run for this command
		if len(split) > 1 {
			hookCommands := strings.Split(split[1], ",")
			shouldRun := false
			for _, cmd := range hookCommands {
				if strings.ToLower(cmd) == command {
					shouldRun = true
				}
			}

			if !shouldRun {
				log.Debugf("hook %s should not run for command %s", hook.Type, command)
				continue
			}
		}

	}

	return nil
}

func executeCachedCommand()
