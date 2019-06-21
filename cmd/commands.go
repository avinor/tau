package cmd

// Command that is available
type Command struct {
	Use              string
	ShortDescription string
	LongDescription  string
	Example          string
	PassThrough bool
}

var (
	validCommands map[string]Command
)

func init() {
	validCommands = map[string]Command{
		"apply": Command{
			Use: "apply",
			ShortDescription: "Builds or changes infrastructure",
			LongDescription: "Builds or changes infrastructure",
			PassThrough: true,
		},
		"init": Command{
			Use: "init",
			ShortDescription: "Initialize a Terraform working directory",
			LongDescription: "Initialize a Terraform working directory",
		},
		"plan": Command{
			Use: "plan",
			ShortDescription: "Generate and show an execution plan",
			LongDescription: "Generate and show an execution plan",
			PassThrough: true,
		},
	}
}