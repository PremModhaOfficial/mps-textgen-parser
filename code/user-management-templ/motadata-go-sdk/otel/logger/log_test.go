package logger

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"

	"github.com/stretchr/testify/assert"
)

const (
	failedToCreateLogger     = "Failed to create logger"
	failedToInitGlobalLogger = "Failed to initialize global logger"
	failedToCreateLoggerFmt  = "Failed to create logger: %v"
	shouldNotLog             = "should not log"
	testModModule            = "test-mod"
	pkgTestModule            = "pkg-test"
)

func TestNewLogger(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	assertions.NotNil(loggerInstance, "Logger should not be nil")
	defer loggerInstance.Close()

	testContext := context.Background()
	loggerInstance.Info(testContext, "test message", String("key", "value"))
}

func TestLogLevels(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "debug"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	testContext := context.Background()

	loggerInstance.Debug(testContext, "debug message")
	loggerInstance.Info(testContext, "info message")
	loggerInstance.Warn(testContext, "warn message")
	loggerInstance.Error(testContext, "error message")
}

func TestContextFields(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	testContext := context.Background()
	testContext = WithTenantID(testContext, "tenant-123")
	testContext = WithRequestID(testContext, "req-456")
	testContext = WithUserID(testContext, "user-789")

	// Log should include tenant_id, request_id, user_id
	loggerInstance.Info(testContext, "request processed",
		String("action", "create"),
		Int("status", 200),
	)
}

func TestNamedLogger(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	testContext := context.Background()

	databaseLogger := loggerInstance.Named("database")
	assertions.NotNil(databaseLogger, "Named logger should not be nil")
	databaseLogger.Info(testContext, "query executed", Duration("duration", 45*time.Millisecond))

	httpRequestLogger := loggerInstance.Named("http")
	assertions.NotNil(httpRequestLogger, "Named logger should not be nil")
	httpRequestLogger.Info(testContext, "request received", String("path", "/api/users"))
}

func TestLoggerWith(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	testContext := context.Background()

	orderComponentLogger := loggerInstance.With(
		String("component", "orders"),
		String("region", "us-east-1"),
	)
	assertions.NotNil(orderComponentLogger, "Logger with fields should not be nil")

	orderComponentLogger.Info(testContext, "order created", String("order_id", "ord-123"))
	orderComponentLogger.Info(testContext, "order shipped", String("order_id", "ord-123"))
}

func TestFileOutput(t *testing.T) {
	assertions := assert.New(t)

	// Create temp directory
	temporaryDirectory := t.TempDir()
	logFilePath := temporaryDirectory + "/test.log"

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = true
	loggerConfig.FilePath = logFilePath

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, "Failed to create logger with file output")

	testContext := context.Background()
	loggerInstance.Info(testContext, "test file log", String("key", "value"))

	loggerInstance.Close()

	// Check file exists
	_, statError := os.Stat(logFilePath)
	assertions.False(os.IsNotExist(statError), "Log file should have been created")
}

func TestSetLevel(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "info"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	assertions.Equal(InfoLevel, loggerInstance.Level(), "Initial level should be InfoLevel")

	loggerInstance.SetLevel(DebugLevel)

	assertions.Equal(DebugLevel, loggerInstance.Level(), "Level should be DebugLevel after SetLevel")
}

