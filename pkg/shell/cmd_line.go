package shell

import (
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/go-cmd/cmd"
	"github.com/go-errors/errors"
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

	log.Debugf("environment variables: %#v", execCmd.Env)

	// Print STDOUT and STDERR lines streaming from Cmd
	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				processLine(options.Stdout, line)
			case line := <-execCmd.Stderr:
				processLine(options.Stderr, line)
			}
		}
	}()

	status := <-execCmd.Start()

	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

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
