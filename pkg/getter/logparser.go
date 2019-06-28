package getter

import (
	"regexp"

	"github.com/apex/log"
	"github.com/go-errors/errors"
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
	strLevel := matches[0][1]
	switch strLevel {
	case "TRACE":
		strLevel = "DEBUG"
	case "ERR":
		strLevel = "ERROR"
	case "":
		strLevel = "INFO"
	}

	level, err := log.ParseLevel(strLevel)
	if err != nil {
		return 0, err
	}

	switch level {
	case log.DebugLevel:
		l.Logger.Debug(message)
	case log.InfoLevel:
		l.Logger.Info(message)
	case log.WarnLevel:
		l.Logger.Warn(message)
	case log.ErrorLevel:
		l.Logger.Error(message)
	case log.FatalLevel:
		l.Logger.Fatal(message)
	default:
		return 0, errors.Errorf("invalid log level")
	}

	return len(message), nil
}
