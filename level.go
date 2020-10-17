package loggroup

import (
	"strings"
)

// Level defines the log level key
type Level int

// Defines the log level
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

func (lv Level) String() string {
	switch lv {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	}
	return "unknown"
}

// ParseLevel parse Level from string
func ParseLevel(s string) Level {
	s = strings.ToUpper(s)
	switch s {
	case "DEBUG":
		return LevelDebug
	case "INFO", "INFORMATION":
		return LevelInfo
	case "WARN", "WARNING":
		return LevelWarn
	case "ERR", "ERROR":
		return LevelError
	case "None", "NIL", "NULL":
		return LevelNone
	default:
		return LevelDebug
	}
}
