package cmd

import (
	"github.com/spf13/cobra"
)

type ptCmd struct {
	meta
}

func newPtCmd(command Command) *cobra.Command {
	pt := &ptCmd{}

	ptCmd := &cobra.Command{
		Use:   command.Use,
		Short: command.ShortDescription,
		Long:  command.LongDescription,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pt.processArgs(args); err != nil {
				return err
			}

			return pt.run(args)
		},
	}

	return ptCmd
}

func (pt *ptCmd) run(args []string) error {
	// _, err := utils.GetSourceArg(args)
	// if err != nil {
	// 	return err
	// }

	// loader, err := config.Load(source, &config.LoadOptions{})
	// if err != nil {
	// 	return err
	// }

	// for _, source := range loader.Sources {
	// 	extraArgs := utils.GetExtraArgs(args, "-backend-config")
	// 	options := &shell.Options{
	// 		WorkingDirectory: source.ModuleDirectory(),
	// 	}

	// 	if err := shell.Execute("terraform", options, extraArgs...); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
