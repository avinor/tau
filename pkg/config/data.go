package config

import (
	"fmt"

	"github.com/avinor/tau/pkg/config/comp"

	"github.com/hashicorp/hcl2/hcl"
)

// Data sources as defined by terraform. This is just a copy of the terraform model and is
// written to the dependency file as defined, no extra processing done.
type Data struct {
	Type string `hcl:"type,label"`
	Name string `hcl:"name,label"`

	Config hcl.Body `hcl:",remain"`

	comp.Remainer
}

// Merge data source with src data source only if type and name matches.
func (d *Data) Merge(src *Data) error {
	if src == nil {
		return nil
	}

	// do not merge data sources that do not match
	if d.Type != src.Type && d.Name != src.Name {
		return nil
	}

	d.Config = hcl.MergeBodies([]hcl.Body{d.Config, src.Config})

	return nil
}

// mergeDatas merges the data arrays into destination config.
func mergeDatas(dest *Config, srcs []*Config) error {
	datas := map[string]*Data{}

	for _, src := range srcs {
		for _, data := range src.Datas {
			key := fmt.Sprintf("%s.%s", data.Type, data.Name)

			if _, ok := datas[key]; !ok {
				datas[key] = data
				continue
			}

			if err := datas[key].Merge(data); err != nil {
				return err
			}
		}
	}

	for _, data := range datas {
		dest.Datas = append(dest.Datas, data)
	}

	return nil
}
