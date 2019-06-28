package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	// loaded is a map of already loaded Sources. Will always be checked so same file is
	// not loaded twice. Map key is absolute path of file
	loaded = map[string]*Source{}
)

// Source information about one file loaded from disk. Includes hcl tag for name
// because it is needed when saving SourceFile.
type Source struct {
	Name         string `hcl:"name,label"`
	File         string
	Content      []byte
	Config       *Config
	Env          map[string]string
	Dependencies map[string]*Source
}

// NewSource creates a new source from a file and parse the configuration.
func NewSource(file string, content []byte) (*Source, error) {
	config := &Config{}
	if err := Parse(content, file, config); err != nil {
		return nil, err
	}

	name := filepath.Base(file)

	env := map[string]string{}
	if config.Environment != nil {
		parsed, err := ParseBody(config.Environment.Config)
		if err != nil {
			return nil, err
		}

		for k, v := range parsed {
			env[k] = v.AsString()
		}
	}

	return &Source{
		Name:         name,
		File:         file,
		Content:      content,
		Config:       config,
		Env:          env,
		Dependencies: map[string]*Source{},
	}, nil
}

// GetSourceFromFile returns the Source for file (should be absolute path). If file exists
// in cache it will return the cached item, otherwise it will create the Source and return
// a pointer to new Source.
func GetSourceFromFile(file string) (*Source, error) {
	if _, already := loaded[file]; already {
		return loaded[file], nil
	}

	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	source, err := NewSource(file, b)
	if err != nil {
		return nil, err
	}

	loaded[file] = source
	return source, nil
}

// IsAlreadyLoaded checks if file is already loaded and returns true if is is.
func IsAlreadyLoaded(file string) bool {
	if _, already := loaded[file]; already {
		return true
	}

	return false
}
