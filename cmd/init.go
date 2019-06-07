package cmd

import (
	"github.com/avinor/tau/pkg/api"
	"github.com/avinor/tau/pkg/executor"
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

	source, err := getSourceArg(args)
	if err != nil {
		return err
	}

	catalog, err := api.NewCatalog(source, &api.Config{
		LoadSources: true,
	})
	if err != nil {
		return err
	}

	for _, module := range catalog.Modules {
		//extraArgs := module.GetBackendArgs()

		shell, err := executor.NewShell(&executor.Config{})
		if err != nil {
			return err
		}

		extraArgs := getExtraArgs(args, "-backend-config")
		extraArgs = append(extraArgs, module.GetBackendArgs()...)

		if err := shell.ExecuteTerraform("init", extraArgs...); err != nil {
			return err
		}
	}

	return catalog.Save()
}
