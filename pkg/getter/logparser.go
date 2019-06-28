package getter

import (
	"regexp"

	"github.com/apex/log"
)

// LogParser parse the registry log messages. Since the registry client uses standard log
// package it is difficult to intercept the messages. This LogParser will register a handler
// for standard log and parse log level and message and forward to apex/log
type LogParser struct {
	Logger log.Interface
}

const (
	logPattern = "\\[(\\w+)\\] (.*)"
)

var (
	logRegex = regexp.MustCompile(logPattern)
)

// Write implements the io.Writer interface, needed for registering a handler for log
func (l LogParser) Write(p []byte) (n int, err error) {
	str := string(p)

	matches := logRegex.FindAllStringSubmatch(str, -1)

	if len(matches) < 1 {
		return 0, nil
	}

	message := matches[0][2]

	level, err := log.ParseLevel(matches[0][1])
	if err != nil {
		l.Logger.Warn(message)
		return len(message), nil
	}

	switch level {
	case log.DebugLevel:
		l.Logger.Debug(message)
	case log.InfoLevel:
		l.Logger.Info(message)
	default:
		l.Logger.Warn(message)
	}

	return len(message), nil
}
