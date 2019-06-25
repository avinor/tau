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
	resolver := Resolver{}

	processor := Processor{
		ctx:      context,
		executor: &executor,
		resolver: &resolver,
	}

	generator := Generator{
		ctx:       context,
		processor: &processor,
		resolver:  &resolver,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator:     generator,
		Processor:     processor,
		Executor:      executor,
	}
}
