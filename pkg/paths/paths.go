package paths

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/strings"
)

const (
	// TauPath directory where all processing is done
	TauPath = ".tau"

	// ModulePath is directory where modules are downloaded
	ModulePath = "module"

	// DependencyPath directory where to process dependencies
	DependencyPath = "dep"
)

// Remove the directory and all its subdirectories
func Remove(path string) {
	os.RemoveAll(path)
}

// Abs returns the absolute path for file, based on
func Abs(pwd, path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	return filepath.Join(pwd, path)
}

// TempDir generates a temporary directory for src using pwd as base directory.
// If no src directory is set it will randomly generate a new temp directory.
func TempDir(pwd, src string) string {
	if src == "" {
		src = strings.SecureRandomAlphaString(16)
	}

	return joinAndCreate(pwd, TauPath, strings.Hash(src))
}

// ModuleDir returns the directory where modules should be stored. It will not
// create the directory as that has to be done by go-getter that is using this
// directory.
func ModuleDir(tempDir, name string) string {
	return join(tempDir, ModulePath, name, false)
}

// DependencyDir creates and return the directory for dependencies.
func DependencyDir(tempDir, name string) string {
	return joinAndCreate(tempDir, DependencyPath, name)
}

// join all parts to a new directory path. If root dir is not set it will use
// working directory. Part is required and will fail if not set.
// If folder is not set or is "." it will use "root" as name. Check for "." is
// because its possible to run "tau init ." and this command will fail when trying
// to create a folder called ".". So it uses root instead
//
// Since its difficult to test log.Fatal (it does os.Exit(1)) it just takes an
// additional parameter if it instead should panic. Set shouldPanic to true
// for tests, and false for real usage
func join(dir, part, folder string, shouldPanic bool) string {
	if dir == "" {
		dir = WorkingDir
	}

	if part == "" {
		if shouldPanic {
			panic("oh no...")
		}
		log.Fatal("Part directory must be set")
	}

	if folder == "" || folder == "." {
		folder = "root"
	}

	path := filepath.Join(dir, part, filepath.Base(folder))
	return path
}

// joinAndCreate does same as join, but creates the directory if it does not exist
func joinAndCreate(dir, part, folder string) string {
	path := join(dir, part, folder, false)
	ensureDirectoryExists(path)

	return path
}

// ensureDirectoryExists makes sure the entire path exists, all parent folders too.
func ensureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.WithField("path", path).Debug("Creating directory")

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("%v", err)
		}
	}
}
