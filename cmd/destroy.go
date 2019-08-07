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

type destroyCmd struct {
	meta

	loader *loader.Loader

	autoApprove bool
}

var (
	// destroyLong is long description of destroy command
	destroyLong = templates.LongDesc(`Destroy resources managed by a module. It
		can either destroy a single resource or all of them. Requires that the 
		module have been initialized first.
		`)

	// destroyExample is examples for destroy command
	destroyExample = templates.Examples(`
		# Destroy all resources from local folder
		tau destroy

		# Destroy all resources in file
		tau destroy -f module.hcl
	`)
)

// newDestroyCmd creates a new destroy command
func newDestroyCmd() *cobra.Command {
	dc := &destroyCmd{}

	destroyCmd := &cobra.Command{
		Use:                   "destroy [-f SORUCE]",
		Short:                 "Destroy tau managed infrastructure",
		Long:                  destroyLong,
		Example:               destroyExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := dc.processArgs(args); err != nil {
				return err
			}

			dc.init()

			return dc.run(args)
		},
	}

	f := destroyCmd.Flags()
	f.BoolVar(&dc.autoApprove, "auto-approve", false, "auto approve destruction")

	dc.addMetaFlags(destroyCmd)

	return destroyCmd
}

func (dc *destroyCmd) init() {
	{
		options := &loader.Options{
			WorkingDirectory: paths.WorkingDir,
			TauDirectory:     dc.TauDir,
			MaxDepth:         1,
		}

		dc.loader = loader.New(options)
	}
}

func (dc *destroyCmd) run(args []string) error {
	files, err := dc.loader.Load(dc.file)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		ui.NewLine()
		ui.Warn("No sources found")
		return nil
	}

	// Want to destroy them in reverse order
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}

	// Verify all modules have been initialized
	if err := files.IsAllInitialized(); err != nil {
		return err
	}

	if err := hooks.RunAll(files, "prepare", "destroy"); err != nil {
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
		dc.resolveDependencies(files)
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

		extraArgs := getExtraArgs(dc.Engine.Compatibility.GetInvalidArgs("destroy")...)

		if dc.autoApprove {
			extraArgs = append(extraArgs, "-auto-approve")
		}

		if err := dc.Engine.Executor.Execute(options, "destroy", extraArgs...); err != nil {
			return err
		}

		paths.Remove(file.VariableFile())
	}

	ui.Separator("")

	if err := hooks.RunAll(files, "finish", "destroy"); err != nil {
		return err
	}

	return nil
}
