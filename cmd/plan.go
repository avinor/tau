package cmd

import (
	"fmt"

	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
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
	files, err := pc.Loader.Load(pc.file)
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

	if err := pc.Runner.RunAll(files, "prepare", "plan"); err != nil {
		return err
	}

	if err := pc.resolveDependencies(files); err != nil {
		return err
	}

	for _, file := range files {
		ui.Separator(file.Name)

		if !paths.IsFile(file.VariableFile()) {
			ui.Warn("Cannot create a plan for %s", file.Name)
			continue
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
	}

	ui.Separator("")

	if err := pc.Runner.RunAll(files, "finish", "plan"); err != nil {
		return err
	}

	return nil
}
