package command

import (
	"sync"

	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

// Executor can execute a command
type Executor struct {
	Command    string
	Arguments  []string
	WorkingDir string

	output string
	hasRun bool
	lock   sync.Mutex
}

// HasRun checks if command has already run
func (e *Executor) HasRun() bool {
	return e.hasRun
}

// Output returns output string once command has been executed. If command has
// not been executed it will always return empty string.
func (e *Executor) Output() string {
	return e.output
}

// Run the command and store result in output
func (e *Executor) Run(env map[string]string) error {
	e.lock.Lock()
	defer e.lock.Unlock()

	buffer := &processors.Buffer{}
	logp := processors.NewUI(ui.Error)

	options := &shell.Options{
		Stdout:           shell.Processors(buffer),
		Stderr:           shell.Processors(logp),
		WorkingDirectory: e.WorkingDir,
		Env:              env,
	}

	args := []string{}
	if e.Arguments != nil {
		args = append(args, e.Arguments...)
	}

	if err := shell.Execute(options, e.Command, args...); err != nil {
		e.hasRun = true
		return err
	}

	e.output = buffer.String()
	e.hasRun = true

	return nil
}
