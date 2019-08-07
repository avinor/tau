package ui

import "os"

var (
	// handler is a singleton handler, override with SetHandler function
	handler Handler = &CliHandler{
		Reader:      os.Stdin,
		Writer:      os.Stderr,
		ErrorWriter: os.Stderr,
	}

	// level is the current output level, defaults to Info
	level = InfoLevel
)

// SetHandler assigns a new handler overriding current (default is cli handler)
func SetHandler(hnd Handler) {
	handler = hnd
}

// SetLevel changes the current output level
func SetLevel(newLevel Level) {
	level = newLevel
}

// Ask for user input
func Ask(query string) (string, error) {
	return handler.Ask(query)
}

// AskSecret will ask for a secret input, not showing what is typed
func AskSecret(query string) (string, error) {
	return handler.AskSecret(query)
}

// Debug will print a debug message, if debugging is activated
func Debug(msg string, args ...interface{}) {
	if level > DebugLevel {
		return
	}

	handler.Debug(msg, args...)
}

// Info will print an information message. Should not be formatted any special way
func Info(msg string, args ...interface{}) {
	if level > InfoLevel {
		return
	}

	handler.Info(msg, args...)
}

// Warn prints a warning
func Warn(msg string, args ...interface{}) {
	if level > WarnLevel {
		return
	}

	handler.Warn(msg, args...)
}

// Error prints an error message
func Error(msg string, args ...interface{}) {
	if level > ErrorLevel {
		return
	}

	handler.Error(msg, args...)
}

// Fatal prints a fatal message and should always exit!
func Fatal(msg string, args ...interface{}) {
	if level > FatalLevel {
		return
	}

	handler.Fatal(msg, args...)

	// Make sure it always exists.. if handler have not implemented that
	// it will make sure to exit here
	os.Exit(1)
}

// Header prints a header, can be bold etc. Implementation can decide how a header
// should be made
func Header(msg string) {
	handler.Header(msg)
}

// Separator between elements
func Separator(title string) {
	handler.Separator(title)
}

// NewLine adds a new line
func NewLine() {
	handler.NewLine()
}
