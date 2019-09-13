package def

import (
	"github.com/avinor/tau/pkg/config"
)

type ExecutorCreator interface {
	CanCreate(hook *config.Hook) bool
	Create(hook *config.Hook) Executor
}
