package paths

import (
	"os"

	"github.com/avinor/tau/pkg/helpers/ui"
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
		ui.Fatal("unable to get working directory")
	}

	WorkingDir = pwd
}
