package cmd

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/sources"
	"github.com/avinor/tau/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	initName             = "init"
	initShortDescription = "Initialize a Terraform working directory"
	initLongDescription  = "Initialize a Terraform working directory"
)

type initCmd struct {
	maxDependencyDepth int
}

func newInitCmd() *cobra.Command {
	ic := &initCmd{}

	initCmd := &cobra.Command{
		Use:   initName,
		Short: initShortDescription,
		Long:  initLongDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return ic.run(args)
		},
	}

	f := initCmd.Flags()
	f.IntVar(&ic.maxDependencyDepth, "max-dependency-depth", 1, "defines max dependency depth when traversing dependencies")

	return initCmd
}

func (ic *initCmd) init(args)

func (ic *initCmd) run(args []string) error {
	source, err := utils.GetSourceArg(args)
	if err != nil {
		return err
	}

	client := config.New(source, &config.Options{
		LoadSources:  true,
		CleanTempDir: true,
	})

	loaded, err := client.Load(source, nil)
	if err != nil {
		return err
	}

	for _, source := range loaded {
		modClient := sources.New(&sources.Options{})
		if err := modClient.Get(source.Config.Module.Source, source.ModuleDirectory(), source.Config.Module.Version); err != nil {
			return err
		}

		options := &shell.Options{
			WorkingDirectory: source.ModuleDirectory(),
		}

		// if err := source.CreateOverrides(); err != nil {
		// 	return err
		// }

		extraArgs := append([]string{"init"}, utils.GetExtraArgs(args, "-backend-config", "-from-module")...)
		if err := shell.Execute("terraform", options, extraArgs...); err != nil {
			return err
		}

		// if err := source.CreateInputVariables(); err != nil {
		// 	return err
		// }
	}

	return nil
}