func TestParseLevel(t *testing.T) {
	testCases := []struct {
		inputLevel    string
		expectedLevel Level
		expectError   bool
	}{
		{"debug", DebugLevel, false},
		{"DEBUG", DebugLevel, false},
		{"info", InfoLevel, false},
		{"INFO", InfoLevel, false},
		{"warn", WarnLevel, false},
		{"warning", WarnLevel, false},
		{"error", ErrorLevel, false},
		{"fatal", FatalLevel, false},
		{"invalid", InfoLevel, true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.inputLevel, func(t *testing.T) {
			assertions := assert.New(t)

			parsedLevel, parseError := ParseLevel(testCase.inputLevel)
			if testCase.expectError {
				assertions.Error(parseError, "ParseLevel(%q) should return error", testCase.inputLevel)
			} else {
				assertions.NoError(parseError, "ParseLevel(%q) should not return error", testCase.inputLevel)
				assertions.Equal(testCase.expectedLevel, parsedLevel, "ParseLevel(%q) returned wrong level", testCase.inputLevel)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	testCases := []struct {
		logLevel       Level
		expectedString string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedString, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(testCase.expectedString, testCase.logLevel.String(), "Level.String() should return correct string")
		})
	}
}

func TestCorrelationID(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	testContext := context.Background()
	testContext = WithCorrelationID(testContext, "corr-abc-123")

	loggerInstance.Info(testContext, "correlated event")
}

func TestNilContext(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, createError := New(loggerConfig)
	assertions.NoError(createError, failedToCreateLogger)
	defer loggerInstance.Close()

	// Should not panic with nil context
	assertions.NotPanics(func() {
		loggerInstance.Info(nil, "message with nil context") //nolint:staticcheck // intentional nil context test
	}, "Logger should handle nil context without panic")
}

/* ---------------------------------------- Global Logger Tests -------------------------------------------------- */

func TestGlobalLoggerInit(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.ServiceName = "test-global"

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	testContext := context.Background()
	Info(testContext, "global logger test", String("test", "init"))
}

func TestGlobalLoggerL(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	// L() should return the global logger
	globalLoggerInstance := L()
	assertions.NotNil(globalLoggerInstance, "L() should not return nil")

	testContext := context.Background()
	globalLoggerInstance.Info(testContext, "using L() accessor")
}

func TestGlobalLoggerPackageFunctions(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "debug"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	testContext := context.Background()
	testContext = WithTenantID(testContext, "tenant-global")
	testContext = WithRequestID(testContext, "req-global")

	// Test all package-level functions
	Debug(testContext, "global debug")
	Info(testContext, "global info", String("key", "value"))
	Warn(testContext, "global warn")
	Error(testContext, "global error")
}

func TestGlobalLoggerWith(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	testContext := context.Background()

	// Use package-level With
	authComponentLogger := With(String("component", "auth"))
	assertions.NotNil(authComponentLogger, "With() should return a logger")
	authComponentLogger.Info(testContext, "auth event", String("action", "login"))
}

func TestGlobalLoggerNamed(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	testContext := context.Background()

	// Use package-level Named
	databaseLogger := Named("database")
	assertions.NotNil(databaseLogger, "Named() should return a logger")
	databaseLogger.Info(testContext, "db query", Duration("latency", 10*time.Millisecond))
}

func TestGlobalLoggerSetLevel(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "info"
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	// Check initial level
	assertions.Equal(InfoLevel, L().Level(), "Initial level should be InfoLevel")

	// Change level dynamically
	SetGlobalLevel(DebugLevel)

	assertions.Equal(DebugLevel, L().Level(), "Level should be DebugLevel after SetGlobalLevel")
}

func TestGlobalLoggerReplaceGlobal(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	firstLoggerConfig := DefaultConfig()
	firstLoggerConfig.ServiceName = "service-1"
	firstLoggerConfig.ConsoleEnabled = true
	firstLoggerConfig.OTELEnabled = false
	firstLoggerConfig.FileEnabled = false

	firstLoggerInstance, firstCreateError := New(firstLoggerConfig)
	assertions.NoError(firstCreateError, "Failed to create first logger")

	secondLoggerConfig := DefaultConfig()
	secondLoggerConfig.ServiceName = "service-2"
	secondLoggerConfig.ConsoleEnabled = true
	secondLoggerConfig.OTELEnabled = false
	secondLoggerConfig.FileEnabled = false

	secondLoggerInstance, secondCreateError := New(secondLoggerConfig)
	assertions.NoError(secondCreateError, "Failed to create second logger")

	// Set first logger as global
	ReplaceGlobal(firstLoggerInstance)

	testContext := context.Background()
	Info(testContext, "from service-1")

	// Replace with second logger, get restore function
	restoreFunction := ReplaceGlobal(secondLoggerInstance)
	assertions.NotNil(restoreFunction, "ReplaceGlobal should return restore function")
	Info(testContext, "from service-2")

	// Restore first logger
	restoreFunction()
	Info(testContext, "back to service-1")

	// Cleanup
	firstLoggerInstance.Close()
	secondLoggerInstance.Close()
	Close()
}

func TestGlobalLoggerDefaultWhenUninitialized(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	// L() should return a default logger when uninitialized
	globalLoggerInstance := L()
	assertions.NotNil(globalLoggerInstance, "L() should not return nil when uninitialized")

	testContext := context.Background()
	// Should not panic
	assertions.NotPanics(func() {
		globalLoggerInstance.Info(testContext, "using default logger")
	}, "Default logger should handle logging without panic")

	// Cleanup
	Close()
}

func TestGlobalLoggerSync(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global logger
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initError := Init(loggerConfig)
	assertions.NoError(initError, failedToInitGlobalLogger)
	defer Close()

	testContext := context.Background()
	Info(testContext, "before sync")

	// Sync should not error (note: may return error on some systems for stdout)
	syncError := Sync()
	// Note: Sync() to stdout may return an error on some systems, which is ok
	if syncError != nil {
		t.Logf("Sync returned: %v (may be expected)", syncError)
	}
}

/* ---------------------------------------- Concurrent Safety Tests -------------------------------------------------- */

func TestConcurrentLoggingSafety(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()
	var wg sync.WaitGroup

	// Spawn 100 goroutines each logging 1000 messages
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				loggerInstance.Info(testContext, "concurrent test",
					Int("goroutine", goroutineID),
					Int("iteration", j),
				)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentLevelChangeSafety(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				loggerInstance.Info(testContext, "message")
			}
		}()
	}

	// Level changers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				if j%2 == 0 {
					loggerInstance.SetLevel(DebugLevel)
				} else {
					loggerInstance.SetLevel(InfoLevel)
				}
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentModuleLoggerSafety(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()
	var wg sync.WaitGroup
	modules := []string{"database", "http", "auth", "cache", "queue"}

	for _, moduleName := range modules {
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(mod string) {
				defer wg.Done()
				moduleLogger := loggerInstance.Module(mod)
				for j := 0; j < 100; j++ {
					moduleLogger.Info(testContext, "module message")
				}
			}(moduleName)
		}
	}

	wg.Wait()
}

