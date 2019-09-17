package cmd

import (
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type applyCmd struct {
	meta

	autoApprove bool
	deletePlan  bool
}

var (
	// applyLong is long description of apply command
	applyLong = templates.LongDesc(`Apply an execution plan where its possible. It will
		loop through all plans generated from plan command and execute them. It will only
		execute for those modules that successfully generated a plan.
		`)

	// applyExample is examples for apply command
	applyExample = templates.Examples(`
		# Apply on current folder
		tau apply

		# Apply a single module
		tau apply -f module.hcl

		# Apply a single module and auto approve
		tau apply -f module.hcl --no-input
	`)
)

// newApplyCmd creates a new apply command
func newApplyCmd() *cobra.Command {
	ac := &applyCmd{}

	applyCmd := &cobra.Command{
		Use:                   "apply [-f SORUCE]",
		Short:                 "Builds or changes infrastructure",
		Long:                  applyLong,
		Example:               applyExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ac.meta.init(args); err != nil {
				return err
			}

			return ac.run(args)
		},
	}

	f := applyCmd.Flags()
	f.BoolVar(&ac.autoApprove, "auto-approve", false, "auto approve deployment")
	f.BoolVar(&ac.deletePlan, "delete-plan", true, "delete terraform plan on success")

	ac.addMetaFlags(applyCmd)

	return applyCmd
}

func (ac *applyCmd) run(args []string) error {
	// load all sources
	files, err := ac.load()
	if err != nil {
		return err
	}

	// Verify all modules have been initialized
	if err := files.IsAllInitialized(); err != nil {
		return err
	}

	// Check if any plans exist, if not then run plan first
	noPlansExists := true
	for _, file := range files {
		if paths.IsFile(file.PlanFile()) {
			noPlansExists = false
			continue
		}
	}

	if !noPlansExists {
		ui.Header("Found tau.plan files, only applying valid plans...")
	}

	for _, file := range files {
		if err := ac.runFile(file, !noPlansExists); err != nil {
			return err
		}
	}

	ui.NewLine()

	return nil
}

func (ac *applyCmd) runFile(file *loader.ParsedFile, onlyPlans bool) error {
	ui.Separator(file.Name)

	// Running prepare hook

	ui.Header("Executing prepare hooks...")

	if err := ac.Runner.Run(file, "prepare", "apply"); err != nil {
		return err
	}

	// Resolving dependencies

	if !paths.IsFile(file.VariableFile()) {
		success, err := ac.resolveDependencies(file)
		if err != nil {
			return err
		}

		if !success {
			return nil
		}
	}

	planFileExists := paths.IsFile(file.PlanFile())

	if !planFileExists && onlyPlans {
		ui.Warn("No plan exists")
		return nil
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

	extraArgs := getExtraArgs(ac.Engine.Compatibility.GetInvalidArgs("apply")...)
	extraArgs = append(extraArgs, "-input=false")

	if ac.autoApprove {
		extraArgs = append(extraArgs, "-auto-approve")
	}

	if planFileExists {
		extraArgs = append(extraArgs, file.PlanFile())
	}

	if err := ac.Engine.Executor.Execute(options, "apply", extraArgs...); err != nil {
		return err
	}

	if ac.deletePlan {
		paths.Remove(file.PlanFile())
	}

	// Executing finish hook

	ui.Header("Executing finish hooks...")

	if err := ac.Runner.Run(file, "finish", "apply"); err != nil {
		return err
	}

	return nil
}
