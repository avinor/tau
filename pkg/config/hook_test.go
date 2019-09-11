package config

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	hookTest1 = `
		hook "name" {
			command = "command"
			trigger_on = "prepare"
		}
	`

	hookTest2 = `
		hook "name" {
			command = "overwrite"
			args = ["arg1", "arg2"]
		}
	`

	hookTest3 = `
		hook "name" {
			args = ["arg3"]
		}
	`

	hookTest4 = `
		hook "name" {
			command = "command"
			trigger_on = "invalid"
		}
	`

	hookTest5 = `
		hook "name" {
			command = ""
			trigger_on = "prepare"
		}
	`

	hookTest6 = `
		hook "name" {
			command = "test"
			script = "path"
			trigger_on = "prepare"
		}
	`

	hookTest7 = `
		hook "name" {
			script = "path"
			trigger_on = "prepare"
		}
	`
)

var (
	hookFile1 = fileFromString("hook1", hookTest1)
	hookFile2 = fileFromString("hook2", hookTest2)
	hookFile3 = fileFromString("hook3", hookTest3)
	hookFile4 = fileFromString("hook4", hookTest4)
	hookFile5 = fileFromString("hook5", hookTest5)
	hookFile6 = fileFromString("hook6", hookTest6)
	hookFile7 = fileFromString("hook7", hookTest7)
)

func TestHookMerge(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected []*Hook
		Error    error
	}{
		{
			[]*File{hookFile1},
			[]*Hook{
				{
					Type:      "name",
					Command:   stringPointer("command"),
					TriggerOn: stringPointer("prepare"),
				},
			},
			nil,
		},
		{
			[]*File{hookFile1, hookFile2},
			[]*Hook{
				{
					Type:      "name",
					Command:   stringPointer("overwrite"),
					TriggerOn: stringPointer("prepare"),
					Arguments: &[]string{"arg1", "arg2"},
				},
			},
			nil,
		},
		{
			[]*File{hookFile1, hookFile2, hookFile3},
			[]*Hook{
				{
					Type:      "name",
					Command:   stringPointer("overwrite"),
					TriggerOn: stringPointer("prepare"),
					Arguments: &[]string{"arg1", "arg2", "arg3"},
				},
			},
			nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeHooks(config, getConfigFromFiles(t, test.Files))

			if test.Error != nil {
				assert.EqualError(t, err, test.Error.Error())
				return
			}

			actual := map[string]string{}
			for _, dep := range config.Dependencies {
				actual[dep.Name] = dep.Source
			}

			assert.ElementsMatch(t, test.Expected, config.Hooks)
		})
	}
}

func TestHookValidation(t *testing.T) {
	tests := []struct {
		Files    []*File
		Expected map[string]ValidationResult
	}{
		{
			[]*File{hookFile1},
			map[string]ValidationResult{
				"name": {Result: true, Error: nil},
			},
		},
		{
			[]*File{hookFile7},
			map[string]ValidationResult{
				"name": {Result: true, Error: nil},
			},
		},
		{
			[]*File{hookFile4},
			map[string]ValidationResult{
				"name": {Result: false, Error: triggerOnValueIncorrect},
			},
		},
		{
			[]*File{hookFile5},
			map[string]ValidationResult{
				"name": {Result: false, Error: scriptOrCommandIsRequired},
			},
		},
		{
			[]*File{hookFile6},
			map[string]ValidationResult{
				"name": {Result: false, Error: scriptAndCommandBothDefined},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			config := &Config{}
			err := mergeHooks(config, getConfigFromFiles(t, test.Files))
			assert.NoError(t, err)

			for _, hook := range config.Hooks {
				result, err := hook.Validate()
				expect := test.Expected[hook.Type]

				assert.Equal(t, expect.Result, result)
				if expect.Error != nil {
					assert.EqualError(t, err, expect.Error.Error())
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}
