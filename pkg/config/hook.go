package config

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
	TriggerOn    string    `hcl:"trigger_on,attr"`
	Command      string    `hcl:"command,attr"`
	Arguments    *[]string `hcl:"args,attr"`
	SetEnv       *bool     `hcl:"set_env,attr"`
	FailOnError  *bool     `hcl:"fail_on_error,attr"`
	DisableCache *bool     `hcl:"disable_cache,attr"`
}
