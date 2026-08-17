package tracer

import (
	"context"
	"fmt"
	"os"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

/* ---------------------------------------- Provider Constructor -------------------------------------------------- */

// newTracerProvider creates a new OTEL tracer provider with the configured exporters.
// This function sets up the OpenTelemetry tracing pipeline including:
//   - Resource identification (service name, version, environment)
//   - OTLP exporter for sending traces to an OTEL collector
//   - Batch span processor for efficient span batching
//   - Sampler configuration for controlling trace volume
//   - Optional stdout exporter for debugging
//
// The provider uses a batch span processor to aggregate spans before export,
// improving performance by reducing network calls and enabling backpressure handling.
//
// Parameters:
//   - tracerConfiguration: The tracer configuration containing OTEL settings
//
// Returns:
//   - *sdktrace.TracerProvider: The configured provider
//   - error: Non-nil if provider creation fails
func newTracerProvider(tracerConfiguration Config) (*sdktrace.TracerProvider, error) {

	// Use background context for provider creation
	// This is appropriate as provider setup happens at application startup
	backgroundContext := context.Background()

	// Create OTEL resource with service identification using common utility
	// This attaches service metadata to all traces for filtering and grouping
	otelResource, resourceError := common.NewOTELResourceFromConfig(
		tracerConfiguration.ServiceName,
		tracerConfiguration.ServiceVersion,
		tracerConfiguration.Environment,
	)

	if resourceError != nil {

		return nil, fmt.Errorf("create resource: %w", resourceError)
	}

	// Collect span processors for the provider
	// Multiple processors can be used (e.g., OTLP batch + stdout debug)
	var spanProcessors []sdktrace.SpanProcessor

	// Create OTLP exporter based on configured protocol (gRPC or HTTP)
	otlpExporter, exporterError := createOTLPTraceExporter(backgroundContext, tracerConfiguration)

	if exporterError != nil {

		return nil, fmt.Errorf("%w: %v", utils.ErrOTELExporterCreationFailed, exporterError)
	}

	// Apply defaults for batch processor settings
	// These values optimize for throughput while limiting memory usage
	maxQueueSize := tracerConfiguration.MaxQueueSize
	if maxQueueSize <= 0 {
		maxQueueSize = int(utils.OTELDefaultMaxQueueSize)
	}

	maxBatchSize := tracerConfiguration.MaxExportBatch
	if maxBatchSize <= 0 {
		maxBatchSize = int(utils.OTELDefaultMaxExportBatch)
	}

	batchTimeout := tracerConfiguration.BatchTimeout
	if batchTimeout <= 0 {
		batchTimeout = utils.OTELDefaultBatchTimeout
	}

	// Create batch span processor for efficient trace export
	// Batching reduces network overhead and improves throughput
	// The processor queues spans and exports them in batches at regular intervals
	batchProcessor := sdktrace.NewBatchSpanProcessor(otlpExporter,
		// Maximum number of spans to queue before blocking
		// Higher values use more memory but handle burst traffic better
		sdktrace.WithMaxQueueSize(maxQueueSize),

		// Maximum number of spans per export batch
		// Larger batches are more efficient but increase latency
		sdktrace.WithMaxExportBatchSize(maxBatchSize),

		// How often to export partial batches and timeout for exports
		sdktrace.WithBatchTimeout(batchTimeout),
	)

	spanProcessors = append(spanProcessors, batchProcessor)

	// Add debug stdout exporter if enabled
	// This is useful for development and troubleshooting
	if tracerConfiguration.OTELDebug {

		stdoutExporter, stdoutError := stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
			stdouttrace.WithPrettyPrint(), // Human-readable JSON output
		)

		if stdoutError != nil {

			return nil, fmt.Errorf("create stdout exporter: %w", stdoutError)
		}

		// Use simple processor for immediate stdout output (no batching)
		// This ensures debug traces appear immediately
		spanProcessors = append(spanProcessors, sdktrace.NewSimpleSpanProcessor(stdoutExporter))
	}

	// Configure sampler based on sampling ratio
	// This determines which traces are recorded and exported
	sampler := createSampler(tracerConfiguration.SamplingRatio)

	// Build provider options starting with resource and sampler
	providerOptions := []sdktrace.TracerProviderOption{
		// Attach resource for service identification in all spans
		sdktrace.WithResource(otelResource),

		// Use parent-based sampling for trace consistency
		// Child spans inherit parent's sampling decision
		sdktrace.WithSampler(sdktrace.ParentBased(sampler)),
	}

	// Add all span processors to the provider
	for _, processor := range spanProcessors {

		providerOptions = append(providerOptions, sdktrace.WithSpanProcessor(processor))
	}

	// Create and return the provider
	tracerProvider := sdktrace.NewTracerProvider(providerOptions...)

	return tracerProvider, nil
}

