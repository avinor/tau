package hooks

import (
	"strings"
	"sync"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/getter"
	pstrings "github.com/avinor/tau/pkg/helpers/strings"
	"github.com/avinor/tau/pkg/helpers/ui"
)

type Options struct {
	Getter   *getter.Client
	CacheDir string
}

type Runner struct {
	options *Options

	// cacheLock makes sure only one command can be generated at a time. For thread safety
	hookLock sync.Mutex

	// cache of all created commands
	hooks map[string]*Command
}

func New(options *Options) *Runner {
	if options == nil {
		options = &Options{}
	}

	return &Runner{
		options: options,
	}
}

// Run all hooks in source for a specific event. Command input can filter hooks that should only be run
// got specific terraform commands.
func (r *Runner) Run(file *loader.ParsedFile, event, command string) error {
	for _, hook := range file.Config.Hooks {
		hook, err := r.getHook(file, hook)
		if err != nil {
			return err
		}

		if hook.HasRun() && !hook.Config.DisableCache {
			ui.Debug("%s has already run")
			continue
		}

		if !hook.ShouldRun(event, command) {
			ui.Debug("%s should not run for command %s", hook.Type, command)
			continue
		}

		if err := hook.Run(); err != nil {
			return err
		}

		for key, value := range hook.Output {
			ui.Debug("setting env %s", key)
			file.Env[key] = value
		}
	}

	return nil
}

func (r *Runner) getHook(file *loader.ParsedFile, event string) error {

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
