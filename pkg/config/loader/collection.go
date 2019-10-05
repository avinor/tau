package loader

import (
	"sync"

	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/tfdiags"
	"github.com/pkg/errors"
)

var (
	// moduleNotInitError is returned when a module is not initialized
	moduleNotInitError = errors.Errorf("module is not initialized")

	lock = sync.Mutex{}
)

// ParsedFileCollection is a collection of parsed files. Using this it is easier to perform
// actions on entire collection
type ParsedFileCollection []*ParsedFile

// WalkFunc is called when walking the collection
type WalkFunc func(file *ParsedFile) error

// IsAllInitialized checks if all modules have been initilized
func (c ParsedFileCollection) IsAllInitialized() error {
	for _, file := range c {
		if !paths.IsDir(file.ModuleDir()) {
			return moduleNotInitError
		}
	}

	return nil
}

// Walk travers the files in collection and execute them in correct
// order depending on dependencies. It could do it in parallell but has
// been limited to do one at the time to not mess up output now
func (c ParsedFileCollection) Walk(walkerFunc WalkFunc) error {
	graph := &dag.AcyclicGraph{}

	for _, file := range c {
		graph.Add(file)
	}

	for _, file := range c {
		for _, dep := range file.Dependencies {
			if contains(c, dep) {
				graph.Connect(dag.BasicEdge(file, dep))
			}
		}
	}

	return graph.Walk(func(vertex dag.Vertex) tfdiags.Diagnostics {
		var diags tfdiags.Diagnostics

		lock.Lock()
		defer lock.Unlock()
		if err := walkerFunc(vertex.(*ParsedFile)); err != nil {
			return diags.Append(err)
		}

		return diags
	}).Err()
}

func contains(list []*ParsedFile, item *ParsedFile) bool {
	for _, file := range list {
		if file == item {
			return true
		}
	}

	return false
}
