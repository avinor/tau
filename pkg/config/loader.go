package config

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/avinor/tau/pkg/paths"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/go-errors/errors"
)

var (
	sourcePathNotFoundError   = errors.Errorf("source path not found")
	dependencySingleFileError = errors.Errorf("dependency must be a single file, cannot be directory")

	moduleRegexp = regexp.MustCompile("(?i).*(\\.hcl|\\.tau)$")
	autoRegexp   = regexp.MustCompile("(?i).*_auto(\\.hcl|\\.tau)")

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
	TempDirectory    string

	// MaxDepth to search for dependencies. Should be enough with 1.
	MaxDepth int
}

// NewLoader creates a new loader client
func NewLoader(options *Options) *Loader {
	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	if options.TempDirectory == "" {
		options.TempDirectory = paths.TempDir(paths.WorkingDir, "")
	}

	return &Loader{
		options: options,
	}
}

// Load sources from path and return list of all Sources found at path. Path can either
// be a single file or a directory, in which case it will load all files found in
// directory.
func (l *Loader) Load(path string) ([]*Source, error) {
	if path == "" {
		return nil, sourcePathNotFoundError
	}

	log.Info(color.New(color.Bold).Sprint("Loading sources..."))

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

// loadFromPath loads all source files matching path pattern and returns the Source
// structs for sources. It does not load dependencies, call loadDependencies on return
// value to load the dependency tree.
func (l *Loader) loadFromPath(path string) ([]*Source, error) {
	path = paths.Abs(l.options.WorkingDirectory, path)

	files, err := findFiles(path, moduleMatchFunc)
	if err != nil {
		return nil, err
	}

	sources := []*Source{}
	for _, file := range files {
		source, err := GetSourceFromFile(file)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

// loadDependencies searches all dependencies for source and recursively loads them into
// sources dependency map. A dependency can only be a single file, it will fail if trying
// to load a dependency that is a directory or resolves to multiple files.
func (l *Loader) loadDependencies(sources []*Source, depth int) error {
	if depth >= l.options.MaxDepth {
		return nil
	}

	for _, source := range sources {
		dir := filepath.Dir(source.SourceFile.File)

		for _, dep := range source.Config.Dependencies {
			path := filepath.Join(dir, dep.Source)
			deps, err := l.loadFromPath(path)

			if err != nil {
				return err
			}

			if len(deps) > 1 {
				return dependencySingleFileError
			}

			if len(deps) == 0 {
				continue
			}

			source.Dependencies[dep.Name] = deps[0]

			if err := l.loadDependencies(deps, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// findFiles searches in path for files matching against a custom matching function. For all files
// that match, and are not directories, it will return as result.
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

	log.Debugf("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}

// sortSources sorts the sources in order of dependencies.
func sortSources(sources []*Source) {
	sort.SliceStable(sources, func(i, j int) bool {
		return lessCompareSources(sources[i], sources[j])
	})
}

// lessCompareSources takes 2 sources and compare them to see if j is less than
// i. It runs recursively to check dependencies of dependency against i. Otherwise
// it will sort them incorrectly
//
// If there is no dependency that it can sort by it sorts alphabetically on name.
// This is just to get a consistent order of all elements.
func lessCompareSources(i, j *Source) bool {
	for _, dep := range j.Dependencies {
		if dep == i || lessCompareSources(i, dep) {
			return true
		}
	}

	return i.Name < j.Name
}
