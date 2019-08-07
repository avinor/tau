package cmd

import (
	"net/http"
	"time"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initCmd with arguments and clients used to initialize a module
type initCmd struct {
	meta

	getter *getter.Client
	loader *loader.Loader

	maxDependencyDepth int
	purge              bool
	noOverrides        bool
	source             string
	sourceVersion      string
	reconfigure        bool
}

var (
	// sourceMustBeAFile is returned if source argument is defined and file reference a directory
	// Can only define source arguments if deploying a single module
	sourceMustBeAFile = errors.Errorf("file cannot be a directory when source is overridden")

	// sourceArgumentRequired required source argument when source-version is defined
	sourceArgumentRequired = errors.Errorf("source has to be set when source version is set")

	// purgeAndReconfigureTogether is returned if purge and reconfigure are both defined
	purgeAndReconfigureTogether = errors.Errorf("purge and reconfigure arguments cannot be defined together")

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
		tau init -f module.hcl

		# Initialize a module and send additional argument to terraform
		tau init -f module.hcl --args input=false
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
	f.BoolVar(&ic.noOverrides, "no-overrides", false, "do not create any overrides (backend config)")
	f.BoolVar(&ic.reconfigure, "reconfigure", false, "reconfigure the backend")
	f.StringVar(&ic.source, "source", "", "override module source location")
	f.StringVar(&ic.sourceVersion, "source-version", "", "override module source version, only valid together with source override")

	ic.addMetaFlags(initCmd)

	return initCmd
}

// init initializes the clients required for initCmd
func (ic *initCmd) init() {
	{
		timeout := time.Duration(ic.timeout) * time.Second

		ui.Debug("http timeout: %s", timeout)

		options := &getter.Options{
			HttpClient: &http.Client{
				Timeout: timeout,
			},
			WorkingDirectory: paths.WorkingDir,
		}

		ic.getter = getter.New(options)
	}

	{
		options := &loader.Options{
			WorkingDirectory: paths.WorkingDir,
			TauDirectory:     ic.TauDir,
			MaxDepth:         ic.maxDependencyDepth,
		}

		ui.Debug("max dependency depth: %s", ic.maxDependencyDepth)

		ic.loader = loader.New(options)
	}
}

// processArgs process arguments and checks for invalid options or combination of arguments
func (ic *initCmd) processArgs(args []string) error {

	// if source-version is defined then source is also required
	if ic.source == "" && ic.sourceVersion != "" {
		return sourceArgumentRequired
	}

	if ic.reconfigure && ic.purge {
		return purgeAndReconfigureTogether
	}

	// Always purge if source is overwritten
	if ic.source != "" {
		ic.purge = true
	}

	return nil
}

// run initialization command
func (ic *initCmd) run(args []string) error {
	if ic.purge {
		ui.Debug("Purging temporary folder")
		paths.Remove(ic.TauDir)
	}

	// load all sources
	files, err := ic.loader.Load(ic.file)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		ui.NewLine()
		ui.Warn("No sources found in path")
		return nil
	}

	// if source defined then it can only deploy a single file, not folder
	if len(files) > 1 {
		if ic.source != "" && paths.IsDir(ic.file) {
			return sourceMustBeAFile
		}
	}

	// Load module files usign go-getter
	if !ic.reconfigure {
		ui.Header("Loading modules...")
		for _, file := range files {
			module := file.Config.Module
			source := module.Source
			version := module.Version

			if ic.source != "" {
				source = ic.source
				version = &ic.sourceVersion
			}

			if err := ic.getter.Get(source, file.ModuleDir(), version); err != nil {
				return err
			}
		}
	}

	if err := hooks.RunAll(files, "prepare", "init"); err != nil {
		return err
	}

	for _, file := range files {
		ui.Separator(file.Name)

		if !ic.noOverrides {
			if err := ic.Engine.CreateOverrides(file); err != nil {
				return err
			}
		}

		ui.NewLine()

		options := &shell.Options{
			WorkingDirectory: file.ModuleDir(),
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              file.Env,
		}

		extraArgs := getExtraArgs(ic.Engine.Compatibility.GetInvalidArgs("init")...)

		if ic.reconfigure {
			extraArgs = append(extraArgs, "-reconfigure", "-force-copy")
		}

		if err := ic.Engine.Executor.Execute(options, "init", extraArgs...); err != nil {
			return err
		}
	}

	ui.Separator("")

	if err := hooks.RunAll(files, "finish", "init"); err != nil {
		return err
	}

	return nil
}
