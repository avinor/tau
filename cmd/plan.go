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
	loader *api.Loader
}

func newPlanCmd() *cobra.Command {
	pc := planCmd{}

	planCmd := &cobra.Command{
		Use:   planName,
		Short: planShortDescription,
		Long:  planLongDescription,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pc.load(args); err != nil {
				return err
			}

			return pc.run(args)
		},
	}

	return planCmd
}

func (pc *planCmd) load(args []string) error {
	pc.loader = api.NewLoader(args[0])

	return pc.loader.Load()
}

func (pc *planCmd) run(args []string) error {
	// Load file -> return config

	// Check dependencies

	plan := pc.loader.GetExecutionPlan()

	plan.CreatePreModule()
	plan.ReadOutputValues()
	plan.RunInit()
	plan.RunPlan()

	// Create pre-module (based on config)

	// Run pre-module

	// Read output from pre-module

	// Run terraform init -from-module

	// Create terraform.tfvars

	// Run terraform plan

	return nil
}
