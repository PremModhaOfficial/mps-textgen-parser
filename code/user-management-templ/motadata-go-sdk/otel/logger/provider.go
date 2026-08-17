package logger

import (
	"context"
	"fmt"
	"os"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

/* ---------------------------------------- Provider Constructor -------------------------------------------------- */

// newOTELProvider creates a new OTEL logger provider with the configured exporters.
// This function sets up the OpenTelemetry logging pipeline including:
//   - Resource identification (service name, version, environment)
//   - OTLP exporter for sending logs to an OTEL collector
//   - Batch processor for efficient log batching
//   - Optional stdout exporter for debugging
//
// The provider uses a batch processor to aggregate logs before export,
// improving performance by reducing network calls and enabling backpressure handling.
//
// Parameters:
//   - loggerConfiguration: The logger configuration containing OTEL settings
//
// Returns:
//   - *sdklog.LoggerProvider: The configured provider
//   - error: Non-nil if provider creation fails
func newOTELProvider(loggerConfiguration Config) (*sdklog.LoggerProvider, error) {

	// Use background context for provider creation
	// This is appropriate as provider setup happens at application startup
	backgroundContext := context.Background()

	// Create OTEL resource with service identification
	// This attaches service metadata to all logs for filtering and grouping
	otelResource, resourceError := common.NewOTELResourceFromConfig(
		loggerConfiguration.ServiceName,
		loggerConfiguration.ServiceVersion,
		loggerConfiguration.Environment,
	)

	if resourceError != nil {

		return nil, fmt.Errorf("create resource: %w", resourceError)
	}

	// Collect processors for the provider
	// Multiple processors can be used (e.g., OTLP batch + stdout debug)
	var processors []sdklog.Processor

	// Create OTLP exporter based on configured protocol (gRPC or HTTP)
	otlpExporter, exporterError := createOTLPLogExporter(backgroundContext, loggerConfiguration)

	if exporterError != nil {

		return nil, fmt.Errorf("%w: %v", utils.ErrOTELExporterCreationFailed, exporterError)
	}

	// Apply defaults for batch processor settings
	// These values optimize for throughput while limiting memory usage
	maxQueueSize := loggerConfiguration.MaxQueueSize
	if maxQueueSize <= 0 {
		maxQueueSize = utils.OTELDefaultMaxQueueSize
	}

	maxBatchSize := loggerConfiguration.MaxExportBatch
	if maxBatchSize <= 0 {
		maxBatchSize = utils.OTELDefaultMaxExportBatch
	}

	exportInterval := loggerConfiguration.ExportInterval
	if exportInterval <= 0 {
		exportInterval = utils.OTELDefaultExportInterval
	}

	batchTimeout := loggerConfiguration.BatchTimeout
	if batchTimeout <= 0 {
		batchTimeout = utils.OTELDefaultBatchTimeout
	}

	// Create batch processor for efficient log export
	// Batching reduces network overhead and improves throughput
	// The processor queues logs and exports them in batches at regular intervals
	processors = append(processors, sdklog.NewBatchProcessor(otlpExporter,
		// Maximum number of logs to queue before blocking
		// Higher values use more memory but handle burst traffic better
		sdklog.WithMaxQueueSize(maxQueueSize),

		// Maximum number of logs per export batch
		// Larger batches are more efficient but increase latency
		sdklog.WithExportMaxBatchSize(maxBatchSize),

		// How often to export partial batches
		// Shorter intervals reduce latency at the cost of efficiency
		sdklog.WithExportInterval(exportInterval),

		// Maximum time to wait for batch export to complete
		// Prevents hung exports from blocking the processor
		sdklog.WithExportTimeout(batchTimeout),
	))

	// Add debug stdout exporter if enabled
	// This is useful for development and troubleshooting
	if loggerConfiguration.OTELDebug {

		stdoutExporter, stdoutError := stdoutlog.New(
			stdoutlog.WithWriter(os.Stdout),
			stdoutlog.WithPrettyPrint(), // Human-readable JSON output
		)

		if stdoutError != nil {

			return nil, fmt.Errorf("create stdout exporter: %w", stdoutError)
		}

		// Use simple processor for immediate stdout output (no batching)
		// This ensures debug logs appear immediately
		processors = append(processors, sdklog.NewSimpleProcessor(stdoutExporter))
	}

	// Build provider options starting with the resource
	providerOptions := []sdklog.LoggerProviderOption{
		sdklog.WithResource(otelResource),
	}

	// Add all processors to the provider
	for _, processor := range processors {

		providerOptions = append(providerOptions, sdklog.WithProcessor(processor))
	}

	// Create and return the provider
	loggerProvider := sdklog.NewLoggerProvider(providerOptions...)

	return loggerProvider, nil
}

/* ---------------------------------------- Exporter Factory Functions -------------------------------------------------- */

// createOTLPLogExporter creates an OTLP exporter based on the configured protocol.
// This factory function selects between gRPC and HTTP transports based on configuration.
//
// Protocol selection:
//   - "grpc": Uses gRPC transport (recommended for high-throughput)
//   - "http" or "http/protobuf": Uses HTTP transport with protobuf encoding
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - loggerConfiguration: Configuration containing OTEL settings
//
// Returns:
//   - sdklog.Exporter: The configured exporter
//   - error: Non-nil if exporter creation fails
func createOTLPLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	// Resolve and validate the protocol string
	protocol, protocolError := common.ResolveProtocol(loggerConfiguration.OTELProtocol)

	if protocolError != nil {

		return nil, protocolError
	}

	// Create exporter based on resolved protocol
	switch protocol {

	case common.ProtocolGRPC:
		// gRPC transport - high performance, bidirectional streaming
		return createGRPCLogExporter(ctx, loggerConfiguration)

	case common.ProtocolHTTP:
		// HTTP transport - simpler, firewall-friendly
		return createHTTPLogExporter(ctx, loggerConfiguration)

	default:
		// Should not reach here if ResolveProtocol works correctly
		return nil, fmt.Errorf("%w: %s", utils.ErrOTELUnsupportedProtocol, protocol)
	}
}

