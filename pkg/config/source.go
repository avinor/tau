package config

import (
	"github.com/avinor/tau/pkg/utils"
	"path"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

// Source for one file loaded
type Source struct {
	Hash         string
	File         string
	Content      []byte
	Dependencies map[string]*Source

	Config *Config
}

// ByDependencies sorts a list of sources by their dependencies
type ByDependencies []*Source

func (a ByDependencies) Len() int {
	return len(a)
}

func (a ByDependencies) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDependencies) Less(i, j int) bool {

	for _, dep := range a[j].Dependencies {
		if dep == a[i] {
			return true
		}
	}

	return false
}

// NewSource creates a new source from a file
func NewSource(file string) (*Source, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config, err := Parser.Parse(b, file)
	if err != nil {
		return nil, err
	}

	log.WithField("indent", 1).Infof("%v loaded", path.Base(file))

	return &Source{
		Hash: utils.Hash(file),
		File:    file,
		Content: b,
		Config:  config,
		Dependencies:    map[string]*Source{},
	}, nil
}

func (src *Source) ModuleDirectory() string {
	return ""
}

func (src *Source) CreateBackendFile() error {
	return nil
}
