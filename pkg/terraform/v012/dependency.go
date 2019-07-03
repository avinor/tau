package v012

import (
	"path/filepath"
	"strings"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

type dependencyProcessor struct {
	Source *config.Source
	File   *hclwrite.File

	executor *Executor
	resolver *Resolver

	// acceptApplyFailure should be set if its acceptable that apply fails. Should be set if
	// no backend is found or unsupported attribute, most probably means a dependency is not deployed
	acceptApplyFailure bool
}

func NewDependencyProcessor(source *config.Source, executor *Executor, resolver *Resolver) *dependencyProcessor {
	f := hclwrite.NewEmptyFile()

	return &dependencyProcessor{
		Source: source,
		File:   f,

		executor: executor,
		resolver: resolver,
	}
}

func (d *dependencyProcessor) Name() string {
	return d.Source.Name
}

func (d *dependencyProcessor) Content() []byte {
	return d.File.Bytes()
}

func (d *dependencyProcessor) Process(dest string) (map[string]cty.Value, bool, error) {
	debugLog := processors.NewUI(ui.Debug)
	errorLog := processors.NewUI(ui.Error)

	if err := hooks.Run(d.Source, "prepare", "init"); err != nil {
		return nil, false, err
	}

	options := &shell.Options{
		Stdout:           shell.Processors(debugLog),
		Stderr:           shell.Processors(d, errorLog),
		WorkingDirectory: dest,
		Env:              d.Source.Env,
	}

	base := filepath.Base(dest)

	ui.Info("- %s", base)

	ui.Debug("running terraform init on %s", base)
	if err := d.executor.Execute(options, "init", "-input=false"); err != nil {
		return nil, false, err
	}

	ui.Debug("running terraform apply on %s", base)
	if err := d.executor.Execute(options, "apply", "-auto-approve", "-input=false"); err != nil {
		// If it accepts failure then just exit with no error, but create = false
		if d.acceptApplyFailure {
			return nil, false, nil
		}

		return nil, false, err
	}

	buffer := &processors.Buffer{}
	options.Stdout = shell.Processors(buffer)

	ui.Debug("reading output from %s", base)
	if err := d.executor.Execute(options, "output", "-json"); err != nil {
		return nil, false, err
	}

	values, err := d.resolver.ResolveStateOutput([]byte(buffer.String()))
	if err != nil {
		return nil, false, err
	}

	return values, true, nil
}

func (d *dependencyProcessor) Write(line string) bool {
	if strings.Contains(line, "Unable to find remote state") {
		d.acceptApplyFailure = true
	}

	if strings.Contains(line, "Unsupported attribute") {
		d.acceptApplyFailure = true
	}

	return true
}
