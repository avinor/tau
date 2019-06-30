package config

import (
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

var (
	// filePathMustBeAbsError is returned when a file path is relative
	filePathMustBeAbsError = errors.Errorf("file path must be absolute")

	// loaded is a map of already loaded Sources. Will always be checked so same file is
	// not loaded twice. Map key is absolute path of file
	loaded = map[string]*Source{}
)

// Source information about one file loaded from disk. Includes hcl tag for name
// because it is needed when saving SourceFile.
type Source struct {
	SourceFile

	Name         string `hcl:"name,label"`
	Config       *Config
	Env          map[string]string
	Dependencies map[string]*Source
}

// GetSourceFromFile returns the Source for file (should be absolute path). If file exists
// in cache it will return the cached item, otherwise it will create the Source and return
// a pointer to new Source.
func GetSourceFromFile(file string) (*Source, error) {
	if isAlreadyLoaded(file) {
		return loaded[file], nil
	}

	if !filepath.IsAbs(file) {
		return nil, filePathMustBeAbsError
	}

	sf, err := GetSourceFile(file)
	if err != nil {
		return nil, err
	}

	source, err := newSource(sf)
	if err != nil {
		return nil, err
	}

	loaded[file] = source
	return source, nil
}

// newSource creates a new Source struct from a SourceFile. It will parse the config,
// merge together with auto files and read environment variables.
func newSource(source *SourceFile) (*Source, error) {
	name := filepath.Base(source.File)

	auto, err := source.GetAutoImport()
	if err != nil {
		return nil, err
	}

	config, err := source.ConfigMergedWith(auto)
	if err != nil {
		return nil, err
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	env, err := parseEnvironmentVariables(config, source.EvalContext())
	if err != nil {
		return nil, err
	}

	return &Source{
		SourceFile:   *source,
		Name:         name,
		Config:       config,
		Env:          env,
		Dependencies: map[string]*Source{},
	}, nil
}

// parseEnvironmentVariables parses the config and returns all the environment variables
// defined in config.
func parseEnvironmentVariables(config *Config, context *hcl.EvalContext) (map[string]string, error) {
	if config != nil && config.Environment == nil {
		return nil, nil
	}

	values := map[string]cty.Value{}
	if err := ParseBody(config.Environment.Config, context, &values); err != nil {
		return nil, err
	}

	env := map[string]string{}
	for key, value := range values {
		env[key] = value.AsString()
	}

	return env, nil
}

// isAlreadyLoaded checks if file is already loaded and returns true if is is.
func isAlreadyLoaded(file string) bool {
	if _, already := loaded[file]; already {
		return true
	}

	return false
}

func validateConfig(config *Config) error {
	return nil
}
