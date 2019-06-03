package api

import (
	"fmt"
)

const (
	maxDependencyDepth = 5
)

type Package struct {
	Source

	modules map[string]*Module
}

func LoadPackage(src, pwd string) (*Package, error) {
	pkg := &Package{
		Source:  getSource(src, pwd),
		modules: map[string]*Module{},
	}

	modules, err := pkg.loadModules(Root)
	if err != nil {
		return nil, err
	}

	for _, module := range modules {
		pkg.modules[module.Hash()] = module
	}

	if err := pkg.resolveRemainingDependencies(0); err != nil {
		return nil, err
	}

	return pkg, nil
}

func (pkg *Package) resolveRemainingDependencies(depth int) error {
	if depth >= maxDependencyDepth {
		return fmt.Errorf("Max dependency depth reached (%v)", maxDependencyDepth)
	}

	mods := []*Module{}

	for _, m := range pkg.modules {
		if m.config == nil {
			mods = append(mods, m)
		}
	}

	if len(mods) == 0 {
		return nil
	}

	for _, mod := range mods {
		if err := mod.resolveDependencies(pkg.modules); err != nil {
			return err
		}
	}

	return pkg.resolveRemainingDependencies(depth + 1)
}
