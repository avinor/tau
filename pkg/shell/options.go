package shell

// Options for shell command
type Options struct {
	WorkingDirectory string
	Stdout           []OutputProcessor
	Stderr           []OutputProcessor
}

type OutputProcessor interface {
	WriteStdout(line string)
	WriteStderr(line string)
}

func Processors(pros ...OutputProcessor) []OutputProcessor {
	return pros
}
