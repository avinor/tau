package ui

// Handler for ui commands. Can ask for input, output messages and enforce a
// specific format.
type Handler interface {
	// Ask for input from user
	Ask(query string) (string, error)

	// AskSecret will ask for a secret input, not showing what is typed
	AskSecret(query string) (string, error)

	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})

	// Output writes to stdout so it can be piped to subsequent commands
	Output(msg string, args ...interface{})

	// Header prints a header, can be bold etc. Implementation can decide how a header
	// should be made. Should always be printed as info message
	Header(msg string)

	// Separator between elements. Should always be printed as info message
	Separator(title string)

	// NewLine adds a new line if necessary
	NewLine()
}
