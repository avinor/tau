package utils

import (
	"github.com/go-errors/errors"
	"strings"
)

// GetSourceArg finds argument that does not start with dash (-)
func GetSourceArg(args []string) (string, error) {
	source := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			if source != "" {
				return "", errors.Errorf("Only one source argument should be defined")
			}

			source = arg
		}
	}

	return source, nil
}

// GetExtraArgs returns all arguments starting with dash (-), but filters out invalid arguments
func GetExtraArgs(args []string, invalidArgs ...string) []string {
	extraArgs := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			invalidArg := false

			for _, ia := range invalidArgs {
				if strings.HasPrefix(arg, ia) {
					invalidArg = true
				}
			}

			if !invalidArg {
				extraArgs = append(extraArgs, arg)
			}
		}
	}

	return extraArgs
}
