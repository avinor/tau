package cmd

import (
	"github.com/avinor/tau/internal/templates"
)

// Command that is available
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
		"plan": Command{
			Use:              "plan SOURCE [terraform options]",
			ShortDescription: "Generate and show an execution plan",
			LongDescription: templates.LongDesc(
				`Generate and show an execution plan for one or more modules.
				`),
		},
	}
)
