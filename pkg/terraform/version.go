package terraform

import (
	"regexp"

	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

const (
	versionPattern = "Terraform v(\\d+.\\d+)"
)

var (
	versionRegex = regexp.MustCompile(versionPattern)
)

func Version() string {
	buffer := &processors.Buffer{}
	logp := &processors.Log{}

	options := &shell.Options{
		Stdout: shell.Processors(buffer),
		Stderr: shell.Processors(logp),
	}

	if err := shell.Execute(options, "terraform", "version"); err != nil {
		return ""
	}

	matches := versionRegex.FindAllStringSubmatch(buffer.Stdout(), -1)

	if len(matches) < 1 && len(matches[0]) < 2 {
		return ""
	}

	return matches[0][1]
}
