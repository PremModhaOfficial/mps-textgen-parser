package tracer

import (
	"context"
	"fmt"
	"os"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

/* ---------------------------------------- Provider Constructor -------------------------------------------------- */

// newTracerProvider creates a new OTEL tracer provider.
func newTracerProvider(tracerConfiguration Config) (*sdktrace.TracerProvider, error) {

	backgroundContext := context.Background()

	// Create resource using common utility
	otelResource, resourceError := common.NewOTELResourceFromConfig(
		tracerConfiguration.ServiceName,
		tracerConfiguration.ServiceVersion,
		tracerConfiguration.Environment,
	)

	if resourceError != nil {

		return nil, fmt.Errorf("create resource: %w", resourceError)
	}

	// Collect span processors
	var spanProcessors []sdktrace.SpanProcessor

	// Create OTLP exporter based on protocol
	otlpExporter, exporterError := createOTLPTraceExporter(backgroundContext, tracerConfiguration)

	if exporterError != nil {

		return nil, fmt.Errorf("create OTLP exporter: %w", exporterError)
	}

	// Add batch span processor for OTLP exporter
	batchProcessor := sdktrace.NewBatchSpanProcessor(otlpExporter,
		sdktrace.WithMaxQueueSize(tracerConfiguration.MaxQueueSize),
		sdktrace.WithMaxExportBatchSize(tracerConfiguration.MaxExportBatch),
		sdktrace.WithBatchTimeout(tracerConfiguration.BatchTimeout),
	)
	spanProcessors = append(spanProcessors, batchProcessor)

	// Add stdout exporter for debugging if enabled
	if tracerConfiguration.OTELDebug {

		stdoutExporter, stdoutError := stdouttrace.New(
			stdouttrace.WithWriter(os.Stdout),
			stdouttrace.WithPrettyPrint(),
		)

		if stdoutError != nil {

			return nil, fmt.Errorf("create stdout exporter: %w", stdoutError)
		}

		// Use simple processor for immediate stdout output
		spanProcessors = append(spanProcessors, sdktrace.NewSimpleSpanProcessor(stdoutExporter))
	}

	// Configure sampler
	var sampler sdktrace.Sampler

	switch {

	case tracerConfiguration.SamplingRatio <= 0:

		sampler = sdktrace.NeverSample()

	case tracerConfiguration.SamplingRatio >= 1:

		sampler = sdktrace.AlwaysSample()

	default:

		sampler = sdktrace.TraceIDRatioBased(tracerConfiguration.SamplingRatio)
	}

	// Build provider options
	providerOptions := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(otelResource),
		sdktrace.WithSampler(sdktrace.ParentBased(sampler)),
	}

	for _, processor := range spanProcessors {

		providerOptions = append(providerOptions, sdktrace.WithSpanProcessor(processor))
	}

	// Create provider
	tracerProvider := sdktrace.NewTracerProvider(providerOptions...)

	return tracerProvider, nil
}

/* ---------------------------------------- Exporter Factory Functions -------------------------------------------------- */

// createOTLPTraceExporter creates an OTLP exporter based on the configured protocol.
func createOTLPTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	protocol, protocolError := common.ResolveProtocol(tracerConfiguration.OTELProtocol)

	if protocolError != nil {

		return nil, protocolError
	}

	switch protocol {

	case common.ProtocolGRPC:

		return createGRPCTraceExporter(ctx, tracerConfiguration)

	case common.ProtocolHTTP:

		return createHTTPTraceExporter(ctx, tracerConfiguration)

	default:

		return nil, fmt.Errorf("unsupported OTEL protocol: %s", protocol)
	}
}

// createGRPCTraceExporter creates a gRPC-based OTLP trace exporter.
func createGRPCTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	exporterOptions := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(tracerConfiguration.OTELEndpoint),
	}

	if tracerConfiguration.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {

			exporterOptions = append(exporterOptions, otlptracegrpc.WithDialOption(dialOption))
		}
	}

	return otlptracegrpc.New(ctx, exporterOptions...)
}

// createHTTPTraceExporter creates an HTTP-based OTLP trace exporter.
func createHTTPTraceExporter(ctx context.Context, tracerConfiguration Config) (sdktrace.SpanExporter, error) {

	exporterOptions := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(tracerConfiguration.OTELEndpoint),
	}

	if tracerConfiguration.OTELInsecure {

		exporterOptions = append(exporterOptions, otlptracehttp.WithInsecure())
	}

	return otlptracehttp.New(ctx, exporterOptions...)
}
