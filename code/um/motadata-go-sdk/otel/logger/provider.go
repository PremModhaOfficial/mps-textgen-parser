package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

/* ---------------------------------------- Provider Constructor -------------------------------------------------- */

// newOTELProvider creates a new OTEL logger provider
func newOTELProvider(loggerConfiguration Config) (*sdklog.LoggerProvider, error) {

	backgroundContext := context.Background()

	otelResource, resourceError := common.NewOTELResourceFromConfig(
		loggerConfiguration.ServiceName,
		loggerConfiguration.ServiceVersion,
		loggerConfiguration.Environment,
	)

	if resourceError != nil {
		return nil, fmt.Errorf("create resource: %w", resourceError)
	}

	var processors []sdklog.Processor

	otlpExporter, exporterError := createOTLPLogExporter(backgroundContext, loggerConfiguration)

	if exporterError != nil {
		return nil, fmt.Errorf("create OTLP exporter: %w", exporterError)
	}

	maxQueueSize := loggerConfiguration.MaxQueueSize

	if maxQueueSize <= 0 {
		maxQueueSize = 2048
	}

	maxBatchSize := loggerConfiguration.MaxExportBatch

	if maxBatchSize <= 0 {
		maxBatchSize = 512
	}

	exportInterval := loggerConfiguration.ExportInterval

	if exportInterval <= 0 {
		exportInterval = time.Second
	}

	batchTimeout := loggerConfiguration.BatchTimeout

	if batchTimeout <= 0 {
		batchTimeout = 30 * time.Second
	}

	processors = append(processors, sdklog.NewBatchProcessor(otlpExporter,
		sdklog.WithMaxQueueSize(maxQueueSize),
		sdklog.WithExportMaxBatchSize(maxBatchSize),
		sdklog.WithExportInterval(exportInterval),
		sdklog.WithExportTimeout(batchTimeout),
	))

	if loggerConfiguration.OTELDebug {
		stdoutExporter, stdoutError := stdoutlog.New(
			stdoutlog.WithWriter(os.Stdout),
			stdoutlog.WithPrettyPrint(),
		)

		if stdoutError != nil {
			return nil, fmt.Errorf("create stdout exporter: %w", stdoutError)
		}

		processors = append(processors, sdklog.NewSimpleProcessor(stdoutExporter))
	}

	providerOptions := []sdklog.LoggerProviderOption{
		sdklog.WithResource(otelResource),
	}

	for _, processor := range processors {
		providerOptions = append(providerOptions, sdklog.WithProcessor(processor))
	}

	loggerProvider := sdklog.NewLoggerProvider(providerOptions...)

	return loggerProvider, nil
}

/* ---------------------------------------- Exporter Factory Functions -------------------------------------------------- */

// createOTLPLogExporter creates an OTLP exporter based on the configured protocol
func createOTLPLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	protocol, protocolError := common.ResolveProtocol(loggerConfiguration.OTELProtocol)

	if protocolError != nil {

		return nil, protocolError
	}

	switch protocol {

	case common.ProtocolGRPC:

		return createGRPCLogExporter(ctx, loggerConfiguration)

	case common.ProtocolHTTP:

		return createHTTPLogExporter(ctx, loggerConfiguration)

	default:

		return nil, fmt.Errorf("unsupported OTEL protocol: %s", protocol)
	}
}

// createGRPCLogExporter creates a gRPC-based OTLP log exporter
func createGRPCLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	exporterOptions := []otlploggrpc.Option{
		otlploggrpc.WithEndpoint(loggerConfiguration.OTELEndpoint),
	}

	if loggerConfiguration.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {

			exporterOptions = append(exporterOptions, otlploggrpc.WithDialOption(dialOption))
		}
	}

	return otlploggrpc.New(ctx, exporterOptions...)
}

// createHTTPLogExporter creates an HTTP-based OTLP log exporter
func createHTTPLogExporter(ctx context.Context, loggerConfiguration Config) (sdklog.Exporter, error) {

	exporterOptions := []otlploghttp.Option{
		otlploghttp.WithEndpoint(loggerConfiguration.OTELEndpoint),
	}

	if loggerConfiguration.OTELInsecure {

		exporterOptions = append(exporterOptions, otlploghttp.WithInsecure())
	}

	return otlploghttp.New(ctx, exporterOptions...)
}
