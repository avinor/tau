package getter

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
	"github.com/otiai10/copy"
)

// LocalGetter is a wrapper around the go-getter implementation of FileGetter
// This is due to the go-getter implementation does not support copying a local
// directory, it will only create a symlink. That causes some problems for
// modules that should be copied into temp folder
type LocalGetter struct {
	getter.FileGetter
}

// Get overwrites FileGetter.Get if Copy is set to true. Otherwise it will just
// call the FileGetter.Get function and use standard functionallity
func (g *LocalGetter) Get(dst string, u *url.URL) error {
	if !g.Copy {
		return g.FileGetter.Get(dst, u)
	}

	path := u.Path
	if u.RawPath != "" {
		path = u.RawPath
	}

	// The source path must exist and be a directory to be usable.
	if fi, err := os.Stat(path); err != nil {
		return fmt.Errorf("source path error: %s", err)
	} else if !fi.IsDir() {
		return fmt.Errorf("source path must be a directory")
	}

	_, err := os.Lstat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// If the destination already exists, just delete it
	if err == nil {
		// Remove the destination
		if err := os.Remove(dst); err != nil {
			return err
		}
	}

	// Create all the parent directories
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	return copy.Copy(path, dst)
}

// GetFile just wraps the FileGetter function
func (g *LocalGetter) GetFile(dst string, u *url.URL) error {
	return g.FileGetter.GetFile(dst, u)
}
