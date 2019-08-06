package v012

// Compatibility implements the def.VersionCompatibility interface
type Compatibility struct{}

// GetValidCommands returns valid commands for terraform 0.12
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

// GetInvalidArgs returns invalid arguments for command. These are invalid because they
// are usually set by tau and should not be configured by user.
func (c *Compatibility) GetInvalidArgs(command string) []string {
	switch command {
	case "init":
		return []string{"-backend-config", "-from-module"}
	case "plan":
		return []string{"-detailed-exitcode", "-out"}
	case "apply":
		return []string{"-auto-approve", "-input"}
	}

	return []string{}
}
