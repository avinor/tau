package config

import (
	"github.com/avinor/tau/pkg/config/comp"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
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

	comp.Remainer
}

// Merge current backend with source. Source will overwrite settings from current backend
func (b *Backend) Merge(src *Backend) error {
	if src == nil {
		return nil
	}

	if b.Type != src.Type {
		return differentBackendTypes
	}

	body := &hclsyntax.Body{
		Attributes: hclsyntax.Attributes{},
	}
	bb := b.Config.(*hclsyntax.Body)
	sb := src.Config.(*hclsyntax.Body)

	for k, v := range bb.Attributes {
		body.Attributes[k] = v
	}

	for k, v := range sb.Attributes {
		body.Attributes[k] = v
	}

	b.Config = body

	return nil
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
