package cmd

import (
	"os"
	"path/filepath"

	"github.com/apex/log"
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type applyCmd struct {
	meta

	loader *config.Loader

	noInput bool
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

// newApplyCmd creates a new plan command
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
	f.BoolVar(&ac.noInput, "no-input", false, "do not ask for any input")

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

	log.Info("")

	if len(loaded) == 0 {
		log.Warn("No sources found")
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
	log.Info(color.New(color.Bold).Sprint("Executing prepare hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "prepare", "plan"); err != nil {
			return err
		}
	}
	log.Info("")

	for _, source := range loaded {
		moduleDir := paths.ModuleDir(ac.TempDir, source.Name)

		log.Info("")
		log.Info("------------------------------------------------------------------------")
		log.Info("")

		if !paths.IsFile(filepath.Join(moduleDir, "tau.tfplan")) {
			log.Warnf(color.YellowString("No plan exists for %s", source.Name))
			continue
		}

		options := &shell.Options{
			WorkingDirectory: moduleDir,
			Stdout:           shell.Processors(&processors.Log{Level: log.InfoLevel}),
			Stderr:           shell.Processors(&processors.Log{Level: log.ErrorLevel}),
			Env:              source.Env,
		}

		extraArgs := getExtraArgs(ac.Engine.Compatibility.GetInvalidArgs("plan")...)
		extraArgs = append(extraArgs, "-out=tau.tfplan")
		if err := ac.Engine.Executor.Execute(options, "plan", extraArgs...); err != nil {
			return err
		}
	}

	log.Info("")
	log.Info("------------------------------------------------------------------------")
	log.Info("")

	log.Info(color.New(color.Bold).Sprint("Executing finish hook..."))
	for _, source := range loaded {
		if err := hooks.Run(source, "finish", "plan"); err != nil {
			return err
		}
	}
	log.Info("")

	return nil
}
