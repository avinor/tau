package shell

import (
	"github.com/go-cmd/cmd"
	log "github.com/sirupsen/logrus"
)

type Options struct {
	WorkingDirectory string
}

func Execute(command string, options *Options, args ...string) error {
	execCmd := cmd.NewCmd(command, args...)
	if options.WorkingDirectory != "" {
		execCmd.Dir = options.WorkingDirectory
	}

	statusChan := execCmd.Start()

	finalStatus := <-statusChan

	log.Debugf("%s", finalStatus.Stdout)
	log.Debugf("%s", finalStatus.Stderr)

	return finalStatus.Error
}
