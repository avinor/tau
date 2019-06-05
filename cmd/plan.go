package cmd

import (
	"github.com/avinor/tau/pkg/api"
	"github.com/spf13/cobra"
)

const (
	planName             = "plan"
	planShortDescription = ""
	planLongDescription  = ""
)

type planCmd struct {
	noPrepare bool
}

func newPlanCmd() *cobra.Command {
	pc := planCmd{}

	planCmd := &cobra.Command{
		Use:   planName,
		Short: planShortDescription,
		Long:  planLongDescription,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pc.run(args)
		},
	}

	return planCmd
}

func (pc *planCmd) run(args []string) error {
	modules, err := api.Load(args[0])
	if err != nil {
		return err
	}

	for _, mod := range modules {
		executor := api.NewExecutor(mod)

		if err := executor.Prepare(mod); err != nil {
			return err
		}

		if err := executor.Run("plan"); err != nil {
			return err
		}
	}

	return nil

	// Check dependencies

	// exec.Prepare()

	// exec.Run("init")
	// exec.Run("plan")

	// Create pre-module (based on config)

	// Run pre-module

	// Read output from pre-module

	// Run terraform init -from-module

	// Create terraform.tfvars

	// Run terraform plan

	return nil
}
