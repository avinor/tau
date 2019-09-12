package hooks

import (
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/config/loader"
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
func GetCommand(file *loader.ParsedFile, hook *config.Hook) (*Command, error) {
	parsedCommand, err := hook.GetExecutingCommand(filepath.Dir(file.FullPath))
	if err != nil {
		return nil, err
	}

	key := getCacheKey(parsedCommand, hook)
	cacheLock.Lock()
	defer cacheLock.Unlock()

	if _, exists := cache[key]; exists {
		return cache[key], nil
	}

	split := strings.Split(*hook.TriggerOn, ":")
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
		parsedCommand: parsedCommand,
		event:         hookEvent,
		commands:      hookCommands,
	}

	return cache[key], nil
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

func getParsedCommand(file *loader.ParsedFile, hook *config.Hook) (string, error) {
	if hook.HasCommand() {
		if strings.HasPrefix(*hook.Command, ".") {
			workingDir := filepath.Dir(file.FullPath)
			return filepath.Join(workingDir, *hook.Command), nil
		}
	}

	if !hook.HasScript() {
		return "", errors.Errorf("hook command or script is required")
	}

	client := getter.New(paths.WorkingDir)

	if err := client.Get(*hook.Script, file.ModuleDir(), nil); err != nil {
		return "", err
	}

	return "", nil
}