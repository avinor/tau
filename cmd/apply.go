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

type applyCmd struct {
	meta

	loader *config.Loader

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
			if err := ac.processArgs(args); err != nil {
				return err
			}

			ac.init()

			return ac.run(args)
		},
	}

	f := applyCmd.Flags()
	f.BoolVar(&ac.autoApprove, "auto-approve", false, "auto approve deployment")
	f.BoolVar(&ac.deletePlan, "delete-plan", true, "delete terraform plan on success")

	ac.addMetaFlags(applyCmd)

	return applyCmd
}

func (ac *applyCmd) init() {
	{
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    ac.TempDir,
			MaxDepth:         1,
		}

		ac.loader = config.NewLoader(options)
	}
}

func (ac *applyCmd) run(args []string) error {
	loaded, err := ac.loader.Load(ac.file)
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
		moduleDir := paths.ModuleDir(ac.TempDir, source.Name)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			return moduleNotInitError
		}
	}

	// Execute prepare hook to make sure we are logged in etc.
	ui.Header("Executing prepare hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "apply"); err != nil {
			return err
		}
	}

	// Check if any plans exist, if not then run plan first
	noPlansExists := true
	for _, source := range loaded {
		moduleDir := paths.ModuleDir(ac.TempDir, source.Name)
		planFile := filepath.Join(moduleDir, "tau.tfplan")

		if paths.IsFile(planFile) {
			noPlansExists = false
			continue
		}
	}

	if noPlansExists {
		ac.resolveDependencies(loaded)
	} else {
		ui.Header("Found tau.plan files, only applying valid plans...")
	}

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(ac.TempDir, source.Name)
		planFile := filepath.Join(moduleDir, "tau.tfplan")
		planFileExists := paths.IsFile(planFile)

		ui.Separator()

		if !planFileExists && !noPlansExists {
			ui.Warn("No plan exists for %s", source.Name)
			continue
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              source.Env,
		}

		extraArgs := getExtraArgs(ac.Engine.Compatibility.GetInvalidArgs("apply")...)
		extraArgs = append(extraArgs, "-input=false")

		if planFileExists {
			extraArgs = append(extraArgs, "tau.tfplan")
		}

		if ac.autoApprove {
			extraArgs = append(extraArgs, "-auto-approve")
		}

		if err := ac.Engine.Executor.Execute(options, "apply", extraArgs...); err != nil {
			return err
		}

		if ac.deletePlan {
			paths.Remove(planFile)
		}
	}

	ui.Separator()

	ui.Header("Executing finish hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "apply"); err != nil {
			return err
		}
	}

	return nil
}
