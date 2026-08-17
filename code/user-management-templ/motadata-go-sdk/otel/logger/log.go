// Package logger provides a high-performance structured logging system with OpenTelemetry integration.
//
// This package wraps uber-go/zap for logging with the OpenTelemetry zap bridge for exporting
// logs to OTEL collectors. It provides:
//   - Structured logging with strongly-typed fields
//   - Context-aware logging with automatic field extraction (trace_id, tenant_id, etc.)
//   - Per-module log levels for fine-grained control
//   - OTEL integration for centralized log collection
//   - File output with rotation (via lumberjack)
//   - Zero-allocation field pooling for high-throughput scenarios
//
// Basic Usage:
//
//	// Initialize the global logger
//	logger.Init(logger.DefaultConfig())
//	defer logger.Close()
//
//	// Log with context
//	ctx := logger.WithTenantID(context.Background(), "tenant-123")
//	logger.Info(ctx, "processing request", logger.String("endpoint", "/api/users"))
//
// Module Usage:
//
//	// Create a module logger with independent level
//	dbLogger := logger.Module("database")
//	dbLogger.SetLevel(logger.DebugLevel)
//	dbLogger.Debug(ctx, "query executed")
//
// Note: This package cannot use dot imports for core/types because the zap field
// constructors (String, Int, etc.) conflict with SDK type names.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Config is an alias for config.LoggerConfig for backward compatibility.
// This allows existing code using logger.Config to continue working.
type Config = config.LoggerConfig

// Field is a log field (re-export zap.Field).
// Fields are key-value pairs that provide structured context to log entries.
// Use the provided constructors (String, Int, etc.) to create fields.
type Field = zap.Field

// Logger is the main logger type providing structured logging with OTEL integration.
// It wraps zap.Logger and adds:
//   - Context-aware field extraction
//   - OTEL log provider management
//   - Graceful shutdown with timeout
//   - Dynamic level changes
//
// Thread Safety: Logger is safe for concurrent use. Multiple goroutines may
// call methods on the same Logger simultaneously.
type Logger struct {
	// zap is the underlying zap logger that performs the actual logging
	zap *zap.Logger

	// level allows dynamic level changes at runtime
	// Using atomic level avoids locks on the hot path
	level zap.AtomicLevel

	// provider is the OTEL logger provider (nil if OTEL is disabled)
	// This is used for graceful shutdown of the OTEL pipeline
	provider *sdklog.LoggerProvider

	// closers holds resources that need to be closed (e.g., file writers)
	closers []io.Closer

	// shutdownTimeout is the maximum time to wait for graceful shutdown
	shutdownTimeout time.Duration

	// closed tracks whether the logger has been closed
	// Using atomic bool for thread-safe idempotent close
	closed atomic.Bool
}

/* ---------------------------------------- Field Constructors -------------------------------------------------- */

