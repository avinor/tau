package config

import (
	hcl2 "github.com/hashicorp/hcl2/hcl"
)

// Config structure for file describing deployment. This includes the module source, inputs
// dependencies, backend etc. One config element is connected to a single deployment
type Config struct {
	Datas        []Data       `hcl:"data,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
	Hooks        []Hook       `hcl:"hook,block"`
	Environment  *Environment `hcl:"environment_variables,block"`
	Backend      *Backend     `hcl:"backend,block"`
	Module       *Module      `hcl:"module,block"`
	Inputs       *Inputs      `hcl:"inputs,block"`
}

// Data sources as defined by terraform. This is just a copy of the terraform model and is
// written to the dependency file as defined, no extra processing done.
type Data struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`

	Config hcl2.Body `hcl:",remain"`
}

// Dependency towards another tau deployment. Source can either be a relative / absolute path
// (start with . or / in that case) to a file or a directory.
//
// For each dependency it will create a remote_state data source to retrieve the values from
// dependency. Backend configuration will be read from the dependency file. To override attributes
// define the backend block in dependency and only define the attributes that should be overridden.
// For instance it can be useful to override token attribute if current module and dependency module
// use different token's for authentication
type Dependency struct {
	Name   string `hcl:"name,label"`
	Source string `hcl:"source,attr"`

	Backend *Backend `hcl:"backend,block"`
}

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

// Backend for remote state storage. This will be added to an override file before running terraform
// init. Any existing backend configuration in module will therefore be overridden.
//
// Supports same attributes as terraform backend configuration. For dependencies this is used to
// configure remote state data source.
type Backend struct {
	Type   string    `hcl:"type,label"`
	Config hcl2.Body `hcl:",remain"`
}

// Module to import and deploy. Uses go-getter to download source, so supports git repos, http(s)
// sources etc. If version is defined it will assume it is a terraform registry source and try
// to download from registry.
type Module struct {
	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}

// Environment variables that should be added to the shell commands run (specfifically terraform).
// Define variables using attributes, blocks not supported
type Environment struct {
	Config hcl2.Body `hcl:",remain"`
}

// Inputs that are sent to module for deployment. Config is converted to a map of attributes.
// Supports all functions supported by terraform.
type Inputs struct {
	Config hcl2.Body `hcl:",remain"`
}
