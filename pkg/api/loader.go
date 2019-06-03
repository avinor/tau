package api

import (
	"crypto/sha1"
	"encoding/hex"
	"os"

	log "github.com/sirupsen/logrus"
)

type Loader struct {
	pwd string
}

func NewLoader() *Loader {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get working directory: %s", err)
	}

	return &Loader{
		pwd: dir,
	}
}

func (l *Loader) Load(src string) error {

	pkg, err := LoadPackage(src, l.pwd)
	if err != nil {
		return err
	}

	for _, m := range pkg.modules {
		log.Infof("%s", m.src)
	}

	// Load src into dst
	// Resolve all .hcl / .tau files
	// Create config / module for all files

	// Check dependencies
	// Load unresolved dependencies into dst
	// Create config for dependencies

	return nil
}

func hash(src string) string {
	h := sha1.New()
	h.Write([]byte(src))
	return hex.EncodeToString(h.Sum(nil))
}
