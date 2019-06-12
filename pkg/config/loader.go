package config

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/avinor/tau/pkg/utils"
	"github.com/go-errors/errors"
	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

// Loader contains all sources loaded
type Loader struct {
	TempDir string
	Pwd     string
	Sources []*Source
	loaded  map[string]*Source
}

// LoadOptions are options when loading modules
type LoadOptions struct {
	LoadSources      bool
	CleanTempDir     bool
	WorkingDirectory string
}

// Load all modules from source
func Load(src string, options *LoadOptions) (*Loader, error) {
	if src == "" {
		return nil, errors.Errorf("Source is empty")
	}

	if options.WorkingDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		options.WorkingDirectory = pwd
	}
	log.Debugf("Current working directory: %v", options.WorkingDirectory)

	tempDir := filepath.Join(options.WorkingDirectory, ".tau", utils.Hash(src))

	if !options.LoadSources {
		loader, err := getLoader(tempDir)
		if err != nil {
			return nil, err
		}

		return loader, nil
	}

	loader := &Loader{
		TempDir: tempDir,
		Pwd:     options.WorkingDirectory,
	}

	if err := loader.loadAllSources(src); err != nil {
		return nil, err
	}

	return loader, nil
}

func getLoader(tempDir string) (*Loader, error) {
	return nil, nil
}

func (l *Loader) Save() error {
	return nil
}

func (l *Loader) loadAllSources(src string) error {
	log.WithField("blank_before", true).Info("Loading modules...")

	sources, err := l.loadSource(src)
	if err != nil {
		return err
	}

	log.WithField("blank_before", true).Info("Loading dependencies...")
	if err := l.loadDependencies(sources, 0); err != nil {
		return err
	}

	sort.Sort(ByDependencies(sources))

	// log.WithField("blank_before", true).Info("Preparing modules...")
	// for _, module := range sources {
	// 	if err := module.Prepare(); err != nil {
	// 		return err
	// 	}
	// }

	l.Sources = sources

	return nil
}

func (l *Loader) loadSource(src string) ([]*Source, error) {
	dst := filepath.Join(l.TempDir, "init", utils.Hash(src))

	if err := l.getSources(src, dst); err != nil {
		return nil, err
	}

	files, err := l.findModuleFiles(dst)
	if err != nil {
		return nil, err
	}

	sources := []*Source{}
	for _, file := range files {
		source, err := NewModule(file)
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
		if _, ok := l.loaded[source.Hash]; !ok {
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
	l.loaded[source.Hash] = source
	deps := []*Source{}

	for _, dep := range source.Config.Dependencies {
		modules, err := l.loadSource(dep.Source)
		if err != nil {
			return nil, err
		}

		for _, mod := range modules {
			hash := mod.Hash()

			if _, ok := l.loaded[hash]; !ok {
				deps = append(deps, mod)
			} else {
				mod = l.loaded[hash]
			}

			mod.deps[dep.Name] = mod
		}
	}

	return deps, nil
}

func (l *Loader) getSources(src, dst string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	log.Debugf("Loading sources for %v", src)

	client := &getter.Client{
		Ctx:  ctx,
		Src:  src,
		Dst:  dst,
		Pwd:  l.pwd,
		Mode: getter.ClientModeAny,
	}

	return client.Get()
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

	log.Debugf("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}
