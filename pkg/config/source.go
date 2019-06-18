package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/avinor/tau/pkg/utils"

	log "github.com/sirupsen/logrus"
)

// Source for one file loaded
type Source struct {
	Hash         string
	File         string
	Content      []byte
	Dependencies map[string]*Source
	Config       *Config

	client *Client
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
func NewSource(file string, client *Client) (*Source, error) {
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

	// TODO Potensial error in hash. Should use src, not full path to file
	return &Source{
		Hash:         utils.Hash(file),
		File:         file,
		Content:      b,
		Config:       config,
		Dependencies: map[string]*Source{},
		client: client,
	}, nil
}

// ModuleDirectory where module should be installed, also creates if does not exist
func (src *Source) ModuleDirectory() string {
	path := filepath.Join(src.client.TempDir, "module", src.Hash)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debugf("Creating module directory")
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	return path
}

// DependencyDirectory where dependencies should be resolved
func (src *Source) DependencyDirectory() string {
	path := filepath.Join(src.client.TempDir, "deps", src.Hash)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Debugf("Creating dependency directory")
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	return path
}

func (src *Source) CreateOverrides() error {
	b, err := GetTerraformOverride(src.Config)
	if err != nil {
		return err
	}

	filename := filepath.Join(src.ModuleDirectory(), "tau_override.tf")
	if err := ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (src *Source) CreateInputVariables() error {
	return nil
}
