package command

import (
	"sync"

	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

type Executor struct {
	Command    string
	Arguments  []string
	WorkingDir string

	output string
	hasRun bool
	lock   sync.Mutex
}

func (e *Executor) HasRun() bool {
	return e.hasRun
}

func (e *Executor) Output() string {
	return e.output
}

// Run the command hook and parse output
func (e *Executor) Run() error {
	e.lock.Lock()
	defer e.lock.Unlock()

	buffer := &processors.Buffer{}
	logp := processors.NewUI(ui.Error)

	options := &shell.Options{
		Stdout:           shell.Processors(buffer),
		Stderr:           shell.Processors(logp),
		WorkingDirectory: e.WorkingDir,
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
