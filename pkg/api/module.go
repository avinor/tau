package api

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"

	"github.com/avinor/tau/pkg/config"
	log "github.com/sirupsen/logrus"
)

// Module is a single module that should be deployed or a dependency for another module
type Module struct {
	file    string
	content []byte
	deps    map[string]*Module

	config     *config.Config
	initConfig *config.InitConfig
	values     *config.ValuesConfig
}

// ByDependencies sorts a list of modules by their dependencies
type ByDependencies []*Module

func (a ByDependencies) Len() int {
	return len(a)
}

func (a ByDependencies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDependencies) Less(i, j int) bool {

	for _, dep := range a[j].deps {
		if dep == a[i] {
			return true
		}
	}

	return false
}

// NewModule creates a new module from a file
func NewModule(file string) (*Module, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config, err := config.Parser.Parse(b, file)
	if err != nil {
		return nil, err
	}

	log.WithField("indent", 1).Infof("%v loaded", path.Base(file))

	return &Module{
		file:    file,
		content: b,
		config:  config,
		deps:    map[string]*Module{},
	}, nil
}

// Prepare module
func (m *Module) Prepare() error {
	log.WithField("indent", 1).Infof("%v", path.Base(m.file))

	return nil
}

func (m *Module) GetBackendArgs() []string {
	return nil
}

// Hash generates a hash of modules content file
func (m *Module) Hash() string {
	md5Ctx := md5.New()
	md5Ctx.Write(m.content)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
