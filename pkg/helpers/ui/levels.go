package ui

import (
	"strings"
)

// Level of output
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// ParseLevel returns the level defined by string. If it is not able to parse the
// level string it will return Info as default level
//
// Since there is no trace level that is interpreted as debug
func ParseLevel(level string) Level {
	level = strings.ToLower(level)

	switch level {
	case "fatal":
		return FatalLevel
	case "err":
		fallthrough
	case "error":
		return ErrorLevel
	case "warn":
		return WarnLevel
	case "info":
		return InfoLevel
	case "trace":
		fallthrough
	case "debug":
		return DebugLevel
	default:
		return InfoLevel
	}
}
