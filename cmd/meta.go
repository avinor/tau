package cmd

import (
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/config"
	"github.com/go-errors/errors"
	"strings"
)

type Meta struct {
	Loader *config.Loader
	Getter *getter.Client
	TempDir string
	SourceDir string
	SourceFile string
}

func (m *Meta) initMeta(args []string) error {
	tempDir := filepath.Join(options.WorkingDirectory, ".tau", utils.Hash(profile))

	loader, err := config.NewLoader(tempDir)
	if err != nil {
		return err
	}

	source, err := utils.GetSourceArg(args)
	if err != nil {
		return err
	}

	client := config.New(source, &config.Options{
		LoadSources:  true,
		CleanTempDir: true,
	})

	return nil
}

// GetSourceArg finds argument that does not start with dash (-)
func GetSourceArg(args []string) (string, error) {
	source := ""
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			if source != "" {
				return "", errors.Errorf("Only one source argument should be defined")
			}

			source = arg
		}
	}

	return source, nil
}

// GetExtraArgs returns all arguments starting with dash (-), but filters out invalid arguments
func GetExtraArgs(args []string, invalidArgs ...string) []string {
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