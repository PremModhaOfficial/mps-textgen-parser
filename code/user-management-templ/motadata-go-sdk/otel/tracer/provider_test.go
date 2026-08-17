package tracer

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

	"github.com/stretchr/testify/assert"
)

/* ---------------------------------------- Provider Creation Tests -------------------------------------------------- */

// TestNewTracerProviderWithGRPC tests tracer provider creation with gRPC protocol
func TestNewTracerProviderWithGRPC(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.SamplingRatio = 1.0

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewTracerProviderWithHTTP tests tracer provider creation with HTTP protocol
func TestNewTracerProviderWithHTTP(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.SamplingRatio = 1.0

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewTracerProviderWithHTTPProtobuf tests tracer provider creation with HTTP/Protobuf protocol
func TestNewTracerProviderWithHTTPProtobuf(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http/protobuf"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

// TestNewTracerProviderWithDefaultProtocol tests tracer provider creation with empty protocol (defaults to gRPC)
func TestNewTracerProviderWithDefaultProtocol(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "" // Empty should default to gRPC
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Protocol Resolution Tests -------------------------------------------------- */

// TestCreateOTLPTraceExporterProtocolResolution tests protocol resolution for trace exporters
func TestCreateOTLPTraceExporterProtocolResolution(t *testing.T) {
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
			config.Enabled = true
			config.OTELEndpoint = "localhost:4317"
			config.OTELProtocol = tc.protocol
			config.OTELInsecure = true
			config.ServiceName = "test-tracer"

			ctx := context.Background()
			exporter, err := createOTLPTraceExporter(ctx, config)

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

/* ---------------------------------------- Sampler Configuration Tests -------------------------------------------------- */

// TestNewTracerProviderSamplerConfiguration tests sampler configuration
func TestNewTracerProviderSamplerConfiguration(t *testing.T) {
	testCases := []struct {
		name          string
		samplingRatio float64
		description   string
	}{
		{
			name:          "always_sample",
			samplingRatio: 1.0,
			description:   "SamplingRatio >= 1.0 should always sample",
		},
		{
			name:          "never_sample",
			samplingRatio: 0.0,
			description:   "SamplingRatio <= 0.0 should never sample",
		},
		{
			name:          "negative_sample",
			samplingRatio: -0.5,
			description:   "Negative SamplingRatio should never sample",
		},
		{
			name:          "ratio_based_50_percent",
			samplingRatio: 0.5,
			description:   "SamplingRatio 0.5 should sample 50%",
		},
		{
			name:          "ratio_based_10_percent",
			samplingRatio: 0.1,
			description:   "SamplingRatio 0.1 should sample 10%",
		},
		{
			name:          "ratio_based_99_percent",
			samplingRatio: 0.99,
			description:   "SamplingRatio 0.99 should sample 99%",
		},
		{
			name:          "above_one_always_samples",
			samplingRatio: 1.5,
			description:   "SamplingRatio > 1.0 should always sample",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			config := DefaultConfig()
			config.Enabled = true
			config.OTELEndpoint = "localhost:4317"
			config.OTELProtocol = "grpc"
			config.OTELInsecure = true
			config.ServiceName = "test-tracer"
			config.SamplingRatio = tc.samplingRatio

			provider, err := newTracerProvider(config)
			assertions.NoError(err, tc.description)
			assertions.NotNil(provider, tc.description)

			// Clean up
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = provider.Shutdown(ctx)
		})
	}
}

/* ---------------------------------------- Batch Processor Configuration Tests -------------------------------------------------- */

// TestNewTracerProviderBatchConfiguration tests batch processor configuration
func TestNewTracerProviderBatchConfiguration(t *testing.T) {
	t.Run("default_batch_settings", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.Enabled = true
		config.OTELEndpoint = "localhost:4317"
		config.OTELProtocol = "grpc"
		config.OTELInsecure = true
		config.ServiceName = "test-tracer"
		// Leave batch settings at zero to test defaults

		provider, err := newTracerProvider(config)
		assertions.NoError(err)
		assertions.NotNil(provider)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = provider.Shutdown(ctx)
	})

	t.Run("custom_batch_settings", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.Enabled = true
		config.OTELEndpoint = "localhost:4317"
		config.OTELProtocol = "grpc"
		config.OTELInsecure = true
		config.ServiceName = "test-tracer"
		config.MaxQueueSize = 4096
		config.MaxExportBatch = 1024
		config.BatchTimeout = 60 * time.Second

		provider, err := newTracerProvider(config)
		assertions.NoError(err)
		assertions.NotNil(provider)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = provider.Shutdown(ctx)
	})
}

/* ---------------------------------------- Debug Mode Tests -------------------------------------------------- */

// TestNewTracerProviderWithDebugMode tests tracer provider creation with debug stdout exporter
func TestNewTracerProviderWithDebugMode(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.OTELDebug = true // Enable debug mode
	config.ServiceName = "test-tracer"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.SamplingRatio = 1.0

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Clean up
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Exporter Creation Tests -------------------------------------------------- */

// TestCreateGRPCTraceExporter tests gRPC trace exporter creation
func TestCreateGRPCTraceExporter(t *testing.T) {
	t.Run("with_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4317"
		config.OTELInsecure = true

		ctx := context.Background()
		exporter, err := createGRPCTraceExporter(ctx, config)
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
		exporter, err := createGRPCTraceExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})
}

// TestCreateHTTPTraceExporter tests HTTP trace exporter creation
func TestCreateHTTPTraceExporter(t *testing.T) {
	t.Run("with_insecure", func(t *testing.T) {
		assertions := assert.New(t)

		config := DefaultConfig()
		config.OTELEndpoint = "localhost:4318"
		config.OTELInsecure = true

		ctx := context.Background()
		exporter, err := createHTTPTraceExporter(ctx, config)
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
		exporter, err := createHTTPTraceExporter(ctx, config)
		assertions.NoError(err)
		assertions.NotNil(exporter)

		_ = exporter.Shutdown(ctx)
	})
}

/* ---------------------------------------- Resource Creation Tests -------------------------------------------------- */

// TestTracerOTELResourceCreation tests that OTEL resource is created correctly for tracer
func TestTracerOTELResourceCreation(t *testing.T) {
	testCases := []struct {
		name           string
		serviceName    string
		serviceVersion string
		environment    string
	}{
		{
			name:           "standard_config",
			serviceName:    "my-tracer",
			serviceVersion: "1.0.0",
			environment:    "production",
		},
		{
			name:           "empty_version",
			serviceName:    "my-tracer",
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

/* ---------------------------------------- Error Handling Tests -------------------------------------------------- */

// TestNewTracerProviderInvalidProtocol tests error handling for invalid protocol
func TestNewTracerProviderInvalidProtocol(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "invalid-protocol"
	config.ServiceName = "test-tracer"

	provider, err := newTracerProvider(config)
	assertions.Error(err)
	assertions.Nil(provider)
	assertions.Contains(err.Error(), "unsupported")
}

/* ---------------------------------------- Provider Shutdown Tests -------------------------------------------------- */

// TestTracerProviderShutdown tests graceful shutdown of tracer provider
func TestTracerProviderShutdown(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"

	provider, err := newTracerProvider(config)
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

// TestTracerProviderShutdownWithTimeout tests provider shutdown respects timeout
func TestTracerProviderShutdownWithTimeout(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"

	provider, err := newTracerProvider(config)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// Shutdown with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_ = provider.Shutdown(ctx)
}

/* ---------------------------------------- Integration Tests -------------------------------------------------- */

// TestTracerInitWithProvider tests tracer initialization with provider
func TestTracerInitWithProvider(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global tracer
	_ = Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer-init"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.SamplingRatio = 1.0

	tracer, err := Init(config)
	assertions.NoError(err)
	assertions.NotNil(tracer)

	// Create a span
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-span")
	assertions.NotNil(span)
	span.End()

	// Clean up
	err = Shutdown(ctx)
	assertions.NoError(err)
}

// TestTracerInitWithHTTPProvider tests tracer initialization with HTTP provider
func TestTracerInitWithHTTPProvider(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global tracer
	_ = Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4318"
	config.OTELProtocol = "http"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer-http"
	config.ServiceVersion = "1.0.0"
	config.Environment = "test"
	config.SamplingRatio = 1.0

	tracer, err := Init(config)
	assertions.NoError(err)
	assertions.NotNil(tracer)

	// Create a span
	ctx := context.Background()
	ctx, span := tracer.Start(ctx, "test-http-span")
	assertions.NotNil(span)
	span.End()

	// Clean up
	err = Shutdown(ctx)
	assertions.NoError(err)
}

/* ---------------------------------------- Concurrent Access Tests -------------------------------------------------- */

// TestConcurrentTracerProviderCreation tests concurrent creation of tracer providers
func TestConcurrentTracerProviderCreation(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"

	// Create multiple providers concurrently
	const numProviders = 10
	providers := make(chan interface{}, numProviders)
	errors := make(chan error, numProviders)

	for i := 0; i < numProviders; i++ {
		go func() {
			provider, err := newTracerProvider(config)
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

/* ---------------------------------------- Propagator Configuration Tests -------------------------------------------------- */

// TestTracerWithCustomPropagators tests tracer initialization with custom propagators
func TestTracerWithCustomPropagators(t *testing.T) {
	testCases := []struct {
		name        string
		propagators []string
	}{
		{
			name:        "default_propagators",
			propagators: nil,
		},
		{
			name:        "tracecontext_only",
			propagators: []string{"tracecontext"},
		},
		{
			name:        "baggage_only",
			propagators: []string{"baggage"},
		},
		{
			name:        "b3_single",
			propagators: []string{"b3"},
		},
		{
			name:        "b3_multi",
			propagators: []string{"b3multi"},
		},
		{
			name:        "jaeger",
			propagators: []string{"jaeger"},
		},
		{
			name:        "all_propagators",
			propagators: []string{"tracecontext", "baggage", "b3", "b3multi", "jaeger"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			// Clean up any existing global tracer
			_ = Shutdown(context.Background())

			config := DefaultConfig()
			config.Enabled = true
			config.OTELEndpoint = "localhost:4317"
			config.OTELProtocol = "grpc"
			config.OTELInsecure = true
			config.ServiceName = "test-tracer-propagators"
			config.Propagators = tc.propagators

			tracer, err := Init(config)
			assertions.NoError(err)
			assertions.NotNil(tracer)

			// Create a span
			ctx := context.Background()
			_, span := tracer.Start(ctx, "test-propagator-span")
			assertions.NotNil(span)
			span.End()

			// Clean up
			_ = Shutdown(ctx)
		})
	}
}

/* ---------------------------------------- Tracer Instance Methods Tests -------------------------------------------------- */

// TestTracerInstanceShutdownIdempotent tests that Tracer.Shutdown is idempotent
func TestTracerInstanceShutdownIdempotent(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global tracer
	_ = Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"

	tracer, err := Init(config)
	assertions.NoError(err)
	assertions.NotNil(tracer)

	ctx := context.Background()

	// First shutdown
	err = tracer.Shutdown(ctx)
	assertions.NoError(err)

	// Second shutdown should be idempotent
	err = tracer.Shutdown(ctx)
	assertions.NoError(err)

	// Third shutdown should also work
	err = tracer.Shutdown(ctx)
	assertions.NoError(err)
}

// TestTracerStartMethodsWithProvider tests all Start* methods on Tracer instance
func TestTracerStartMethodsWithProvider(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing global tracer
	_ = Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = true
	config.OTELEndpoint = "localhost:4317"
	config.OTELProtocol = "grpc"
	config.OTELInsecure = true
	config.ServiceName = "test-tracer"
	config.SamplingRatio = 1.0

	tracer, err := Init(config)
	assertions.NoError(err)
	assertions.NotNil(tracer)
	defer func() { _ = Shutdown(context.Background()) }()

	ctx := context.Background()

	// Test Start
	_, span := tracer.Start(ctx, "internal-span")
	assertions.NotNil(span)
	span.End()

	// Test StartServer
	_, span = tracer.StartServer(ctx, "server-span")
	assertions.NotNil(span)
	span.End()

	// Test StartClient
	_, span = tracer.StartClient(ctx, "client-span")
	assertions.NotNil(span)
	span.End()

	// Test StartProducer
	_, span = tracer.StartProducer(ctx, "producer-span")
	assertions.NotNil(span)
	span.End()

	// Test StartConsumer
	_, span = tracer.StartConsumer(ctx, "consumer-span")
	assertions.NotNil(span)
	span.End()
}
