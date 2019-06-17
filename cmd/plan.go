package cmd

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	planName             = "plan"
	planShortDescription = "Generate and show an execution plan"
	planLongDescription  = "Generate and show an execution plan"
)

type planCmd struct {
}

func newPlanCmd() *cobra.Command {
	pc := &planCmd{}

	planCmd := &cobra.Command{
		Use:   planName,
		Short: planShortDescription,
		Long:  planLongDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.run(args)
		},
	}

	return planCmd
}

func (pc *planCmd) run(args []string) error {
	source, err := utils.GetSourceArg(args)
	if err != nil {
		return err
	}

	loader, err := config.Load(source, &config.LoadOptions{})
	if err != nil {
		return err
	}

	for _, source := range loader.Sources {
		extraArgs := utils.GetExtraArgs(args, "-backend-config")
		options := &shell.Options{
			WorkingDirectory: source.ModuleDirectory(),
		}

		if err := shell.Execute("terraform", options, extraArgs...); err != nil {
			return err
		}
	}

	return nil
}
