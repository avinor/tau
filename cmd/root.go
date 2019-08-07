package cmd

import (
	"fmt"
	"strings"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/spf13/cobra"
)

var (
	rootLong = templates.LongDesc(`Tau is a thin wrapper on top of terraform to manage module deployments.
		It can deploy either a single module or all modules in a folder, taking into consideration the
		dependencies between modules.

		All arguments that are not handled by tau will be forwarded to terraform.
		`)
)

var (
	debug         bool
	workingDir    string
	terraformArgs []string
)

// NewRootCmd returns the root command for TAU.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tau",
		Short: "Tau (Terraform Avinor Utility) manages terraform deployments",
		Long:  rootLong,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if debug {
				ui.SetLevel(ui.DebugLevel)
			}
		},
	}

	p := rootCmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "enable verbose debug logs")
	p.StringVar(&workingDir, "working-directory", "", "working directory (default to current directory)")
	p.StringArrayVarP(&terraformArgs, "args", "a", []string{}, "arguments to forward to terraform in key=value format")

	rootCmd.AddCommand(newInitCmd())
	rootCmd.AddCommand(newPlanCmd())
	rootCmd.AddCommand(newApplyCmd())
	rootCmd.AddCommand(newDestroyCmd())
	rootCmd.AddCommand(newOutputCmd())
	rootCmd.AddCommand(newVersionCmd())

	for name, cmd := range passThroughCommands {
		rootCmd.AddCommand(newPtCmd(name, cmd))
	}

	return rootCmd
}

// getExtraArgs returns extra terraform arguments, but filters out invalid arguments
func getExtraArgs(invalidArgs ...string) []string {
	extraArgs := []string{}
	for _, arg := range terraformArgs {
		invalidArg := false
		arg = fmt.Sprintf("-%s", arg)

		for _, ia := range invalidArgs {
			if strings.HasPrefix(arg, ia) {
				invalidArg = true
			}
		}

		if !invalidArg {
			extraArgs = append(extraArgs, arg)
		}
	}

	return extraArgs
}
