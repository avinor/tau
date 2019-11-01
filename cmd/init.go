package cmd

import (
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// initCmd with arguments and clients used to initialize a module
type initCmd struct {
	meta

	options *initOptions
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
	ic := &initCmd{
		options: &initOptions{
			source: &config.Module{},
		},
	}

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
			if err := ic.meta.init(args); err != nil {
				return err
			}

			if err := ic.processArgs(args); err != nil {
				return err
			}

			return ic.run(args)
		},
	}

	f := initCmd.Flags()
	f.BoolVar(&ic.options.purge, "purge", true, "purge temporary folder before init")
	f.BoolVar(&ic.options.noOverrides, "no-overrides", false, "do not create any overrides (backend config)")
	f.BoolVar(&ic.options.reconfigure, "reconfigure", false, "reconfigure the backend")
	f.StringVar(&ic.options.source.Source, "source", "", "override module source location")
	f.StringVar(&ic.options.source.Version, "source-version", "", "override module source version, only valid together with source override")

	ic.addMetaFlags(initCmd)

	return initCmd
}

// processArgs process arguments and checks for invalid options or combination of arguments
func (ic *initCmd) processArgs(args []string) error {

	// if source-version is defined then source is also required
	if ic.options.source.Source == "" && ic.options.source.Version != "" {
		return sourceArgumentRequired
	}

	if ic.options.reconfigure && ic.options.purge {
		return purgeAndReconfigureTogether
	}

	// Always purge if source is overwritten
	if ic.options.source.Source != "" {
		ic.options.purge = true
	}

	return nil
}

// run initialization command
func (ic *initCmd) run(args []string) error {
	if ic.options.purge {
		ui.Debug("Purging temporary folder")
		paths.Remove(ic.TauDir)
	}

	// load all sources
	files, err := ic.load()
	if err != nil {
		return err
	}

	// if source defined then it can only deploy a single file, not folder
	if len(files) > 1 && ic.options.source.Source != "" {
		return sourceMustBeAFile
	}

	if err := files.Walk(ic.runFile); err != nil {
		return err
	}

	ui.NewLine()

	return nil
}

func (ic *initCmd) runFile(file *loader.ParsedFile) error {
	ui.Separator(file.Name)

	// Running prepare hook

	ui.Header("Executing prepare hooks...")

	if err := ic.Runner.Run(file, "prepare", "init"); err != nil {
		return err
	}

	// Executing terraform command

	ic.runInit(file, ic.options)

	// Executing finish hook

	ui.Header("Executing finish hooks...")

	if err := ic.Runner.Run(file, "finish", "init"); err != nil {
		return err
	}

	return nil
}
