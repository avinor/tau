package loader

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestModuleRegexp(t *testing.T) {
	tests := []struct {
		Name  string
		Match bool
	}{
		{"test.hcl", true},
		{"/tmp/test.hcl", true},
		{"test.tau", true},
		{"/tmp/test.tau", true},
		{"test_auto.hcl", false},
		{"test_auto.tau", false},
		{"test.hclr", false},
		{"test.taur", false},
		{"/tmp/hcl", false},
		{"/tmp/hcl.mp3", false},
		{"/tmp/test.HCL", true},
		{"/tmp/test.TAU", true},
		{"/tmp/TEST.TAU", true},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			assert.Equal(t, test.Match, moduleMatchFunc(test.Name), test.Name)
		})
	}
}