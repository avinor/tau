package config

import (
	"github.com/avinor/tau/pkg/strings"
	"io/ioutil"
	"os"
)

// Source for one file loaded
type Source struct {
	File         string
	Content      []byte
	ContentHash string
	Dependencies map[string]*Source
	Config       *Config
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

	config, err := Parse(b, file)
	if err != nil {
		return nil, err
	}

	// log.WithField("indent", 1).Infof("%v loaded", path.Base(file))

	return &Source{
		File:         file,
		Content:      b,
		ContentHash: strings.HashFromBytes(b),
		Config:       config,
		Dependencies: map[string]*Source{},
	}, nil
}

// ModuleDirectory where module should be installed, also creates if does not exist
// func (src *Source) ModuleDirectory() string {
// 	path := filepath.Join(src.client.TempDir, "module", src.Hash)

// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		log.Debugf("Creating module directory")
// 		if err := os.MkdirAll(path, os.ModePerm); err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	return path
// }

// // DependencyDirectory where dependencies should be resolved
// func (src *Source) DependencyDirectory() string {
// 	path := filepath.Join(src.client.TempDir, "deps", src.Hash)

// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		log.Debugf("Creating dependency directory")
// 		if err := os.MkdirAll(path, os.ModePerm); err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	return path
// }

// func (src *Source) CreateOverrides() error {
// 	b, err := GetTerraformOverride(src.Config)
// 	if err != nil {
// 		return err
// 	}

// 	filename := filepath.Join(src.ModuleDirectory(), "tau_override.tf")
// 	if err := ioutil.WriteFile(filename, b, os.ModePerm); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (src *Source) CreateInputVariables() error {
// 	return nil
// }
