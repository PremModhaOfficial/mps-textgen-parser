package otel

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"

	"github.com/stretchr/testify/assert"
)

const (
	otelInstanceShouldNotBeNil = "OTEL instance should not be nil"
	invalidEndpoint            = "invalid://endpoint"
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
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
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
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
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
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
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

/* ---------------------------------------- Additional Coverage Tests -------------------------------------------------- */

func TestInitWithTracerEnabled(t *testing.T) {
	t.Run("tracer enabled with debug stdout exporter", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELDebug = true
		cfg.Tracer.OTELProtocol = "grpc"
		cfg.Tracer.OTELInsecure = true
		cfg.Tracer.ServiceName = "test-tracer-service"
		cfg.Tracer.SamplingRatio = 1.0

		otelInstance, err := Init(cfg)

		assertions.NoError(err, "Init with tracer enabled should succeed")
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
		assertions.NotNil(otelInstance.Tracer, "Tracer should not be nil when enabled")
		assertions.NotNil(otelInstance.Logger, "Logger should still be initialized")
	})

	t.Run("tracer and metrics both enabled", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELDebug = true
		cfg.Tracer.OTELProtocol = "grpc"
		cfg.Tracer.OTELInsecure = true
		cfg.Tracer.SamplingRatio = 1.0
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = false
		cfg.Metrics.PrometheusEnabled = false

		otelInstance, err := Init(cfg)

		assertions.NoError(err, "Init with tracer and metrics enabled should succeed")
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
		assertions.NotNil(otelInstance.Tracer, "Tracer should not be nil")
		assertions.True(otelInstance.metrics, "metrics flag should be true")
	})
}

func TestInitMetricsFailureWithCleanup(t *testing.T) {
	t.Run("metrics failure cleans up logger", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		// Enable metrics with OTEL exporter pointing to an invalid endpoint
		// and use an invalid protocol to trigger a failure
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = true
		cfg.Metrics.OTELProtocol = "invalid_protocol"
		cfg.Metrics.OTELEndpoint = invalidEndpoint

		otelInstance, err := Init(cfg)

		assertions.Error(err, "Init should return error when metrics init fails")
		assertions.Nil(otelInstance, "OTEL instance should be nil on metrics failure")
		assertions.Contains(err.Error(), "initialize metrics")
	})
}

func TestInitTracerFailureWithCleanup(t *testing.T) {
	t.Run("tracer failure cleans up logger only (no metrics)", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELProtocol = "invalid_protocol"
		cfg.Tracer.OTELEndpoint = invalidEndpoint
		cfg.Metrics.Enabled = false

		otelInstance, err := Init(cfg)

		assertions.Error(err, "Init should return error when tracer init fails")
		assertions.Nil(otelInstance, "OTEL instance should be nil on tracer failure")
		assertions.Contains(err.Error(), "initialize tracer")
	})

	t.Run("tracer failure cleans up logger and metrics", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		// Enable metrics (no-op, will succeed)
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = false
		cfg.Metrics.PrometheusEnabled = false
		// Enable tracer with invalid protocol to cause failure
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELProtocol = "invalid_protocol"
		cfg.Tracer.OTELEndpoint = invalidEndpoint

		otelInstance, err := Init(cfg)

		assertions.Error(err, "Init should return error when tracer init fails")
		assertions.Nil(otelInstance, "OTEL instance should be nil on tracer failure")
		assertions.Contains(err.Error(), "initialize tracer")
	})
}

