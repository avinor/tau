package config

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

func (c *Config) Merge(src *Config) error {
	if src == nil {
		return nil
	}

	return nil
}

func (c Config) Validate() (bool, error) {
	return true, nil
}
