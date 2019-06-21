package getter

import (
	"regexp"
	"github.com/apex/log"
)

// LogParser parse the registry log messages
type LogParser struct {}

const (
	logPattern = "\\[(\\w+)\\] (.*)"
)

var (
	logRegex = regexp.MustCompile(logPattern)
)

func (l LogParser) Write(p []byte) (n int, err error) {
	str := string(p)

	matches := logRegex.FindAllStringSubmatch(str, -1)

	if len(matches) < 1 {
		return 0, nil
	}

	message := matches[0][2]

	level, err := log.ParseLevel(matches[0][1])
	if err != nil {
		log.Warn(message)
		return len(message), nil
	}

	switch level {
	case log.DebugLevel:
		log.Debug(message)
	case log.InfoLevel:
		log.Info(message)
	default:
		log.Warn(message)
	}

	return len(message), nil
}