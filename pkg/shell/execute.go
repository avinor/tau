package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-cmd/cmd"
	"github.com/go-errors/errors"

	"github.com/avinor/tau/pkg/helpers/ui"
)

// Execute a shell command
func Execute(options *Options, command string, args ...string) error {

	if options == nil {
		options = &Options{}
	}

	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  false,
		Streaming: true,
	}

	execCmd := cmd.NewCmdOptions(cmdOptions, command, args...)
	if options.WorkingDirectory != "" {
		execCmd.Dir = options.WorkingDirectory
	}

	execCmd.Env = os.Environ()

	if len(options.Env) > 0 {
		for k, v := range options.Env {
			variable := fmt.Sprintf("%s=%s", k, v)
			execCmd.Env = append(execCmd.Env, variable)
		}
	}

	ui.Debug("environment variables: %#v", execCmd.Env)
	ui.Debug("command: %s %s", execCmd.Name, strings.Join(execCmd.Args, " "))

	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {

		defer close(doneChan)

		for execCmd.Stdout != nil || execCmd.Stderr != nil {
			select {
			case line, open := <-execCmd.Stdout:
				if !open {
					execCmd.Stdout = nil
					continue
				}
				processLine(options.Stdout, line)
			case line, open := <-execCmd.Stderr:
				if !open {
					execCmd.Stderr = nil
					continue
				}
				processLine(options.Stderr, line)
			}
		}
	}()

	status := <-execCmd.Start()

	<-doneChan

	if status.Error != nil {
		return status.Error
	}

	if status.Exit != 0 {
		return errors.Errorf("%s command exited with exit code %v", command, status.Exit)
	}

	return nil
}

func processLine(processors []OutputProcessor, line string) {
	for _, out := range processors {
		if !out.Write(line) {
			return
		}
	}
}