/* ---------------------------------------- Edge Case Tests -------------------------------------------------- */

func TestEmptyMessage(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()

	// Should not panic
	loggerInstance.Info(testContext, "")
	loggerInstance.Info(testContext, "", String("key", "value"))
}

func TestLargeMessage(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()

	// Create a large message (1MB)
	largeMessage := make([]byte, 1024*1024)
	for i := range largeMessage {
		largeMessage[i] = 'a'
	}

	// Should not panic
	loggerInstance.Info(testContext, string(largeMessage))
}

func TestManyFields(t *testing.T) {
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		t.Fatalf(failedToCreateLoggerFmt, err)
	}
	defer loggerInstance.Close()

	testContext := context.Background()

	// Create 100 fields
	fields := make([]Field, 100)
	for i := 0; i < 100; i++ {
		fields[i] = Int("field", i)
	}

	// Should not panic
	loggerInstance.Info(testContext, "many fields", fields...)
}

/* ---------------------------------------- IsInitialized Tests -------------------------------------------------- */

func TestIsInitialized(t *testing.T) {
	assertions := assert.New(t)

	// Clean up first
	Close()

	// Should return false when not initialized
	assertions.False(IsInitialized(), "Should return false when not initialized")

	// Initialize
	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, initErr := Init(loggerConfig)
	assertions.NoError(initErr)
	defer Close()

	// Should return true after initialization
	assertions.True(IsInitialized(), "Should return true after initialization")
}

/* ---------------------------------------- MustInit Tests -------------------------------------------------- */

func TestMustInit(t *testing.T) {
	t.Run("successful MustInit", func(t *testing.T) {
		assertions := assert.New(t)
		Close()
		defer Close()

		loggerConfig := DefaultConfig()
		loggerConfig.ConsoleEnabled = false
		loggerConfig.OTELEnabled = false
		loggerConfig.FileEnabled = false

		assertions.NotPanics(func() {
			loggerInstance := MustInit(loggerConfig)
			assertions.NotNil(loggerInstance)
		})
	})

	t.Run("MustInit panics on invalid config", func(t *testing.T) {
		assertions := assert.New(t)
		Close()
		defer Close()

		loggerConfig := DefaultConfig()
		loggerConfig.Level = "invalid_level"

		assertions.Panics(func() {
			_ = MustInit(loggerConfig)
		})
	})
}

/* ---------------------------------------- InitFromUnifiedConfig Tests -------------------------------------------------- */

