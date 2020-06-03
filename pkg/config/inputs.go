package config

import (
	"github.com/hashicorp/hcl/v2"

	"github.com/avinor/tau/pkg/config/comp"
	helperhcl "github.com/avinor/tau/pkg/helpers/hcl"
)

// Inputs that are sent to module for deployment. Config is converted to a map of attributes.
// Supports all functions supported by terraform.
type Inputs struct {
	Config hcl.Body `hcl:",remain"`

	comp.Remainer
}

// Merge current inputs with config from source
func (i *Inputs) Merge(src *Inputs) error {
	if src == nil {
		return nil
	}

	i.Config = helperhcl.MergeBodiesWithOverides([]hcl.Body{i.Config, src.Config})

	return nil
}

// mergeInputs merges only the inputs from all configurations in srcs into dest
func mergeInputs(dest *Config, srcs []*Config) error {
	for _, src := range srcs {
		if src.Inputs == nil {
			continue
		}

		if dest.Inputs == nil {
			dest.Inputs = src.Inputs
			continue
		}

		if err := dest.Inputs.Merge(src.Inputs); err != nil {
			return err
		}
	}

	return nil
}