// createGRPCLogExporter creates a gRPC-based OTLP log exporter.
// gRPC provides high performance through HTTP/2 multiplexing and binary encoding.
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - loggerConfiguration: Configuration containing endpoint and security settings
//
// Returns:
//   - sdklog.Exporter: The configured gRPC exporter
//   - error: Non-nil if exporter creation fails
func createGRPCLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	// Build exporter options starting with the endpoint
	exporterOptions := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(loggerConfiguration.OTELEndpoint),
	}

	// Add insecure transport credentials if configured
	// WARNING: Only use insecure mode for local development
	if loggerConfiguration.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {

			exporterOptions = append(exporterOptions, otlploggrpc.WithDialOption(dialOption))
		}
	}

	return otlploggrpc.New(ctx, exporterOptions...)
}

// createHTTPLogExporter creates an HTTP-based OTLP log exporter.
// HTTP transport uses REST conventions with protobuf encoding.
//
// Parameters:
//   - ctx: Context for the exporter creation
//   - loggerConfiguration: Configuration containing endpoint and security settings
//
// Returns:
//   - sdklog.Exporter: The configured HTTP exporter
//   - error: Non-nil if exporter creation fails
func createHTTPLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	// Build exporter options starting with the endpoint
	exporterOptions := []otlploghttp.Option{
		otlploghttp.WithEndpoint(loggerConfiguration.OTELEndpoint),
	}

	// Add insecure option if configured
	// WARNING: Only use insecure mode for local development
	if loggerConfiguration.OTELInsecure {

		exporterOptions = append(exporterOptions, otlploghttp.WithInsecure())
	}

	return otlploghttp.New(ctx, exporterOptions...)
}
