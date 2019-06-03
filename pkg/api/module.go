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

type Level int

const (
	Root Level = 1 << iota
	Dependency
)

type Module struct {
	Source

	content []byte
	level   Level

	config *config.Config
	deps   map[string]*Module
}

func NewModule(src, pwd string, level Level) (*Module, error) {
	if _, err := os.Stat(src); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(src)
	if err != nil {
		return nil, err
	}

	config, err := config.Parser.Parse(b, src)
	if err != nil {
		return nil, err
	}

	log.WithField("indent", 1).Infof("%v loaded", path.Base(src))

	return &Module{
		Source:  getSource(src, pwd),
		content: b,
		level:   level,
		config:  config,
	}, nil
}

func (m *Module) resolveDependencies(loaded map[string]*Module) error {
	m.deps = map[string]*Module{}

	for _, dep := range m.config.Dependencies {
		source := getSource(dep.Source, m.pwd)
		_, err := source.loadModules(Dependency)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Module) Hash() string {
	md5Ctx := md5.New()
	md5Ctx.Write(m.content)
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
