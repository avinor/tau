package cmd

import (
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
			return pc.run()
		},
	}

	return planCmd
}

func (pc *planCmd) run() error {
	// catalog, err := api.LoadCatalog(&api.LoaderConfig{})
	// if err != nil {
	// 	return err
	// }

	// shell := api.NewExecutor()

	// for _, module := range catalog.Modules {
	// 	config := &api.ShellConfig{}

	// 	if err := shell.ExecuteTerraform(config, "plan", extraArgs); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
