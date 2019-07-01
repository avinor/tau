package hooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldRun(t *testing.T) {
	tests := []struct {
		cmd     *Command
		event   string
		command string
		expects bool
	}{
		{&Command{event: "prepare", commands: []string{}}, "prepare", "init", true},
		{&Command{event: "prepare", commands: []string{}}, "finish", "init", false},
		{&Command{event: "prepare", commands: []string{"init"}}, "prepare", "init", true},
		{&Command{event: "prepare", commands: []string{"plan"}}, "prepare", "init", false},
		{&Command{event: "prepare", commands: []string{"init"}}, "prepare", "INIT", true},
		{&Command{event: "prepare", commands: []string{}}, "Prepare", "init", true},
		{&Command{event: "prepare", commands: []string{"init"}}, "prepare", "plan", false},
		{&Command{event: "finish", commands: []string{}}, "prepare", "init", false},
		{&Command{event: "finish", commands: []string{"init"}}, "finish", "init", true},
		{&Command{event: "finish", commands: []string{"init"}}, "finish", "INIT", true},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			expects := test.cmd.ShouldRun(test.event, test.command)

			assert.Equal(t, test.expects, expects)
		})
	}
}

func TestParseOutput(t *testing.T) {
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
			expects := parseOutput(test.output)

			assert.Equal(t, test.expects, expects)
		})
	}
}
