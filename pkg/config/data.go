package config

import (
	"github.com/hashicorp/hcl2/hcl"
)

// Data sources as defined by terraform. This is just a copy of the terraform model and is
// written to the dependency file as defined, no extra processing done.
type Data struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`

	Config hcl.Body `hcl:",remain"`
}
