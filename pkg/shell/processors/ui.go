package processors

// UI forwards messages to the ui handler. Set writer to the function on ui handler
// it should call, Debug, Info, Error etc depending on the level it should log at
type UI struct {
	writer func(string, ...interface{})
}

// NewUI returns a new UI processor
func NewUI(writer func(string, ...interface{})) *UI {
	return &UI{
		writer: writer,
	}
}

// Write line to ui level defined by writer func
func (u *UI) Write(line string) bool {
	u.writer(line)

	return true
}