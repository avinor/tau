package config

import (
	"github.com/hashicorp/hcl2/hcl"
	"github.com/pkg/errors"
)

var (
	differentBackendTypes = errors.Errorf("cannot merge backends with different types")
)

// Backend for remote state storage. This will be added to an override file before running terraform
// init. Any existing backend configuration in module will therefore be overridden.
//
// Supports same attributes as terraform backend configuration. For dependencies this is used to
// configure remote state data source.
type Backend struct {
	Type   string   `hcl:"type,label"`
	Config hcl.Body `hcl:",remain"`
}

// Merge current backend with source. Source will overwrite settings from current backend
func (b *Backend) Merge(src *Backend) error {
	if src == nil {
		return nil
	}

	if b.Type != src.Type {
		return differentBackendTypes
	}

	b.Config = hcl.MergeBodies([]hcl.Body{b.Config, src.Config})

	return nil
}

// Validate that all settings in backend are correct and required settings are configured
func (b Backend) Validate() (bool, error) {
	return true, nil
}

// mergeBackends merges only the backends from all configurations in srcs into dest
func mergeBackends(dest *Config, srcs []*Config) error {
	for _, src := range srcs {
		if src.Backend == nil {
			continue
		}

		if dest.Backend == nil {
			dest.Backend = src.Backend
			continue
		}

		if err := dest.Backend.Merge(src.Backend); err != nil {
			return err
		}
	}

	return nil
}