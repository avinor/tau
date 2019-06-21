package v012

type Compatibility struct {
}

func (c *Compatibility) GetValidCommands() []string {
	return nil
}

func (c *Compatibility) GetInvalidArgs(command string) []string {
	return nil
}