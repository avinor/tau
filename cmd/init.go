package cmd

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/avinor/tau/pkg/terraform"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type initCmd struct {
	meta

	maxDependencyDepth int
}

var (
	initCommand = validCommands["init"]
)

func newInitCmd() *cobra.Command {
	ic := &initCmd{}

	initCmd := &cobra.Command{
		Use:   initCommand.Use,
		Short: initCommand.ShortDescription,
		Long:  initCommand.LongDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ic.processArgs(args); err != nil {
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

	log.Info(color.New(color.Bold).Sprint("Loading modules..."))
	for _, source := range loaded {
		module := source.Config.Module
		moduleDir := dir.Module(ic.TempDir, source.File)

		if module == nil {
			log.WithField("file", source.File).Fatal("No module defined in source")
			continue
		}

		if err := ic.Getter.Get(module.Source, moduleDir, module.Version); err != nil {
			return err
		}
	}

	log.Info("")

	for _, source := range loaded {
		moduleDir := dir.Module(ic.TempDir, source.File)
		depsDir := dir.Dependency(ic.TempDir, source.File)

		if err := ic.Engine.CreateOverrides(source, moduleDir); err != nil {
			return err
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(new(processors.Log)),
			Stderr:           shell.Processors(new(processors.Log)),
		}

		extraArgs := getExtraArgs(args, ic.Engine.Compatibility.GetInvalidArgs("init")...)
		if err := terraform.Execute(options, "init", extraArgs...); err != nil {
			return err
		}

		if err := ic.Engine.CreateValues(source, depsDir, moduleDir); err != nil {
			return err
		}
	}

	return nil
}
