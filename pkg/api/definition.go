package api

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-errors/errors"
	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

// Definition of the api
type Definition struct {
	tempDir string
	config  *Config
}

// Config parameters for configuring api definition
type Config struct {
	Source             string
	WorkingDirectory   string
	ExtraArguments     []string
	MaxDependencyDepth int
	CleanTempDir    bool
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

	tempDir := filepath.Join(config.WorkingDirectory, ".tau", hash(config.Source))

	if config.CleanTempDir {
		log.Debugf("Cleaning temp directory...")
		// os.RemoveAll(config.WorkingDirectory)
	}

	def := &Definition{
		tempDir: tempDir,
		config:  config,
	}

	return def, nil
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

func getPwd(src string) (string, string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	log.Debugf("Current working directory: %v", pwd)

	getterSource, err := getter.Detect(src, pwd, getter.Detectors)
	if err != nil {
		return "", "", err
	}

	if strings.Index(getterSource, "file://") == 0 {
		log.Debug("File source detected. Changing working directory")
		rootPath := strings.Replace(getterSource, "file://", "", 1)

		fi, err := os.Stat(rootPath)

		if err != nil {
			return "", "", err
		}

		if !fi.IsDir() {
			pwd = path.Dir(rootPath)
			src = path.Base(rootPath)
		} else {
			pwd = rootPath
			src = "."
		}

		log.Debugf("New working directory: %v", pwd)
		log.Debugf("New source: %v", src)
	}

	return pwd, src, nil
}
