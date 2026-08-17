package utils

import (
	"os"
	"time"
)

var (
	CurrentDir, _ = os.Getwd()
)

const (
	Yes = "yes"

	No = "no"

	PathSeparator = string(os.PathSeparator)

	KeySeparator = "^"

	SpaceSeparator = " "

	NewLineSeparator = "\n"

	NotAvailable = -1

	Empty = ""
)

/* ---------------------------------------- OTEL Constants -------------------------------------------------- */

// OTEL Protocol constants - supported exporter protocols
const (
	// OTELProtocolGRPC represents the gRPC protocol for OTEL exporters
	OTELProtocolGRPC = "grpc"

	// OTELProtocolHTTP represents the HTTP protocol for OTEL exporters
	OTELProtocolHTTP = "http"

	// OTELProtocolHTTPProtobuf represents the HTTP/protobuf protocol for OTEL exporters
	OTELProtocolHTTPProtobuf = "http/protobuf"
)

// OTEL Default configuration values - sensible defaults for batch processing
const (
	// OTELDefaultMaxQueueSize is the default maximum queue size for batch processors
	OTELDefaultMaxQueueSize = 2048

	// OTELDefaultMaxExportBatch is the default maximum batch size for exports
	OTELDefaultMaxExportBatch = 512

	// OTELDefaultExportInterval is the default interval between batch exports
	OTELDefaultExportInterval = time.Second

	// OTELDefaultBatchTimeout is the default timeout for batch exports
	OTELDefaultBatchTimeout = 30 * time.Second

	// OTELDefaultShutdownTimeout is the default timeout for graceful shutdown
	OTELDefaultShutdownTimeout = 10 * time.Second
)

// OTEL Pool size constants - maximum sizes for object pools
const (
	// OTELFieldPoolMaxSize is the maximum size of the field slice pool for logging
	OTELFieldPoolMaxSize = 1024

	// OTELTimerPoolMaxSize is the maximum size of the timer pool for metrics
	OTELTimerPoolMaxSize = 1024

	// OTELLabelsPoolMaxSize is the maximum size of the labels pool for metrics
	OTELLabelsPoolMaxSize = 1024

	// OTELAttributesPoolMaxSize is the maximum size of the attributes pool for metrics
	OTELAttributesPoolMaxSize = 1024

	// OTELMaxContextFields is the maximum number of context fields expected in logs
	OTELMaxContextFields = 6

	// OTELMaxLabels is the maximum number of labels expected in metrics
	OTELMaxLabels = 16
)
