package ui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLevelParser(t *testing.T) {
	tests := []struct {
		Level       string
		ExpectLevel Level
	}{
		{"TRACE", DebugLevel},
		{"DEBUG", DebugLevel},
		{"INFO", InfoLevel},
		{"WARN", WarnLevel},
		{"ERR", ErrorLevel},
		{"ERROR", ErrorLevel},
		{"SOMETHING", InfoLevel},
		{"", InfoLevel},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			level := ParseLevel(test.Level)

			assert.Equal(t, test.ExpectLevel, level)
		})
	}
}
