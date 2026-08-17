package logger

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModuleLogger(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "info"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	testContext := context.Background()

	// Create module loggers
	databaseModuleLogger := Module("database")
	httpModuleLogger := Module("http")

	assertions.NotNil(databaseModuleLogger, "Database module logger should not be nil")
	assertions.NotNil(httpModuleLogger, "HTTP module logger should not be nil")

	// Both should inherit default level (info)
	databaseModuleLogger.Info(testContext, "database info message")
	httpModuleLogger.Info(testContext, "http info message")
}

func TestModuleLevelFromConfig(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "info"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.ModuleLevels = map[string]string{
		"database": "debug",
		"http":     "warn",
		"auth":     "error",
	}

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	// Verify module levels
	assertions.Equal(DebugLevel, ModuleLevel("database"), "Database module level should be debug")
	assertions.Equal(WarnLevel, ModuleLevel("http"), "HTTP module level should be warn")
	assertions.Equal(ErrorLevel, ModuleLevel("auth"), "Auth module level should be error")
}

func TestModuleDynamicLevelChange(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "info"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	testContext := context.Background()

	// Create module logger
	databaseModuleLogger := Module("database")

	// Initial level should be info (inherited)
	assertions.Equal(InfoLevel, databaseModuleLogger.Level(), "Initial level should be InfoLevel")

	// Change level dynamically via module logger
	databaseModuleLogger.SetLevel(DebugLevel)
	assertions.Equal(DebugLevel, databaseModuleLogger.Level(), "Level should be DebugLevel after SetLevel")

	// Change level via global function
	SetModuleLevel("database", WarnLevel)
	assertions.Equal(WarnLevel, databaseModuleLogger.Level(), "Level should be WarnLevel after SetModuleLevel")

	// Create another logger for same module - should share level
	secondDatabaseModuleLogger := Module("database")
	assertions.Equal(WarnLevel, secondDatabaseModuleLogger.Level(), "New module logger should share level with existing")

	// Log messages
	databaseModuleLogger.Info(testContext, "this should be filtered (info < warn)")
	databaseModuleLogger.Warn(testContext, "this should appear (warn >= warn)")
}

func TestModuleLevelFiltering(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "debug" // Base level allows all
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	testContext := context.Background()

	// Set module to warn level
	SetModuleLevel("database", WarnLevel)

	databaseModuleLogger := Module("database")
	assertions.NotNil(databaseModuleLogger, "Database module logger should not be nil")

	// These should be filtered
	databaseModuleLogger.Debug(testContext, "debug - should be filtered")
	databaseModuleLogger.Info(testContext, "info - should be filtered")

	// These should appear
	databaseModuleLogger.Warn(testContext, "warn - should appear")
	databaseModuleLogger.Error(testContext, "error - should appear")
}

func TestModuleLevelFromYAML(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	temporaryDirectory := t.TempDir()
	configFilePath := filepath.Join(temporaryDirectory, "config.yaml")

	yamlConfigContent := `
service:
  name: test-service

logger:
  level: info
  console_enabled: true
  modules:
    database: debug
    http: warn
    cache: error
`
	writeError := os.WriteFile(configFilePath, []byte(yamlConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	_, initError := InitFromConfig(configFilePath)
	assertions.NoError(initError, "InitFromConfig should not fail")
	defer Close()

	// Verify module levels
	assertions.Equal(DebugLevel, ModuleLevel("database"), "Database module level should be debug")
	assertions.Equal(WarnLevel, ModuleLevel("http"), "HTTP module level should be warn")
	assertions.Equal(ErrorLevel, ModuleLevel("cache"), "Cache module level should be error")
}

func TestModuleLevelFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	// Set module levels via env
	os.Setenv("LOG_MODULE_DATABASE", "debug")
	os.Setenv("LOG_MODULE_HTTP", "warn")
	os.Setenv("LOG_MODULE_AUTH", "error")
	defer func() {
		os.Unsetenv("LOG_MODULE_DATABASE")
		os.Unsetenv("LOG_MODULE_HTTP")
		os.Unsetenv("LOG_MODULE_AUTH")
	}()

	_, initError := InitFromEnv()
	assertions.NoError(initError, "InitFromEnv should not fail")
	defer Close()

	// Verify module levels (env vars are lowercased)
	assertions.Equal(DebugLevel, ModuleLevel("database"), "Database module level should be debug")
	assertions.Equal(WarnLevel, ModuleLevel("http"), "HTTP module level should be warn")
	assertions.Equal(ErrorLevel, ModuleLevel("auth"), "Auth module level should be error")
}

func TestModuleLoggerWith(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	testContext := context.Background()

	// Create module logger with preset fields
	databaseModuleLoggerWithFields := Module("database").With(
		String("component", "postgres"),
		String("host", "localhost"),
	)
	assertions.NotNil(databaseModuleLoggerWithFields, "Module logger with fields should not be nil")

	databaseModuleLoggerWithFields.Info(testContext, "query executed", Duration("latency", 10))
}

func TestSetModuleLevels(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	// Set multiple levels at once
	SetModuleLevels(map[string]Level{
		"database": DebugLevel,
		"http":     WarnLevel,
		"cache":    ErrorLevel,
	})

	// Verify
	assertions.Equal(DebugLevel, ModuleLevel("database"), "Database module level should be debug")
	assertions.Equal(WarnLevel, ModuleLevel("http"), "HTTP module level should be warn")
	assertions.Equal(ErrorLevel, ModuleLevel("cache"), "Cache module level should be error")
}

func TestListModuleLevels(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.ModuleLevels = map[string]string{
		"database": "debug",
		"http":     "warn",
	}

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	moduleLevelsMap := ListModuleLevels()

	assertions.Equal(2, len(moduleLevelsMap), "ListModuleLevels should return 2 modules")
	assertions.Equal(DebugLevel, moduleLevelsMap["database"], "Database module level should be debug")
	assertions.Equal(WarnLevel, moduleLevelsMap["http"], "HTTP module level should be warn")
}

func TestModuleLoggerName(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	databaseModuleLogger := Module("database")
	assertions.Equal("database", databaseModuleLogger.ModuleName(), "ModuleName() should return correct module name")
}

func TestModuleLoggerFromLogger(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, "Failed to create logger")
	defer loggerInstance.Close()

	testContext := context.Background()

	// Create module logger from instance
	databaseModuleLoggerFromInstance := loggerInstance.Module("database")
	assertions.NotNil(databaseModuleLoggerFromInstance, "Module logger from instance should not be nil")
	databaseModuleLoggerFromInstance.Info(testContext, "from instance module logger")

	// Set level on instance module
	databaseModuleLoggerFromInstance.SetLevel(DebugLevel)
	databaseModuleLoggerFromInstance.Debug(testContext, "debug from instance module logger")
}

func TestModuleLevelInheritance(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "warn" // Default level is warn
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	// Module without explicit level should inherit default
	databaseModuleLogger := Module("database")
	assertions.Equal(WarnLevel, databaseModuleLogger.Level(), "Module should inherit default warn level")
}

func TestModuleContextFields(t *testing.T) {
	assertions := assert.New(t)

	// Clean up
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, "Failed to init logger")
	defer Close()

	testContext := context.Background()
	testContext = WithTenantID(testContext, "tenant-123")
	testContext = WithRequestID(testContext, "req-456")

	databaseModuleLogger := Module("database")
	assertions.NotNil(databaseModuleLogger, "Module logger should not be nil")
	databaseModuleLogger.Info(testContext, "with context fields", String("query", "SELECT *"))
}
