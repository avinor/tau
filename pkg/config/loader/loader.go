package loader

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/pkg/errors"
)

var (
	// sourcePathNotFoundError is returned when source could not find any modules
	sourcePathNotFoundError = errors.Errorf("source path not found")

	// dependencySingleFileError is returned if a dependency resolves to multiple files
	dependencySingleFileError = errors.Errorf("dependency must be a single file, cannot be directory")

	// moduleRegexp is regular expression to match module files
	moduleRegexp = regexp.MustCompile("(?i).*(\\.hcl|\\.tau)$")

	// moduleMatchFunc for checking if a filename match module pattern. Since regexp
	// does not support look ahead this function makes sure to also check that it does
	// not contain _auto keyword
	moduleMatchFunc = func(str string) bool {
		return moduleRegexp.MatchString(str) && !strings.Contains(strings.ToLower(str), "_auto")
	}
)

// Loader client for loading sources
type Loader struct {
	options *Options

	// loaded is a map of already loaded ParsedFile. Will always be checked so same file is
	// not loaded twice. Map key is absolute path of file
	loaded map[string]*ParsedFile
}

// Options when loading modules. WorkingDirectory is directory where it will search for
// files. If this is not set it will default to current working directory
type Options struct {
	WorkingDirectory string

	// TauDirectory is the directory where tau will store temporary files. This is required
	// for ParsedFile to know location where to download and store files.
	TauDirectory string

	// CacheDirectory where to download cached resources and scripts for processing
	CacheDirectory string

	// Getter to retrieve source code with
	Getter *getter.Client

	// MaxDepth to search for dependencies. Should be enough with 1.
	MaxDepth int
}

// New creates a new loader client with options
func New(options *Options) *Loader {
	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	if options.TauDirectory == "" {
		path.Join(options.WorkingDirectory, paths.TauPath)
	}

	return &Loader{
		options: options,
		loaded:  map[string]*ParsedFile{},
	}
}

// Load files from path and return list of all ParsedFile found at path. Path can either
// be a single file or a directory, in which case it will load all files found in
// directory.
func (l *Loader) Load(paths []string) (ParsedFileCollection, error) {
	files := make([]*ParsedFile, 0)

	for _, path := range paths {
		if path == "" {
			return nil, sourcePathNotFoundError
		}

		loaded, err := l.loadFromPath(path)
		if err != nil {
			return nil, err
		}

		files = append(files, loaded...)
	}

	if err := l.loadDependencies(files, 0); err != nil {
		return nil, err
	}

	return files, nil
}

// loadFromPath loads all files matching path pattern and returns the ParsedFile
// structs for files. It does not load dependencies, call loadDependencies on return
// value to load the dependency tree.
func (l *Loader) loadFromPath(path string) ([]*ParsedFile, error) {
	path = paths.Abs(l.options.WorkingDirectory, path)

	sources, err := findFiles(path, moduleMatchFunc)
	if err != nil {
		return nil, err
	}

	files := []*ParsedFile{}
	for _, file := range sources {
		file, err := l.getParsedFile(file)
		if err != nil {
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}

// loadDependencies searches all dependencies for files and recursively loads them into
// sources dependency map. A dependency can only be a single file, it will fail if trying
// to load a dependency that is a directory or resolves to multiple files.
func (l *Loader) loadDependencies(files []*ParsedFile, depth int) error {
	if depth >= l.options.MaxDepth {
		return nil
	}

	for _, file := range files {
		dir := filepath.Dir(file.FullPath)

		for _, dep := range file.Config.Dependencies {
			path := filepath.Join(dir, dep.Source)
			deps, err := l.loadFromPath(path)

			if err != nil {
				return err
			}

			if len(deps) > 1 {
				return dependencySingleFileError
			}

			if len(deps) == 0 {
				continue
			}

			file.Dependencies[dep.Name] = deps[0]

			if err := l.loadDependencies(deps, depth+1); err != nil {
				return err
			}
		}
	}

	return nil
}

// getParsedFile checks if the file has already been parsed and returns previous parsed file
// or loads the file if not already loaded. Next time this is called with same source file
// it will return a reference to the previous loaded file.
func (l *Loader) getParsedFile(file string) (*ParsedFile, error) {
	if _, already := l.loaded[file]; already {
		return l.loaded[file], nil
	}

	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	parsed, err := NewParsedFile(file, content, l.options.TauDirectory, l.options.CacheDirectory)
	if err != nil {
		return nil, err
	}

	l.loaded[file] = parsed
	return parsed, nil
}
