package loader

import (
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/helpers/ui"
)

// findFiles searches in path for files matching against a custom matching function. All files
// that match, and are not directories, will be return as result.
func findFiles(path string, matchFunc func(string) bool) ([]string, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		if matchFunc(fi.Name()) {
			return []string{path}, nil
		}
		return nil, nil
	}

	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return nil, err
	}

	matches := []string{}
	for _, file := range files {
		if matchFunc(file) {
			fi, err := os.Stat(file)

			if err != nil {
				return nil, err
			}

			if !fi.IsDir() {
				matches = append(matches, file)
			}
		}
	}

	ui.Debug("Found %v template file(s): %v", len(matches), matches)

	return matches, nil
}
