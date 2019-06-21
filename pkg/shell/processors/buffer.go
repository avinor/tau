package processors

import (
	"strings"
)

type Buffer struct {
	stdoutBuilder strings.Builder
	stderrBuilder strings.Builder
}

func (b *Buffer) WriteStdout(line string) {
	b.stdoutBuilder.WriteString(line)
	b.stdoutBuilder.WriteByte('\n')
}

func (b *Buffer) WriteStderr(line string) {
	b.stderrBuilder.WriteString(line)
	b.stdoutBuilder.WriteByte('\n')
}

func (b *Buffer) Stdout() string {
	return b.stdoutBuilder.String()
}

func (b *Buffer) Stderr() string {
	return b.stderrBuilder.String()
}
