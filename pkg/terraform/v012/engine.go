package v012

import (
	"github.com/avinor/tau/pkg/terraform/lang"
)

type Engine struct {
	Compatibility
	Generator
	Processor
}

func NewEngine() *Engine {
	context := lang.EvalContext()

	processor := Processor{
		ctx: context,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator: Generator{
			ctx:       context,
			processor: &processor,
			resolver:  &Resolver{},
		},
		Processor: processor,
	}
}
