# OTEL Common Package

This package provides shared utilities and types used across the OpenTelemetry (OTEL) integration modules (logger, metrics, tracer).

## Overview

The `common` package contains:

- **Protocol Resolution**: Utilities for resolving OTEL exporter protocols (gRPC/HTTP)
- **Resource Creation**: Factory functions for creating OTEL resources with service identification
- **Lifecycle Management**: Error aggregation and shutdown coordination utilities
- **gRPC Configuration**: Helper functions for gRPC dial options

## Package Structure

```
common/
├── exporter.go     # Protocol resolution and gRPC configuration
├── lifecycle.go    # Shutdown error aggregation
├── resource.go     # OTEL resource creation
├── common_test.go  # Unit and benchmark tests
└── README.md       # This file
```

## Components

### Protocol Resolution

Resolves and normalizes OTEL exporter protocol strings:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

// Resolve protocol from configuration
protocol, err := common.ResolveProtocol("grpc")  // Returns ProtocolGRPC
protocol, err := common.ResolveProtocol("http")  // Returns ProtocolHTTP
protocol, err := common.ResolveProtocol("")      // Defaults to ProtocolGRPC
```

**Supported Protocols:**
- `grpc` - gRPC protocol (default)
- `http`, `http/protobuf` - HTTP protocol

### OTEL Resource Creation

Creates standardized OTEL resources with service identification:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

// Create resource with service info
serviceInfo := common.ServiceInfo{
    Name:        "my-service",
    Version:     "1.0.0",
    Environment: "production",
}

resource, err := common.NewOTELResource(serviceInfo)

// Or use the convenience function
resource, err := common.NewOTELResourceFromConfig("my-service", "1.0.0", "production")
```

### Shutdown Error Aggregation

Coordinates shutdown of multiple components and aggregates errors:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

// Aggregate multiple errors
errors := []error{err1, err2, err3}
aggregatedErr := common.AggregateErrors(errors)

// Or use ShutdownCollector for cleaner code
collector := common.NewShutdownCollector()
collector.Collect(tracer.Shutdown(ctx))
collector.Collect(metrics.Shutdown(ctx))
collector.Collect(logger.Close())

if collector.HasErrors() {
    return collector.Error()
}
```

### gRPC Configuration

Helper functions for configuring gRPC connections:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"

// Get dial options for insecure connection
options := common.GRPCDialOptions(true)  // Returns []grpc.DialOption

// Get single dial option
option := common.GRPCDialOption(true)    // Returns grpc.DialOption
```

## Types

### ExporterProtocol

```go
type ExporterProtocol string

const (
    ProtocolGRPC ExporterProtocol = "grpc"
    ProtocolHTTP ExporterProtocol = "http"
)
```

### ServiceInfo

```go
type ServiceInfo struct {
    Name        string  // Service name
    Version     string  // Service version
    Environment string  // Deployment environment
}
```

### ShutdownCollector

```go
type ShutdownCollector struct {
    // ... internal fields
}

// Methods:
func NewShutdownCollector() *ShutdownCollector
func (c *ShutdownCollector) Collect(err error)
func (c *ShutdownCollector) CollectFunc(fn func() error)
func (c *ShutdownCollector) Error() error
func (c *ShutdownCollector) HasErrors() bool
func (c *ShutdownCollector) Count() int
```

## Testing

Run tests with coverage:

```bash
go test -cover ./otel/common/...
```

Run benchmarks:

```bash
go test -bench=. ./otel/common/...
```

## Dependencies

- `go.opentelemetry.io/otel/sdk/resource` - OTEL resource SDK
- `go.opentelemetry.io/otel/semconv` - Semantic conventions
- `google.golang.org/grpc` - gRPC framework

## Usage in Other Packages

This package is used internally by:
- `otel/logger` - For resource creation and protocol resolution
- `otel/metrics` - For resource creation and protocol resolution
- `otel/tracer` - For resource creation and protocol resolution
- `otel` (main) - For shutdown coordination

External usage is supported but typically not necessary as the main OTEL packages provide higher-level abstractions.
