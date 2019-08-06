package loader

import (
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/pkg/errors"
)

var (
	// moduleNotInitError is returned when a module is not initialized
	moduleNotInitError = errors.Errorf("module is not initialized")
)

// ParsedFileCollection is a collection of parsed files. Using this it is easier to perform
// actions on entire collection
type ParsedFileCollection []*ParsedFile

// IsAllInitialized checks if all modules have been initilized
func (c ParsedFileCollection) IsAllInitialized() error {
	for _, file := range c {
		if !paths.IsDir(file.ModuleDir()) {
			return moduleNotInitError
		}
	}

	return nil
}
