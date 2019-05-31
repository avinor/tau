package api

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/avinor/tau/pkg/config"
	getter "github.com/hashicorp/go-getter"
)

const (
	maxDependencyDepth = 5
)

type Loader struct {
	hashing func(string) string
	pwd     string
}

type Module struct {
	src string // ./examples
	dst string // .tau/dfaf

	config *config.Config
	// Executor
}

func NewLoader() *Loader {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %s", err)
	}

	return &Loader{
		hashing: getDstPath,
		pwd:     dir,
	}
}

func (l *Loader) Load(src string) error {
	dst, err := l.loadSources(src)

	if err != nil {
		return err
	}

	l.resolveModules(dst)

	// Load src into dst
	// Resolve all .hcl / .tau files
	// Create config / module for all files

	// Check dependencies
	// Load unresolved dependencies into dst
	// Create config for dependencies

	return nil
}

func (l *Loader) GetModules() ([]*Module, error) {
	return nil, nil
}

func (l *Loader) loadSources(src string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	dst := fmt.Sprintf(".tau/%s", l.hashing(src))
	var err error

	// Try with .hcl and .tau extension if cannot find file
	for _, extsrc := range []string{src, fmt.Sprintf("%s.tau", src), fmt.Sprintf("%s.hcl", src)} {
		client := &getter.Client{
			Ctx:  ctx,
			Src:  extsrc,
			Dst:  dst,
			Pwd:  l.pwd,
			Mode: getter.ClientModeAny,
		}

		if err = client.Get(); err == nil {
			return dst, nil
		}
	}

	return dst, err
}

func (l *Loader) resolveModules(src string) error {

	matches, err := resolveModuleExt(src, []string{"*.hcl", "*.tau"})

	if err != nil {
		return err
	}

	for _, match := range matches {
		log.Info(match)
	}

	return nil
}

func resolveModuleExt(src string, exts []string) ([]string, error) {
	matches := []string{}

	for _, ext := range exts {
		m, err := filepath.Glob(filepath.Join(src, ext))
		if err != nil {
			return nil, err
		}

		matches = append(matches, m...)
	}

	return matches, nil
}

func getDstPath(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
