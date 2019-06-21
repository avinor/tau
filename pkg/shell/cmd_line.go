package shell

import (
	"time"

	"github.com/go-cmd/cmd"
)

// Execute a shell command
func Execute(command string, options *Options, args ...string) error {
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

	// Print STDOUT and STDERR lines streaming from Cmd
	go func() {
		for {
			select {
			case line := <-execCmd.Stdout:
				options.ProcessStdout(line)
			case line := <-execCmd.Stderr:
				options.ProcessStderr(line)
			}
		}
	}()

	status := <-execCmd.Start()

	for len(execCmd.Stdout) > 0 || len(execCmd.Stderr) > 0 {
		time.Sleep(10 * time.Millisecond)
	}

	return status.Error
}
