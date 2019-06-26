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
	Dependencies map[string]*Source
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

	config := &Config{}
	if err := Parse(b, file, config); err != nil {
		return nil, err
	}

	name := filepath.Base(file)

	return &Source{
		Name:         name,
		File:         file,
		Content:      b,
		ContentHash:  strings.HashFromBytes(b),
		Config:       config,
		Dependencies: map[string]*Source{},
	}, nil
}
