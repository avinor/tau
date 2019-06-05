package api

type ModuleExecutor struct {
	module *Module
}

type Executor struct {
	modules []*ModuleExecutor
}

func NewExecutor(module []*Module) (*Executor, error) {
	return &Executor{}, nil
}
