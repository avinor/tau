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
	switch command {
	case "init":
		return []string{"-backend-config", "-from-module"}
	case "plan":
		return []string{"-detailed-exitcode", "-out"}
	}

	return []string{}
}
