package logger

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

	"github.com/stretchr/testify/assert"
)

/* ---------------------------------------- Provider Creation Tests -------------------------------------------------- */

// TestNewOTELProviderWithGRPC tests OTEL provider creation with gRPC protocol
func TestNewOTELProviderWithGRPC(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-service"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewOTELProviderWithHTTP tests OTEL provider creation with HTTP protocol
func TestNewOTELProviderWithHTTP(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http"
	config.OTELInsecure = true
	config.ServiceName = "test-service"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewOTELProviderWithHTTPProtobuf tests OTEL provider creation with HTTP/Protobuf protocol
func TestNewOTELProviderWithHTTPProtobuf(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http/protobuf"
	config.OTELInsecure = true
	config.ServiceName = "test-service"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewOTELProviderWithDefaultProtocol tests OTEL provider creation with empty protocol (defaults to gRPC)
func TestNewOTELProviderWithDefaultProtocol(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "" // Empty should default to gRPC
	config.OTELInsecure = true
	config.ServiceName = "test-service"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Protocol Resolution Tests -------------------------------------------------- */

// TestCreateOTLPLogExporterProtocolResolution tests protocol resolution for log exporters
func TestCreateOTLPLogExporterProtocolResolution(t *testing.T) {
	testCases := []struct {
		name         string
		protocol     string
		expectError  bool
		errorContain string
	}{
		{
			name:        "grpc_lowercase",
			protocol:    "grpc",
			expectError: false,
		},
		{
			name:        "grpc_uppercase",
			protocol:    "GRPC",
			expectError: false,
		},
		{
			name:        "http_lowercase",
			protocol:    "http",
			expectError: false,
		},
		{
			name:        "http_uppercase",
			protocol:    "HTTP",
			expectError: false,
		},
		{
			name:        "http_protobuf",
			protocol:    "http/protobuf",
			expectError: false,
		},
		{
			name:        "empty_defaults_to_grpc",
			protocol:    "",
			expectError: false,
		},
		{
			name:         "unsupported_protocol",
			protocol:     "invalid",
			expectError:  true,
			errorContain: "unsupported",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			config := DefaultConfig()
			config.OTELEnabled = true
			config.OTELEndpoint = "localhost:4317"
			config.OTELProtocol = tc.protocol
			config.OTELInsecure = true
			config.ServiceName = "test-service"

			ctx := context.Background()
			exporter, err := createOTLPLogExporter(ctx, config)

			if tc.expectError {
				assertions.Error(err)
				if tc.errorContain != "" {
					assertions.Contains(err.Error(), tc.errorContain)
				}
				assertions.Nil(exporter)
			} else {
				assertions.NoError(err)
				assertions.NotNil(exporter)
				// Clean up
				_ = exporter.Shutdown(ctx)
			}
		})
	}
}

/* ---------------------------------------- Batch Processor Configuration Tests -------------------------------------------------- */

// TestNewOTELProviderBatchConfiguration tests batch processor configuration defaults
func TestNewOTELProviderBatchConfiguration(t *testing.T) {
	t.Run("default_batch_settings", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEnabled = true
		config.OTELEndpoint = "localhost:4317"
		config.OTELProtocol = "grpc"
		config.OTELInsecure = true
		config.ServiceName = "test-service"
		// Leave batch settings at zero to test defaults

		provider, err := newOTELProvider(config)
		assertions.NoError(err)
		assertions.NotNil(provider)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = provider.Shutdown(ctx)
	})

	t.Run("custom_batch_settings", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEnabled = true
		config.OTELEndpoint = "localhost:4317"
		config.OTELProtocol = "grpc"
		config.OTELInsecure = true
		config.ServiceName = "test-service"
		config.MaxQueueSize = 4096
		config.MaxExportBatch = 1024
		config.ExportInterval = 2 * time.Second
		config.BatchTimeout = 60 * time.Second

		provider, err := newOTELProvider(config)
		assertions.NoError(err)
		assertions.NotNil(provider)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = provider.Shutdown(ctx)
	})
}

/* ---------------------------------------- Debug Mode Tests -------------------------------------------------- */

// TestNewOTELProviderWithDebugMode tests OTEL provider creation with debug stdout exporter
func TestNewOTELProviderWithDebugMode(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.OTELDebug = true // Enable debug mode
	config.ServiceName = "test-service"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Exporter Creation Tests -------------------------------------------------- */

// TestCreateGRPCLogExporter tests gRPC log exporter creation
func TestCreateGRPCLogExporter(t *testing.T) {
	t.Run("with_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4317"
		config.OTELInsecure = true

		ctx := context.Background()
		exporter, err := createGRPCLogExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})

	t.Run("without_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4317"
		config.OTELInsecure = false

		ctx := context.Background()
		exporter, err := createGRPCLogExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})
}

// TestCreateHTTPLogExporter tests HTTP log exporter creation
func TestCreateHTTPLogExporter(t *testing.T) {
	t.Run("with_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4318"
		config.OTELInsecure = true

		ctx := context.Background()
		exporter, err := createHTTPLogExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})

	t.Run("without_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4318"
		config.OTELInsecure = false

		ctx := context.Background()
		exporter, err := createHTTPLogExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})
}

