package api

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	maxDependencyDepth = 5
)

type Loader struct {
	Source
	modules map[string]*Module
}

func NewLoader(src string) *Loader {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %s", err)
	}

	return &Loader{
		Source:  getSource(src, pwd),
		modules: map[string]*Module{},
	}
}

func (l *Loader) Load() error {
	modules, err := l.loadModules(Root)
	if err != nil {
		return err
	}

	for _, module := range modules {
		l.modules[module.Hash()] = module
	}

	return l.resolveRemainingDependencies(0)
}

func (l *Loader) resolveRemainingDependencies(depth int) error {
	if depth >= maxDependencyDepth {
		return fmt.Errorf("Max dependency depth reached (%v)", maxDependencyDepth)
	}

	mods := []*Module{}

	for _, m := range l.modules {
		if m.deps == nil {
			mods = append(mods, m)
		}
	}

	if len(mods) == 0 {
		return nil
	}

	for _, mod := range mods {
		if err := mod.resolveDependencies(l.modules); err != nil {
			return err
		}
	}

	return l.resolveRemainingDependencies(depth + 1)
}