func TestInitFromUnifiedConfig(t *testing.T) {
	assertions := assert.New(t)
	Close()
	defer Close()

	cfg := config.DefaultConfig()
	cfg.Logger.ConsoleEnabled = false
	cfg.Logger.OTELEnabled = false
	cfg.Logger.FileEnabled = false

	loggerInstance, err := InitFromUnifiedConfig(cfg)
	assertions.NoError(err)
	assertions.NotNil(loggerInstance)
}

/* ---------------------------------------- ModuleLevels Additional Tests -------------------------------------------------- */

func TestModuleLevelsSetDefaultLevel(t *testing.T) {
	assertions := assert.New(t)
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, err := Init(loggerConfig)
	assertions.NoError(err)
	defer Close()

	// Get the module levels manager
	manager := getModuleLevels()

	// Set default level
	manager.SetDefaultLevel(DebugLevel)
	assertions.Equal(DebugLevel, manager.DefaultLevel())

	// Change it back
	manager.SetDefaultLevel(InfoLevel)
	assertions.Equal(InfoLevel, manager.DefaultLevel())
}

func TestModuleLevelsLevel(t *testing.T) {
	assertions := assert.New(t)
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, err := Init(loggerConfig)
	assertions.NoError(err)
	defer Close()

	manager := getModuleLevels()

	// Level for non-existent module should return default
	level := manager.Level("nonexistent-module")
	assertions.Equal(InfoLevel, level)

	// Set level for a module
	manager.SetLevel("test-module", DebugLevel)
	level = manager.Level("test-module")
	assertions.Equal(DebugLevel, level)
}

/* ---------------------------------------- Level String Edge Cases -------------------------------------------------- */

func TestLevelStringUnknown(t *testing.T) {
	assertions := assert.New(t)

	// Create an invalid level (out of range)
	invalidLevel := Level(99)
	levelStr := invalidLevel.String()

	// Should return some string representation
	assertions.NotEmpty(levelStr)
}

/* ---------------------------------------- Console Format Tests -------------------------------------------------- */

func TestConsoleFormatConsole(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = true
	loggerConfig.ConsoleFormat = "console" // Not JSON
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	testContext := context.Background()
	loggerInstance.Info(testContext, "console format test")
}

/* ---------------------------------------- InitFromConfig Error Test -------------------------------------------------- */

func TestInitFromConfigError(t *testing.T) {
	assertions := assert.New(t)
	Close()
	defer Close()

	// Try to load from non-existent file
	_, err := InitFromConfig("/nonexistent/path/config.yaml")
	assertions.Error(err)
}

/* ---------------------------------------- Log Level Bailout Tests -------------------------------------------------- */

func TestLogLevelBailout(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "error" // Only error and above
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	testContext := context.Background()

	// These should bail out early without logging
	assertions.NotPanics(func() {
		loggerInstance.Debug(testContext, shouldNotLog)
		loggerInstance.Info(testContext, shouldNotLog)
		loggerInstance.Warn(testContext, shouldNotLog)
	})

	// This should log
	assertions.NotPanics(func() {
		loggerInstance.Error(testContext, "should log")
	})
}

/* ---------------------------------------- Close Edge Cases -------------------------------------------------- */

func TestLoggerDoubleClose(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)

	// First close
	err = loggerInstance.Close()
	assertions.NoError(err)

	// Second close should be idempotent
	err = loggerInstance.Close()
	assertions.NoError(err)
}

/* ---------------------------------------- Additional Coverage Tests -------------------------------------------------- */

func TestZapLevelUnknown(t *testing.T) {
	assertions := assert.New(t)

	// Test zapLevel() with an unknown Level value (default branch)
	unknownLevel := Level(99)
	zapLvl := unknownLevel.zapLevel()

	// Unknown levels should default to InfoLevel in zapLevel()
	assertions.Equal(zapLvl, InfoLevel.zapLevel(), "Unknown level should map to zapcore.InfoLevel")

	// Also test with negative unknown value
	negativeUnknown := Level(-5)
	zapLvlNeg := negativeUnknown.zapLevel()
	assertions.Equal(zapLvlNeg, InfoLevel.zapLevel(), "Negative unknown level should map to zapcore.InfoLevel")
}

func TestPutFieldSliceNil(t *testing.T) {
	// putFieldSlice(nil) should not panic
	assertions := assert.New(t)
	assertions.NotPanics(func() {
		putFieldSlice(nil)
	}, "putFieldSlice(nil) should not panic")
}

