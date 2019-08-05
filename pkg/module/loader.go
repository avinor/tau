package module

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/pkg/errors"
)

var (
	// sourcePathNotFoundError is returned when source could not find any modules
	sourcePathNotFoundError = errors.Errorf("source path not found")

	// moduleRegexp is regular expression to match module files
	moduleRegexp = regexp.MustCompile("(?i).*(\\.hcl|\\.tau)$")

	// autoRegexp is regular expression to match _auto import files
	autoRegexp = regexp.MustCompile("(?i).*_auto(\\.hcl|\\.tau)")

	// moduleMatchFunc for checking if a filename match module pattern. Since regexp
	// does not support look ahead this function makes sure to also check that it does
	// not contain _auto keyword
	moduleMatchFunc = func(str string) bool {
		return moduleRegexp.MatchString(str) && !strings.Contains(str, "_auto")
	}

	// autoMatchFunc checks that the filename is an auto import file (contains _auto)
	autoMatchFunc = func(str string) bool {
		return autoRegexp.MatchString(str)
	}
)

// Loader client for loading sources
type Loader struct {
	options *Options
}

// Options when loading modules. If TempDirectory is not set it will create a random
// temporary directory. This is not advised as it will create a new temporary directory
// for every call to load.
type Options struct {
	WorkingDirectory string

	// MaxDepth to search for dependencies. Should be enough with 1.
	MaxDepth int
}

// NewLoader creates a new loader client
func NewLoader(options *Options) *Loader {
	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	return &Loader{
		options: options,
	}
}

// Load modules from path and return list of all Modules found at path. Path can either
// be a single file or a directory, in which case it will load all files found in
// directory.
func (l *Loader) Load(path string) ([]*Config, error) {
	if path == "" {
		return nil, sourcePathNotFoundError
	}

	ui.Header("Loading sources...")

	sources, err := l.loadFromPath(path)
	if err != nil {
		return nil, err
	}

	if err := l.loadDependencies(sources, 0); err != nil {
		return nil, err
	}

	sortSources(sources)

	return sources, nil
}

// loadFromPath loads all source files matching path pattern and returns the Config
// structs for sources. It does not load dependencies, call loadDependencies on return
// value to load the dependency tree.
func (l *Loader) loadFromPath(path string) ([]*Config, error) {
	path = paths.Abs(l.options.WorkingDirectory, path)

	files, err := findFiles(path, moduleMatchFunc)
	if err != nil {
		return nil, err
	}

	sources := []*Config{}
	for _, file := range files {
		source, err := GetSourceFromFile(file)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

// findFiles searches in path for files matching against a custom matching function. All files
// that match, and are not directories, will be return as result.
func findFiles(path string, matchFunc func(string) bool) ([]string, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		if matchFunc(fi.Name()) {
			return []string{path}, nil
		}
		return nil, nil
	}

	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	matches := []string{}
	for _, file := range files {
		if matchFunc(file) {
			fi, err := os.Stat(file)

			if err != nil {
				return nil, err
			}

			if !fi.IsDir() {
				matches = append(matches, file)
			}
		}
	}

	ui.Debug("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}
