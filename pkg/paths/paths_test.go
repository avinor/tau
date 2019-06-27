package paths

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestJoinPaths tests various combinations of paths to see that the function
// join correctly can create paths that are not invalid
func TestJoinPaths(t *testing.T) {
	tests := []struct {
		dir        string
		part       string
		folder     string
		expects    string
		shouldFail bool
	}{
		{WorkingDir, TauPath, "test", filepath.Join(WorkingDir, TauPath, "test"), false},
		{WorkingDir, TauPath, ".", filepath.Join(WorkingDir, TauPath, "root"), false},
		{WorkingDir, TauPath, "", filepath.Join(WorkingDir, TauPath, "root"), false},
		{WorkingDir, "", "", filepath.Join(WorkingDir, TauPath, "root"), true},
		{WorkingDir, TauPath, "./example", filepath.Join(WorkingDir, TauPath, "example"), false},
		{WorkingDir, TauPath, "../example", filepath.Join(WorkingDir, TauPath, "example"), false},
		{WorkingDir, TauPath, "../example/test", filepath.Join(WorkingDir, TauPath, "test"), false},
		{WorkingDir, TauPath, "example", filepath.Join(WorkingDir, TauPath, "example"), false},
		{WorkingDir, TauPath, "example/test", filepath.Join(WorkingDir, TauPath, "test"), false},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			if test.shouldFail {
				assert.Panics(t, func() { join(test.dir, test.part, test.folder, true) })
				return
			}

			assert.Equal(t, test.expects, join(test.dir, test.part, test.folder, false))
		})
	}
}
