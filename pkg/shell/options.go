package shell

// Options when running shell command
type Options struct {
	WorkingDirectory string
	Stdout           []OutputProcessor
	Stderr           []OutputProcessor
	Env              map[string]string
}

// OutputProcessor can process a line from command output, does not separate between
// stdout and stderr. If it returns true it will proceed to next OutputProcessor in line.
// With false it will not proceed to other processors
type OutputProcessor interface {
	Write(line string) bool
}

// Processors is a helper function to create a list of processors
func Processors(pros ...OutputProcessor) []OutputProcessor {
	return pros
}
