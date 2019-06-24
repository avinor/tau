package config

import (
	"io/ioutil"
	"os"

	"github.com/avinor/tau/pkg/strings"
)

// Source for one file loaded
type Source struct {
	File         string
	Content      []byte
	ContentHash  string
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

	return &Source{
		File:         file,
		Content:      b,
		ContentHash:  strings.HashFromBytes(b),
		Config:       config,
		Dependencies: map[string]*Source{},
	}, nil
}
