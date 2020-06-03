package hooks

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/avinor/tau/pkg/config"
	"github.com/avinor/tau/pkg/helpers/strings"
)

func TestShouldRun(t *testing.T) {
	tests := []struct {
		hook    *config.Hook
		event   string
		command string
		expects bool
	}{
		{&config.Hook{TriggerOn: strings.ToPointer("prepare")}, "prepare", "init", true},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare")}, "prepare", "init", true},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare")}, "finish", "init", false},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare:init")}, "prepare", "init", true},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare:plan")}, "prepare", "init", false},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare:init")}, "prepare", "INIT", true},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare")}, "Prepare", "init", true},
		{&config.Hook{TriggerOn: strings.ToPointer("prepare:init")}, "prepare", "plan", false},
		{&config.Hook{TriggerOn: strings.ToPointer("finish")}, "prepare", "init", false},
		{&config.Hook{TriggerOn: strings.ToPointer("finish:init")}, "finish", "init", true},
		{&config.Hook{TriggerOn: strings.ToPointer("finish:init")}, "finish", "INIT", true},
		{&config.Hook{TriggerOn: strings.ToPointer("finish:init,plan")}, "finish", "plan", true},
		{&config.Hook{TriggerOn: strings.ToPointer("finish:init,plan")}, "finish", "init", true},
	}

	runner := Runner{}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			expects := runner.ShouldRun(test.hook, test.event, test.command)

			assert.Equal(t, test.expects, expects)
		})
	}
}
