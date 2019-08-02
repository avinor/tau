package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/helpers/hclcontext"
	"github.com/avinor/tau/pkg/helpers/ui"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/zclconf/go-cty/cty"
)

var (
	parser = hclparse.NewParser()
)

// File is a config file that is not parsed. It will read the content and store for later processing
// Used when same config files should be read with different evaluation contexts.
//
// Children is used to add files it is dependent on that should be processed together with this one.
// When parsing configuration it will merge config from all children together with File. Current file
// takes presedence and will overwrite all children files.
type File struct {
	Name     string
	FullPath string
	Content  []byte

	Children []*File
}

// NewFile returns a new File. It will check that it exists and read content, but not parse it
func NewFile(file string) (*File, error) {
	name := filepath.Base(file)
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}

	ui.Info("- loading %s", name)

	if _, err := os.Stat(absPath); err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:     name,
		FullPath: absPath,
		Content:  content,
		Children: []*File{},
	}, nil
}

// AddChild adds a child File.
func (f *File) AddChild(file *File) {
	f.Children = append(f.Children, file)
}

// Config returns the full configuration for file. This includes the merged configuration from
// all children. Should only call this once as it will do full parsing of file and all children
func (f *File) Config() (*Config, error) {
	configs := []*Config{}
	context := f.getEvalContext()

	for _, file := range append(f.Children, f) {
		parsed, err := file.parse(context)
		if err != nil {
			return nil, err
		}

		configs = append(configs, parsed)
	}

	config := &Config{}
	if err := config.Merge(configs); err != nil {
		return nil, err
	}

	return config, nil
}

// parse the file using evaluation context from input. It will add source variables to the context
// variables if not set that makes it possible to retrieve file name etc
func (f *File) parse(context *hcl.EvalContext) (*Config, error) {
	hclFile, diags := parser.ParseHCL(f.Content, f.FullPath)
	if diags.HasErrors() {
		return nil, diags
	}

	config := &Config{}
	bodyDiags := gohcl.DecodeBody(hclFile.Body, context, config)

	if bodyDiags.HasErrors() {
		return nil, diags
	}

	return config, nil
}

// getEvalContext gets the context for this file. Adding variables for source to default context
func (f File) getEvalContext() *hcl.EvalContext {
	ext := filepath.Ext(f.Name)

	values := map[string]cty.Value{
		"source": cty.ObjectVal(map[string]cty.Value{
			"path":     cty.StringVal(f.FullPath),
			"name":     cty.StringVal(f.Name[0 : len(f.Name)-len(ext)]),
			"filename": cty.StringVal(f.Name),
		}),
	}

	return hclcontext.WithVariables(hclcontext.Default, values)
}
