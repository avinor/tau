package cmd

import (
	"path"
	"os"
	"github.com/avinor/tau/pkg/dir"
	"github.com/avinor/tau/pkg/config"
	"github.com/go-errors/errors"
	"strings"
	"github.com/hashicorp/go-getter"
)

type Meta struct {
	Loader *config.Loader
	TempDir string
	SourceDir string
	SourceFile string
}

func (m *Meta) initMeta(args []string) error {
	{
		src, err := getSourceArg(args)
		if err != nil {
			return err
		}

		src, wd, err := splitSource(src)
		if err != nil {
			return err
		}

		m.SourceFile = src
		m.SourceDir = wd
	}

	m.TempDir = dir.TempDir(dir.Working, m.SourceFile)

	{
		options := &config.Options{
			WorkingDirectory: dir.Working,
			TempDirectory: m.TempDir,
		}

		m.Loader = config.NewLoader(options)
	}

	return nil
}

// getSourceArg finds argument that does not start with dash (-)
func getSourceArg(args []string) (string, error) {
	source := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			if source != "" {
				return "", errors.Errorf("Only one source argument should be defined")
			}

			source = arg
		}
	}

	if source == "" {
		return "", errors.Errorf("No source argument found")
	}

	return source, nil
}

// getExtraArgs returns all arguments starting with dash (-), but filters out invalid arguments
func getExtraArgs(args []string, invalidArgs ...string) []string {
	extraArgs := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			invalidArg := false

			for _, ia := range invalidArgs {
				if strings.HasPrefix(arg, ia) {
					invalidArg = true
				}
			}

			if !invalidArg {
				extraArgs = append(extraArgs, arg)
			}
		}
	}

	return extraArgs
}

// Split the source directory into working directory and source directory
func splitSource(src string) (string, string, error) {
	pwd := ""

	getterSource, err := getter.Detect(src, dir.Working, getter.Detectors)
	if err != nil {
		return "", "", errors.Errorf("Failed to detect source")
	}

	if strings.Index(getterSource, "file://") == 0 {
		rootPath := strings.Replace(getterSource, "file://", "", 1)

		fi, err := os.Stat(rootPath)
		if err != nil {
			return "", "", errors.Errorf("Failed to read %v", rootPath)
		}

		if !fi.IsDir() {
			pwd = path.Dir(rootPath)
			src = path.Base(rootPath)
		} else {
			pwd = rootPath
			src = "."
		}

		// log.Debugf("New source directory: %v", pwd)
		// log.Debugf("New source: %v", src)
	}

	return src, pwd, nil
}