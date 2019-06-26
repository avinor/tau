package cmd

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/avinor/tau/internal/templates"
)

type initCmd struct {
	meta

	maxDependencyDepth int
}

var (
	initLong = templates.LongDesc(`Initialize tau working folder based on SOURCE argument.
		SOURCE can either be a single file or a folder. If it is a folder it will initialize
		all modules in the folder, ordering them by dependencies.
		`)

	initExample = templates.Examples(`
		# Initialize a single module
		tau init module.hcl

		# Initialize a folder
		tau init folder

		# Initialize a module and send additional argument to terraform
		tau init module.hcl -input=false
	`)
)

func newInitCmd() *cobra.Command {
	ic := &initCmd{}

	initCmd := &cobra.Command{
		Use:   "init SOURCE [terraform options]",
		Short: "Initialize a Tau working directory",
		Long:  initLong,
		Example: initExample,
		Args:  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage: true,
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
		moduleDir := dir.Module(ic.TempDir, source.Name)

		if module == nil {
			log.WithField("file", source.Name).Fatal("No module defined in source")
			continue
		}

		if err := ic.Getter.Get(module.Source, moduleDir, module.Version); err != nil {
			return err
		}
	}
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Resolving dependencies..."))
	for _, source := range loaded {
		if source.Config.Inputs == nil {
			continue
		}

		moduleDir := dir.Module(ic.TempDir, source.Name)
		depsDir := dir.Dependency(ic.TempDir, source.Name)

		vars, err := ic.Engine.ResolveDependencies(source, depsDir)
		if err != nil {
			return err
		}

		if len(vars) == 0 {
			continue
		}

		if err := ic.Engine.WriteInputVariables(source, moduleDir, vars); err != nil {
			return err
		}
	}
	log.Info("")

	for _, source := range loaded {
		moduleDir := dir.Module(ic.TempDir, source.Name)

		if err := ic.Engine.CreateOverrides(source, moduleDir); err != nil {
			return err
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(new(processors.Log)),
			Stderr:           shell.Processors(new(processors.Log)),
		}

		extraArgs := getExtraArgs(args, ic.Engine.Compatibility.GetInvalidArgs("init")...)
		if err := ic.Engine.Executor.Execute(options, "init", extraArgs...); err != nil {
			return err
		}
	}

	return config.SaveSourcesFile(loaded, ic.TempDir)
}
