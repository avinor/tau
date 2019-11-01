package loader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldDeleteFile(t *testing.T) {
	tests := []struct {
		Name    string
		Delete  bool
		Altered string
	}{
		{"test.hcl", false, "test.hcl"},
		{"/tmp/test.hcl", false, "/tmp/test.hcl"},
		{"/tmp/delete_test.hcl", true, "/tmp/test.hcl"},
		{"delete_test.hcl", true, "test.hcl"},
		{"/tmp/destroy_test.hcl", true, "/tmp/test.hcl"},
		{"/tmp/DELETE_test.hcl", true, "/tmp/test.hcl"},
		{"/tmp/Destroy_test.hcl", true, "/tmp/test.hcl"},
		{"/tmp/DELETEtest.hcl", true, "/tmp/test.hcl"},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			del, altered := shouldDeleteFile(test.Name)

			assert.Equal(t, test.Delete, del)
			if del {
				assert.Equal(t, test.Altered, altered)
			}
		})
	}
}
