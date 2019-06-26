package config

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/getter"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/go-errors/errors"
)

// Loader client for loading sources
type Loader struct {
	options *Options
	loaded  map[string]*Source
	getter  *getter.Client
}

// Options when loading modules
type Options struct {
	WorkingDirectory string
	TempDirectory    string
	Getter           *getter.Client
}

// NewLoader creates a new loader client
func NewLoader(options *Options) *Loader {
	if options.WorkingDirectory == "" {
		options.WorkingDirectory = dir.Working
	}

	if options.Getter == nil {
		options.Getter = getter.New(nil)
	}

	return &Loader{
		options: options,
		loaded:  map[string]*Source{},
		getter:  options.Getter,
	}
}

// Load all sources from src
func (l *Loader) Load(src string, version *string) ([]*Source, error) {
	if src == "" {
		return nil, errors.Errorf("Source is empty")
	}

	log.Info(color.New(color.Bold).Sprint("Loading sources..."))

	sources, err := l.loadSource(src, nil)
	if err != nil {
		return nil, err
	}

	log.Info("")
	log.Info(color.New(color.Bold).Sprint("Loading dependencies..."))
	if err := l.loadDependencies(sources, 0); err != nil {
		return nil, err
	}

	sort.Sort(ByDependencies(sources))

	log.Info("")

	return sources, nil
}

func (l *Loader) loadSource(src string, version *string) ([]*Source, error) {
	dst := dir.Source(l.options.TempDirectory, src)

	if err := l.getter.Get(src, dst, version); err != nil {
		return nil, err
	}

	files, err := l.findModuleFiles(dst)
	if err != nil {
		return nil, err
	}

	sources := []*Source{}
	for _, file := range files {
		source, err := NewSource(file)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

func (l *Loader) loadDependencies(sources []*Source, depth int) error {
	remaining := []*Source{}

	for _, source := range sources {
		if _, ok := l.loaded[source.ContentHash]; !ok {
			remaining = append(remaining, source)
		}
	}

	for _, source := range remaining {
		deps, err := l.loadModuleDependencies(source)
		if err != nil {
			return err
		}

		if err := l.loadDependencies(deps, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) loadModuleDependencies(source *Source) ([]*Source, error) {
	l.loaded[source.ContentHash] = source
	deps := []*Source{}

	for _, dep := range source.Config.Dependencies {
		sources, err := l.loadSource(dep.Source, dep.Version)
		if err != nil {
			return nil, err
		}

		for _, src := range sources {
			if _, ok := l.loaded[src.ContentHash]; !ok {
				deps = append(deps, src)
			} else {
				src = l.loaded[src.ContentHash]
			}

			source.Dependencies[dep.Name] = src
		}
	}

	return deps, nil
}

func (l *Loader) findModuleFiles(dst string) ([]string, error) {

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

			if !fi.IsDir() {
				matches = append(matches, match)
			}
		}
	}

	// log.Debugf("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}
