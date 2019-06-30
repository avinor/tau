package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/hclcontext"

	"github.com/apex/log"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/zclconf/go-cty/cty"
)

var (
	loadedSourceFiles = map[string]*SourceFile{}
)

// SourceFile is a file to load config source from
type SourceFile struct {
	File    string
	Content []byte
}

// GetSourceFile returns the SourceFile for input file. If file is already loaded
// before it will return a pointer to that, otherwise it will load the file and
// return a new pointer
func GetSourceFile(file string) (*SourceFile, error) {
	if _, already := loadedSourceFiles[file]; already {
		return loadedSourceFiles[file], nil
	}

	log.Infof("- loading %s", filepath.Base(file))

	if _, err := os.Stat(file); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sf := &SourceFile{
		File:    file,
		Content: content,
	}

	loadedSourceFiles[file] = sf
	return sf, nil
}

// Config returns the configuration struct for SourceFile. It will parse the content
// of file.
func (sf *SourceFile) Config() (*Config, error) {
	return sf.configWithContext(sf.EvalContext())
}

// GetAutoImport returns a list of all files that should be auto imported together
// with this sourcefile. That is all files that have _auto in filename and resides in
// same folder as sourcefile
//
// Returned list will not include the sourcefile itself, this is just a list of the
// additional files that should be loaded.
func (sf *SourceFile) GetAutoImport() ([]*SourceFile, error) {
	sources := []*SourceFile{}
	path := filepath.Dir(sf.File)

	autoFiles, err := findFiles(path, autoMatchFunc)
	if err != nil {
		return nil, err
	}

	for _, file := range autoFiles {
		source, err := GetSourceFile(file)
		if err != nil {
			return nil, err
		}

		sources = append(sources, source)
	}

	return sources, nil
}

// ConfigMergedWith returns the config from sourcefile merged together with the
// config from input sources.
func (sf *SourceFile) ConfigMergedWith(sources []*SourceFile) (*Config, error) {
	config, err := sf.Config()
	if err != nil {
		return nil, err
	}

	configs := []*Config{config}

	for _, source := range sources {
		config, err := source.configWithContext(sf.EvalContext())
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}

	return sf.mergeConfigs(configs), nil
}

// EvalContext returns hcl eval context for SourceFile. It will add a variable source.file
// that is full path for file, and source.name that is base name of file
func (sf *SourceFile) EvalContext() *hcl.EvalContext {
	values := map[string]cty.Value{
		"source": cty.ObjectVal(map[string]cty.Value{
			"file": cty.StringVal(sf.File),
			"name": cty.StringVal(filepath.Base(sf.File)),
		}),
	}

	return hclcontext.WithVariables(hclcontext.Default, values)
}

// configWithContext returns the parsed context using a custom evalcontext
func (sf *SourceFile) configWithContext(context *hcl.EvalContext) (*Config, error) {
	file, err := ParseFile(sf)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := ParseBody(file.Body, context, config); err != nil {
		return nil, err
	}

	return config, nil
}

// mergeConfig takes each element and merge them together. Instead of merging entire file
// this will ensure to merge attribute by attribute for input block for instance. A full file
// merge will replace entire block. For data, dependencies and hooks it use the last definition
// of element, overriding previous definition.
func (sf *SourceFile) mergeConfigs(configs []*Config) *Config {
	new := &Config{}
	datas := map[string]Data{}
	dependencies := map[string]Dependency{}
	hooks := map[string]Hook{}

	// Merge all fields together
	for _, config := range configs {
		for _, data := range config.Datas {
			datas[fmt.Sprintf("%s.%s", data.Type, data.Name)] = data
		}
		for _, dep := range config.Dependencies {
			dependencies[dep.Name] = dep
		}
		for _, hook := range config.Hooks {
			hooks[hook.Type] = hook
		}

		if new.Environment == nil {
			new.Environment = config.Environment
		} else if config.Environment != nil {
			new.Environment.Config = hcl.MergeBodies([]hcl.Body{new.Environment.Config, config.Environment.Config})
		}

		if new.Backend == nil {
			new.Backend = config.Backend
		} else if config.Backend != nil {
			new.Backend.Config = hcl.MergeBodies([]hcl.Body{new.Backend.Config, config.Backend.Config})
		}

		if config.Module != nil {
			new.Module = config.Module
		}

		if new.Inputs == nil {
			new.Inputs = config.Inputs
		} else if config.Inputs != nil {
			new.Inputs.Config = hcl.MergeBodies([]hcl.Body{new.Inputs.Config, config.Inputs.Config})
		}
	}

	// Set array elements
	for _, data := range datas {
		new.Datas = append(new.Datas, data)
	}
	for _, dep := range dependencies {
		new.Dependencies = append(new.Dependencies, dep)
	}
	for _, hook := range hooks {
		new.Hooks = append(new.Hooks, hook)
	}

	return new
}
