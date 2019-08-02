package config

import (
	"github.com/hashicorp/hcl2/hcl"
)

// Inputs that are sent to module for deployment. Config is converted to a map of attributes.
// Supports all functions supported by terraform.
type Inputs struct {
	Config hcl.Body `hcl:",remain"`
}
