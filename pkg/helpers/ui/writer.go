package ui

import (
	"regexp"
)

// Writer implements io.Writer
type Writer struct {}

const (
	// logPattern to match "[LEVEL] message"
	logPattern = "\\[(\\w+)\\] (.*)"
)

var (
	// logRegex to check against
	logRegex = regexp.MustCompile(logPattern)
)

// Write tries to parse the incoming string. It will try to read the log level
// at start of string as [LEVEL]. If not found it will just use InfoLevel
func (u Writer) Write(p []byte) (n int, err error) {
	str := string(p)

	matches := logRegex.FindAllStringSubmatch(str, -1)

	if len(matches) < 1 {
		return 0, nil
	}

	message := matches[0][2]
	strLevel := matches[0][1]

	level := ParseLevel(strLevel)
	switch level {
	case DebugLevel:
		Debug(message)
	case InfoLevel:
		Info(message)
	case WarnLevel:
		Warn(message)
	case ErrorLevel:
		Error(message)
	case FatalLevel:
		Fatal(message)
	default:
		Info(message)
	}

	return len(message), nil
}