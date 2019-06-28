package cmd

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/getter"
	"github.com/avinor/tau/pkg/paths"
	"github.com/avinor/tau/pkg/terraform"
	"github.com/fatih/color"
	"github.com/go-errors/errors"
	"github.com/spf13/pflag"
)

type meta struct {
	timeout int

	Getter *getter.Client
	Loader *config.Loader
	Engine *terraform.Engine

	TempDir    string
	SourceDir  string
	SourceFile string
}

func (m *meta) processArgs(args []string) error {
	log.Debug(color.New(color.Bold).Sprint("Processing arguments..."))

	if workingDir == "" {
		workingDir = paths.WorkingDir
	}

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

		log.Debugf("- Source file: %s", m.SourceFile)
		log.Debugf("- Source dir: %s", m.SourceDir)
	}

	m.TempDir = paths.TempDir(workingDir, m.SourceFile)

	log.Debugf("- Temp dir: %s", m.TempDir)

	{
		timeout := time.Duration(m.timeout) * time.Second

		log.Debugf("- Http timeout: %s", timeout)

		options := &getter.Options{
			HttpClient: &http.Client{
				Timeout: timeout,
			},
			WorkingDirectory: m.SourceDir,
		}

		m.Getter = getter.New(options)
	}

	{
		options := &config.Options{
			WorkingDirectory: paths.WorkingDir,
			TempDirectory:    m.TempDir,
			Getter:           m.Getter,
		}

		m.Loader = config.NewLoader(options)
	}

	log.Debug("")

	{
		m.Engine = terraform.NewEngine()
	}

	return nil
}

func (m *meta) addMetaFlags(f *pflag.FlagSet) {
	f.IntVar(&m.timeout, "timeout", 10, "timeout for http client when retrieving sources")
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
	pwd := paths.WorkingDir
	client := getter.New(nil)

	getterSource, err := client.Detect(src)
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
	}

	return src, pwd, nil
}
