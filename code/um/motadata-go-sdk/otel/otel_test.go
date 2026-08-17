package otel

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
	"github.com/stretchr/testify/assert"
)

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

// testConfig returns a minimal config for testing with all external connections disabled
func testConfig() config.Config {
	cfg := config.DefaultConfig()

	cfg.Service.Name = "test-service"
	cfg.Service.Version = "1.0.0"
	cfg.Service.Environment = "test"

	// Disable all external exporters
	cfg.Logger.ConsoleEnabled = false
	cfg.Logger.OTELEnabled = false
	cfg.Logger.FileEnabled = false

	cfg.Metrics.Enabled = false
	cfg.Metrics.OTELEnabled = false
	cfg.Metrics.PrometheusEnabled = false

	cfg.Tracer.Enabled = false

	return cfg
}

// cleanupGlobals ensures a clean state for testing
func cleanupGlobals() {
	_ = logger.Close()
	_ = metrics.Shutdown(context.Background())
	_ = tracer.Shutdown(context.Background())
}

/* ---------------------------------------- Init Tests -------------------------------------------------- */

func TestInit(t *testing.T) {
	t.Run("successful initialization with all disabled", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		otelInstance, err := Init(cfg)

		assertions.NoError(err, "Init should not return error")
		assertions.NotNil(otelInstance, "OTEL instance should not be nil")
		assertions.NotNil(otelInstance.Logger, "Logger should not be nil")
		assertions.Nil(otelInstance.Tracer, "Tracer should be nil when disabled")
	})

	t.Run("logger initialization failure propagates error", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Logger.Level = "invalid_level" // This will cause logger init to fail

		otelInstance, err := Init(cfg)

		assertions.Error(err, "Init should return error for invalid logger config")
		assertions.Nil(otelInstance, "OTEL instance should be nil on error")
		assertions.Contains(err.Error(), "initialize logger")
	})

	t.Run("initialization with metrics enabled (no-op)", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		// Enable metrics but without OTEL or Prometheus exporters
		// This should create a no-op provider internally
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = false
		cfg.Metrics.PrometheusEnabled = false

		otelInstance, err := Init(cfg)

		assertions.NoError(err, "Init with no-op metrics should succeed")
		assertions.NotNil(otelInstance, "OTEL instance should not be nil")
	})
}

func TestInitFromEnv(t *testing.T) {
	t.Run("initializes from environment", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		// Set minimal environment for testing
		t.Setenv("LOG_CONSOLE_ENABLED", "false")
		t.Setenv("LOG_OTEL_ENABLED", "false")
		t.Setenv("LOG_FILE_ENABLED", "false")
		t.Setenv("METRICS_ENABLED", "false")
		t.Setenv("TRACER_ENABLED", "false")

		otelInstance, err := InitFromEnv()

		assertions.NoError(err, "InitFromEnv should not return error")
		assertions.NotNil(otelInstance, "OTEL instance should not be nil")
	})
}

func TestMustInit(t *testing.T) {
	t.Run("successful MustInit", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()

		assertions.NotPanics(func() {
			otelInstance := MustInit(cfg)
			assertions.NotNil(otelInstance)
		}, "MustInit should not panic with valid config")
	})

	t.Run("MustInit panics on invalid config", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Logger.Level = "invalid_level"

		assertions.Panics(func() {
			_ = MustInit(cfg)
		}, "MustInit should panic with invalid config")
	})
}

func TestMustInitFromEnv(t *testing.T) {
	t.Run("successful MustInitFromEnv", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		t.Setenv("LOG_CONSOLE_ENABLED", "false")
		t.Setenv("LOG_OTEL_ENABLED", "false")
		t.Setenv("LOG_FILE_ENABLED", "false")
		t.Setenv("METRICS_ENABLED", "false")
		t.Setenv("TRACER_ENABLED", "false")

		assertions.NotPanics(func() {
			otelInstance := MustInitFromEnv()
			assertions.NotNil(otelInstance)
		}, "MustInitFromEnv should not panic with valid env")
	})
}

