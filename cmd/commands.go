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
		"apply": Command{
			Use:              "apply [-f SOURCE]",
			ShortDescription: "Builds or changes infrastructure",
			LongDescription:  "Builds or changes infrastructure",
		},
		"destroy": Command{
			Use:              "destroy [-f SOURCE",
			ShortDescription: "Destroy Tau-managed infrastructure",
			LongDescription:  "Destroy Tau-managed infrastructure",
		},
		"get": Command{
			Use:              "get [-f SOURCE]",
			ShortDescription: "Download and install modules for the configuration",
			LongDescription:  "Download and install modules for the configuration",
		},
		"import": Command{
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
		"plan": {
			Use:              "plan [-f SOURCE]",
			ShortDescription: "Generate and show an execution plan",
			LongDescription: templates.LongDesc(
				`Generate and show an execution plan for one or more modules.
				`),
		},
		"refresh": Command{
			Use:              "refresh [-f SOURCE]",
			ShortDescription: "Update local state file against real resources",
			LongDescription:  "Update local state file against real resources",
		},
		"show": Command{
			Use:              "show [-f SOURCE]",
			ShortDescription: "Inspect Terraform state or plan",
			LongDescription:  "Inspect Terraform state or plan",
		},
		"taint": Command{
			Use:              "taint -f SOURCE",
			ShortDescription: "Manually mark a resource for recreation",
			LongDescription:  "Manually mark a resource for recreation",
			SingleResource:   true,
		},
		"untaint": Command{
			Use:              "untaint -f SOURCE",
			ShortDescription: "Manually unmark a resource as tainted",
			LongDescription:  "Manually unmark a resource as tainted",
			SingleResource:   true,
		},
	}
)
