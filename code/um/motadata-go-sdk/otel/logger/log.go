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
	"go.opentelemetry.io/contrib/bridges/otelzap"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Config is an alias for config.LoggerConfig for backward compatibility
type Config = config.LoggerConfig

// Field is a log field (re-export zap.Field)
type Field = zap.Field

// Logger is the main logger type
type Logger struct {
	zap *zap.Logger

	level zap.AtomicLevel

	provider *sdklog.LoggerProvider

	closers []io.Closer

	shutdownTimeout time.Duration

	closed atomic.Bool
}

/* ---------------------------------------- Field Constructors -------------------------------------------------- */

// Field constructors (re-exported from zap for convenience)
var (
	String   = zap.String
	Strings  = zap.Strings
	Int      = zap.Int
	Int64    = zap.Int64
	Float64  = zap.Float64
	Bool     = zap.Bool
	Time     = zap.Time
	Duration = zap.Duration
	Err      = zap.Error
	Any      = zap.Any
	Binary   = zap.Binary
)

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global logger state
var (
	globalMu sync.RWMutex

	globalLogger *Logger

	explicitlyInited atomic.Bool
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// DefaultConfig returns the default logger configuration with service defaults
func DefaultConfig() Config {
	defaultConfig := config.DefaultLoggerConfig()

	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

// New creates a new Logger instance
func New(config Config) (*Logger, error) {
	level, err := ParseLevel(config.Level)

	if err != nil {
		return nil, err
	}

	atomicLevel := zap.NewAtomicLevelAt(level.zapLevel())

	encoderConfig := zapcore.EncoderConfig{
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

	var cores []zapcore.Core
	var closers []io.Closer
	var provider *sdklog.LoggerProvider

	// Console core setup
	if config.ConsoleEnabled {
		var encoder zapcore.Encoder

		if config.ConsoleFormat == "console" {
			encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		}

		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), atomicLevel))
	}

	// OTEL core setup
	if config.OTELEnabled {
		var err error

		provider, err = newOTELProvider(config)

		if err != nil {
			return nil, fmt.Errorf("otel provider: %w", err)
		}

		otelCore := otelzap.NewCore(config.ServiceName, otelzap.WithLoggerProvider(provider))
		cores = append(cores, otelCore)
	}

	// File core setup
	if config.FileEnabled {
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.FileMaxSizeMB,
			MaxBackups: config.FileMaxBackups,
			MaxAge:     config.FileMaxAgeDays,
			Compress:   config.FileCompress,
			LocalTime:  true,
		}

		closers = append(closers, fileWriter)

		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			atomicLevel,
		)

		cores = append(cores, fileCore)
	}

	// Use noop core if no cores configured (all outputs disabled)
	// This provides minimal overhead for benchmarking and disabled logging scenarios
	if len(cores) == 0 {
		cores = append(cores, zapcore.NewNopCore())
	}

	core := zapcore.NewTee(cores...)

	options := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel)}

	if config.AddCaller {
		options = append(options, zap.AddCaller(), zap.AddCallerSkip(config.CallerSkip))
	}

	zapLogger := zap.New(core, options...)

	zapLogger = zapLogger.With(
		zap.String("service", config.ServiceName),
		zap.String("version", config.ServiceVersion),
		zap.String("env", config.Environment),
	)

	initModuleLevels(level)

	if len(config.ModuleLevels) > 0 {
		moduleLevels := getModuleLevels()

		for moduleName, levelString := range config.ModuleLevels {
			if parsedLevel, err := ParseLevel(levelString); err == nil {
				moduleLevels.SetLevel(moduleName, parsedLevel)
			}
		}
	}

	timeout := config.ShutdownTimeout

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Logger{
		zap:             zapLogger,
		level:           atomicLevel,
		provider:        provider,
		closers:         closers,
		shutdownTimeout: timeout,
	}, nil
}

/* ---------------------------------------- Logger Instance Methods -------------------------------------------------- */