func TestExtractFieldsPooledNilContext(t *testing.T) {
	assertions := assert.New(t)

	// extractFieldsPooled with nil context should return nil
	wrapper := extractFieldsPooled(nil) //nolint:staticcheck // intentional nil context test
	assertions.Nil(wrapper, "extractFieldsPooled(nil) should return nil")
}

func TestExtractFieldsPooledEmptyContext(t *testing.T) {
	assertions := assert.New(t)

	// Context with no fields should return nil (fast path)
	wrapper := extractFieldsPooled(context.Background())
	assertions.Nil(wrapper, "extractFieldsPooled with empty context should return nil")
}

func TestExtractFieldsWithContext(t *testing.T) {
	assertions := assert.New(t)

	// extractFields (allocating version) with context values
	ctx := context.Background()
	ctx = WithTenantID(ctx, "tenant-1")
	ctx = WithRequestID(ctx, "req-1")
	ctx = WithUserID(ctx, "user-1")
	ctx = WithCorrelationID(ctx, "corr-1")

	fields := extractFields(ctx)
	assertions.NotNil(fields, "extractFields should return fields when context has values")
	assertions.Len(fields, 4, "Should have 4 fields from context")
}

func TestExtractFieldsNilContext(t *testing.T) {
	assertions := assert.New(t)

	// extractFields with nil context should return nil
	fields := extractFields(nil) //nolint:staticcheck // intentional nil context test
	assertions.Nil(fields, "extractFields(nil) should return nil")
}

func TestGetFieldSliceAndReturn(t *testing.T) {
	assertions := assert.New(t)

	// Exercise the pool get/put cycle
	wrapper := getFieldSlice()
	assertions.NotNil(wrapper, "getFieldSlice should return a wrapper")
	assertions.Empty(wrapper.Fields, "Fields should start empty")

	// Use it
	wrapper.Fields = append(wrapper.Fields, String("key", "value"))
	assertions.Len(wrapper.Fields, 1)

	// Return to pool
	putFieldSlice(wrapper)
}

func TestSyncWhenGlobalLoggerNil(t *testing.T) {
	assertions := assert.New(t)

	// Close global logger so globalLogger becomes nil
	Close()

	// Sync() when globalLogger is nil should return nil
	err := Sync()
	assertions.NoError(err, "Sync() with nil globalLogger should return nil")
}

func TestCloseWhenGlobalLoggerNil(t *testing.T) {
	assertions := assert.New(t)

	// Ensure global logger is nil
	Close()

	// Close() again when globalLogger is nil should return nil
	err := Close()
	assertions.NoError(err, "Close() with nil globalLogger should return nil")
}

func TestGetModuleLevelsSlowPath(t *testing.T) {
	assertions := assert.New(t)

	// Force moduleLevels to nil to exercise the slow path (double-check locking)
	moduleLevelsMu.Lock()
	savedModuleLevels := moduleLevels
	moduleLevels = nil
	moduleLevelsMu.Unlock()

	// getModuleLevels should lazily initialize
	levels := getModuleLevels()
	assertions.NotNil(levels, "getModuleLevels should create a new instance when nil")

	// Verify default level is InfoLevel
	assertions.Equal(InfoLevel, levels.DefaultLevel(), "Default level should be InfoLevel")

	// Restore previous state
	moduleLevelsMu.Lock()
	moduleLevels = savedModuleLevels
	moduleLevelsMu.Unlock()
}

func TestGetModuleLevelsSlowPathConcurrent(t *testing.T) {
	assertions := assert.New(t)

	// Force moduleLevels to nil
	moduleLevelsMu.Lock()
	savedModuleLevels := moduleLevels
	moduleLevels = nil
	moduleLevelsMu.Unlock()

	// Multiple goroutines calling getModuleLevels concurrently when nil
	var wg sync.WaitGroup
	results := make([]*ModuleLevels, 10)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = getModuleLevels()
		}(i)
	}
	wg.Wait()

	// All should get the same instance
	for i := 1; i < len(results); i++ {
		assertions.Same(results[0], results[i], "All goroutines should get the same ModuleLevels instance")
	}

	// Restore
	moduleLevelsMu.Lock()
	moduleLevels = savedModuleLevels
	moduleLevelsMu.Unlock()
}

