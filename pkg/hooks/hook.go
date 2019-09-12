package hooks

import (
	"strings"

	"github.com/avinor/tau/pkg/config"
	pstrings "github.com/avinor/tau/pkg/helpers/strings"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

type Hook struct {
	executor *Executor
}

func (h *Hook) HasRun()

// ShouldRun checks if for a given event and command this Command should run
func (h *Hook) ShouldRun(event string, command string) bool {
	event = strings.ToLower(event)
	command = strings.ToLower(command)

	if h.event != event {
		return false
	}

	if len(h.commands) == 0 {
		return true
	}

	for _, cmd := range h.commands {
		if cmd == command {
			return true
		}
	}

	return false
}

// Run the command hook and parse output
func (h *Hook) Run() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	buffer := &processors.Buffer{}
	logp := processors.NewUI(ui.Error)

	options := &shell.Options{
		Stdout: shell.Processors(buffer),
		Stderr: shell.Processors(logp),
	}

	if c.Hook.WorkingDir != nil && *c.Hook.WorkingDir != "" {
		options.WorkingDirectory = *c.Hook.WorkingDir
	}

	args := []string{}
	if c.Hook.Arguments != nil {
		args = append(args, *c.Hook.Arguments...)
	}

	ui.Info("- %s", c.Hook.Type)

	if err := shell.Execute(options, c.parsedCommand, args...); err != nil {
		if c.Hook.FailOnError != nil && !*c.Hook.FailOnError {
			c.HasRun = true
			return nil
		}

		return err
	}

	if c.Hook.SetEnv != nil && *c.Hook.SetEnv {
		for key, value := range parseOutput(buffer.String()) {
			c.Output[key] = value
		}
	}

	c.HasRun = true

	return nil
}

// parseOutput as key=value. If a line cannot be parsed as key=value it will be ignored
func parseOutput(output string) map[string]string {
	matches := outputRegexp.FindAllStringSubmatch(output, -1)
	values := map[string]string{}

	if len(matches) == 0 {
		return values
	}

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		values[match[1]] = match[2]
	}

	return values
}

// getCacheKey returns a unique cache key for a given command with arguments. If disable_cache
// is set it will generate a random key to make sure it creates new instances
func getCacheKey(command string, hook *config.Hook) string {
	if hook.DisableCache != nil && *hook.DisableCache {
		return pstrings.SecureRandomAlphaString(16)
	}

	if hook.Arguments == nil {
		return command
	}

	return strings.Join(append([]string{command}, *hook.Arguments...), "_")
}