// Field constructors (re-exported from zap for convenience).
// These provide strongly-typed field creation for structured logging.
// Using these constructors ensures type safety and optimal performance.
//
// Usage:
//
//	logger.Info(ctx, "user action",
//	    String("user_id", "123"),
//	    Int("count", 42),
//	    Duration("elapsed", time.Second),
//	    Err(someError),
//	)
var (
	// String constructs a field with the given key and string value
	String = zap.String

	// Strings constructs a field with the given key and string slice value
	Strings = zap.Strings

	// Int constructs a field with the given key and int value
	Int = zap.Int

	// Int64 constructs a field with the given key and int64 value
	Int64 = zap.Int64

	// Float64 constructs a field with the given key and float64 value
	Float64 = zap.Float64

	// Bool constructs a field with the given key and bool value
	Bool = zap.Bool

	// Time constructs a field with the given key and time.Time value
	Time = zap.Time

	// Duration constructs a field with the given key and time.Duration value
	Duration = zap.Duration

	// Err is shorthand for Error("error", err)
	Err = zap.Error

	// Any constructs a field with the given key and arbitrary value
	// Use this sparingly as it requires reflection
	Any = zap.Any

	// Binary constructs a field with the given key and []byte value
	Binary = zap.Binary
)

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global logger state.
// The global logger provides a convenient default for applications that don't need
// multiple logger instances. It's protected by a RWMutex for thread-safe access.
var (
	// globalMu protects access to globalLogger
	globalMu sync.RWMutex

	// globalLogger is the default logger instance
	globalLogger *Logger

	// explicitlyInited tracks whether Init() was called (vs. auto-initialization)
	explicitlyInited atomic.Bool
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// DefaultConfig returns the default logger configuration with service defaults.
// This provides sensible defaults for development and can be customized as needed.
//
// Default values:
//   - Level: "info"
//   - ConsoleEnabled: true
//   - ConsoleFormat: "json"
//   - ServiceName: "app"
//   - ServiceVersion: "0.0.0"
//   - Environment: "development"
//
// Returns:
//   - Config: A logger configuration with defaults applied
func DefaultConfig() Config {

	// Start with the base default config from the config package
	defaultConfig := config.DefaultLoggerConfig()

	// Apply service identification defaults
	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

// New creates a new Logger instance with the given configuration.
// This is the primary constructor for creating loggers.
//
// The function configures:
//   - Log level parsing and atomic level creation
//   - Encoder configuration (JSON or console format)
//   - Console output core (if enabled)
//   - OTEL output core via otelzap bridge (if enabled)
//   - File output core with rotation (if enabled)
//   - Caller information (if enabled)
//   - Service identification fields
//   - Module-level configuration
//
// Parameters:
//   - config: The logger configuration
//
// Returns:
//   - *Logger: The configured logger instance
//   - error: Non-nil if configuration is invalid
func New(config Config) (*Logger, error) {

	// Parse the configured log level
	level, err := ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	// Create atomic level for dynamic level changes without locks
	atomicLevel := zap.NewAtomicLevelAt(level.zapLevel())
	encoderConfig := newEncoderConfig()

	// Build cores for all enabled outputs
	cores, closers, provider, err := buildCores(config, encoderConfig, atomicLevel)
	if err != nil {
		return nil, err
	}

	// Create the zap logger with all cores
	zapLogger := buildZapLogger(config, cores)

	// Initialize module levels from configuration
	initModuleLevels(level)
	applyModuleLevels(config.ModuleLevels)

	// Configure shutdown timeout with default fallback
	timeout := config.ShutdownTimeout
	if timeout <= 0 {
		timeout = utils.OTELDefaultShutdownTimeout
	}

	return &Logger{
		zap:             zapLogger,
		level:           atomicLevel,
		provider:        provider,
		closers:         closers,
		shutdownTimeout: timeout,
	}, nil
}

// newEncoderConfig creates the standard encoder configuration for structured logging.
func newEncoderConfig() zapcore.EncoderConfig {

	return zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// buildCores creates all enabled output cores (console, OTEL, file).
func buildCores(config Config, encoderConfig zapcore.EncoderConfig, atomicLevel zap.AtomicLevel) ([]zapcore.Core, []io.Closer, *sdklog.LoggerProvider, error) {

	var cores []zapcore.Core
	var closers []io.Closer
	var provider *sdklog.LoggerProvider

	// Add console core if enabled
	if config.ConsoleEnabled {
		cores = append(cores, newConsoleCore(config.ConsoleFormat, encoderConfig, atomicLevel))
	}

	// Add OTEL core if enabled
	if config.OTELEnabled {
		otelCore, otelProvider, err := newOTELCore(config)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("otel provider: %w", err)
		}
		cores = append(cores, otelCore)
		provider = otelProvider
	}

	// Add file core if enabled
	if config.FileEnabled {
		fileCore, fileCloser := newFileCore(config, encoderConfig, atomicLevel)
		cores = append(cores, fileCore)
		closers = append(closers, fileCloser)
	}

	// Use noop core if no outputs configured
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewNopCore())
	}

	return cores, closers, provider, nil
}

