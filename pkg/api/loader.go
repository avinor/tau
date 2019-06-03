package api

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/go-getter"
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
	source, err := getRootSource(src)
	if err != nil {
		log.Fatalf("Unable to get root source: %s", err)
	}

	return &Loader{
		Source:  *source,
		modules: map[string]*Module{},
	}
}

func (l *Loader) Load() error {
	log.WithField("blank_before", true).Info("Loading modules...")

	modules, err := l.loadModules(Root)
	if err != nil {
		return err
	}

	for _, module := range modules {
		l.modules[module.Hash()] = module
	}

	log.WithField("blank_before", true).Info("Loading dependencies...")
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

func getRootSource(src string) (*Source, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	log.Debugf("Current working directory: %v", pwd)

	getterSource, err := getter.Detect(src, pwd, getter.Detectors)
	if err != nil {
		return nil, err
	}

	if strings.Index(getterSource, "file://") == 0 {
		log.Debug("File source detected. Changing working directory")
		rootPath := strings.Replace(getterSource, "file://", "", 1)

		fi, err := os.Stat(rootPath)

		if err != nil {
			return nil, err
		}

		if !fi.IsDir() {
			pwd = path.Dir(rootPath)
			src = path.Base(rootPath)
		} else {
			pwd = rootPath
			src = "."
		}

		log.Debugf("New working directory: %v", pwd)
		log.Debugf("New source: %v", src)
	}

	source := getSource(src, pwd)
	return &source, nil
}
