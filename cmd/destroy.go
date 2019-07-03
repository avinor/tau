package cmd

import (
	"os"
	"path/filepath"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/spf13/cobra"
)

type destroyCmd struct {
	meta

	loader *config.Loader

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
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    dc.TempDir,
			MaxDepth:         1,
		}

		dc.loader = config.NewLoader(options)
	}
}

func (dc *destroyCmd) run(args []string) error {
	loaded, err := dc.loader.Load(dc.file)
	if err != nil {
		return err
	}

	if len(loaded) == 0 {
		ui.NewLine()
		ui.Warn("No sources found")
		return nil
	}

	// Want to destroy them in reverse order
	for i, j := 0, len(loaded)-1; i < j; i, j = i+1, j-1 {
        loaded[i], loaded[j] = loaded[j], loaded[i]
    }

	// Verify all modules have been initialized
	for _, source := range loaded {
		moduleDir := paths.ModuleDir(dc.TempDir, source.Name)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			return moduleNotInitError
		}
	}

	// Execute prepare hook to make sure we are logged in etc.
	ui.Header("Executing prepare hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "destroy"); err != nil {
			return err
		}
	}

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(dc.TempDir, source.Name)

		ui.Separator()

		if !paths.IsFile(filepath.Join(moduleDir, "tau.tfplan")) {
			ui.Warn("No plan exists for %s", source.Name)
			continue
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              source.Env,
		}

		extraArgs := getExtraArgs(dc.Engine.Compatibility.GetInvalidArgs("destroy")...)

		if dc.autoApprove {
			extraArgs = append(extraArgs, "-auto-approve")
		}

		if err := dc.Engine.Executor.Execute(options, "destroy", extraArgs...); err != nil {
			return err
		}
	}

	ui.Separator()

	ui.Header("Executing finish hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "destroy"); err != nil {
			return err
		}
	}

	return nil
}
