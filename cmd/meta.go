package cmd

import (
	"time"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	hooksdef "github.com/avinor/tau/pkg/hooks/def"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/avinor/tau/pkg/terraform"
	"github.com/avinor/tau/pkg/terraform/def"
	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	// noSourceInPath is returned when there are no source files in path
	noSourceInPath = errors.Errorf("no source files found in path")
)

type meta struct {
	timeout            int
	maxDependencyDepth int
	files              []string
	noAutoInit         bool

	Engine *terraform.Engine
	Getter *getter.Client
	Loader *loader.Loader
	Runner *hooks.Runner

	TauDir   string
	CacheDir string
}

type initOptions struct {
	purge       bool
	noOverrides bool
	source      *config.Module
	reconfigure bool
}

func (m *meta) init(args []string) error {
	if workingDir == "" {
		workingDir = paths.WorkingDir
	}

	m.TauDir = paths.JoinAndCreate(workingDir, paths.TauPath)
	m.CacheDir = paths.JoinAndCreate(workingDir, paths.CachePath)

	{
		timeout := time.Duration(m.timeout) * time.Second

		options := &getter.Options{
			Timeout:          timeout,
			WorkingDirectory: workingDir,
		}

		m.Getter = getter.New(options)
	}

	{
		options := &loader.Options{
			WorkingDirectory: workingDir,
			TauDirectory:     m.TauDir,
			CacheDirectory:   m.CacheDir,
			MaxDepth:         m.maxDependencyDepth,
			Getter:           m.Getter,
		}

		m.Loader = loader.New(options)
	}

	{
		m.Runner = hooks.New(&hooksdef.Options{
			Getter:   m.Getter,
			CacheDir: m.CacheDir,
		})
	}

	{
		m.Engine = terraform.NewEngine(&def.Options{
			Runner: m.Runner,
		})
	}

	ui.Debug("tau dir: %s", m.TauDir)
	ui.Debug("http timeout: %s", m.timeout)
	ui.Debug("max dependency depth: %s", m.maxDependencyDepth)

	return nil
}

// addMetaFlags adds common arguments defined on meta to the command cmd.
// All commands that include meta have to call this!
func (m *meta) addMetaFlags(cmd *cobra.Command) {
	f := cmd.Flags()
	f.IntVar(&m.timeout, "timeout", 10, "timeout for http client when retrieving sources")
	f.StringArrayVarP(&m.files, "file", "f", []string{"."}, "file or directory to run configuration for")
	f.BoolVar(&m.noAutoInit, "no-auto-init", false, "disable auto init")
	f.IntVar(&m.maxDependencyDepth, "max-dependency-depth", 1, "defines max dependency depth when traversing dependencies") //nolint:lll
}

// load wraps the Loader.Load function to load all files and return to caller.
// Also prints some helpful messages and checks that there are loaded files.
func (m *meta) load() (loader.ParsedFileCollection, error) {
	ui.Header("Loading files...")

	// load all sources
	files, err := m.Loader.Load(m.files)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		ui.Info("- Loaded %s", file)
	}

	if len(files) == 0 {
		return nil, noSourceInPath
	}

	return files, nil
}

// resolveDependencies resolves the dependencies for all files
func (m *meta) resolveDependencies(file *loader.ParsedFile) (bool, error) {
	if file.Config.Inputs == nil {
		return true, nil
	}

	ui.Header("Resolving dependencies...")

	success, err := m.Engine.ResolveDependencies(file)
	if err != nil {
		return false, err
	}

	if !success {
		ui.NewLine()
		ui.Info(color.GreenString("Some of the dependencies failed to resolve. This can be because dependency"))
		ui.Info(color.GreenString("have not been applied yet, and therefore it cannot read remote-state."))
		ui.NewLine()

		return false, nil
	}

	if err := m.Engine.WriteInputVariables(file); err != nil {
		return false, err
	}

	return true, nil
}

// autoInit can be called by any command to auto initialize the module
func (m *meta) autoInit(file *loader.ParsedFile) error {
	if file.IsInitialized() {
		return nil
	}

	if m.noAutoInit {
		ui.Debug("no-auto-init set, not initializing module")
		return nil
	}

	return m.runInit(file, nil)
}

// runInit initializes the parsed file
func (m *meta) runInit(file *loader.ParsedFile, options *initOptions) error {
	if options == nil {
		options = &initOptions{}
	}

	ui.Header("Initializing tau...")

	// Loading module

	if !options.reconfigure {
		module := file.Config.Module

		if options.source != nil && options.source.Source != "" {
			module = options.source
		}

		if module.Version != "" {
			ui.Info("- Loading module from terraform registry %s, version %s", module.Source, module.Version)
		} else {
			ui.Info("- Loading module from %s", module.Source)
		}

		if err := m.Getter.Get(module.GetSource(), file.ModuleDir()); err != nil {
			return err
		}
	}

	// Creating overrides

	if !options.noOverrides {
		ui.Info("- Creating overrides for backend")

		if err := m.Engine.CreateOverrides(file); err != nil {
			return err
		}
	}

	// Executing terraform command

	ui.NewLine()
	ui.Info(color.New(color.FgGreen, color.Bold).Sprint("Tau has been successfully initialized!"))
	ui.NewLine()

	shellOptions := &shell.Options{
		WorkingDirectory: file.ModuleDir(),
		Stdout:           shell.Processors(processors.NewUI(ui.Info)),
		Stderr:           shell.Processors(processors.NewUI(ui.Error)),
		Env:              file.Env,
	}

	extraArgs := getExtraArgs(m.Engine.Compatibility.GetInvalidArgs("init")...)

	if options.reconfigure {
		extraArgs = append(extraArgs, "-reconfigure", "-force-copy")
	}

	if err := m.Engine.Executor.Execute(shellOptions, "init", extraArgs...); err != nil {
		return err
	}

	return nil
}
