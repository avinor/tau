package def

// Executor can execute a hook and return the output from hook execution.
type Executor interface {
	HasRun() bool
	Run() error
	Output() string
}