/* ---------------------------------------- Sampler Configuration -------------------------------------------------- */

// createSampler creates a sampler based on the sampling ratio.
// This determines which traces are recorded and exported.
//
// Sampling strategies:
//   - ratio <= 0: Never sample (no traces exported)
//   - ratio >= 1: Always sample (all traces exported)
//   - 0 < ratio < 1: Probability-based sampling (e.g., 0.1 = 10%)
//
// Parameters:
//   - ratio: Sampling probability (0.0 to 1.0)
//
// Returns:
//   - sdktrace.Sampler: The configured sampler
func createSampler(ratio float64) sdktrace.Sampler {

	switch {

	case ratio <= 0:
		// Never sample - discard all traces
		// Useful when tracing should be completely disabled
		return sdktrace.NeverSample()

	case ratio >= 1:
		// Always sample - record all traces
		// Useful for development or low-traffic services
		return sdktrace.AlwaysSample()

	default:
		// Probability-based sampling using trace ID
		// This ensures consistent sampling across distributed services
		return sdktrace.TraceIDRatioBased(ratio)
	}
}

/* ---------------------------------------- Exporter Factory Functions -------------------------------------------------- */

// createOTLPTraceExporter creates an OTLP exporter based on the configured protocol.
// This factory function selects between gRPC and HTTP transports based on configuration.
//
// Protocol selection:
//   - "grpc": Uses gRPC transport (recommended for high-throughput)
//   - "http" or "http/protobuf": Uses HTTP transport with protobuf encoding
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - tracerConfiguration: Configuration containing OTEL settings
//
// Returns:
//   - sdktrace.SpanExporter: The configured exporter
//   - error: Non-nil if exporter creation fails
func createOTLPTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	// Resolve and validate the protocol string using common utility
	protocol, protocolError := common.ResolveProtocol(tracerConfiguration.OTELProtocol)

	if protocolError != nil {

		return nil, protocolError
	}

	// Create exporter based on resolved protocol
	switch protocol {

	case common.ProtocolGRPC:
		// gRPC transport - high performance, bidirectional streaming
		return createGRPCTraceExporter(ctx, tracerConfiguration)

	case common.ProtocolHTTP:
		// HTTP transport - simpler, firewall-friendly
		return createHTTPTraceExporter(ctx, tracerConfiguration)

	default:
		// Should not reach here if ResolveProtocol works correctly
		return nil, fmt.Errorf("%w: %s", utils.ErrOTELUnsupportedProtocol, protocol)
	}
}

// createGRPCTraceExporter creates a gRPC-based OTLP trace exporter.
// gRPC provides high performance through HTTP/2 multiplexing and binary encoding.
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - tracerConfiguration: Configuration containing endpoint and security settings
//
// Returns:
//   - sdktrace.SpanExporter: The configured gRPC exporter
//   - error: Non-nil if exporter creation fails
func createGRPCTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	// Build exporter options starting with the endpoint
	exporterOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(tracerConfiguration.OTELEndpoint),
	}

	// Add insecure transport credentials if configured
	// WARNING: Only use insecure mode for local development
	if tracerConfiguration.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {

			exporterOptions = append(exporterOptions, otlptracegrpc.WithDialOption(dialOption))
		}
	}

	return otlptracegrpc.New(ctx, exporterOptions...)
}

// createHTTPTraceExporter creates an HTTP-based OTLP trace exporter.
// HTTP transport uses REST conventions with protobuf encoding.
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - tracerConfiguration: Configuration containing endpoint and security settings
//
// Returns:
//   - sdktrace.SpanExporter: The configured HTTP exporter
//   - error: Non-nil if exporter creation fails
func createHTTPTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	// Build exporter options starting with the endpoint
	exporterOptions := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(tracerConfiguration.OTELEndpoint),
	}

	// Add insecure option if configured
	// WARNING: Only use insecure mode for local development
	if tracerConfiguration.OTELInsecure {

		exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
	}

	return otlptracehttp.New(ctx, exporterOptions...)
}
