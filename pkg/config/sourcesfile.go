package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-errors/errors"
	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

type SourceFile struct {
	Sources []*Source `hcl:"source,block"`
}

const (
	tauFileName = "tau.hcl"
)

// LoadTempDir loads sources from temp dir. Temp dir already has to be initialized
func LoadSourcesFile(tempdir string) ([]*Source, error) {
	file := filepath.Join(tempdir, tauFileName)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil, errors.Errorf("Temporary tau file does not exist. Make sure init is run first.")
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	sf := &SourceFile{}
	if err := Parse(b, file, sf); err != nil {
		return nil, err
	}

	return sf.Sources, nil
}

// Save loaded sources in temp directory
func SaveSourcesFile(sources []*Source, tempdir string) error {
	sf := &SourceFile{
		Sources: sources,
	}

	f := hclwrite.NewEmptyFile()
	body := f.Body()

	gohcl.EncodeIntoBody(sf, body)

	file := filepath.Join(tempdir, tauFileName)

	return ioutil.WriteFile(file, f.Bytes(), os.ModePerm)
}
