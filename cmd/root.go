package cmd

import (
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootName             = "tau"
	rootShortDescription = "TAU (Terraform Avinor Utility) manages terraform deployments"
	rootLongDescription  = "TAU (Terraform Avinor Utility) manages terraform deployments"
)

var (
	debug      bool
	workingDir string
)

// NewRootCmd returns the root command for TAU.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   rootName,
		Short: rootShortDescription,
		Long:  rootLongDescription,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				log.SetLevel(log.DebugLevel)
			}
		},
	}

	p := rootCmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "enable verbose debug logs")
	p.StringVar(&workingDir, "working-directory", "", "working directory (default to current directory)")

	rootCmd.AddCommand(newInitCmd())

	return rootCmd
}

func getSourceArg(args []string) (string, error) {
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

func getExtraArgs(args []string, invalidArgs ...string) []string {
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