func TestInitFromFileWithValidFile(t *testing.T) {
	t.Run("valid YAML config file", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		// Create a temporary valid YAML config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		yamlContent := `service:
  name: test-file-service
  version: "1.0.0"
  environment: test
logger:
  level: info
  console:
    enabled: false
  otel:
    enabled: false
  file:
    enabled: false
metrics:
  enabled: false
tracer:
  enabled: false
`
		err := os.WriteFile(configPath, []byte(yamlContent), 0644)
		assertions.NoError(err, "Should be able to write temp config file")

		otelInstance, initErr := InitFromFile(configPath)

		assertions.NoError(initErr, "InitFromFile should succeed with valid YAML")
		assertions.NotNil(otelInstance, otelInstanceShouldNotBeNil)
		assertions.Equal("test-file-service", otelInstance.config.Service.Name)

		_ = otelInstance.Shutdown(context.Background())
	})

	t.Run("invalid YAML content", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		// Create a temporary file with invalid YAML
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "bad_config.yaml")

		invalidYAML := `{{{not valid yaml: [[[`
		err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
		assertions.NoError(err, "Should be able to write temp file")

		otelInstance, initErr := InitFromFile(configPath)

		assertions.Error(initErr, "InitFromFile should return error for invalid YAML")
		assertions.Nil(otelInstance, "OTEL instance should be nil on error")
		assertions.Contains(initErr.Error(), "load config")
	})
}

func TestMustInitFromFileSuccess(t *testing.T) {
	t.Run("success case with valid file", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		// Create a temporary valid YAML config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		yamlContent := `service:
  name: must-init-file-service
  version: "2.0.0"
  environment: test
logger:
  level: info
  console:
    enabled: false
  otel:
    enabled: false
  file:
    enabled: false
metrics:
  enabled: false
tracer:
  enabled: false
`
		err := os.WriteFile(configPath, []byte(yamlContent), 0644)
		assertions.NoError(err, "Should be able to write temp config file")

		assertions.NotPanics(func() {
			otelInstance := MustInitFromFile(configPath)
			assertions.NotNil(otelInstance)
			assertions.Equal("must-init-file-service", otelInstance.config.Service.Name)
			_ = otelInstance.Shutdown(context.Background())
		}, "MustInitFromFile should not panic with valid file")
	})
}

func TestMustInitFromEnvPanic(t *testing.T) {
	t.Run("panics on invalid env config", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		// Set an invalid log level to trigger a failure
		t.Setenv("LOG_LEVEL", "invalid_level")
		t.Setenv("LOG_CONSOLE_ENABLED", "false")
		t.Setenv("LOG_OTEL_ENABLED", "false")
		t.Setenv("LOG_FILE_ENABLED", "false")
		t.Setenv("METRICS_ENABLED", "false")
		t.Setenv("TRACER_ENABLED", "false")

		assertions.Panics(func() {
			_ = MustInitFromEnv()
		}, "MustInitFromEnv should panic with invalid env config")
	})
}

func TestShutdownWithPartialComponents(t *testing.T) {
	t.Run("shutdown with only logger", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		otelInstance, err := Init(cfg)
		assertions.NoError(err)

		// otelInstance has Logger but no Tracer and metrics=false
		assertions.NotNil(otelInstance.Logger)
		assertions.Nil(otelInstance.Tracer)
		assertions.False(otelInstance.metrics)

		ctx := context.Background()
		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr, "Shutdown with only logger should succeed")
	})

	t.Run("shutdown with logger and metrics only", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = false
		cfg.Metrics.PrometheusEnabled = false

		otelInstance, err := Init(cfg)
		assertions.NoError(err)

		assertions.NotNil(otelInstance.Logger)
		assertions.Nil(otelInstance.Tracer)
		assertions.True(otelInstance.metrics)

		ctx := context.Background()
		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr, "Shutdown with logger and metrics should succeed")
	})

	t.Run("shutdown with logger and tracer only", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELDebug = true
		cfg.Tracer.OTELProtocol = "grpc"
		cfg.Tracer.OTELInsecure = true
		cfg.Tracer.SamplingRatio = 1.0

		otelInstance, err := Init(cfg)
		assertions.NoError(err)

		assertions.NotNil(otelInstance.Logger)
		assertions.NotNil(otelInstance.Tracer)
		assertions.False(otelInstance.metrics)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr, "Shutdown with logger and tracer should succeed")
	})

	t.Run("shutdown with all three components", func(t *testing.T) {
		assertions := assert.New(t)
		cleanupGlobals()
		defer cleanupGlobals()

		cfg := testConfig()
		cfg.Metrics.Enabled = true
		cfg.Metrics.OTELEnabled = false
		cfg.Metrics.PrometheusEnabled = false
		cfg.Tracer.Enabled = true
		cfg.Tracer.OTELDebug = true
		cfg.Tracer.OTELProtocol = "grpc"
		cfg.Tracer.OTELInsecure = true
		cfg.Tracer.SamplingRatio = 1.0

		otelInstance, err := Init(cfg)
		assertions.NoError(err)

		assertions.NotNil(otelInstance.Logger)
		assertions.NotNil(otelInstance.Tracer)
		assertions.True(otelInstance.metrics)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownErr := otelInstance.Shutdown(ctx)
		assertions.NoError(shutdownErr, "Shutdown with all components should succeed")
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
