package def

import (
	"github.com/avinor/tau/pkg/config"
)

// ExecutorCreator defines a creator that can make executors. Should always call
// CanCreate first to check if this specific creator can create an executor for
// hook.
type ExecutorCreator interface {
	CanCreate(hook *config.Hook) bool
	Create(hook *config.Hook) (Executor, error)
}
