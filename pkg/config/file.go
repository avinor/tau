package config

import (
	"path/filepath"

	"github.com/avinor/tau/pkg/helpers/hclcontext"
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

	children []*File

	// context to evaluate expressions with. New variables can be added to this by calling AddToContext()
	context *hcl.EvalContext
}

// NewFile returns a new File. It will check that it exists and read content, but not parse it
func NewFile(filename string, content []byte) (*File, error) {
	name := filepath.Base(filename)
	absPath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	return &File{
		Name:     name,
		FullPath: absPath,
		Content:  content,
		children: []*File{},
		context:  getNewEvalContext(filename),
	}, nil
}

// AddChild adds a child File.
func (f *File) AddChild(file *File) {
	f.children = append(f.children, file)
}

// AddToContext adds a variable to the evaluation context of this file
func (f *File) AddToContext(key string, value cty.Value) {
	f.context.Variables[key] = value
}

// EvalContext returns the evaluation context for this file
func (f *File) EvalContext() *hcl.EvalContext {
	return f.context
}

// Config returns the full configuration for file. This includes the merged configuration from
// all children. Should only call this once as it will do full parsing of file and all children
func (f *File) Config() (*Config, error) {
	configs := []*Config{}

	for _, file := range append(f.children, f) {
		parsed, err := file.parse(f.context)
		if err != nil {
			return nil, err
		}

		configs = append(configs, parsed)
	}

	config := &Config{}
	if err := config.Merge(configs); err != nil {
		return nil, err
	}

	config.PostProcess(f)

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
		return nil, bodyDiags
	}

	return config, nil
}

// GetEvalContext gets the context for this file. Adding variables for source to default context
func getNewEvalContext(fullPath string) *hcl.EvalContext {
	name := filepath.Base(fullPath)
	ext := filepath.Ext(name)

	value := cty.ObjectVal(map[string]cty.Value{
		"path":     cty.StringVal(fullPath),
		"name":     cty.StringVal(name[0 : len(name)-len(ext)]),
		"filename": cty.StringVal(name),
	})

	context := hclcontext.NewContext()
	context.Variables["source"] = value

	return context
}
