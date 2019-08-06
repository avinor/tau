package config

// Module to import and deploy. Uses go-getter to download source, so supports git repos, http(s)
// sources etc. If version is defined it will assume it is a terraform registry source and try
// to download from registry.
type Module struct {
	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}

// Merge current module with config from source
func (m *Module) Merge(src *Module) error {
	if src == nil {
		return nil
	}

	if src.Source != "" {
		m.Source = src.Source
	}

	if src.Version != nil {
		m.Version = src.Version
	}

	return nil
}

// mergeModules merges only the modules from all configurations in srcs into dest
func mergeModules(dest *Config, srcs []*Config) error {
	for _, src := range srcs {
		if src.Module == nil {
			continue
		}

		if dest.Module == nil {
			dest.Module = src.Module
			continue
		}

		if err := dest.Module.Merge(src.Module); err != nil {
			return err
		}
	}

	return nil
}
