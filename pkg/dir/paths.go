package dir

import (
	"log"
	"path/filepath"
	"github.com/avinor/tau/pkg/strings"
	"os"
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
func Module(tempDir, module string) string {
	return joinAndCreate(tempDir, ModulePath, filepath.Base(module))
}

// Source generates a source directory
func Source(tempDir, source string) string {
	return joinAndCreate(tempDir, SourcePath, filepath.Base(source))
}

// Dependency generates a dependency directory
func Dependency(tempDir, dep string) string {
	return joinAndCreate(tempDir, SourcePath, dep)
}

func joinAndCreate(dir, part, folder string) string {
	if dir == "" {
		dir = Working
	}

	if part == "" {
		log.Fatal("Part directory must be set")
	}

	if folder == "" {
		log.Fatal("Folder directory must be set")
	}

	path := filepath.Join(dir, part, folder)
	ensureDirectoryExists(path)

	return path
}

func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}