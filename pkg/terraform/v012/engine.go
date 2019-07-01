package v012

type Engine struct {
	Compatibility
	Generator
	Processor
	Executor
}

func NewEngine() *Engine {
	executor := Executor{}
	resolver := Resolver{}

	processor := Processor{
		executor: &executor,
	}

	generator := Generator{
		processor: &processor,
		resolver:  &resolver,
		executor:  &executor,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator:     generator,
		Processor:     processor,
		Executor:      executor,
	}
}
