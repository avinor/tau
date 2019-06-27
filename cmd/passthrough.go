package cmd

import (
	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ptCmd struct {
	meta
	name    string
	command Command
}

func newPtCmd(name string, command Command) *cobra.Command {
	pt := &ptCmd{
		name:    name,
		command: command,
	}

	ptCmd := &cobra.Command{
		Use:                   command.Use,
		Short:                 command.ShortDescription,
		Long:                  command.LongDescription,
		Args:                  cobra.MinimumNArgs(1),
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pt.processArgs(args); err != nil {
				return err
			}

			return pt.run(args)
		},
	}

	if pt.command.Example != "" {
		ptCmd.Example = pt.command.Example
	}

	f := ptCmd.Flags()
	pt.addMetaFlags(f)

	return ptCmd
}

func (pt *ptCmd) run(args []string) error {
	log.Info(color.New(color.Bold).Sprint("Loading initialized sources..."))

	loaded, err := config.LoadSourcesFile(pt.TempDir)
	if err != nil {
		return err
	}

	log.Info("")

	for _, source := range loaded {
		moduleDir := dir.Module(pt.TempDir, source.Name)

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(new(processors.Log)),
			Stderr:           shell.Processors(new(processors.Log)),
			Env:              source.Env,
		}

		log.Info("------------------------------------------------------------------------")

		extraArgs := getExtraArgs(args, pt.Engine.Compatibility.GetInvalidArgs(pt.name)...)
		if err := pt.Engine.Executor.Execute(options, pt.name, extraArgs...); err != nil {
			return err
		}
	}

	return nil
}
