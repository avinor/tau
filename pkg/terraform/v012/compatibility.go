package v012

type Compatibility struct {
}

func (c *Compatibility) GetValidCommands() []string {
	return []string{
		"apply",
		"destroy",
		"env",
		"get",
		"graph",
		"import",
		"init",
		"output",
		"plan",
		"providers",
		"refresh",
		"show",
		"taint",
		"untaint",
		"workspace",
		"force-unlock",
		"state",
	}
}

func (c *Compatibility) GetInvalidArgs(command string) []string {
	if command == "init" {
		return []string{"-backend-config", "-from-module"}
	}

	return []string{}
}
