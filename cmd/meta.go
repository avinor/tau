package cmd

import (
	"time"

	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/helpers/paths"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/terraform"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type meta struct {
	timeout            int
	maxDependencyDepth int
	file               string

	Engine *terraform.Engine
	Getter *getter.Client
	Loader *loader.Loader
	Runner *hooks.Runner

	TauDir   string
	CacheDir string
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
			MaxDepth:         m.maxDependencyDepth,
		}

		m.Loader = loader.New(options)
	}

	{
		m.Engine = terraform.NewEngine()
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
	f.StringVarP(&m.file, "file", "f", ".", "file or directory to run configuration for")
	f.IntVar(&m.maxDependencyDepth, "max-dependency-depth", 1, "defines max dependency depth when traversing dependencies") //nolint:lll
}

// resolveDependencies resolves the dependencies for all files
func (m *meta) resolveDependencies(files []*loader.ParsedFile) error {
	showDepFailureInfo := false
	ui.Header("Resolving dependencies...")
	for _, file := range files {
		if file.Config.Inputs == nil {
			continue
		}

		success, err := m.Engine.ResolveDependencies(file)
		if err != nil {
			return err
		}

		if !success {
			showDepFailureInfo = true
			continue
		}

		if err := m.Engine.WriteInputVariables(file); err != nil {
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
