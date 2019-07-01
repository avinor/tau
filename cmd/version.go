package cmd

import (
	"fmt"

	"github.com/avinor/tau/pkg/terraform"
	"github.com/spf13/cobra"
)

var (
	BuildVersion = "0.1.0"

	BuildTag = "unset"

	BuildSha = "unset"

	GitTreeState = "unset"
)

type versionCmd struct {
	client bool
	//noCheck bool
}

func newVersionCmd() *cobra.Command {
	ver := &versionCmd{}

	verCmd := &cobra.Command{
		Use:   "version",
		Short: "Shows version of tau and terraform",
		Long:  "Shows version of tau and terraform",
		Run: func(cmd *cobra.Command, args []string) {
			ver.printVersion()
		},
	}

	f := verCmd.Flags()
	f.BoolVar(&ver.client, "client", false, "check only version for tau client")
	//f.BoolVar(&ver.noCheck, "no-check", false, "do not check for upgrades")

	return verCmd
}

func (ver *versionCmd) printVersion() {
	fmt.Printf("tau v%s\n", BuildVersion)

	if !ver.client {
		version := terraform.Version()

		if version == "" {
			fmt.Println("terraform not found!")
		} else {
			fmt.Printf("terraform v%s\n", version)
		}
	}
}
