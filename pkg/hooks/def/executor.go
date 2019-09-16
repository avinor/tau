package def

// Executor can execute a hook and return the output from hook execution.
type Executor interface {
	HasRun() bool
	Run(env map[string]string) error
	Output() string
}
