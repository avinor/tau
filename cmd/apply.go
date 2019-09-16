package cmd

import (
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
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
	files, err := ac.Loader.Load(ac.file)
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

	if err := ac.Runner.RunAll(files, "prepare", "apply"); err != nil {
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

	if noPlansExists {
		//ac.resolveDependencies(files)
	} else {
		ui.Header("Found tau.plan files, only applying valid plans...")
	}

	for _, file := range files {
		planFileExists := paths.IsFile(file.PlanFile())

		ui.Separator(file.Name)

		if !planFileExists && !noPlansExists {
			ui.Warn("No plan exists for %s", file.Name)
			continue
		}

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
	}

	ui.Separator("")

	if err := ac.Runner.RunAll(files, "finish", "apply"); err != nil {
		return err
	}

	return nil
}
