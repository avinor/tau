package config

import (
	"sort"
	"github.com/avinor/tau/pkg/utils"
	"path/filepath"
	"os"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
)

// LoadOptions are options when loading modules
type LoadOptions struct {
	LoadSources bool
	CleanTempDir bool
	WorkingDirectory string
}

// Load all modules from source
func Load(src string, options *LoadOptions) (*Loader, error) {
	if src == "" {
		return nil, errors.Errorf("Source is empty")
	}

	pwd := options.WorkingDirectory
	if pwd == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		pwd = wd
	}
	log.Debugf("Current working directory: %v", pwd)

	tempDir := filepath.Join(pwd, ".tau", utils.Hash(src))
	
	if !options.LoadSources {

	}

	if err := loader.loadModules(); err != nil {
		return nil, err
	}

	return loader, nil
}

func (l *Loader) Save() error {
	return nil
}

func (l *Loader) get

func (l *Loader) loadModules() error {
	log.WithField("blank_before", true).Info("Loading modules...")

	modules, err := l.loadSource(l.src)
	if err != nil {
		return err
	}

	log.WithField("blank_before", true).Info("Loading dependencies...")
	if err := l.loadDependencies(modules, 0); err != nil {
		return err
	}

	sort.Sort(ByDependencies(modules))

	log.WithField("blank_before", true).Info("Preparing modules...")
	for _, module := range modules {
		if err := module.Prepare(); err != nil {
			return err
		}
	}

	return modules, nil
}

func (l *Loader) loadSource(src string) ([]*Module, error) {
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

func (l *Loader) loadDependencies(modules []*Module, depth int) error {
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

func (l *Loader) loadModuleDependencies(module *Module) ([]*Module, error) {
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