/* ---------------------------------------- Shutdown Tests -------------------------------------------------- */

func TestOTELShutdown(t *testing.T) {
	t.Run("shutdown with all components", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		otelInstance, err := Init(cfg)
		assertions.NoError(err)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr, "Shutdown should not return error")
	})

	t.Run("shutdown with nil components", func(t *testing.T) {
		assertions := assert.New(t)

		otelInstance := &OTEL{
			Logger:  nil,
			Tracer:  nil,
			metrics: false,
		}

		ctx := context.Background()
		err := otelInstance.Shutdown(ctx)

		assertions.NoError(err, "Shutdown should handle nil components")
	})
}

func TestPackageLevelShutdown(t *testing.T) {
	t.Run("package level shutdown", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		_, initErr := Init(cfg)
		assertions.NoError(initErr)

		ctx := context.Background()
		err := Shutdown(ctx)

		assertions.NoError(err, "Package level Shutdown should not return error")
	})

	t.Run("shutdown without initialization", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()

		ctx := context.Background()
		err := Shutdown(ctx)

		assertions.NoError(err, "Shutdown should not error when not initialized")
	})
}

/* ---------------------------------------- InitFromFile Tests -------------------------------------------------- */

func TestInitFromFile(t *testing.T) {
	t.Run("file not found error", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		_, err := InitFromFile("/nonexistent/path/config.yaml")

		assertions.Error(err, "InitFromFile should return error for nonexistent file")
		assertions.Contains(err.Error(), "load config")
	})
}

func TestMustInitFromFile(t *testing.T) {
	t.Run("panics on file not found", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		assertions.Panics(func() {
			_ = MustInitFromFile("/nonexistent/path/config.yaml")
		}, "MustInitFromFile should panic for nonexistent file")
	})
}

/* ---------------------------------------- OTEL Struct Tests -------------------------------------------------- */

func TestOTELStructFields(t *testing.T) {
	t.Run("config is stored", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Service.Name = "my-test-service"

		otelInstance, err := Init(cfg)
		assertions.NoError(err)
		defer func() { _ = otelInstance.Shutdown(context.Background()) }()

		// The config should be stored
		assertions.Equal("my-test-service", otelInstance.config.Service.Name)
	})
}

/* ---------------------------------------- Integration Tests -------------------------------------------------- */

func TestFullLifecycle(t *testing.T) {
	t.Run("init and shutdown cycle", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Metrics.Enabled = true

		otelInstance, initErr := Init(cfg)
		assertions.NoError(initErr)
		assertions.NotNil(otelInstance)

		// Use the logger
		testContext := context.Background()
		otelInstance.Logger.Info(testContext, "test message")

		// Shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr)
	})

	t.Run("multiple init-shutdown cycles", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()

		for i := 0; i < 3; i++ {
			otelInstance, err := Init(cfg)
			assertions.NoError(err, "Init cycle %d should succeed", i)
			assertions.NotNil(otelInstance)

			ctx := context.Background()
			shutdownErr := otelInstance.Shutdown(ctx)
			assertions.NoError(shutdownErr, "Shutdown cycle %d should succeed", i)

			cleanupGlobals()
		}
	})
}

/* ---------------------------------------- Benchmark Tests -------------------------------------------------- */

func BenchmarkInit(b *testing.B) {
	cfg := testConfig()
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		otelInstance, _ := Init(cfg)
		if otelInstance != nil {
			_ = otelInstance.Shutdown(context.Background())
		}
		cleanupGlobals()
	}
}

func BenchmarkShutdown(b *testing.B) {
	cfg := testConfig()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		otelInstance, _ := Init(cfg)
		b.StartTimer()

		if otelInstance != nil {
			_ = otelInstance.Shutdown(context.Background())
		}
		cleanupGlobals()
	}
}
