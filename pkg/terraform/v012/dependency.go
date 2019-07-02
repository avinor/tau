package v012

import (
	"path/filepath"
	"strings"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
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

	// unableToFindRemoteState is set to true if it gets this error when running terraform plan.
	// This probably means that a dependency has not been run and needs to run first.
	unableToFindRemoteState bool
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
	debugLog := &processors.Log{Level: log.DebugLevel}
	errorLog := &processors.Log{Level: log.ErrorLevel}

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

	log.Infof("- %s", base)

	log.Debugf("running terraform init on %s", base)
	if err := d.executor.Execute(options, "init"); err != nil {
		return nil, false, err
	}

	log.Debugf("running terraform apply on %s", base)
	if err := d.executor.Execute(options, "apply"); err != nil {
		if d.unableToFindRemoteState {
			return nil, false, nil
		}

		return nil, false, err
	}

	buffer := &processors.Buffer{}
	options.Stdout = shell.Processors(buffer)

	log.Debugf("reading output from %s", base)
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
		d.unableToFindRemoteState = true
	}

	if d.unableToFindRemoteState == true {
		return false
	}

	return true
}
