package cmd

import (
	"github.com/avinor/tau/internal/templates"
)

// Command description for a command that should just be passed through to terraform
type Command struct {
	Use              string
	ShortDescription string
	LongDescription  string
	Example          string
	SingleResource   bool
}

var (
	validCommands = map[string]Command{
		"apply": {
			Use:              "apply [-f SOURCE]",
			ShortDescription: "Builds or changes infrastructure",
			LongDescription:  "Builds or changes infrastructure",
		},
		"destroy": {
			Use:              "destroy [-f SOURCE",
			ShortDescription: "Destroy Tau-managed infrastructure",
			LongDescription:  "Destroy Tau-managed infrastructure",
		},
		"get": {
			Use:              "get [-f SOURCE]",
			ShortDescription: "Download and install modules for the configuration",
			LongDescription:  "Download and install modules for the configuration",
		},
		"import": {
			Use:              "import -f SOURCE",
			ShortDescription: "Import existing infrastructure into Terraform",
			LongDescription:  "Import existing infrastructure into Terraform",
			SingleResource:   true,
		},
		"output": {
			Use:              "output [-f SOURCE]",
			ShortDescription: "Read an output from a state file",
			LongDescription: templates.LongDesc(
				`Read an output from the state file for module deployed.
				`),
		},
		"refresh": {
			Use:              "refresh [-f SOURCE]",
			ShortDescription: "Update local state file against real resources",
			LongDescription:  "Update local state file against real resources",
		},
		"show": {
			Use:              "show [-f SOURCE]",
			ShortDescription: "Inspect Terraform state or plan",
			LongDescription:  "Inspect Terraform state or plan",
		},
		"taint": {
			Use:              "taint -f SOURCE",
			ShortDescription: "Manually mark a resource for recreation",
			LongDescription:  "Manually mark a resource for recreation",
			SingleResource:   true,
		},
		"untaint": {
			Use:              "untaint -f SOURCE",
			ShortDescription: "Manually unmark a resource as tainted",
			LongDescription:  "Manually unmark a resource as tainted",
			SingleResource:   true,
		},
	}
)
