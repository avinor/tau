package v012

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/avinor/tau/pkg/config/loader"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/avinor/tau/pkg/hooks"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// DependencyProcessor implements the def.DepdendencyProcessor interface
type DependencyProcessor struct {
	ParsedFile *loader.ParsedFile
	File       *hclwrite.File

	executor *Executor
	resolver *Resolver

	// acceptApplyFailure should be set if its acceptable that apply fails. Should be set if
	// no backend is found or unsupported attribute, most probably means a dependency is not deployed
	acceptApplyFailure bool
}

// NewDependencyProcessor creates a new dependencyProcessor structure from input arguments
func NewDependencyProcessor(file *loader.ParsedFile, executor *Executor, resolver *Resolver) *DependencyProcessor {
	f := hclwrite.NewEmptyFile()

	return &DependencyProcessor{
		ParsedFile: file,
		File:       f,

		executor: executor,
		resolver: resolver,
	}
}

// WriteContent writes the context of main.tf
func (d *DependencyProcessor) WriteContent(dest string) error {
	file := filepath.Join(dest, "main.tf")
	if err := ioutil.WriteFile(file, d.File.Bytes(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// Process the dependency and return the variables from output.
func (d *DependencyProcessor) Process() (map[string]cty.Value, bool, error) {
	dest := d.ParsedFile.DependencyDir(d.ParsedFile.Name)
	if err := d.WriteContent(dest); err != nil {
		return nil, false, err
	}

	debugLog := processors.NewUI(ui.Debug)
	errorLog := processors.NewUI(ui.Error)

	if err := hooks.Run(d.ParsedFile, "prepare", "init"); err != nil {
		return nil, false, err
	}

	options := &shell.Options{
		Stdout:           shell.Processors(debugLog),
		Stderr:           shell.Processors(d, errorLog),
		WorkingDirectory: dest,
		Env:              d.ParsedFile.Env,
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

// Write implements the shell.OutputProcessor interface so it can use DependencyProcessor
// as a processer when executing commands, and therefore set acceptApplyFailure if it detects
// acceptable error messages in output
func (d *DependencyProcessor) Write(line string) bool {
	if strings.Contains(line, "Unable to find remote state") {
		d.acceptApplyFailure = true
	}

	if strings.Contains(line, "Unsupported attribute") {
		d.acceptApplyFailure = true
	}

	return true
}
