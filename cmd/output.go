package cmd

import (
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/spf13/cobra"
)

type outputCmd struct {
	meta

	loader *loader.Loader

	json bool
}

var (
	// outputLong is long description of output command
	outputLong = templates.LongDesc(`Print all the output variables from a module.
		If including the --json flag it will print output in json format.
		`)

	// outputExample is examples for output command
	outputExample = templates.Examples(`
		# Combine output from all files
		tau output

		# Print output from module.hcl in json format
		tau output -f module.hcl --json
	`)
)

// newOutputCmd creates a new output command
func newOutputCmd() *cobra.Command {
	oc := &outputCmd{}

	outputCmd := &cobra.Command{
		Use:                   "output [-f SORUCE]",
		Short:                 "Print output from module",
		Long:                  outputLong,
		Example:               outputExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := oc.processArgs(args); err != nil {
				return err
			}

			oc.init()

			return oc.run(args)
		},
	}

	f := outputCmd.Flags()
	f.BoolVar(&oc.json, "json", false, "print in json format")

	oc.addMetaFlags(outputCmd)

	return outputCmd
}

func (oc *outputCmd) init() {
	{
		options := &loader.Options{
			WorkingDirectory: paths.WorkingDir,
			TauDirectory:     oc.TauDir,
			MaxDepth:         1,
		}

		oc.loader = loader.New(options)
	}
}

func (oc *outputCmd) run(args []string) error {
	files, err := oc.loader.Load(oc.file)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		ui.NewLine()
		ui.Warn("No sources found")
		return nil
	}

	// Verify all modules have been initialized
	if err := files.IsAllInitialized(); err != nil {
		return err
	}

	if err := hooks.RunAll(files, "prepare", "output"); err != nil {
		return err
	}

	// Check if any plans exist, if not then run plan first
	noVariablesExists := true
	for _, file := range files {
		if paths.IsFile(file.VariableFile()) {
			noVariablesExists = false
			continue
		}
	}

	if noVariablesExists {
		oc.resolveDependencies(files)
	}

	for _, file := range files {
		ui.Separator(file.Name)

		if !paths.IsFile(file.VariableFile()) {
			ui.Warn("No values file exists for %s", file.Name)
			continue
		}

		options := &shell.Options{
			WorkingDirectory: file.ModuleDir(),
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              file.Env,
		}

		extraArgs := getExtraArgs(oc.Engine.Compatibility.GetInvalidArgs("output")...)

		if err := oc.Engine.Executor.Execute(options, "output", extraArgs...); err != nil {
			return err
		}

		paths.Remove(file.VariableFile())
	}

	ui.Separator("")

	if err := hooks.RunAll(files, "finish", "output"); err != nil {
		return err
	}

	return nil
}
