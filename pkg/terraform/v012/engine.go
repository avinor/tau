package v012

// Engine is an engine for terraform 0.12 versions.
type Engine struct {
	Compatibility
	Generator
	Executor
}

// NewEngine creates a new engine and returns reference
func NewEngine() *Engine {
	executor := Executor{}

	generator := Generator{
		executor: &executor,
	}

	return &Engine{
		Compatibility: Compatibility{},
		Generator:     generator,
		Executor:      executor,
	}
}
