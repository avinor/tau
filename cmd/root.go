package cmd

import (
	"strings"

	"github.com/avinor/tau/pkg/api"
	"github.com/go-errors/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootName             = "tau"
	rootShortDescription = "TAU (Terraform Avinor Utility) manages terraform deployments"
	rootLongDescription  = "TAU (Terraform Avinor Utility) manages terraform deployments"
)

type rootCmd struct {
	debug              bool
	noPrepare          bool
	maxDependencyDepth int
	workingDir         string

	api *api.Definition
}

// NewRootCmd returns the root command for TAU.
func NewRootCmd() *cobra.Command {
	rc := &rootCmd{}

	rootCmd := &cobra.Command{
		Use:   rootName,
		Short: rootShortDescription,
		Long:  rootLongDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if rc.debug {
				log.SetLevel(log.DebugLevel)
			}

			if err := rc.load(args); err != nil {
				return err
			}

			return rc.api.Run(args[0])
		},
	}

	f := rootCmd.Flags()
	f.BoolVar(&rc.debug, "debug", false, "enable verbose debug logs")
	f.BoolVar(&rc.noPrepare, "no-prepare", false, "do not prepare tau folder before running terraform")
	f.IntVar(&rc.maxDependencyDepth, "max-dependency-depth", 1, "defines max dependency depth when traversing dependencies")
	f.StringVar(&rc.workingDir, "working-directory", "", "working directory (default to current directory)")

	return rootCmd
}

func (rc *rootCmd) load(args []string) error {
	if len(args) < 1 {
		return errors.Errorf("Too few arguments defined")
	}

	if strings.HasPrefix(args[0], "-") {
		return errors.Errorf("First argument should not start with a -")
	}

	source := ""
	extraArgs := []string{}
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "-") {
			extraArgs = append(extraArgs, arg)
		} else if source != "" {
			return errors.Errorf("Only one source argument should be defined")
		} else {
			source = arg
		}
	}

	def, err := api.New(&api.Config{
		Source:             source,
		WorkingDirectory:   rc.workingDir,
		ExtraArguments:     extraArgs,
		MaxDependencyDepth: rc.maxDependencyDepth,
		CleanTempDir:       !rc.noPrepare,
	})
	if err != nil {
		return err
	}

	rc.api = def

	return nil
}
