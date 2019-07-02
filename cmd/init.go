package cmd

import (
	"net/http"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/paths"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initCmd with arguments and clients used to initialize a module
type initCmd struct {
	meta

	getter *getter.Client
	loader *config.Loader

	maxDependencyDepth int
	purge              bool
	source             string
	sourceVersion      string
}

var (
	// sourceMustBeAFile is returned if source argument is defined and file reference a directory
	// Can only define source arguments if deploying a single module
	sourceMustBeAFile = errors.Errorf("file cannot be a directory when source is overridden")

	// sourceArgumentRequired required source argument when source-version is defined
	sourceArgumentRequired = errors.Errorf("source has to be set when source version is set")

	// initLong is long description of init command
	initLong = templates.LongDesc(`Initialize tau working folder based on file argument or by
		default the current folder. If file argument or no argument is defined it will deploy
		all resources in folder, file can also reference a single file to deploy. When
		initializing entire folder it will sort them in order of dependency.
		`)

	// initExample is examples for init command
	initExample = templates.Examples(`
		# Initialize current folder
		tau init

		# Initialize a single module
		tau init module.hcl

		# Initialize a module and send additional argument to terraform
		tau init module.hcl --args input=false
	`)
)

func newInitCmd() *cobra.Command {
	ic := &initCmd{}

	initCmd := &cobra.Command{
		Use:                   "init [-f SOURCE]",
		Short:                 "Initialize a tau working directory",
		Long:                  initLong,
		Example:               initExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ic.meta.processArgs(args); err != nil {
				return err
			}

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
	f.StringVar(&ic.source, "source", "", "override module source location")
	f.StringVar(&ic.sourceVersion, "source-version", "", "override module source version, only valid together with source override")

	ic.addMetaFlags(initCmd)

	return initCmd
}

// init initializes the clients required for initCmd
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

// processArgs process arguments and checks for invalid options or combination of arguments
func (ic *initCmd) processArgs(args []string) error {

	// if source defined then it can only deploy a single file, not folder
	if ic.source != "" {
		fi, err := os.Stat(ic.file)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			return sourceMustBeAFile
		}
	}

	// if source-version is defined then source is also required
	if ic.source == "" && ic.sourceVersion != "" {
		return sourceArgumentRequired
	}

	return nil
}

// run initialization command
func (ic *initCmd) run(args []string) error {
	if ic.purge {
		log.Debug("Purging temporary folder")
		log.Debug("")
		paths.Remove(ic.TempDir)
	}

	// load all sources
	loaded, err := ic.loader.Load(ic.file)
	if err != nil {
		return err
	}

	log.Info("")

	if len(loaded) == 0 {
		log.Warn("No sources found in path")
		return nil
	}

	// Load module files usign go-getter
	log.Info(color.New(color.Bold).Sprint("Loading modules..."))
	for _, source := range loaded {
		module := source.Config.Module
		moduleDir := paths.ModuleDir(ic.TempDir, source.Name)

		if module == nil {
			log.WithField("file", source.Name).Fatal("No module defined in source")
			continue
		}

		source := module.Source
		version := module.Version

		if ic.source != "" {
			source = ic.source
			version = &ic.sourceVersion
		}

		if err := ic.getter.Get(source, moduleDir, version); err != nil {
			return err
		}
	}
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Executing prepare hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "init"); err != nil {
			return err
		}
	}

	// log.Info(color.New(color.Bold).Sprint("Resolving dependencies..."))
	// for _, source := range loaded {
	// 	if source.Config.Inputs == nil {
	// 		continue
	// 	}

	// 	moduleDir := paths.ModuleDir(ic.TempDir, source.Name)
	// 	depsDir := paths.DependencyDir(ic.TempDir, source.Name)

	// 	vars, err := ic.Engine.ResolveDependencies(source, depsDir)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if err := ic.Engine.WriteInputVariables(source, moduleDir, vars); err != nil {
	// 		return err
	// 	}
	// }
	// log.Info("")

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

		log.Info("")
		log.Info("------------------------------------------------------------------------")
		log.Info("")

		extraArgs := getExtraArgs(ic.Engine.Compatibility.GetInvalidArgs("init")...)
		if err := ic.Engine.Executor.Execute(options, "init", extraArgs...); err != nil {
			return err
		}
	}

	log.Info("")
	log.Info("------------------------------------------------------------------------")
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Executing finish hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "init"); err != nil {
			return err
		}
	}
	log.Info("")

	return nil
}
