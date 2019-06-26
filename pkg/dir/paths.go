package dir

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/strings"
	"os"
	"path/filepath"
)

const (
	// TauPath directory where all processing is done
	TauPath = ".tau"

	// ModulePath is directory where modules are downloaded
	ModulePath = "module"

	// SourcePath directory to download / copy sources
	SourcePath = "source"

	// DependencyPath directory where to process dependencies
	DependencyPath = "dep"
)

// Remove removes the directory
func Remove(path string) {
	os.RemoveAll(path)
}

// TempDir generates a temporary directory for src
func TempDir(pwd, src string) string {
	if src == "" {
		src = strings.SecureRandomAlphaString(16)
	}

	return joinAndCreate(pwd, TauPath, strings.Hash(src))
}

// Module generates a module directory
func Module(tempDir, name string) string {
	return join(tempDir, ModulePath, name)
}

// Source generates a source directory
func Source(tempDir, name string) string {
	return join(tempDir, SourcePath, name)
}

// Dependency generates a dependency directory
func Dependency(tempDir, name string) string {
	return joinAndCreate(tempDir, DependencyPath, name)
}

func join(dir, part, folder string) string {
	if dir == "" {
		dir = Working
	}

	if part == "" {
		log.Fatal("Part directory must be set")
	}

	if folder == "" || folder == "." {
		folder = "root"
	}

	path := filepath.Join(dir, part, folder)
	return path
}

func joinAndCreate(dir, part, folder string) string {
	path := join(dir, part, folder)
	ensureDirectoryExists(path)

	return path
}

func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithField("path", path).Debug("Creating directory")

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
