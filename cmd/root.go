package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	rootName             = "tau"
	rootShortDescription = "TAU (Terraform Avinor Utility) manages terraform deployments"
	rootLongDescription  = "TAU (Terraform Avinor Utility) manages terraform deployments"
)

var (
	debug bool
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

	return rootCmd
}