// newConsoleCore creates a console output core with the specified format.
func newConsoleCore(format string, encoderConfig zapcore.EncoderConfig, level zap.AtomicLevel) zapcore.Core {

	var encoder zapcore.Encoder
	if format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
}

// newOTELCore creates an OTEL output core with the configured provider.
func newOTELCore(config Config) (zapcore.Core, *sdklog.LoggerProvider, error) {

	provider, err := newOTELProvider(config)
	if err != nil {
		return nil, nil, err
	}

	core := otelzap.NewCore(config.ServiceName, otelzap.WithLoggerProvider(provider))
	return core, provider, nil
}

// newFileCore creates a file output core with rotation using lumberjack.
func newFileCore(config Config, encoderConfig zapcore.EncoderConfig, level zap.AtomicLevel) (zapcore.Core, io.Closer) {

	fileWriter := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.FileMaxSizeMB,
		MaxBackups: config.FileMaxBackups,
		MaxAge:     config.FileMaxAgeDays,
		Compress:   config.FileCompress,
		LocalTime:  true,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(fileWriter),
		level,
	)

	return core, fileWriter
}

// buildZapLogger creates the zap logger with service fields and options.
func buildZapLogger(config Config, cores []zapcore.Core) *zap.Logger {

	core := zapcore.NewTee(cores...)

	options := []zap.Option{
		zap.AddStacktrace(zapcore.ErrorLevel),
	}

	if config.AddCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}

	zapLogger := zap.New(core, options...)

	return zapLogger.With(
		zap.String("service", config.ServiceName),
		zap.String("version", config.ServiceVersion),
		zap.String("env", config.Environment),
	)
}

// applyModuleLevels applies module-level overrides from configuration.
func applyModuleLevels(moduleLevelConfig map[string]string) {

	if len(moduleLevelConfig) == 0 {
		return
	}

	moduleLevels := getModuleLevels()
	for moduleName, levelString := range moduleLevelConfig {
		if parsedLevel, parseErr := ParseLevel(levelString); parseErr == nil {
			moduleLevels.SetLevel(moduleName, parsedLevel)
		}
	}
}

/* ---------------------------------------- Logger Instance Methods -------------------------------------------------- */

