package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
)

type SourceFile struct {
	Sources []*Source `hcl:"source,block"`
}

// LoadTempDir loads sources from temp dir. Temp dir already has to be initialized
func LoadSourcesFile() ([]*Source, error) {
	return nil, nil
}

// Save loaded sources in temp directory
func SaveSourcesFile(sources []*Source, tempdir string) error {
	sf := &SourceFile{
		Sources: sources,
	}

	f := hclwrite.NewEmptyFile()
	body := f.Body()

	gohcl.EncodeIntoBody(sf, body)

	file := filepath.Join(tempdir, "tau.hcl")

	return ioutil.WriteFile(file, f.Bytes(), os.ModePerm)
}
