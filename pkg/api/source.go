package api

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-getter"
)

// Source is parent type for all types that need to load sources
type Source struct {
	src string
	dst string
	pwd string
}

func getSource(src, pwd string) Source {
	return Source{
		src: src,
		dst: filepath.Join(pwd, ".tau", hash(src)),
		pwd: pwd,
	}
}

func (src *Source) loadModules(level Level) ([]*Module, error) {
	if err := src.loadSources(); err != nil {
		return nil, err
	}

	files, err := src.findModuleFiles()
	if err != nil {
		return nil, err
	}

	modules := []*Module{}
	for _, file := range files {
		module, err := NewModule(file, src.pwd, level)
		if err != nil {
			return nil, err
		}

		modules = append(modules, module)
	}

	return modules, nil
}

func (src *Source) loadSources() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	client := &getter.Client{
		Ctx:  ctx,
		Src:  src.src,
		Dst:  src.dst,
		Pwd:  src.pwd,
		Mode: getter.ClientModeAny,
	}

	return client.Get()
}

func (src *Source) findModuleFiles() ([]string, error) {

	matches := []string{}

	for _, ext := range []string{"*.hcl", "*.tau"} {
		m, err := filepath.Glob(filepath.Join(src.dst, ext))
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

	return matches, nil
}

func hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
