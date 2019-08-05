package module

import "github.com/avinor/tau/pkg/config"

type Config struct {
	Name string

	File         *config.File
	Env          map[string]string
	Dependencies map[string]*Config
}
