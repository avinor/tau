package config

import (
	"strings"

	"github.com/pkg/errors"
)

var (
	// ValidHookTriggers is a list of valid values for trigger_on
	ValidHookTriggers = []string{"prepare", "finish"}

	// scriptOrCommandIsRequired is returned if command or script is not set
	scriptOrCommandIsRequired = errors.Errorf("hook command or script is required")

	// scriptAndCommandBothDefined is returned if both command and script attribute is defined
	scriptAndCommandBothDefined = errors.Errorf("script and command can not be both defined in hook")

	// triggerOnValueIncorrect is returned if the trigger_on value is incorrect value
	triggerOnValueIncorrect = errors.Errorf("trigger_on has to be one of: %s", strings.Join(ValidHookTriggers, ", "))
)

// Hook describes a hook that should be run at specific time during deployment.
// Can be used to set environment variables or prepare environment before deployment
//
// TriggerOn decides at which event this hook should trigger. On event command specified
// in Command will run. If read_output is set to true it will try to parse the output
// from command (stdout) as key=value pairs and add them to list of environment
// variables that are sent to terraform commands
//
// To prevent same command from running multiple times it will assume that running same command
// multiple times always produce same result and therefore cache output. To prevent this
// set disable_cache = true. It will force the command to run for every source including hook
//
// By default it will fail command if hook fails. To prevent this set fail_on_error = false
type Hook struct {
	Type         string    `hcl:"type,label"`
	TriggerOn    *string   `hcl:"trigger_on,attr"`
	Command      *string   `hcl:"command,attr"`
	Script       *string   `hcl:"script,attr"`
	Arguments    *[]string `hcl:"args,attr"`
	SetEnv       *bool     `hcl:"set_env,attr"`
	FailOnError  *bool     `hcl:"fail_on_error,attr"`
	DisableCache *bool     `hcl:"disable_cache,attr"`
	WorkingDir   *string   `hcl:"working_dir,attr"`
}

// Merge current hook with config from source
func (h *Hook) Merge(src *Hook) error {
	if src == nil {
		return nil
	}

	// do not merge different hook types
	if h.Type != src.Type {
		return nil
	}

	h.TriggerOn = setFirstStringPointer(src.TriggerOn, h.TriggerOn)
	h.Command = setFirstStringPointer(src.Command, h.Command)
	h.Script = setFirstStringPointer(src.Script, h.Script)
	h.WorkingDir = setFirstStringPointer(src.WorkingDir, h.WorkingDir)
	h.SetEnv = setFirstBoolPointer(src.SetEnv, h.SetEnv)
	h.FailOnError = setFirstBoolPointer(src.FailOnError, h.FailOnError)
	h.DisableCache = setFirstBoolPointer(src.DisableCache, h.DisableCache)

	if src.Arguments != nil {
		if h.Arguments == nil {
			h.Arguments = src.Arguments
		} else {
			for _, arg := range *src.Arguments {
				*h.Arguments = append(*h.Arguments, arg)
			}
		}
	}

	return nil
}

// Validate that all required settings are correct
func (h Hook) Validate() (bool, error) {
	if !h.HasCommand() {
		if !h.HasScript() {
			return false, scriptOrCommandIsRequired
		}
	}

	if h.HasCommand() && h.HasScript() {
		return false, scriptAndCommandBothDefined
	}

	if h.TriggerOn == nil {
		return false, triggerOnValueIncorrect
	}

	validTrigger := false
	for _, trigger := range ValidHookTriggers {
		if trigger == strings.ToLower(*h.TriggerOn) {
			validTrigger = true
		}
	}

	if !validTrigger {
		return false, triggerOnValueIncorrect
	}

	return true, nil
}

// HasScript returns true if script is defined
func (h Hook) HasScript() bool {
	return h.Script != nil && *h.Script != ""
}

// HasCommand returns true is command is defined
func (h Hook) HasCommand() bool {
	return h.Command != nil && *h.Command != ""
}

// setFirstStringPointer returns first string that is not empty
func setFirstStringPointer(args ...*string) *string {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}

	return nil
}

// setFirstBoolPointer returns first bool pointer that has a reference
func setFirstBoolPointer(args ...*bool) *bool {
	for _, arg := range args {
		if arg != nil {
			return arg
		}
	}

	return nil
}

// mergeHooks merges the hooks arrays into destination config.
func mergeHooks(dest *Config, srcs []*Config) error {
	hooks := map[string]*Hook{}

	for _, src := range srcs {
		for _, hook := range src.Hooks {
			if _, ok := hooks[hook.Type]; !ok {
				hooks[hook.Type] = hook
				continue
			}

			if err := hooks[hook.Type].Merge(hook); err != nil {
				return err
			}
		}
	}

	for _, hook := range hooks {
		dest.Hooks = append(dest.Hooks, hook)
	}

	return nil
}
