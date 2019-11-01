package cmd

import (
	"fmt"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type planCmd struct {
	meta
}

var (
	// planLong is long description of plan command
	planLong = templates.LongDesc(`Generate and show an execution plan where its possible.
		Command will resolve dependencies, create input variables and run terraform plan.
		For some dependencies it will not be possible if resources it depends on have not
		been deployed yet. It will not be able to show a plan, but apply will be able to
		apply the resources. 
		`)

	// planExample is examples for plan command
	planExample = templates.Examples(`
		# Plan current folder
		tau plan

		# Plan a single module
		tau plan -f module.hcl
	`)
)

// newPlanCmd creates a new plan command
func newPlanCmd() *cobra.Command {
	pc := &planCmd{}

	planCmd := &cobra.Command{
		Use:                   "plan [-f SORUCE]",
		Short:                 "Generate and show an execution plan",
		Long:                  planLong,
		Example:               planExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := pc.meta.init(args); err != nil {
				return err
			}

			return pc.run(args)
		},
	}

	pc.addMetaFlags(planCmd)

	return planCmd
}

func (pc *planCmd) run(args []string) error {
	// load all sources
	files, err := pc.load()
	if err != nil {
		return err
	}

	// Verify all modules have been initialized
	if pc.meta.noAutoInit {
		if err := files.IsAllInitialized(); err != nil {
			return err
		}
	}

	if err := files.Walk(pc.runFile); err != nil {
		return err
	}

	ui.NewLine()

	return nil
}

func (pc *planCmd) runFile(file *loader.ParsedFile) error {
	ui.Separator(file.Name)

	// Running prepare hook

	ui.Header("Executing prepare hooks...")

	if err := pc.Runner.Run(file, "prepare", "plan"); err != nil {
		return err
	}

	pc.autoInit(file)

	// Resolving dependencies

	success, err := pc.resolveDependencies(file)
	if err != nil {
		return err
	}

	if !success {
		return nil
	}

	// Executing terraform command

	ui.NewLine()
	ui.Info(color.New(color.FgGreen, color.Bold).Sprint("Tau has been successfully initialized!"))
	ui.NewLine()

	if !paths.IsFile(file.VariableFile()) {
		ui.Warn("Cannot create a plan for %s", file.Name)
		return nil
	}

	options := &shell.Options{
		WorkingDirectory: file.ModuleDir(),
		Stdout:           shell.Processors(processors.NewUI(ui.Info)),
		Stderr:           shell.Processors(processors.NewUI(ui.Error)),
		Env:              file.Env,
	}

	extraArgs := getExtraArgs(pc.Engine.Compatibility.GetInvalidArgs("plan")...)
	extraArgs = append(extraArgs, fmt.Sprintf("-out=%s", file.PlanFile()))
	if err := pc.Engine.Executor.Execute(options, "plan", extraArgs...); err != nil {
		return err
	}

	// Executing finish hook

	ui.Header("Executing finish hooks...")

	if err := pc.Runner.Run(file, "finish", "plan"); err != nil {
		return err
	}

	return nil
}
