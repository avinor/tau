package v012

import (
	"github.com/avinor/tau/pkg/terraform/lang"
)

type Engine struct {
	Compatibility
	Generator
	Processor
	Executor
}

func NewEngine() *Engine {
	context := lang.EvalContext()

	executor := Executor{}

	processor := Processor{
		ctx:      context,
		executor: &executor,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator: Generator{
			ctx:       context,
			processor: &processor,
			resolver:  &Resolver{},
		},
		Processor: processor,
		Executor:  executor,
	}
}
