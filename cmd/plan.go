package cmd

import (
	"os"

	"github.com/apex/log"
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/paths"
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
		# Initialize current folder
		tau init

		# Initialize a single module
		tau init module.hcl

		# Initialize a module and send additional argument to terraform
		tau init module.hcl --args input=false
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

	log.Info("")

	if len(loaded) == 0 {
		log.Warn("No sources found")
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
	log.Info(color.New(color.Bold).Sprint("Executing prepare hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "plan"); err != nil {
			return err
		}
	}
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Resolving dependencies..."))
	for _, source := range loaded {
		if source.Config.Inputs == nil {
			continue
		}

		moduleDir := paths.ModuleDir(pc.TempDir, source.Name)
		depsDir := paths.DependencyDir(pc.TempDir, source.Name)

		vars, success, err := pc.Engine.ResolveDependencies(source, depsDir)
		if err != nil {
			return err
		}

		if !success {
			log.Warnf("skipping module, could not resolve dependencies.")
		}

		if err := pc.Engine.WriteInputVariables(source, moduleDir, vars); err != nil {
			return err
		}
	}
	log.Info("")

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(pc.TempDir, source.Name)

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(new(processors.Log)),
			Stderr:           shell.Processors(new(processors.Log)),
			Env:              source.Env,
		}

		log.Info("------------------------------------------------------------------------")

		extraArgs := getExtraArgs(pc.Engine.Compatibility.GetInvalidArgs("plan")...)
		if err := pc.Engine.Executor.Execute(options, "plan", extraArgs...); err != nil {
			return err
		}
	}

	log.Info(color.New(color.Bold).Sprint("Executing finish hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "plan"); err != nil {
			return err
		}
	}
	log.Info("")

	return nil
}
