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
		resolver: &resolver,
	}

	generator := Generator{
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
