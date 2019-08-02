package config

// Module to import and deploy. Uses go-getter to download source, so supports git repos, http(s)
// sources etc. If version is defined it will assume it is a terraform registry source and try
// to download from registry.
type Module struct {
	Source  string  `hcl:"source,attr"`
	Version *string `hcl:"version,attr"`
}