func TestModuleLevelsListModules(t *testing.T) {
	assertions := assert.New(t)

	manager := newModuleLevels(InfoLevel)

	// Initially empty
	modules := manager.ListModules()
	assertions.Empty(modules, "ListModules should return empty map initially")

	// Set some levels
	manager.SetLevel("db", DebugLevel)
	manager.SetLevel("http", WarnLevel)

	modules = manager.ListModules()
	assertions.Len(modules, 2, "ListModules should return 2 modules")
	assertions.Equal(DebugLevel, modules["db"])
	assertions.Equal(WarnLevel, modules["http"])
}

func TestModuleLevelsReset(t *testing.T) {
	assertions := assert.New(t)

	manager := newModuleLevels(InfoLevel)
	manager.SetLevel("db", DebugLevel)
	manager.SetLevel("http", WarnLevel)

	assertions.Len(manager.ListModules(), 2)

	manager.Reset()

	assertions.Empty(manager.ListModules(), "ListModules should be empty after Reset")
}

func TestModuleLevelsUpdateExisting(t *testing.T) {
	assertions := assert.New(t)

	manager := newModuleLevels(InfoLevel)

	// Set level first time (creates new entry)
	manager.SetLevel("db", DebugLevel)
	assertions.Equal(DebugLevel, manager.Level("db"))

	// Update existing entry (uses atomic update path)
	manager.SetLevel("db", ErrorLevel)
	assertions.Equal(ErrorLevel, manager.Level("db"))
}

func TestModuleLevelsGetAtomicLevelDoubleCheck(t *testing.T) {
	assertions := assert.New(t)

	manager := newModuleLevels(InfoLevel)

	// First call creates the level (slow path)
	level1 := manager.getAtomicLevel("new-module")
	assertions.NotNil(level1)

	// Second call returns existing level (fast path)
	level2 := manager.getAtomicLevel("new-module")
	assertions.Equal(level1.Level(), level2.Level(), "Should return the same level")
}

func TestModuleLoggerSetLevel(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	modLogger := loggerInstance.Module(testModModule)
	assertions.NotNil(modLogger)

	// Set level and verify
	modLogger.SetLevel(DebugLevel)
	assertions.Equal(DebugLevel, modLogger.Level())

	modLogger.SetLevel(ErrorLevel)
	assertions.Equal(ErrorLevel, modLogger.Level())

	// Verify module name
	assertions.Equal(testModModule, modLogger.ModuleName())
}

func TestModuleLoggerWithFields(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	modLogger := loggerInstance.Module(testModModule)
	withLogger := modLogger.With(String("key", "value"))
	assertions.NotNil(withLogger)
	assertions.Equal(testModModule, withLogger.ModuleName())

	// Should still be able to log
	testContext := context.Background()
	assertions.NotPanics(func() {
		withLogger.Info(testContext, "test with fields")
	})
}

func TestModuleLoggerAllLevels(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	modLogger := loggerInstance.Module("all-levels")
	modLogger.SetLevel(DebugLevel)

	testContext := context.Background()

	// Exercise all module logger methods
	assertions.NotPanics(func() {
		modLogger.Debug(testContext, "debug msg", String("k", "v"))
		modLogger.Info(testContext, "info msg", String("k", "v"))
		modLogger.Warn(testContext, "warn msg", String("k", "v"))
		modLogger.Error(testContext, "error msg", String("k", "v"))
	})

	// Now set level high so Debug/Info/Warn bail out (disabled level paths)
	modLogger.SetLevel(ErrorLevel)
	assertions.NotPanics(func() {
		modLogger.Debug(testContext, shouldNotLog)
		modLogger.Info(testContext, shouldNotLog)
		modLogger.Warn(testContext, shouldNotLog)
		modLogger.Error(testContext, "should log")
	})
}

func TestPackageLevelModuleFunctions(t *testing.T) {
	assertions := assert.New(t)
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	_, err := Init(loggerConfig)
	assertions.NoError(err)
	defer Close()

	// Test package-level Module function
	modLogger := Module(pkgTestModule)
	assertions.NotNil(modLogger)

	// Test package-level SetModuleLevel
	SetModuleLevel(pkgTestModule, DebugLevel)
	assertions.Equal(DebugLevel, ModuleLevel(pkgTestModule))

	// Test package-level SetModuleLevels
	SetModuleLevels(map[string]Level{
		"mod-a": WarnLevel,
		"mod-b": ErrorLevel,
	})
	assertions.Equal(WarnLevel, ModuleLevel("mod-a"))
	assertions.Equal(ErrorLevel, ModuleLevel("mod-b"))

	// Test ListModuleLevels
	allLevels := ListModuleLevels()
	assertions.NotEmpty(allLevels)

	// Test ResetModuleLevels
	ResetModuleLevels()
	allLevels = ListModuleLevels()
	assertions.Empty(allLevels)
}

