package logger

import (
	"context"
	"testing"
)

/* ---------------------------------------- Memory Allocation Benchmarks -------------------------------------------------- */

func BenchmarkLogInfoNoFields(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerInstance.Info(testContext, "simple message")
	}
}

func BenchmarkLogInfoWithFields(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerInstance.Info(testContext, "message with fields",
			String("key1", "value1"),
			String("key2", "value2"),
			Int("count", i),
		)
	}
}

func BenchmarkLogInfo10Fields(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerInstance.Info(testContext, "message with 10 fields",
			String("field1", "value1"),
			String("field2", "value2"),
			String("field3", "value3"),
			String("field4", "value4"),
			String("field5", "value5"),
			Int("field6", 100),
			Int("field7", 200),
			Int("field8", 300),
			Bool("field9", true),
			Float64("field10", 3.14),
		)
	}
}

func BenchmarkLogWithContextFields(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	testContext := context.Background()
	testContext = WithTenantID(testContext, "tenant-123")
	testContext = WithRequestID(testContext, "req-456")
	testContext = WithUserID(testContext, "user-789")
	testContext = WithCorrelationID(testContext, "corr-abc")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		loggerInstance.Info(testContext, "message with context")
	}
}

func BenchmarkLogDisabledLevel(b *testing.B) {
	loggerConfig := DefaultConfig()
	loggerConfig.Level = "error" // Only error and above
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	loggerInstance, _ := New(loggerConfig)
	defer loggerInstance.Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Debug is disabled, should be minimal overhead
		loggerInstance.Debug(testContext, "debug message",
			String("key", "value"),
			Int("count", i),
		)
	}
}

/* ---------------------------------------- Concurrent Access Benchmarks -------------------------------------------------- */

func BenchmarkLogConcurrent(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loggerInstance.Info(testContext, "concurrent message",
				String("key", "value"),
			)
		}
	})
}

func BenchmarkLogConcurrentWithContext(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		testContext := context.Background()
		testContext = WithTenantID(testContext, "tenant-123")
		testContext = WithRequestID(testContext, "req-456")

		for pb.Next() {
			loggerInstance.Info(testContext, "concurrent message with context",
				String("key", "value"),
			)
		}
	})
}

/* ---------------------------------------- Named Logger Benchmarks -------------------------------------------------- */

func BenchmarkNamedLoggerCreation(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		namedLogger := loggerInstance.Named("database")
		_ = namedLogger
	}
}

func BenchmarkLoggerWithCreation(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		withLogger := loggerInstance.With(
			String("component", "auth"),
			String("version", "v1"),
		)
		_ = withLogger
	}
}

func BenchmarkNamedLoggerLog(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	namedLogger := loggerInstance.Named("database")
	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		namedLogger.Info(testContext, "query executed")
	}
}

/* ---------------------------------------- Global Logger Benchmarks -------------------------------------------------- */

func BenchmarkGlobalLoggerAccess(b *testing.B) {
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	Init(loggerConfig)
	defer Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = L()
	}
}

func BenchmarkGlobalLoggerConcurrentAccess(b *testing.B) {
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false

	Init(loggerConfig)
	defer Close()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = L()
		}
	})
}

func BenchmarkGlobalLogInfo(b *testing.B) {
	Close()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.AddCaller = false

	Init(loggerConfig)
	defer Close()

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Info(testContext, "benchmark",
			String("key", "value"),
			Int("count", i),
		)
	}
}

/* ---------------------------------------- Module Logger Benchmarks -------------------------------------------------- */

func BenchmarkModuleLoggerCreation(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		moduleLogger := loggerInstance.Module("database")
		_ = moduleLogger
	}
}

func BenchmarkModuleLoggerLog(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	moduleLogger := loggerInstance.Module("database")
	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		moduleLogger.Info(testContext, "module message")
	}
}

func BenchmarkModuleLevelChange(b *testing.B) {
	loggerInstance := createDiscardLogger(b)
	defer loggerInstance.Close()

	moduleLogger := loggerInstance.Module("database")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			moduleLogger.SetLevel(DebugLevel)
		} else {
			moduleLogger.SetLevel(InfoLevel)
		}
	}
}

func BenchmarkModuleLevelFiltering(b *testing.B) {
	Close()
	ResetModuleLevels()

	loggerConfig := DefaultConfig()
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.AddCaller = false

	Init(loggerConfig)
	defer Close()

	// Set module level to error - most logs will be filtered
	SetModuleLevel("database", ErrorLevel)

	testContext := context.Background()
	databaseModuleLogger := Module("database")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// These will be filtered - should be very fast
		databaseModuleLogger.Debug(testContext, "filtered", Int("i", i))
		databaseModuleLogger.Info(testContext, "filtered", Int("i", i))
	}
}

/* ---------------------------------------- Field Creation Benchmarks -------------------------------------------------- */

func BenchmarkFieldString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = String("key", "value")
	}
}

func BenchmarkFieldInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Int("count", 42)
	}
}

func BenchmarkFieldErr(b *testing.B) {
	testError := context.DeadlineExceeded

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Err(testError)
	}
}

func BenchmarkFieldAny(b *testing.B) {
	testData := map[string]interface{}{
		"name": "test",
		"age":  30,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Any("data", testData)
	}
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

func createDiscardLogger(b *testing.B) *Logger {
	b.Helper()

	loggerConfig := DefaultConfig()
	loggerConfig.Level = "debug"
	loggerConfig.ConsoleEnabled = false
	loggerConfig.OTELEnabled = false
	loggerConfig.FileEnabled = false
	loggerConfig.AddCaller = false

	loggerInstance, err := New(loggerConfig)
	if err != nil {
		b.Fatalf("Failed to create logger: %v", err)
	}

	return loggerInstance
}
