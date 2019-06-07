package api

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

// Catalog is a collection of modules to deploy
type Catalog struct {
	TempDir      string
	Config       *Config
	SettingsFile string
	Modules      []*Module
}

// Config parameters for configuring api definition
type Config struct {
	WorkingDirectory string
	LoadSources      bool
}

// NewCatalog returns a new catalog from source
func NewCatalog(src string, config *Config) (*Catalog, error) {
	if src == "" {
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

	tempDir := filepath.Join(config.WorkingDirectory, ".tau", hash(src))
	settingsFile := filepath.Join(tempDir, "tau.yaml")

	modules, err := getModules(config.LoadSources)
	if err != nil {
		return nil, err
	}

	return &Catalog{
		TempDir:      tempDir,
		Config:       config,
		SettingsFile: settingsFile,
		Modules:      modules,
	}, nil
}

// Save catalog to temp directory for later processing
func (c *Catalog) Save() error {
	return nil
}

func getModules(loadSources bool) ([]*Module, error) {
	return nil, nil
}

func hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
