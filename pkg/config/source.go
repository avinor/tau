package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/strings"
)

// Source for one file loaded
type Source struct {
	Name         string `hcl:"name,label"`
	ContentHash  string `hcl:"hash,attr"`
	File         string
	Content      []byte
	Config       *Config
	Env          map[string]string
	Dependencies map[string]*Source
}

// NewSource creates a new source from a file
func NewSource(file string, content []byte) (*Source, error) {
	config := &Config{}
	if err := Parse(content, file, config); err != nil {
		return nil, err
	}

	name := filepath.Base(file)

	env := map[string]string{}
	if config.Environment != nil {
		for k, v := range ParseBody(config.Environment.Config) {
			env[k] = v.AsString()
		}
	}

	return &Source{
		Name:         name,
		File:         file,
		Content:      content,
		ContentHash:  strings.HashFromBytes(content),
		Config:       config,
		Env:          env,
		Dependencies: map[string]*Source{},
	}, nil
}

func NewSourceFromFile(file string) (*Source, error) {
	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return NewSource(file, b)
}
