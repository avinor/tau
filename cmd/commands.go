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
}

var (
	validCommands = map[string]Command{
		// "apply": Command{
		// 	Use:              "apply",
		// 	ShortDescription: "Builds or changes infrastructure",
		// 	LongDescription:  "Builds or changes infrastructure",
		// 	PassThrough:      true,
		// },
		"plan": {
			Use:              "plan SOURCE [terraform options]",
			ShortDescription: "Generate and show an execution plan",
			LongDescription: templates.LongDesc(
				`Generate and show an execution plan for one or more modules.
				`),
		},
		"output": {
			Use:              "output SOURCE [terraform options]",
			ShortDescription: "Read an output from a state file",
			LongDescription: templates.LongDesc(
				`Read an output from the state file for module deployed.
				`),
		},
	}
)
