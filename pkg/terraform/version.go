package terraform

import (
	"regexp"

	"github.com/apex/log"
	"github.com/avinor/tau/pkg/shell"
	"github.com/avinor/tau/pkg/shell/processors"
)

const (
	versionPattern = "Terraform v(.*)"
)

var (
	versionRegex = regexp.MustCompile(versionPattern)
)

// Version checks the current terraform version, returns empty string if not found
func Version() string {
	buffer := &processors.Buffer{}
	logp := &processors.Log{Level: log.ErrorLevel}

	options := &shell.Options{
		Stdout: shell.Processors(buffer),
		Stderr: shell.Processors(logp),
	}

	if err := shell.Execute(options, "terraform", "version"); err != nil {
		return ""
	}

	matches := versionRegex.FindAllStringSubmatch(buffer.String(), -1)

	if len(matches) < 1 || len(matches[0]) < 2 {
		return ""
	}

	return matches[0][1]
}
