package processors

import (
	"github.com/apex/log"
)

// Log processor will write lines to log. If Level is not set it will default to
// Info level.
type Log struct {
	Level log.Level
}

// Write line to log level defined during initialization
func (l *Log) Write(line string) bool {
	switch l.Level {
	case log.FatalLevel:
		log.Fatal(line)
	case log.ErrorLevel:
		log.Error(line)
	case log.WarnLevel:
		log.Warn(line)
	case log.DebugLevel:
		log.Debug(line)
	case log.InfoLevel:
		fallthrough
	default:
		log.Info(line)
	}

	return true
}
