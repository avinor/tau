package paths

import (
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/helpers/ui"
)

const (
	// TauPath directory where all processing is done
	TauPath = ".tau"
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

// Join the paths in arguments together and returned the concated path. If any element is empty
// or "." it will replace that with "root", just to ensure it is a valid pathname. If first element
// is empty or not absolute path it will prepend the working directory to path.
func Join(path ...string) string {
	paths := []string{}

	for idx, item := range path {
		if idx == 0 {
			if !filepath.IsAbs(item) {
				item = filepath.Join(WorkingDir, item)
			}
		}

		if item == "" || item == "." {
			item = "root"
		}

		paths = append(paths, item)
	}

	return filepath.Join(paths...)
}

// JoinAndCreate joins the paths together and makes sure the path exists
func JoinAndCreate(path ...string) string {
	joined := Join(path...)
	EnsureDirectoryExists(joined)

	return joined
}

// EnsureDirectoryExists makes sure the entire path exists, all parent folders too.
func EnsureDirectoryExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ui.Debug("Creating directory %s", path)

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			ui.Fatal("%v", err)
		}
	}
}
