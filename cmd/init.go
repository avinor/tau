package cmd

import (
	"net/http"
	"time"

	"github.com/apex/log"
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/paths"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type initCmd struct {
	meta

	getter *getter.Client
	loader *config.Loader

	maxDependencyDepth int
	purge              bool
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
		Use:                   "init SOURCE [terraform options]",
		Short:                 "Initialize a tau working directory",
		Long:                  initLong,
		Example:               initExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ic.processArgs(args); err != nil {
				return err
			}

			ic.init()

			return ic.run(args)
		},
	}

	f := initCmd.Flags()
	f.IntVar(&ic.maxDependencyDepth, "max-dependency-depth", 1, "defines max dependency depth when traversing dependencies") //nolint:lll
	f.BoolVar(&ic.purge, "purge", false, "purge temporary folder before init")

	ic.addMetaFlags(initCmd)

	return initCmd
}

func (ic *initCmd) init() {
	{
		timeout := time.Duration(ic.timeout) * time.Second

		log.Debugf("- Http timeout: %s", timeout)

		options := &getter.Options{
			HttpClient: &http.Client{
				Timeout: timeout,
			},
			WorkingDirectory: paths.WorkingDir,
		}

		ic.getter = getter.New(options)
	}

	{
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    ic.TempDir,
			MaxDepth:         ic.maxDependencyDepth,
		}

		ic.loader = config.NewLoader(options)
	}
}

func (ic *initCmd) run(args []string) error {
	if ic.purge {
		log.Debug("Purging temporary folder")
		log.Debug("")
		paths.Remove(ic.TempDir)
	}

	loaded, err := ic.loader.Load(ic.file)
	if err != nil {
		return err
	}

	log.Info(color.New(color.Bold).Sprint("Loading modules..."))
	for _, source := range loaded {
		module := source.Config.Module
		moduleDir := paths.ModuleDir(ic.TempDir, source.Name)

		if module == nil {
			log.WithField("file", source.Name).Fatal("No module defined in source")
			continue
		}

		if err := ic.getter.Get(module.Source, moduleDir, module.Version); err != nil {
			return err
		}
	}
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Resolving dependencies..."))
	for _, source := range loaded {
		if source.Config.Inputs == nil {
			continue
		}

		moduleDir := paths.ModuleDir(ic.TempDir, source.Name)
		depsDir := paths.DependencyDir(ic.TempDir, source.Name)

		vars, err := ic.Engine.ResolveDependencies(source, depsDir)
		if err != nil {
			return err
		}

		if err := ic.Engine.WriteInputVariables(source, moduleDir, vars); err != nil {
			return err
		}
	}
	log.Info("")

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(ic.TempDir, source.Name)

		if err := ic.Engine.CreateOverrides(source, moduleDir); err != nil {
			return err
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(new(processors.Log)),
			Stderr:           shell.Processors(new(processors.Log)),
			Env:              source.Env,
		}

		log.Info("------------------------------------------------------------------------")

		extraArgs := getExtraArgs(args, ic.Engine.Compatibility.GetInvalidArgs("init")...)
		if err := ic.Engine.Executor.Execute(options, "init", extraArgs...); err != nil {
			return err
		}
	}

	return config.SaveSourcesFile(loaded, ic.TempDir)
}
