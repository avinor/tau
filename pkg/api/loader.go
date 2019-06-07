package api

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-getter"
	log "github.com/sirupsen/logrus"
)

type loader struct {
	tempDir string
	src     string
	pwd     string
	loaded  map[string]*Module
}

func newLoader(src, tempDir string) (*loader, error) {
	loader := &loader{
		tempDir: tempDir,
		src:     src,
		loaded:  map[string]*Module{},
	}

	if err := loader.setSourceDirectory(); err != nil {
		return nil, err
	}

	return loader, nil
}

func (l *loader) loadModules() ([]*Module, error) {
	log.WithField("blank_before", true).Info("Loading modules...")

	modules, err := l.loadSource(l.src)
	if err != nil {
		return nil, err
	}

	log.WithField("blank_before", true).Info("Loading dependencies...")
	if err := l.loadDependencies(modules, 0); err != nil {
		return nil, err
	}

	sort.Sort(ByDependencies(modules))

	log.WithField("blank_before", true).Info("Preparing modules...")
	for _, module := range modules {
		if err := module.Prepare(); err != nil {
			return nil, err
		}
	}

	return modules, nil
}

func (l *loader) loadSource(src string) ([]*Module, error) {
	dst := filepath.Join(l.tempDir, "init", hash(src))

	if err := l.getSources(src, dst); err != nil {
		return nil, err
	}

	files, err := l.findModuleFiles(dst)
	if err != nil {
		return nil, err
	}

	modules := []*Module{}
	for _, file := range files {
		module, err := NewModule(file)
		if err != nil {
			return nil, err
		}

		modules = append(modules, module)
	}

	return modules, nil
}

func (l *loader) loadDependencies(modules []*Module, depth int) error {
	remaining := []*Module{}

	for _, module := range modules {
		if _, ok := l.loaded[module.Hash()]; !ok {
			remaining = append(remaining, module)
		}
	}

	for _, module := range remaining {
		deps, err := l.loadModuleDependencies(module)
		if err != nil {
			return err
		}

		if err := l.loadDependencies(deps, depth+1); err != nil {
			return err
		}
	}

	return nil
}

func (l *loader) loadModuleDependencies(module *Module) ([]*Module, error) {
	l.loaded[module.Hash()] = module
	deps := []*Module{}

	for _, dep := range module.config.Dependencies {
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

func (l *loader) getSources(src, dst string) error {
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

func (l *loader) findModuleFiles(dst string) ([]string, error) {

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

func (l *loader) setSourceDirectory() error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	getterSource, err := getter.Detect(l.src, pwd, getter.Detectors)
	if err != nil {
		return err
	}

	if strings.Index(getterSource, "file://") == 0 {
		log.Debug("File source detected. Changing source directory")
		rootPath := strings.Replace(getterSource, "file://", "", 1)

		fi, err := os.Stat(rootPath)

		if err != nil {
			return err
		}

		if !fi.IsDir() {
			l.pwd = path.Dir(rootPath)
			l.src = path.Base(rootPath)
		} else {
			l.pwd = rootPath
			l.src = "."
		}

		log.Debugf("New source directory: %v", pwd)
		log.Debugf("New source: %v", l.src)
	}

	return nil
}
