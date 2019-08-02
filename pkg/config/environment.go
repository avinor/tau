package config

import (
	"github.com/hashicorp/hcl2/hcl"
)

// Environment variables that should be added to the shell commands run (specfifically terraform).
// Define variables using attributes, blocks not supported
type Environment struct {
	Config hcl.Body `hcl:",remain"`
}
