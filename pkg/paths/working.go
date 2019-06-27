package paths

import (
	"os"

	"github.com/apex/log"
)

var (
	// WorkingDir is directory where user is running command from
	WorkingDir string
)

// init working directory or panic. Need to be able to resolve working directory
// or it should just fail
func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("unable to get working directory")
	}

	WorkingDir = pwd
}
