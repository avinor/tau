package cmd

import (
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/avinor/tau/internal/templates"
)

var (
	rootLong = templates.LongDesc(`TAU is a thin wrapper on top of terraform to manage module deployments.
		It can deploy either a single module or all modules in a folder, taking into consideration the
		dependencies between modules.

		All arguments that are not handled by tau will be forwarded to terraform.
		`)
)

var (
	debug      bool
	workingDir string
)

// NewRootCmd returns the root command for TAU.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tau",
		Short: "TAU (Terraform Avinor Utility) manages terraform deployments",
		Long:  rootLong,
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

	for name, cmd := range validCommands {
		rootCmd.AddCommand(newPtCmd(name, cmd))
	}

	return rootCmd
}