/* ---------------------------------------- Resource Creation Tests -------------------------------------------------- */

// TestOTELResourceCreation tests that OTEL resource is created correctly
func TestOTELResourceCreation(t *testing.T) {
	testCases := []struct {
		name           string
		serviceName    string
		serviceVersion string
		environment    string
	}{
		{
			name:           "standard_config",
			serviceName:    "my-service",
			serviceVersion: "1.0.0",
			environment:    "production",
		},
		{
			name:           "empty_version",
			serviceName:    "my-service",
			serviceVersion: "",
			environment:    "development",
		},
		{
			name:           "minimal_config",
			serviceName:    "app",
			serviceVersion: "0.0.0",
			environment:    "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			resource, err := common.NewOTELResourceFromConfig(
				tc.serviceName,
				tc.serviceVersion,
				tc.environment,
			)
			assertions.NoError(err)
			assertions.NotNil(resource)
		})
	}
}

/* ---------------------------------------- Integration Tests -------------------------------------------------- */

// TestLoggerWithOTELProvider tests creating a logger with OTEL provider enabled
// Note: This test creates exporters but doesn't require an actual collector
func TestLoggerWithOTELProvider(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.ConsoleEnabled = false
	config.FileEnabled = false
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-logger"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.ShutdownTimeout = 100 * time.Millisecond // Short timeout since no collector

	logger, err := New(config)
	assertions.NoError(err)
	assertions.NotNil(logger)

	// Log a message (won't actually be sent without collector)
	ctx := context.Background()
	logger.Info(ctx, "test message from OTEL provider test")

	// Clean up - ignore timeout errors since no collector is running
	_ = logger.Close()
}

// TestLoggerWithOTELProviderAndHTTP tests creating a logger with HTTP OTEL provider
// Note: This test creates exporters but doesn't require an actual collector
func TestLoggerWithOTELProviderAndHTTP(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.ConsoleEnabled = false
	config.FileEnabled = false
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http"
	config.OTELInsecure = true
	config.ServiceName = "test-logger-http"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.ShutdownTimeout = 100 * time.Millisecond // Short timeout since no collector

	logger, err := New(config)
	assertions.NoError(err)
	assertions.NotNil(logger)

	// Log a message (won't actually be sent without collector)
	ctx := context.Background()
	logger.Info(ctx, "test message from HTTP OTEL provider test")

	// Clean up - ignore errors since no collector is running
	_ = logger.Close()
}

/* ---------------------------------------- Error Handling Tests -------------------------------------------------- */

// TestNewOTELProviderInvalidProtocol tests error handling for invalid protocol
func TestNewOTELProviderInvalidProtocol(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "invalid-protocol"
	config.ServiceName = "test-service"

	provider, err := newOTELProvider(config)
	assertions.Error(err)
	assertions.Nil(provider)
	assertions.Contains(err.Error(), "unsupported")
}

// TestLoggerCreationWithInvalidOTELConfig tests logger creation fails with invalid OTEL config
func TestLoggerCreationWithInvalidOTELConfig(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.ConsoleEnabled = false
	config.FileEnabled = false
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "invalid-protocol"
	config.ServiceName = "test-service"

	logger, err := New(config)
	assertions.Error(err)
	assertions.Nil(logger)
}

/* ---------------------------------------- Provider Shutdown Tests -------------------------------------------------- */

// TestOTELProviderShutdown tests graceful shutdown of OTEL provider
func TestOTELProviderShutdown(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-service"

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Shutdown should succeed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = provider.Shutdown(ctx)
	assertions.NoError(err)

	// Second shutdown should also succeed (idempotent)
	err = provider.Shutdown(ctx)
	assertions.NoError(err)
}

// TestOTELProviderShutdownWithTimeout tests provider shutdown respects timeout
func TestOTELProviderShutdownWithTimeout(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-service"
	config.ShutdownTimeout = 1 * time.Second

	provider, err := newOTELProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Shutdown with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Concurrent Access Tests -------------------------------------------------- */

// TestConcurrentOTELProviderCreation tests concurrent creation of OTEL providers
func TestConcurrentOTELProviderCreation(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.OTELEnabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-service"

	// Create multiple providers concurrently
	const numProviders = 10
	providers := make(chan interface{}, numProviders)
	errors := make(chan error, numProviders)

	for i := 0; i < numProviders; i++ {
		go func() {
			provider, err := newOTELProvider(config)
			if err != nil {
				errors <- err
				return
			}
			providers <- provider
		}()
	}

	// Collect results
	var createdProviders []interface{}
	for i := 0; i < numProviders; i++ {
		select {
		case provider := <-providers:
			createdProviders = append(createdProviders, provider)
		case err := <-errors:
			t.Errorf("Error creating provider: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for providers")
		}
	}

	assertions.Len(createdProviders, numProviders)

	// Clean up all providers
	ctx := context.Background()
	for _, p := range createdProviders {
		if provider, ok := p.(interface{ Shutdown(context.Context) error }); ok {
			_ = provider.Shutdown(ctx)
		}
	}
}
