package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/avinor/tau/pkg/hclcontext"
	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

const (
	tauFileName = "tau.hcl"
)

var (
	temporaryTauFileNotExistError = errors.Errorf("Temporary tau file does not exist. Make sure init is run first.")
)

// Checkpoint stores a collection of sources for later reload
type Checkpoint struct {
	Sources []*Source `hcl:"source,block"`
}

// LoadCheckpoint loads a checkpoint file from the temporary directory, that has to
// exist before loading. Use checkpoint to save loaded sources between different commands
// , ie. terraform init -> terraform plan
func LoadCheckpoint(tempdir string) ([]*Source, error) {
	file := filepath.Join(tempdir, tauFileName)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, temporaryTauFileNotExistError
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	f, err := ParseFile(&SourceFile{File: file, Content: b})
	if err != nil {
		return nil, err
	}

	cp := &Checkpoint{}
	if err := ParseBody(f.Body, hclcontext.Default, cp); err != nil {
		return nil, err
	}

	return cp.Sources, nil
}

// SaveCheckpoint saves the sources to a checkpoint file in temporary directory.
// To reload call LoadCheckpoint
func SaveCheckpoint(sources []*Source, tempdir string) error {
	cp := &Checkpoint{
		Sources: sources,
	}

	f := hclwrite.NewEmptyFile()
	body := f.Body()

	gohcl.EncodeIntoBody(cp, body)

	file := filepath.Join(tempdir, tauFileName)

	return ioutil.WriteFile(file, f.Bytes(), os.ModePerm)
}
