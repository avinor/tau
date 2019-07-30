package hooks

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/avinor/tau/pkg/config"
	pstrings "github.com/avinor/tau/pkg/helpers/strings"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

// Command is an internal representation of a hook command. Used to execute the hook
type Command struct {
	Hook   *config.Hook
	HasRun bool
	Output map[string]string

	parsedCommand string
	lock          sync.Mutex
	event         string
	commands      []string
	err           error
}

var (
	// cacheLock makes sure only one command can be generated at a time. For thread safety
	cacheLock = sync.Mutex{}

	// cache of all created commands
	cache = map[string]*Command{}

	outputRegexp = regexp.MustCompile("(?m:^\\s*\"?([^\"=\\s]*)\"?\\s*=\\s*\"?([^\"\\n]*)\"?$)")
)

// GetCommand creates a new command or return command from cache if it has already been
// created before. Cache key is based on full path for command including arguments
func GetCommand(source *config.Source, hook *config.Hook) *Command {
	command := hook.Command
	if strings.HasPrefix(hook.Command, ".") {
		workingDir := filepath.Dir(source.File)
		command = filepath.Join(workingDir, hook.Command)
	}

	key := getCacheKey(command, hook)
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if _, exists := cache[key]; exists {
		return cache[key]
	}

	split := strings.Split(hook.TriggerOn, ":")
	hookEvent := strings.ToLower(split[0])
	hookCommands := []string{}

	if len(split) > 1 {
		for _, cmd := range strings.Split(split[1], ",") {
			hookCommands = append(hookCommands, strings.ToLower(cmd))
		}
	}

	cache[key] = &Command{
		Hook:          hook,
		Output:        map[string]string{},
		parsedCommand: command,
		event:         hookEvent,
		commands:      hookCommands,
	}

	return cache[key]
}

// ShouldRun checks if for a given event and command this Command should run
func (c *Command) ShouldRun(event string, command string) bool {
	event = strings.ToLower(event)
	command = strings.ToLower(command)

	if c.event != event {
		return false
	}

	if len(c.commands) == 0 {
		return true
	}

	for _, cmd := range c.commands {
		if cmd == command {
			return true
		}
	}

	return false
}

// Run the command hook and parse output
func (c *Command) Run() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.HasRun {
		return c.err
	}

	buffer := &processors.Buffer{}
	logp := processors.NewUI(ui.Error)

	options := &shell.Options{
		Stdout: shell.Processors(buffer),
		Stderr: shell.Processors(logp),
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
