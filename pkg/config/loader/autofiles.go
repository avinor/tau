package loader

import (
	"path/filepath"
	"regexp"

	"github.com/avinor/tau/pkg/config"
)

var (
	// autoRegexp is regular expression to match _auto import files
	autoRegexp = regexp.MustCompile("(?i).*_auto(\\.hcl|\\.tau)")

	// autoMatchFunc checks that the filename is an auto import file (contains _auto)
	autoMatchFunc = func(str string) bool {
		return autoRegexp.MatchString(str)
	}

	// autoImportPaths is a cache of auto imported files. Key is the path where to search for auto import
	// files.
	autoImportPaths = map[string][]*config.File{}
)

// AddAutoImports searches in the directory for file for any auto imports and add them to the
// list of children.
func AddAutoImports(file *config.File) error {
	dir := filepath.Dir(file.FullPath)

	if _, exists := autoImportPaths[dir]; exists {
		addAutoChildren(file, autoImportPaths[dir])
		return nil
	}

	autoFiles, err := findFiles(dir, autoMatchFunc)
	if err != nil {
		return err
	}

	cacheList := []*config.File{}
	for _, af := range autoFiles {
		configFile, err := config.NewFile(af)
		if err != nil {
			return err
		}

		cacheList = append(cacheList, configFile)
	}

	autoImportPaths[dir] = cacheList

	addAutoChildren(file, cacheList)
	return nil
}

// addAutoChildren calls AddChild for each children on the config file `file`
func addAutoChildren(file *config.File, children []*config.File) {
	for _, child := range children {
		file.AddChild(child)
	}
}
