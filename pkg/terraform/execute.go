package terraform

import "github.com/avinor/tau/pkg/shell"

// Execute wraps shell.Execute to execute terraform commands
func Execute(options *shell.Options, command string, args ...string) error {
	args = append([]string{command}, args...)

	return shell.Execute(options, "terraform", args...)
}
