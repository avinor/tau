package cmd

import (
	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/shell"
	"github.com/spf13/cobra"
)

const (
	initName             = "init"
	initShortDescription = "Initialize a Terraform working directory"
	initLongDescription  = "Initialize a Terraform working directory"
)

type initCmd struct {
	meta

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
			if err := ic.initMeta(args); err != nil {
				return err
			}

			return ic.run(args)
		},
	}

	f := initCmd.Flags()
	f.IntVar(&ic.maxDependencyDepth, "max-dependency-depth", 2, "defines max dependency depth when traversing dependencies")

	ic.addMetaFlags(f)

	return initCmd
}

func (ic *initCmd) run(args []string) error {
	loaded, err := ic.Loader.Load(ic.SourceFile, nil)
	if err != nil {
		return err
	}

	for _, source := range loaded {
		module := source.Config.Module
		moduleDir := dir.Module(ic.TempDir, source.File)

		if err := ic.Getter.Get(module.Source, moduleDir, module.Version); err != nil {
			return err
		}

		// 	// if err := source.CreateOverrides(); err != nil {
		// 	// 	return err
		// 	// }

		options := &shell.Options{
			WorkingDirectory: moduleDir,
		}

		extraArgs := append([]string{"init"}, getExtraArgs(args, "-backend-config", "-from-module")...)
		if err := shell.Execute("terraform", options, extraArgs...); err != nil {
			return err
		}

		// 	// if err := source.CreateInputVariables(); err != nil {
		// 	// 	return err
		// 	// }
	}

	return nil
}
