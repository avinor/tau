package v012

import (
	"github.com/avinor/tau/pkg/terraform/lang"
)

type Engine struct {
	Backend
	Compatibility
	Generator
	Processor
}

func NewEngine() *Engine {
	context := lang.EvalContext()

	backend := Backend{
		ctx: context,
	}

	return &Engine{
		Backend:       backend,
		Compatibility: Compatibility{},
		Generator: Generator{
			backend: &backend,
		},
		Processor: Processor{},
	}
}
