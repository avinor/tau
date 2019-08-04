package module

import "github.com/avinor/tau/pkg/helpers/paths"

// Loader client for loading sources
type Loader struct {
	options *Options
}

// Options when loading modules. If TempDirectory is not set it will create a random
// temporary directory. This is not advised as it will create a new temporary directory
// for every call to load.
type Options struct {
	WorkingDirectory string

	// MaxDepth to search for dependencies. Should be enough with 1.
	MaxDepth int
}

// NewLoader creates a new loader client
func NewLoader(options *Options) *Loader {
	if options.WorkingDirectory == "" {
		options.WorkingDirectory = paths.WorkingDir
	}

	return &Loader{
		options: options,
	}
}
