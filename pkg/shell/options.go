package shell

// Options for shell command
type Options struct {
	WorkingDirectory string
	Stdout           []OutputProcessor
	Stderr           []OutputProcessor
}

type OutputProcessor interface {
	Stdout(line string)
	Stderr(line string)
}

func NewOptions() *Options {
	return &Options{
		Stdout: []OutputProcessor{},
		Stderr: []OutputProcessor{},
	}
}

func (o *Options) WithWorkingDirectory(wd string) *Options {
	o.WorkingDirectory = wd
	return o
}

func (o *Options) WithStdout(processor OutputProcessor) *Options {
	o.Stdout = append(o.Stdout, processor)
	return o
}

func (o *Options) WithStderr(processor OutputProcessor) *Options {
	o.Stderr = append(o.Stderr, processor)
	return o
}

func (o *Options) ProcessStdout(line string) {
	for _, out := range o.Stdout {
		out.Stdout(line)
	}
}

func (o *Options) ProcessStderr(line string) {
	for _, out := range o.Stderr {
		out.Stderr(line)
	}
}