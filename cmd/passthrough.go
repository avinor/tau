package cmd

import (
	"os"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ptCmd struct {
	meta
	name    string
	command Command

	loader *config.Loader
}

var (
	moduleNotInitError = errors.Errorf("module is not initialized")
)

func newPtCmd(name string, command Command) *cobra.Command {
	pt := &ptCmd{
		name:    name,
		command: command,
	}

	ptCmd := &cobra.Command{
		Use:                   command.Use,
		Short:                 command.ShortDescription,
		Long:                  command.LongDescription,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pt.processArgs(args); err != nil {
				return err
			}

			pt.init()

			return pt.run(args)
		},
	}

	if pt.command.Example != "" {
		ptCmd.Example = pt.command.Example
	}

	if pt.command.SingleResource {
		ptCmd.MarkFlagRequired("file")
	}

	pt.addMetaFlags(ptCmd)

	return ptCmd
}

func (pt *ptCmd) init() {
	{
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    pt.TempDir,
			MaxDepth:         1,
		}

		pt.loader = config.NewLoader(options)
	}
}

func (pt *ptCmd) run(args []string) error {
	loaded, err := pt.loader.Load(pt.file)
	if err != nil {
		return err
	}

	if len(loaded) == 0 {
		ui.NewLine()
		ui.Warn("No sources found")
		return nil
	}

	// Verify all modules have been initialized
	for _, source := range loaded {
		moduleDir := paths.ModuleDir(pt.TempDir, source.Name)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			return moduleNotInitError
		}
	}

	ui.Header("Executing prepare hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", pt.name); err != nil {
			return err
		}
	}

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(pt.TempDir, source.Name)

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              source.Env,
		}

		ui.Separator()

		extraArgs := getExtraArgs(pt.Engine.Compatibility.GetInvalidArgs(pt.name)...)
		if err := pt.Engine.Executor.Execute(options, pt.name, extraArgs...); err != nil {
			return err
		}
	}

	ui.Header("Executing finish hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", pt.name); err != nil {
			return err
		}
	}

	return nil
}
