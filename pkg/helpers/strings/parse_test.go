package strings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseVars(t *testing.T) {
	tests := []struct {
		output  string
		expects map[string]string
	}{
		{
			"key=value",
			map[string]string{
				"key": "value",
			},
		},
		{
			`
key=value
test=success
			`,
			map[string]string{
				"key":  "value",
				"test": "success",
			},
		},
		{
			`
			key=value
			test=success
			`,
			map[string]string{
				"key":  "value",
				"test": "success",
			},
		},
		{
			`
			"key"=value
			test=success
			"quoted" = "value"
			not = "longer string with space"
			`,
			map[string]string{
				"key":    "value",
				"test":   "success",
				"quoted": "value",
				"not":    "longer string with space",
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			expects := ParseVars(test.output)

			assert.Equal(t, test.expects, expects)
		})
	}
}
