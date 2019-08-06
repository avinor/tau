package paths

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJoinPaths tests various combinations of paths to see that the function
// Join correctly can create paths that are not invalid
func TestJoinPaths(t *testing.T) {
	tests := []struct {
		paths   []string
		expects string
	}{
		{[]string{WorkingDir, TauPath, "test"}, filepath.Join(WorkingDir, TauPath, "test")},
		{[]string{WorkingDir, TauPath, "."}, filepath.Join(WorkingDir, TauPath, "root")},
		{[]string{WorkingDir, TauPath, ""}, filepath.Join(WorkingDir, TauPath, "root")},
		{[]string{WorkingDir, "", ""}, filepath.Join(WorkingDir, "root", "root")},
		{[]string{WorkingDir, TauPath, "./example"}, filepath.Join(WorkingDir, TauPath, "example")},
		{[]string{WorkingDir, TauPath, "../example"}, filepath.Join(WorkingDir, "example")},
		{[]string{TauPath, "test"}, filepath.Join(WorkingDir, TauPath, "test")},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			assert.Equal(t, test.expects, Join(test.paths...))
		})
	}
}
