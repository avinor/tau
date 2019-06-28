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
	modulePattern = regexp.MustCompile(".*(\\.hcl|\\.tau)")
	autoPattern   = regexp.MustCompile(".*(\\.hcl|\\.tau)")
)

// Loader client for loading sources
type Loader struct {
	options *Options
}

// Options when loading modules. If TempDirectory is not set it will create a random
// temporary directory. This is not adviced as it will create a new temporary directory
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
		return nil, errors.Errorf("source path is empty")
	}

	log.Info(color.New(color.Bold).Sprint("Loading sources..."))

	sources, err := l.loadFromPath(path)
	if err != nil {
		return nil, err
	}

	if err := l.loadDependencies(sources, 0); err != nil {
		return nil, err
	}

	sort.Sort(ByDependencies(sources))

	log.Info("")

	return sources, nil
}

// loadFromPath loads all source files matching path pattern and returns the Source
// structs for sources. It does not load dependencies, call loadDependencies on return
// value to load the dependency tree.
func (l *Loader) loadFromPath(path string) ([]*Source, error) {
	path = paths.Abs(l.options.WorkingDirectory, path)

	log.Infof("- loading %s", filepath.Base(path))

	files, err := findModuleFiles(path)
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
		dir := filepath.Dir(source.File)

		for _, dep := range source.Config.Dependencies {
			path := filepath.Join(dir, dep.Source)
			deps, err := l.loadFromPath(path)

			if err != nil {
				return err
			}

			if len(deps) > 1 {
				return errors.Errorf("dependency must be a single file, cannot be directory")
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

// findModuleFiles searches in dst for files matching *.hcl or *.tau, or if dst is a single
// file it will just return that as result. Will exclude any files matching _auto as those should
// not be loaded as module sources.
func findModuleFiles(dst string) ([]string, error) {

	fi, err := os.Stat(dst)

	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return []string{dst}, nil
	}

	matches := []string{}
	for _, ext := range []string{"*.hcl", "*.tau"} {
		m, err := filepath.Glob(filepath.Join(dst, ext))
		if err != nil {
			return nil, err
		}

		for _, match := range m {
			fi, err := os.Stat(match)

			if err != nil {
				return nil, err
			}

			if !fi.IsDir() && !strings.Contains(fi.Name(), "_auto") {
				matches = append(matches, match)
			}
		}
	}

	log.Debugf("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}

func findFiles(path, pattern string) ([]string, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		match, err := regexp.MatchString(fi.Name(), pattern)

		if err != nil {
			return nil, err
		} else if match {
			return []string{path}, nil
		}

		return nil, nil
	}

	matches, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		fi, err := os.Stat(match)

		if err != nil {
			return nil, err
		}

		if !fi.IsDir() && !strings.Contains(fi.Name(), "_auto") {
			matches = append(matches, match)
		}
	}
}
