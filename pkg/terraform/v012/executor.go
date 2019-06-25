package v012

import "github.com/avinor/tau/pkg/shell"

type Executor struct {}

// Execute wraps shell.Execute to execute terraform commands
func (e *Executor) Execute(options *shell.Options, command string, args ...string) error {
	args = append([]string{command}, args...)

	return shell.Execute(options, "terraform", args...)
}