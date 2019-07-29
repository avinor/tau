package cmd

import (
	"os"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type ptCmd struct {
	meta
	name    string
	command passThroughCommand

	loader *config.Loader
}

// passThroughCommand description for a command that should just be passed through to terraform
type passThroughCommand struct {
	Use              string
	ShortDescription string
	LongDescription  string
	Example          string
	SingleResource   bool
}

var (
	moduleNotInitError = errors.Errorf("module is not initialized")

	passThroughCommands = map[string]passThroughCommand{
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
		},
		"output": {
			Use:              "output [-f SOURCE]",
			ShortDescription: "Read an output from a state file",
			LongDescription: templates.LongDesc(
				`Read an output from the state file for module deployed.
				`),
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
		"taint": {
			Use:              "taint -f SOURCE",
			ShortDescription: "Manually mark a resource for recreation",
			LongDescription:  "Manually mark a resource for recreation",
			SingleResource:   true,
		},
		"untaint": {
			Use:              "untaint -f SOURCE",
			ShortDescription: "Manually unmark a resource as tainted",
			LongDescription:  "Manually unmark a resource as tainted",
			SingleResource:   true,
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

	ui.Separator()

	ui.Header("Executing finish hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", pt.name); err != nil {
			return err
		}
	}

	return nil
}