func TestFieldsPooledWithContextAndFields(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	// Create context with fields
	ctx := context.Background()
	ctx = WithTenantID(ctx, "t1")

	// fieldsPooled with both context fields and provided fields
	allFields, wrapper := loggerInstance.fieldsPooled(ctx, []Field{String("extra", "val")})
	assertions.NotEmpty(allFields)
	// wrapper should be nil because fields were copied and wrapper was returned to pool
	assertions.Nil(wrapper, "wrapper should be nil when both context and extra fields are present")

	// fieldsPooled with context fields but no extra fields
	allFields2, wrapper2 := loggerInstance.fieldsPooled(ctx, nil)
	assertions.NotEmpty(allFields2)
	// wrapper should be non-nil (returned directly)
	if wrapper2 != nil {
		putFieldSlice(wrapper2)
	}

	// fieldsPooled with nil context (no context fields)
	allFields3, wrapper3 := loggerInstance.fieldsPooled(nil, []Field{String("k", "v")}) //nolint:staticcheck // intentional nil context test
	assertions.Len(allFields3, 1)
	assertions.Nil(wrapper3)
}

func TestCloseWithFileCloserError(t *testing.T) {
	assertions := assert.New(t)

	// Create logger with file output to a temp dir, then close it
	tmpDir := t.TempDir()
	logFilePath := tmpDir + "/close-test.log"

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = true
	loggerConfig.FilePath = logFilePath

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)

	testContext := context.Background()
	loggerInstance.Info(testContext, "before close")

	// Close should work correctly
	err = loggerInstance.Close()
	// May or may not error depending on sync, but should not panic
	if err != nil {
		t.Logf("Close returned: %v (may be expected)", err)
	}
}

func TestNewLoggerInvalidLevel(t *testing.T) {
	assertions := assert.New(t)

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "not_a_real_level"

	loggerInstance, err := New(loggerConfig)
	assertions.Error(err, "New with invalid level should return error")
	assertions.Nil(loggerInstance, "Logger should be nil on error")
}

func TestDefaultLoggerCreation(t *testing.T) {
	assertions := assert.New(t)

	// Clean up everything
	Close()

	// Reset globalLogger to nil to force defaultLogger() path
	globalMu.Lock()
	globalLogger = nil
	globalMu.Unlock()

	// L() should trigger defaultLogger()
	loggerInstance := L()
	assertions.NotNil(loggerInstance, "L() should return a non-nil default logger")

	// Clean up
	Close()
}

func TestApplyModuleLevelsFromConfig(t *testing.T) {
	assertions := assert.New(t)
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.ModuleLevels = map[string]string{
		"database": "debug",
		"http":     "warn",
		"invalid":  "not_a_level", // Should be silently ignored
	}

	loggerInstance, err := New(loggerConfig)
	assertions.NoError(err)
	defer loggerInstance.Close()

	// Verify module levels were applied
	manager := getModuleLevels()
	assertions.Equal(DebugLevel, manager.Level("database"))
	assertions.Equal(WarnLevel, manager.Level("http"))
}

func TestApplyModuleLevelsEmpty(t *testing.T) {
	// applyModuleLevels with nil/empty map should be a no-op
	assertions := assert.New(t)
	assertions.NotPanics(func() {
		applyModuleLevels(nil)
		applyModuleLevels(map[string]string{})
	})
}

func TestInitFromEnv(t *testing.T) {
	assertions := assert.New(t)
	Close()
	defer Close()

	// InitFromEnv should use env-based config (defaults)
	loggerInstance, err := InitFromEnv()
	assertions.NoError(err)
	assertions.NotNil(loggerInstance)
}

func TestConfigFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// ConfigFromEnv should return a valid config
	cfg := ConfigFromEnv()
	assertions.NotEmpty(cfg.Level)
}
