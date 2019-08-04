package module

import "github.com/avinor/tau/pkg/config"

type Module struct {
	Name         string
	Config       *config.Config
	Env          map[string]string
	Dependencies map[string]*Module
}
