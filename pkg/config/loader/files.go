package loader

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/avinor/tau/pkg/helpers/ui"
)

var (
	// deleteFileRegex is regular expression to match files that should be deleted
	deleteFileRegex = regexp.MustCompile("(?i)^(delete|destroy)_?")
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

// shouldDeleteFile checks if the file should be deleted and returns true if it
// should. It will also return an altered filename if file should be deleted.
// Ignore altered filename if file should not be deleted
func shouldDeleteFile(filename string) (bool, string) {
	base := filepath.Base(filename)

	if !deleteFileRegex.MatchString(base) {
		return false, filename
	}

	dir := filepath.Dir(filename)
	altered := deleteFileRegex.ReplaceAllString(base, "")

	return true, filepath.Join(dir, altered)
}
