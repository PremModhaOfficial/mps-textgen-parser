package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Level represents log severity
type Level int8

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	DebugLevel Level = iota - 1

	InfoLevel

	WarnLevel

	ErrorLevel

	FatalLevel
)

/* ---------------------------------------- Level Methods -------------------------------------------------- */

// String returns the level name
func (logLevel Level) String() string {

	switch logLevel {

	case DebugLevel:

		return "debug"

	case InfoLevel:

		return "info"

	case WarnLevel:

		return "warn"

	case ErrorLevel:

		return "error"

	case FatalLevel:

		return "fatal"

	default:

		return "unknown"
	}
}

// zapLevel converts to zapcore.Level
func (logLevel Level) zapLevel() zapcore.Level {
	switch logLevel {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

/* ---------------------------------------- Level Parsing Functions -------------------------------------------------- */

// ParseLevel parses a level string
func ParseLevel(levelString string) (Level, error) {

	switch strings.ToLower(strings.TrimSpace(levelString)) {

	case "debug":

		return DebugLevel, nil

	case "info":

		return InfoLevel, nil

	case "warn", "warning":

		return WarnLevel, nil

	case "error":

		return ErrorLevel, nil

	case "fatal":

		return FatalLevel, nil

	default:

		return InfoLevel, fmt.Errorf("unknown level: %s", levelString)
	}
}
