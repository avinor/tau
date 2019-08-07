package cmd

import (
	"github.com/avinor/tau/internal/templates"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type outputCmd struct {
	meta

	loader *loader.Loader

	json bool
}

var (
	// outputRequiresSingleFile is returned if using json output and it tries to process multiple files
	outputRequiresSingleFile = errors.Errorf("can only process a single file when using json output")

	// outputLong is long description of output command
	outputLong = templates.LongDesc(`Print all the output variables from a module.
		If including the --json flag it will print output in json format.
		`)

	// outputExample is examples for output command
	outputExample = templates.Examples(`
		# Combine output from all files
		tau output

		# Print output from module.hcl in json format
		tau output -f module.hcl --json
	`)
)

// newOutputCmd creates a new output command
func newOutputCmd() *cobra.Command {
	oc := &outputCmd{}

	outputCmd := &cobra.Command{
		Use:                   "output [-f SORUCE]",
		Short:                 "Print output from module",
		Long:                  outputLong,
		Example:               outputExample,
		DisableFlagsInUseLine: true,
		SilenceUsage:          true,
		SilenceErrors:         true,
		Args:                  cobra.MaximumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := oc.processArgs(args); err != nil {
				return err
			}

			oc.init()

			return oc.run(args)
		},
	}

	f := outputCmd.Flags()
	f.BoolVar(&oc.json, "json", false, "print in json format")

	oc.addMetaFlags(outputCmd)

	return outputCmd
}

func (oc *outputCmd) init() {
	{
		options := &loader.Options{
			WorkingDirectory: paths.WorkingDir,
			TauDirectory:     oc.TauDir,
			MaxDepth:         1,
		}

		oc.loader = loader.New(options)
	}
}

func (oc *outputCmd) run(args []string) error {
	files, err := oc.loader.Load(oc.file)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		ui.NewLine()
		ui.Warn("No sources found")
		return nil
	}

	// if source defined then it can only deploy a single file, not folder
	if len(files) > 1 && oc.json {
		return outputRequiresSingleFile
	}

	// Verify all modules have been initialized
	if err := files.IsAllInitialized(); err != nil {
		return err
	}

	if err := hooks.RunAll(files, "prepare", "output"); err != nil {
		return err
	}

	// Check if any plans exist, if not then run plan first
	noVariablesExists := true
	for _, file := range files {
		if paths.IsFile(file.VariableFile()) {
			noVariablesExists = false
			continue
		}
	}

	if noVariablesExists {
		oc.resolveDependencies(files)
	}

	var values map[string]cty.Value

	for _, file := range files {
		ui.Separator(file.Name)

		if !paths.IsFile(file.VariableFile()) {
			ui.Warn("No values file exists for %s", file.Name)
			continue
		}

		outputProcessor := oc.Engine.Executor.NewOutputProcessor()

		options := &shell.Options{
			WorkingDirectory: file.ModuleDir(),
			Stdout:           shell.Processors(outputProcessor),
			Stderr:           shell.Processors(processors.NewUI(ui.Error)),
			Env:              file.Env,
		}

		if !oc.json {
			options.Stdout = append(options.Stdout, processors.NewUI(ui.Info))
		}

		extraArgs := getExtraArgs(oc.Engine.Compatibility.GetInvalidArgs("output")...)

		if oc.json {
			extraArgs = append(extraArgs, "-json")
		}

		if err := oc.Engine.Executor.Execute(options, "output", extraArgs...); err != nil {
			return err
		}

		if oc.json {
			output, err := outputProcessor.GetOutput()
			if err != nil {
				return err
			}
			values = output
		}

		paths.Remove(file.VariableFile())
	}

	ui.Separator("")

	if err := hooks.RunAll(files, "finish", "output"); err != nil {
		return err
	}

	ui.NewLine()

	if oc.json {
		obj := cty.ObjectVal(values)
		bytes, err := ctyjson.Marshal(obj, obj.Type())
		if err != nil {
			return err
		}

		ui.Output("%v", string(bytes))
	}

	return nil
}
