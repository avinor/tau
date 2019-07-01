package hooks

type Event int

const (
	Before Event = iota
	After
)

type Hook struct {
	Event   Event
	Command string
}

func Register(hooks ...Hook) error {
	return nil
}

func Trigger(event Event) error {
	return nil
}
