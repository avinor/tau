package config

import (
	"github.com/avinor/tau/pkg/sources"
	"os"
	"path/filepath"
	"sort"

	"github.com/avinor/tau/pkg/utils"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

// Loader client for loading sources
type Loader struct {
	tempDir string
	loaded  map[string]*Source
	options *Options
}

// Options are options when loading modules
type Options struct {
	WorkingDirectory string
}

// NewLoader creates a new loader client
func NewLoader(tempDir string, options *Options) *Loader {
	if options.WorkingDirectory == "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.WithError(err).Fatal("Unable to get working directory")
		}

		options.WorkingDirectory = pwd
	}
	log.Debugf("Current working directory: %v", options.WorkingDirectory)

	tempDir := filepath.Join(options.WorkingDirectory, ".tau", utils.Hash(profile))

	return &Client{
		TempDir: tempDir,
		loaded:  map[string]*Source{},
		options: options,
	}
}

// Load all sources from src
func (l *Loader) Load(src string, version *string) ([]*Source, error) {
	if src == "" {
		return nil, errors.Errorf("Source is empty")
	}

	cSrc, cPwd := sources.ResolveDirectory(src)
	sClient := sources.New(&sources.Options{
		WorkingDirectory: cPwd,
	})

	log.WithField("blank_before", true).Info("Loading modules...")

	sources, err := c.loadSource(sClient, cSrc, nil)
	if err != nil {
		return nil, err
	}

	log.WithField("blank_before", true).Info("Loading dependencies...")
	if err := c.loadDependencies(sClient, sources, 0); err != nil {
		return nil, err
	}

	sort.Sort(ByDependencies(sources))

	return sources, nil
}

func (l *Loader) LoadTempDir() ([]*Source, error) {
	return nil, nil
}

// Save loaded sources in temp directory
func (c *Client) Save(sources []*Source) error {
	return nil
}

func (l *Loader) CleanTempDir() {
	log.Debugf("Removing temp directory")
	os.RemoveAll(c.TempDir)
}

func (c *Client) loadSource(sClient *sources.Client, src string, version *string) ([]*Source, error) {
	dst := filepath.Join(c.TempDir, "sources", utils.Hash(src))

	if err := sClient.Get(src, dst, version); err != nil {
		return nil, err
	}

	files, err := c.findModuleFiles(dst)
	if err != nil {
		return nil, err
	}

	sources := []*Source{}
	for _, file := range files {
		source, err := NewSource(file, c)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

func (c *Client) loadDependencies(sClient *sources.Client, sources []*Source, depth int) error {
	remaining := []*Source{}

	for _, source := range sources {
		if _, ok := c.loaded[source.Hash]; !ok {
			remaining = append(remaining, source)
		}
	}

	for _, source := range remaining {
		deps, err := c.loadModuleDependencies(sClient, source)
		if err != nil {
			return err
		}

		if err := c.loadDependencies(sClient, deps, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) loadModuleDependencies(sClient *sources.Client, source *Source) ([]*Source, error) {
	c.loaded[source.Hash] = source
	deps := []*Source{}

	for _, dep := range source.Config.Dependencies {
		sources, err := c.loadSource(sClient, dep.Source, dep.Version)
		if err != nil {
			return nil, err
		}

		for _, src := range sources {
			if _, ok := c.loaded[src.Hash]; !ok {
				deps = append(deps, src)
			} else {
				src = c.loaded[src.Hash]
			}

			source.Dependencies[dep.Name] = src
		}
	}

	return deps, nil
}

func (c *Client) findModuleFiles(dst string) ([]string, error) {

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

func (c *Client) getSavedSources(src string) ([]*Source, error) {
	return nil, nil
}