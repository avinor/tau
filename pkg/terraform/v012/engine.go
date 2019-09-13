package v012

import (
	"github.com/avinor/tau/pkg/terraform/def"
)

// Engine is an engine for terraform 0.12 versions.
type Engine struct {
	Compatibility
	Generator
	Executor
}

// NewEngine creates a new engine and returns reference
func NewEngine(options *def.Options) *Engine {
	executor := Executor{}

	generator := Generator{
		executor: &executor,
		runner:   options.Runner,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator:     generator,
		Executor:      executor,
	}
}
