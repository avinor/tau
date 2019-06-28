package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewSource tries to create different sources and see that they are created correctly.
// It will only validate some fields, not testing content for instance as that is not needed.
func TestNewSource(t *testing.T) {
	tests := []struct{
		File string
		Content string
		Expects Source
		Err bool
	}{
		{
			"/tmp/test.hcl",
			"",
			Source{
				Name: "test.hcl",
				File: "/tmp/test.hcl",
				Env: map[string]string{},
				Dependencies: map[string]*Source{},
			},
			false,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			source, err := NewSource(test.File, []byte(test.Content))	

			if test.Err {
				assert.Error(t, err)
			}

			assert.Equal(t, test.Expects.Name, source.Name)
			assert.Equal(t, test.Expects.File, source.File)
			assert.Equal(t, test.Expects.Env, source.Env)
			assert.Equal(t, test.Expects.Dependencies, source.Dependencies)
		})
	}
}
