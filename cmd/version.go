package cmd

import (
	"github.com/spf13/cobra"

	"github.com/avinor/tau/pkg/helpers/ui"
)

var (
	// BuildTag set during build to git tag, if any
	BuildTag = "unset"

	// BuildSha is the git sha set during build
	BuildSha = "unset"

	// GitTreeState is state of git tree during build, dirty if there are uncommitted changes
	GitTreeState = "unset"
)

type versionCmd struct {
	detailed bool
	// noCheck bool
}

func newVersionCmd() *cobra.Command {
	ver := &versionCmd{}

	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Shows version of tau",
		Long:  "Shows version of tau",
		Run: func(cmd *cobra.Command, args []string) {
			ver.printVersion()
		},
	}

	f := verCmd.Flags()
	f.BoolVar(&ver.detailed, "detailed", false, "show detailed information")
	//f.BoolVar(&ver.noCheck, "no-check", false, "do not check for upgrades")

	return verCmd
}

func (ver *versionCmd) printVersion() {
	ui.Info("tau version %s", BuildTag)
	ui.Info("commit: %s", BuildSha)
}
