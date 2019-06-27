package config

import (
	hcl2 "github.com/hashicorp/hcl2/hcl"
)

// Config structure for file describing deployment. This includes the module source, inputs
// dependencies, backend etc. One config element is connected to a single deployment
type Config struct {
	Datas        []Data       `hcl:"data,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
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
// (start with . or / in that case) or can also be any file that can be retrieved with go-getter.
//
// For each dependency it will create a remote_state data source to retrieve the values from
// dependency. Backend configuration will be read from the dependency file. To override attributes
// define the backend block in dependency and only define the attributes that should be overridden.
// For instance it can be useful to override token attribute if current module and dependency module
// use different token's for authentication
type Dependency struct {
	Name string `hcl:"name,label"`

	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`

	Backend *Backend `hcl:"backend,block"`
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

// Module defining module to create
type Module struct {
	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}

// Environment variables
type Environment struct {
	Config hcl2.Body `hcl:",remain"`
}

// Inputs that are converted to terraform.tfvars for module
type Inputs struct {
	Config hcl2.Body `hcl:",remain"`
}
