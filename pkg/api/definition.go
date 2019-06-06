package api

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

// Definition of the api
type Definition struct {
	config *Config
	modules []*Module
}

// Config parameters for configuring api definition
type Config struct {
	Source             string
	WorkingDirectory   string
	ExtraArguments     []string
	MaxDependencyDepth int
	CleanTempDir       bool
}

// New returns a new api definition
func New(config *Config) (*Definition, error) {
	if config.Source == "" {
		return nil, errors.Errorf("Source is empty")
	}

	if config.WorkingDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		config.WorkingDirectory = pwd
		log.Debugf("Current working directory: %v", pwd)
	}

	loader, err := newLoader(config)
	if err != nil {
		return nil, err
	}

	modules, err := loader.load()
	if err != nil {
		return nil, err
	}

	return &Definition{
		modules: modules,
		config: config,
	}, nil
}

// Run a terraform command on loaded modules
func (d *Definition) Run(cmd string) error {
	return nil
}

func hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
