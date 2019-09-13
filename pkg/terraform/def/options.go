package def

import "github.com/avinor/tau/pkg/hooks"

// Options sent to New function when making a new Engine.
type Options struct {
	Runner *hooks.Runner
}
