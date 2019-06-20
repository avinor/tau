package dir

import (
	"os"
)

var (
	// Working is current working directory
	Working string
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("Unable to get working directory")
	}

	Working = pwd
}
