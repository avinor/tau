package cmd

import (
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/shell"
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

func (ic *initCmd) run(args []string) error {
	// TODO Option to clean out
	// if config.CleanTempDir {
	// 	log.Debugf("Cleaning temp directory...")
	// 	os.RemoveAll(config.WorkingDirectory)
	// }

	source, err := utils.GetSourceArg(args)
	if err != nil {
		return err
	}

	loader, err := config.Load(source, &config.LoadOptions{
		LoadSources: true,
	})
	if err != nil {
		return err
	}

	for _, source := range loader.Sources {
		// if err := source.CreateOverrides(); err != nil {
		// 	return err
		// }

		extraArgs := append([]string{"init"}, utils.GetExtraArgs(args, "-backend-config", "-from-module")...)

		initArgs := append(extraArgs, "-from-module", source.Config.Module.Source, "-backend=false")
		options := &shell.Options{
			WorkingDirectory: source.ModuleDirectory(),
		}

		if err := shell.Execute("terraform", options, initArgs...); err != nil {
			return err
		}

		if err := source.CreateInputVariables(); err != nil {
			return err
		}
	}

	return loader.Save()
}
