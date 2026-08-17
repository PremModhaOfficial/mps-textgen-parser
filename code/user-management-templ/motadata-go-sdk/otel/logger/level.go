package logger

import (
	"fmt"
	"strings"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.uber.org/zap/zapcore"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Level represents log severity levels.
// The levels follow standard syslog severity ordering where lower values
// indicate more verbose/less severe messages.
//
// Level ordering (from most to least verbose):
//   - DebugLevel (-1): Detailed debugging information
//   - InfoLevel  (0):  General operational information
//   - WarnLevel  (1):  Warning conditions that should be reviewed
//   - ErrorLevel (2):  Error conditions that need attention
//   - FatalLevel (3):  Fatal errors that terminate the application
//
// The underlying int8 values are compatible with zap's level system,
// enabling efficient level checking and dynamic level changes.
type Level int8

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	// DebugLevel is for detailed debugging information.
	// Use for development-time tracing, variable dumps, and step-by-step execution.
	// Should typically be disabled in production due to high volume.
	DebugLevel Level = iota - 1

	// InfoLevel is for general operational information.
	// Use for application lifecycle events, request processing, and normal operations.
	// This is the recommended default level for production.
	InfoLevel

	// WarnLevel is for warning conditions that should be reviewed.
	// Use for deprecated features, approaching limits, and recoverable issues.
	// These don't require immediate action but indicate potential problems.
	WarnLevel

	// ErrorLevel is for error conditions that need attention.
	// Use for failed operations, unexpected conditions, and issues requiring investigation.
	// These indicate problems that affected functionality but didn't crash the service.
	ErrorLevel

	// FatalLevel is for fatal errors that terminate the application.
	// Logging at this level will call os.Exit(1) after writing the message.
	// Use only for unrecoverable errors where continued operation is impossible.
	FatalLevel
)

/* ---------------------------------------- Level Methods -------------------------------------------------- */

// String returns the human-readable level name.
// This is used for logging output and configuration display.
//
// Returns:
//   - string: The level name in lowercase (e.g., "debug", "info", "warn", "error", "fatal")
//
// Example:
//
//	level := InfoLevel
//	fmt.Println(level.String())  // Output: "info"
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
		// Return "unknown" for any unrecognized level
		// This shouldn't happen with proper usage but provides safety
		return "unknown"
	}
}

// zapLevel converts the SDK Level to zapcore.Level.
// This internal method bridges our Level type to zap's internal representation,
// enabling integration with zap's level-checking and filtering mechanisms.
//
// Returns:
//   - zapcore.Level: The corresponding zap level
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
		// Default to InfoLevel for safety
		// This ensures unknown levels don't accidentally enable verbose logging
		return zapcore.InfoLevel
	}
}

/* ---------------------------------------- Level Parsing Functions -------------------------------------------------- */

// ParseLevel parses a level string into a Level constant.
// This function is case-insensitive and trims whitespace, making it
// suitable for parsing configuration values from various sources.
//
// Supported values:
//   - "debug"          -> DebugLevel
//   - "info"           -> InfoLevel
//   - "warn", "warning" -> WarnLevel
//   - "error"          -> ErrorLevel
//   - "fatal"          -> FatalLevel
//
// Parameters:
//   - levelString: The level string to parse (case-insensitive)
//
// Returns:
//   - Level: The parsed level (defaults to InfoLevel on error)
//   - error: Non-nil if the level string is not recognized
//
// Example:
//
//	level, err := ParseLevel("DEBUG")    // Returns DebugLevel, nil
//	level, err := ParseLevel("warning")  // Returns WarnLevel, nil
//	level, err := ParseLevel("invalid")  // Returns InfoLevel, error
func ParseLevel(levelString string) (Level, error) {

	// Normalize input: trim whitespace and convert to lowercase
	// This allows flexible configuration like " INFO " or "INFO"
	switch strings.ToLower(strings.TrimSpace(levelString)) {

	case "debug":
		return DebugLevel, nil

	case "info":
		return InfoLevel, nil

	case "warn", "warning":
		// Support both "warn" (common) and "warning" (syslog standard)
		return WarnLevel, nil

	case "error":
		return ErrorLevel, nil

	case "fatal":
		return FatalLevel, nil

	default:
		// Return InfoLevel as default with error indicating invalid input
		// This ensures the caller can still function while being notified
		return InfoLevel, fmt.Errorf("%w: %s", utils.ErrOTELInvalidLogLevel, levelString)
	}
}
