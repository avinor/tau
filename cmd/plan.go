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
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type planCmd struct {
	meta

	loader *config.Loader
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
			if err := pc.processArgs(args); err != nil {
				return err
			}

			pc.init()

			return pc.run(args)
		},
	}

	pc.addMetaFlags(planCmd)

	return planCmd
}

func (pc *planCmd) init() {
	{
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    pc.TempDir,
			MaxDepth:         1,
		}

		pc.loader = config.NewLoader(options)
	}
}

func (pc *planCmd) run(args []string) error {
	loaded, err := pc.loader.Load(pc.file)
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
		moduleDir := paths.ModuleDir(pc.TempDir, source.Name)

		if _, err := os.Stat(moduleDir); os.IsNotExist(err) {
			return moduleNotInitError
		}
	}

	// Execute prepare hook to make sure we are logged in etc.
	ui.Header("Executing prepare hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "plan"); err != nil {
			return err
		}
	}

	if err := pc.resolveDependencies(loaded); err != nil {
		return err
	}

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(pc.TempDir, source.Name)

		ui.Separator()

		if !paths.IsFile(filepath.Join(moduleDir, "terraform.tfvars")) {
			ui.Warn("Cannot create a plan for %s", source.Name)
			continue
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(processors.NewUI(ui.Info)),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              source.Env,
		}

		extraArgs := getExtraArgs(pc.Engine.Compatibility.GetInvalidArgs("plan")...)
		extraArgs = append(extraArgs, "-out=tau.tfplan")
		if err := pc.Engine.Executor.Execute(options, "plan", extraArgs...); err != nil {
			return err
		}
	}

	ui.Separator()

	ui.Header("Executing finish hook...")
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "plan"); err != nil {
			return err
		}
	}

	return nil
}

func (m *meta) resolveDependencies(loaded []*config.Source) error {
	showDepFailureInfo := false
	ui.Header("Resolving dependencies...")
	for _, source := range loaded {
		if source.Config.Inputs == nil {
			continue
		}

		moduleDir := paths.ModuleDir(m.TempDir, source.Name)
		depsDir := paths.DependencyDir(m.TempDir, source.Name)

		vars, success, err := m.Engine.ResolveDependencies(source, depsDir)
		if err != nil {
			return err
		}

		if !success {
			showDepFailureInfo = true
			continue
		}

		if err := m.Engine.WriteInputVariables(source, moduleDir, vars); err != nil {
			return err
		}
	}

	if showDepFailureInfo {
		ui.NewLine()
		ui.Info(color.GreenString("Some of the dependencies failed to resolve. This can be because dependency"))
		ui.Info(color.GreenString("have not been applied yet, and therefore it cannot read remote-state."))
		ui.Info(color.GreenString("It will continue to plan those modules that can be applied and skip failed."))
		ui.NewLine()
	}

	return nil
}