// Debug logs at debug level
func (logger *Logger) Debug(ctx context.Context, message string, fields ...Field) {
	// Early bailout for disabled level to avoid allocations
	if !logger.level.Enabled(zapcore.DebugLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Debug(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Info logs at info level
func (logger *Logger) Info(ctx context.Context, message string, fields ...Field) {
	// Early bailout for disabled level to avoid allocations
	if !logger.level.Enabled(zapcore.InfoLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Info(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Warn logs at warn level
func (logger *Logger) Warn(ctx context.Context, message string, fields ...Field) {
	// Early bailout for disabled level to avoid allocations
	if !logger.level.Enabled(zapcore.WarnLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Warn(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Error logs at error level
func (logger *Logger) Error(ctx context.Context, message string, fields ...Field) {
	// Early bailout for disabled level to avoid allocations
	if !logger.level.Enabled(zapcore.ErrorLevel) {
		return
	}

	allFields, pooledPtr := logger.fieldsPooled(ctx, fields)

	logger.zap.Error(message, allFields...)

	if pooledPtr != nil {
		putFieldSlice(pooledPtr)
	}
}

// Fatal logs at fatal level then exits
func (logger *Logger) Fatal(ctx context.Context, message string, fields ...Field) {
	logger.zap.Fatal(message, logger.fields(ctx, fields)...)
}

// With returns a logger with preset fields
func (logger *Logger) With(fields ...Field) *Logger {
	return &Logger{
		zap:             logger.zap.With(fields...),
		level:           logger.level,
		provider:        logger.provider,
		closers:         logger.closers,
		shutdownTimeout: logger.shutdownTimeout,
	}
}

// Named returns a named sub-logger
func (logger *Logger) Named(name string) *Logger {
	return &Logger{
		zap:             logger.zap.Named(name),
		level:           logger.level,
		provider:        logger.provider,
		closers:         logger.closers,
		shutdownTimeout: logger.shutdownTimeout,
	}
}

// SetLevel changes log level dynamically
func (logger *Logger) SetLevel(level Level) {
	logger.level.SetLevel(level.zapLevel())
}

// Level returns current log level
func (logger *Logger) Level() Level {
	return Level(logger.level.Level())
}

// Sync flushes buffered logs
func (logger *Logger) Sync() error {
	return logger.zap.Sync()
}

// Close shuts down the logger (safe to call multiple times)
func (logger *Logger) Close() error {
	if !logger.closed.CompareAndSwap(false, true) {
		return nil
	}

	_ = logger.Sync()

	timeout := logger.shutdownTimeout

	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	if logger.provider != nil {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := logger.provider.Shutdown(ctx); err != nil {
			return err
		}
	}

	for _, closer := range logger.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	return nil
}

/* ---------------------------------------- Field Helper Functions -------------------------------------------------- */

// fieldsPooled combines context and provided fields using pooled allocation
func (logger *Logger) fieldsPooled(ctx context.Context, fields []Field) ([]Field, *FieldSliceWrapper) {
	wrapper := extractFieldsPooled(ctx)

	if wrapper == nil {
		return fields, nil
	}

	if len(fields) == 0 {
		return wrapper.Fields, wrapper
	}

	allFields := make([]Field, 0, len(wrapper.Fields)+len(fields))
	allFields = append(allFields, wrapper.Fields...)
	allFields = append(allFields, fields...)

	putFieldSlice(wrapper)

	return allFields, nil
}

// fields combines context and provided fields (allocating version)
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

// Init initializes the global logger with the given configuration
func Init(config Config) (*Logger, error) {
	logger, err := New(config)

	if err != nil {
		return nil, err
	}

	ReplaceGlobal(logger)
	explicitlyInited.Store(true)

	return logger, nil
}

// L returns the global logger instance
func L() *Logger {
	globalMu.RLock()
	logger := globalLogger
	globalMu.RUnlock()

	if logger == nil {
		return defaultLogger()
	}

	return logger
}

// IsInitialized returns true if the global logger was explicitly initialized
func IsInitialized() bool {
	return explicitlyInited.Load()
}

// MustInit initializes the global logger and panics on failure
func MustInit(config Config) *Logger {
	logger, err := Init(config)

	if err != nil {
		panic(fmt.Sprintf("logger: failed to initialize: %v", err))
	}

	return logger
}

// ReplaceGlobal replaces the global logger with the provided logger
func ReplaceGlobal(logger *Logger) func() {
	globalMu.Lock()
	previous := globalLogger
	globalLogger = logger
	globalMu.Unlock()

	return func() {
		ReplaceGlobal(previous)
	}
}

// defaultLogger creates a minimal default logger for when global is uninitialized
func defaultLogger() *Logger {
	globalMu.Lock()
	defer globalMu.Unlock()

	if globalLogger != nil {
		return globalLogger
	}

	config := DefaultConfig()
	config.CallerSkip = 3

	logger, err := New(config)

	if err != nil {
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

// Debug logs at debug level using the global logger
func Debug(ctx context.Context, msg string, fields ...Field) {
	L().Debug(ctx, msg, fields...)
}

// Info logs at info level using the global logger
func Info(ctx context.Context, msg string, fields ...Field) {
	L().Info(ctx, msg, fields...)
}

// Warn logs at warn level using the global logger
func Warn(ctx context.Context, msg string, fields ...Field) {
	L().Warn(ctx, msg, fields...)
}

// Error logs at error level using the global logger
func Error(ctx context.Context, msg string, fields ...Field) {
	L().Error(ctx, msg, fields...)
}

// Fatal logs at fatal level using the global logger
func Fatal(ctx context.Context, msg string, fields ...Field) {
	L().Fatal(ctx, msg, fields...)
}

// With returns a new logger with preset fields
func With(fields ...Field) *Logger {
	return L().With(fields...)
}

// Named returns a named sub-logger
func Named(name string) *Logger {
	return L().Named(name)
}

// SetGlobalLevel changes the log level of the global logger
func SetGlobalLevel(level Level) {
	L().SetLevel(level)
}

// Sync flushes the global logger's buffered logs
func Sync() error {
	globalMu.RLock()
	logger := globalLogger
	globalMu.RUnlock()

	if logger != nil {
		return logger.Sync()
	}

	return nil
}

// Close shuts down the global logger
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

// LoadConfig loads logger configuration from a YAML file
func LoadConfig(configPath string) (Config, error) {
	return config.LoadLoggerConfig(configPath)
}

// ConfigFromEnv loads logger configuration from environment variables
func ConfigFromEnv() Config {
	return config.LoadLoggerConfigFromEnv()
}

// InitFromConfig loads config from a file and initializes the global logger
func InitFromConfig(configPath string) (*Logger, error) {
	config, err := LoadConfig(configPath)

	if err != nil {
		return nil, err
	}

	return Init(config)
}

// InitFromEnv initializes the global logger using environment variables
func InitFromEnv() (*Logger, error) {
	config := ConfigFromEnv()

	return Init(config)
}

// InitFromUnifiedConfig initializes the global logger from the unified config
func InitFromUnifiedConfig(unifiedConfig config.Config) (*Logger, error) {
	return Init(unifiedConfig.GetLoggerConfig())
}
