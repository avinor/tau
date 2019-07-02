package processors

import (
	"strings"
)

// Buffer collects all lines in buffer that can later be read and processed
type Buffer struct {
	builder strings.Builder
}

// Write line to buffer
func (b *Buffer) Write(line string) bool {
	b.builder.WriteString(line)
	b.builder.WriteByte('\n')

	return true
}

// String returns string from buffer
func (b *Buffer) String() string {
	return b.builder.String()
}
