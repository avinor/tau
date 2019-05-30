package config

import (
	hcl2 "github.com/hashicorp/hcl2/hcl"
)

// Config represents the configuration object to read
type Config struct {
	Datas        []Data       `hcl:"data,block"`
	Dependencies []Dependency `hcl:"dependency,block"`
	Backend      *Backend     `hcl:"backend,block"`
	Module       *Module      `hcl:"module,block"`
	Inputs       *Inputs      `hcl:"inputs,block"`

	Remaining hcl2.Body `hcl:",remain"`
}

// Data sources to read from
type Data struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`

	Config hcl2.Body `hcl:",remain"`
}

// Dependency is another module this depends on
type Dependency struct {
	Name string `hcl:"name,label"`

	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}

// Backend for remote state
type Backend struct {
	Type   string    `hcl:"type,label"`
	Config hcl2.Body `hcl:",remain"`
}

// Module defining module to create
type Module struct {
	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}

// Inputs that are converted to terraform.tfvars for module
type Inputs struct {
	Config hcl2.Body `hcl:",remain"`
}