// Debug logs at debug level.
// Use for detailed debugging information that is too verbose for production.
//
// Parameters:
//   - ctx: Context for automatic field extraction (trace_id, tenant_id, etc.)
//   - message: The log message
//   - fields: Additional structured fields
func (logger *Logger) Debug(ctx context.Context, message string, fields ...Field) {

	// Early bailout if level is disabled to avoid allocation
	// This check is very fast due to the atomic level
	if !logger.level.Enabled(zapcore.DebugLevel) {
		return
	}

	// Extract fields from context using pooled allocation
	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	// Write the log entry
	logger.zap.Debug(message, allFields...)

	// Return pooled slice to pool for reuse
	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Info logs at info level.
// Use for general operational information about application state.
//
// Parameters:
//   - ctx: Context for automatic field extraction
//   - message: The log message
//   - fields: Additional structured fields
func (logger *Logger) Info(ctx context.Context, message string, fields ...Field) {

	if !logger.level.Enabled(zapcore.InfoLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Info(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Warn logs at warn level.
// Use for conditions that should be reviewed but don't require immediate action.
//
// Parameters:
//   - ctx: Context for automatic field extraction
//   - message: The log message
//   - fields: Additional structured fields
func (logger *Logger) Warn(ctx context.Context, message string, fields ...Field) {

	if !logger.level.Enabled(zapcore.WarnLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Warn(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Error logs at error level.
// Use for error conditions that need attention but didn't crash the service.
//
// Parameters:
//   - ctx: Context for automatic field extraction
//   - message: The log message
//   - fields: Additional structured fields
func (logger *Logger) Error(ctx context.Context, message string, fields ...Field) {

	if !logger.level.Enabled(zapcore.ErrorLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Error(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Fatal logs at fatal level then exits.
// Use only for unrecoverable errors where continued operation is impossible.
// WARNING: This calls os.Exit(1) after logging.
//
// Parameters:
//   - ctx: Context for automatic field extraction
//   - message: The log message
//   - fields: Additional structured fields
func (logger *Logger) Fatal(ctx context.Context, message string, fields ...Field) {

	// Fatal always logs and exits - no level check needed
	logger.zap.Fatal(message, logger.fields(ctx, fields)...)
}

// With returns a logger with preset fields.
// These fields are included in all subsequent log entries from the returned logger.
//
// Parameters:
//   - fields: Fields to preset on all future log entries
//
// Returns:
//   - *Logger: A new logger with preset fields
//
// Example:
//
//	reqLogger := logger.With(String("request_id", reqID))
//	reqLogger.Info(ctx, "processing")  // Includes request_id
func (logger *Logger) With(fields ...Field) *Logger {

	return &Logger{
		zap:             logger.zap.With(fields...),
		level:           logger.level,
		provider:        logger.provider,
		closers:         logger.closers,
		shutdownTimeout: logger.shutdownTimeout,
	}
}

// Named returns a named sub-logger.
// The name is appended to the logger's name, creating a hierarchy.
//
// Parameters:
//   - name: The sub-logger name
//
// Returns:
//   - *Logger: A named sub-logger
//
// Example:
//
//	dbLogger := logger.Named("database")
//	dbLogger.Info(ctx, "connected")  // Logger name: "parent.database"
func (logger *Logger) Named(name string) *Logger {

	return &Logger{
		zap:             logger.zap.Named(name),
		level:           logger.level,
		provider:        logger.provider,
		closers:         logger.closers,
		shutdownTimeout: logger.shutdownTimeout,
	}
}

// SetLevel changes log level dynamically.
// This takes effect immediately for all goroutines using this logger.
//
// Parameters:
//   - level: The new log level
func (logger *Logger) SetLevel(level Level) {

	logger.level.SetLevel(level.zapLevel())
}

// Level returns current log level.
//
// Returns:
//   - Level: The current log level
func (logger *Logger) Level() Level {

	return Level(logger.level.Level())
}

// Sync flushes buffered logs.
// Call this before exiting to ensure all logs are written.
//
// Returns:
//   - error: Non-nil if flushing fails
func (logger *Logger) Sync() error {

	return logger.zap.Sync()
}

// Close shuts down the logger gracefully.
// This method is idempotent - multiple calls are safe.
//
// It performs:
//  1. Flush buffered logs
//  2. Shutdown OTEL provider with timeout
//  3. Close file writers
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func (logger *Logger) Close() error {

	// Ensure Close is only executed once using atomic compare-and-swap
	if !logger.closed.CompareAndSwap(false, true) {
		return nil
	}

	// Collect all shutdown errors to report them together
	var errs []error

	// Flush any buffered logs (collect error but continue shutdown)
	if syncErr := logger.Sync(); syncErr != nil {
		errs = append(errs, fmt.Errorf("sync: %w", syncErr))
	}

	// Configure shutdown timeout with fallback
	timeout := logger.shutdownTimeout
	if timeout <= 0 {
		timeout = utils.OTELDefaultShutdownTimeout
	}

	// Shutdown OTEL provider with timeout
	if logger.provider != nil {

		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		if providerErr := logger.provider.Shutdown(ctx); providerErr != nil {
			errs = append(errs, fmt.Errorf("otel provider: %w", providerErr))
		}

		cancel()
	}

	// Close file writers and other closers
	for _, closer := range logger.closers {

		if closeErr := closer.Close(); closeErr != nil {
			errs = append(errs, fmt.Errorf("closer: %w", closeErr))
		}
	}

	// Return aggregated errors if any occurred
	if len(errs) == 1 {
		return errs[0]
	}

	if len(errs) > 1 {
		return fmt.Errorf("%w: %v", utils.ErrOTELShutdownFailed, errs)
	}

	return nil
}

/* ---------------------------------------- Field Helper Functions -------------------------------------------------- */

// fieldsPooled combines context and provided fields using pooled allocation.
// This is the high-performance version used by logging methods.
//
// Parameters:
//   - ctx: Context to extract fields from
//   - fields: Additional fields to include
//
// Returns:
//   - []Field: Combined fields (may be the input fields or a new slice)
//   - *FieldSliceWrapper: Pool wrapper to return (nil if no pooling used)
func (logger *Logger) fieldsPooled(ctx context.Context, fields []Field) ([]Field, *FieldSliceWrapper) {

	// Extract context fields using the pool
	wrapper := extractFieldsPooled(ctx)

	// Fast path: no context fields
	if wrapper == nil {
		return fields, nil
	}

	// Fast path: no additional fields
	if len(fields) == 0 {
		return wrapper.Fields, wrapper
	}

	// Combine context fields and provided fields
	allFields := make([]Field, 0, len(wrapper.Fields)+len(fields))
	allFields = append(allFields, wrapper.Fields...)
	allFields = append(allFields, fields...)

	// Return wrapper to pool since we copied its contents
	putFieldSlice(wrapper)

	return allFields, nil
}

// fields combines context and provided fields (allocating version).
// This is used by Fatal where pooling isn't needed.
//
// Parameters:
//   - ctx: Context to extract fields from
//   - fields: Additional fields to include
//
// Returns:
//   - []Field: Combined fields
func (logger *Logger) fields(ctx context.Context, fields []Field) []Field {

	contextFields := extractFields(ctx)

	if len(contextFields) == 0 {
		return fields
	}

	allFields := make([]Field, 0, len(contextFields)+len(fields))
	allFields = append(allFields, contextFields...)
	allFields = append(allFields, fields...)

	return allFields
}

/* ---------------------------------------- Global Logger Functions -------------------------------------------------- */

// Init initializes the global logger with the given configuration.
// This should be called once at application startup.
//
// Parameters:
//   - config: The logger configuration
//
// Returns:
//   - *Logger: The initialized logger
//   - error: Non-nil if initialization fails
func Init(config Config) (*Logger, error) {

	logger, err := New(config)

	if err != nil {
		return nil, err
	}

	// Replace the global logger and mark as explicitly initialized
	ReplaceGlobal(logger)
	explicitlyInited.Store(true)

	return logger, nil
}

// L returns the global logger instance.
// If no global logger has been initialized, a default one is created.
//
// Returns:
//   - *Logger: The global logger
func L() *Logger {

	globalMu.RLock()
	logger := globalLogger
	globalMu.RUnlock()

	// Auto-initialize if not set
	if logger == nil {
		return defaultLogger()
	}

	return logger
}

// IsInitialized returns true if the global logger was explicitly initialized.
// This distinguishes between explicit Init() calls and auto-initialization.
//
// Returns:
//   - bool: True if Init() was called
func IsInitialized() bool {

	return explicitlyInited.Load()
}

// MustInit initializes the global logger and panics on failure.
// Use this when logger initialization is required for application startup.
//
// Parameters:
//   - config: The logger configuration
//
// Returns:
//   - *Logger: The initialized logger
func MustInit(config Config) *Logger {

	logger, err := Init(config)

	if err != nil {
		panic(fmt.Sprintf("logger: failed to initialize: %v", err))
	}

	return logger
}

// ReplaceGlobal replaces the global logger with the provided logger.
// It returns a function that restores the previous logger.
//
// Parameters:
//   - logger: The new global logger
//
// Returns:
//   - func(): Function to restore the previous logger
func ReplaceGlobal(logger *Logger) func() {

	globalMu.Lock()
	previous := globalLogger
	globalLogger = logger
	globalMu.Unlock()

	return func() {
		ReplaceGlobal(previous)
	}
}

// defaultLogger creates a minimal default logger for when global is uninitialized.
// This uses double-check locking to avoid creating multiple default loggers.
func defaultLogger() *Logger {

	globalMu.Lock()
	defer globalMu.Unlock()

	// Double-check after acquiring write lock
	if globalLogger != nil {
		return globalLogger
	}

	// Create minimal logger with console output
	config := DefaultConfig()
	config.CallerSkip = 3 // Account for extra call stack frames

	logger, err := New(config)

	if err != nil {
		// Fallback to noop logger if creation fails
		nopLogger := zap.NewNop()

		logger = &Logger{
			zap:   nopLogger,
			level: zap.NewAtomicLevelAt(zapcore.InfoLevel),
		}
	}

	globalLogger = logger

	return globalLogger
}

/* ---------------------------------------- Package-level Logging Functions -------------------------------------------------- */

// Debug logs at debug level using the global logger.
func Debug(ctx context.Context, msg string, fields ...Field) {
	L().Debug(ctx, msg, fields...)
}

// Info logs at info level using the global logger.
func Info(ctx context.Context, msg string, fields ...Field) {
	L().Info(ctx, msg, fields...)
}

// Warn logs at warn level using the global logger.
func Warn(ctx context.Context, msg string, fields ...Field) {
	L().Warn(ctx, msg, fields...)
}

// Error logs at error level using the global logger.
func Error(ctx context.Context, msg string, fields ...Field) {
	L().Error(ctx, msg, fields...)
}

// Fatal logs at fatal level using the global logger.
func Fatal(ctx context.Context, msg string, fields ...Field) {
	L().Fatal(ctx, msg, fields...)
}

// With returns a new logger with preset fields.
func With(fields ...Field) *Logger {
	return L().With(fields...)
}

// Named returns a named sub-logger.
func Named(name string) *Logger {
	return L().Named(name)
}

// SetGlobalLevel changes the log level of the global logger.
func SetGlobalLevel(level Level) {
	L().SetLevel(level)
}

// Sync flushes the global logger's buffered logs.
func Sync() error {

	globalMu.RLock()
	logger := globalLogger
	globalMu.RUnlock()

	if logger != nil {
		return logger.Sync()
	}

	return nil
}

// Close shuts down the global logger.
func Close() error {

	globalMu.Lock()
	logger := globalLogger
	globalLogger = nil
	globalMu.Unlock()

	explicitlyInited.Store(false)

	if logger != nil {
		return logger.Close()
	}

	return nil
}

/* ---------------------------------------- Config Loading Functions -------------------------------------------------- */

// LoadConfig loads logger configuration from a YAML file.
//
// Parameters:
//   - configPath: Path to the configuration file
//
// Returns:
//   - Config: The loaded configuration
//   - error: Non-nil if loading fails
func LoadConfig(configPath string) (Config, error) {

	return config.LoadLoggerConfig(configPath)
}

// ConfigFromEnv loads logger configuration from environment variables.
//
// Returns:
//   - Config: Configuration from environment
func ConfigFromEnv() Config {

	return config.LoadLoggerConfigFromEnv()
}

// InitFromConfig loads config from a file and initializes the global logger.
//
// Parameters:
//   - configPath: Path to the configuration file
//
// Returns:
//   - *Logger: The initialized logger
//   - error: Non-nil if initialization fails
func InitFromConfig(configPath string) (*Logger, error) {

	config, err := LoadConfig(configPath)

	if err != nil {
		return nil, err
	}

	return Init(config)
}

// InitFromEnv initializes the global logger using environment variables.
//
// Returns:
//   - *Logger: The initialized logger
//   - error: Non-nil if initialization fails
func InitFromEnv() (*Logger, error) {

	config := ConfigFromEnv()

	return Init(config)
}

// InitFromUnifiedConfig initializes the global logger from the unified config.
//
// Parameters:
//   - unifiedConfig: The unified application configuration
//
// Returns:
//   - *Logger: The initialized logger
//   - error: Non-nil if initialization fails
func InitFromUnifiedConfig(unifiedConfig config.Config) (*Logger, error) {

	return Init(unifiedConfig.GetLoggerConfig())
}
