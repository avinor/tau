package cmd

import (
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type ptCmd struct {
	meta
	name    string
	command passThroughCommand
}

// passThroughCommand description for a command that should just be passed through to terraform
type passThroughCommand struct {
	Use              string
	ShortDescription string
	LongDescription  string
	Example          string
	SingleResource   bool
	MaximumNArgs     int
}

var (
	passThroughCommands = map[string]passThroughCommand{
		"force-unlock": {
			Use:              "force-unlock -f SOURCE ID",
			ShortDescription: "Force unlock remote state lock",
			LongDescription:  "Force unlock remote state lock",
			SingleResource:   true,
			MaximumNArgs:     1,
		},
		"get": {
			Use:              "get [-f SOURCE]",
			ShortDescription: "Download and install modules for the configuration",
			LongDescription:  "Download and install modules for the configuration",
		},
		"import": {
			Use:              "import -f SOURCE",
			ShortDescription: "Import existing infrastructure into Terraform",
			LongDescription:  "Import existing infrastructure into Terraform",
			SingleResource:   true,
			MaximumNArgs:     2,
		},
		"refresh": {
			Use:              "refresh [-f SOURCE]",
			ShortDescription: "Update local state file against real resources",
			LongDescription:  "Update local state file against real resources",
		},
		"show": {
			Use:              "show [-f SOURCE]",
			ShortDescription: "Inspect Terraform state or plan",
			LongDescription:  "Inspect Terraform state or plan",
		},
		"state": {
			Use:              "state -f SOURCE [ARGS...]",
			ShortDescription: "Advanced state management",
			LongDescription:  "Advanced state management",
			SingleResource:   true,
			MaximumNArgs:     10,
		},
		"taint": {
			Use:              "taint -f SOURCE ADDRESS",
			ShortDescription: "Manually mark a resource for recreation",
			LongDescription:  "Manually mark a resource for recreation",
			SingleResource:   true,
			MaximumNArgs:     1,
		},
		"untaint": {
			Use:              "untaint -f SOURCE ADDRESS",
			ShortDescription: "Manually unmark a resource as tainted",
			LongDescription:  "Manually unmark a resource as tainted",
			SingleResource:   true,
			MaximumNArgs:     1,
		},
	}
)

func newPtCmd(name string, command passThroughCommand) *cobra.Command {
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
		Args:                  cobra.MaximumNArgs(command.MaximumNArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pt.meta.init(args); err != nil {
				return err
			}

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

func (pt *ptCmd) run(args []string) error {
	// load all sources
	files, err := pt.load()
	if err != nil {
		return err
	}

	// Verify all modules have been initialized
	if err := files.IsAllInitialized(); err != nil {
		return err
	}

	for _, file := range files {
		if err := pt.runFile(file, args); err != nil {
			return err
		}
	}

	ui.NewLine()

	return nil
}

func (pt *ptCmd) runFile(file *loader.ParsedFile, args []string) error {
	ui.Separator(file.Name)

	// Running prepare hook

	ui.Header("Executing prepare hooks...")

	if err := pt.Runner.Run(file, "prepare", pt.name); err != nil {
		return err
	}

	// Executing terraform command

	ui.NewLine()
	ui.Info(color.New(color.FgGreen, color.Bold).Sprint("Tau has been successfully initialized!"))
	ui.NewLine()

	options := &shell.Options{
		WorkingDirectory: file.ModuleDir(),
		Stdout:           shell.Processors(processors.NewUI(ui.Info)),
		Stderr:           shell.Processors(processors.NewUI(ui.Error)),
		Env:              file.Env,
	}

	ui.Separator(file.Name)

	extraArgs := getExtraArgs(pt.Engine.Compatibility.GetInvalidArgs(pt.name)...)
	extraArgs = append(extraArgs, args...)
	if err := pt.Engine.Executor.Execute(options, pt.name, extraArgs...); err != nil {
		return err
	}

	// Executing finish hook

	ui.Header("Executing finish hooks...")

	if err := pt.Runner.Run(file, "finish", pt.name); err != nil {
		return err
	}

	return nil
}
