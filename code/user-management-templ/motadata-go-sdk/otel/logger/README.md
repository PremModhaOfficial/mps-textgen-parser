# Logger Package

A production-ready, structured logging package built on [Uber's Zap](https://github.com/uber-go/zap) with OpenTelemetry integration for distributed tracing correlation.

## Table of Contents

- [Cross-Module Usage](#cross-module-usage)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
  - [Configuration Options](#configuration-options)
  - [Configuration Methods](#configuration-methods)
  - [Environment Variables](#environment-variables)
  - [YAML Configuration](#yaml-configuration)
- [Usage](#usage)
  - [Basic Logging](#basic-logging)
  - [Structured Fields](#structured-fields)
  - [Context-Aware Logging](#context-aware-logging)
  - [Named Loggers](#named-loggers)
  - [Module Loggers](#module-loggers)
  - [Dynamic Level Changes](#dynamic-level-changes)
- [Global Logger](#global-logger)
- [OpenTelemetry Integration](#opentelemetry-integration)
- [File Logging](#file-logging)
- [Best Practices](#best-practices)
- [API Reference](#api-reference)

---

## Cross-Module Usage

The logger is designed for use across multiple packages in your service. **Initialize once, use everywhere.**

```
myapp/
├── main.go              # Initialize logger here
├── handlers/
│   └── order.go         # Use logger.Module("http") or logger.Info()
├── services/
│   └── order.go         # Use logger.Module("order-service")
└── repository/
    └── order.go         # Use logger.Module("database")
```

```go
// main.go - Initialize once
func main() {
    logger.Init(cfg)
    defer logger.Close()
    // ...
}

// Any other package - Just import and use
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"

var log = logger.Module("mypackage")  // or use logger.Info() directly

func DoSomething(ctx context.Context) {
    log.Info(ctx, "doing something")
}
```

---

## Features

- **High Performance**: Built on Zap, one of the fastest logging libraries for Go
- **Structured Logging**: JSON-formatted logs with typed fields
- **Multiple Outputs**: Console, file, and OpenTelemetry exporters
- **Context Propagation**: Automatic extraction of trace IDs, tenant IDs, and request IDs
- **Module-Level Logging**: Per-module log levels with dynamic runtime updates
- **Log Rotation**: Built-in file rotation with compression support
- **OpenTelemetry Bridge**: Seamless integration with OTEL collectors for observability

---

## Installation

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
```

---

## Quick Start

### Minimal Setup

```go
package main

import (
    "context"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

func main() {
    // Initialize with defaults
    cfg := logger.DefaultConfig()
    cfg.ServiceName = "my-service"
    cfg.ServiceVersion = "1.0.0"
    cfg.Environment = "production"

    log, err := logger.Init(cfg)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    // Start logging
    ctx := context.Background()
    logger.Info(ctx, "Application started", logger.String("port", "8080"))
}
```

### Using Environment Variables

```go
package main

import (
    "context"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

func main() {
    // Load configuration from environment variables
    log, err := logger.InitFromEnv()
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    ctx := context.Background()
    logger.Info(ctx, "Application started")
}
```

---

## Configuration

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Level` | string | `"info"` | Minimum log level: `debug`, `info`, `warn`, `error`, `fatal` |
| `ServiceName` | string | `"app"` | Service identifier added to all logs |
| `ServiceVersion` | string | `"0.0.0"` | Service version added to all logs |
| `Environment` | string | `"development"` | Deployment environment (development, staging, production) |
| `ConsoleEnabled` | bool | `true` | Enable console output |
| `ConsoleFormat` | string | `"json"` | Console format: `json` or `console` |
| `OTELEnabled` | bool | `false` | Enable OpenTelemetry log export |
| `OTELEndpoint` | string | `"localhost:4317"` | OTEL collector gRPC endpoint |
| `OTELInsecure` | bool | `true` | Use insecure gRPC connection |
| `FileEnabled` | bool | `false` | Enable file logging |
| `FilePath` | string | `"logs/app.log"` | Log file path |
| `FileMaxSizeMB` | int | `100` | Maximum log file size in MB before rotation |
| `FileMaxBackups` | int | `5` | Number of rotated files to keep |
| `FileMaxAgeDays` | int | `7` | Days to retain old log files |
| `FileCompress` | bool | `true` | Compress rotated log files |
| `AddCaller` | bool | `true` | Include caller file and line number |
| `CallerSkip` | int | `2` | Stack frames to skip for caller info |
| `ModuleLevels` | map[string]string | `nil` | Per-module log level overrides |

### Configuration Methods

#### 1. Programmatic Configuration

```go
cfg := logger.DefaultConfig()
cfg.Level = "debug"
cfg.ServiceName = "order-service"
cfg.ServiceVersion = "2.1.0"
cfg.Environment = "production"
cfg.ConsoleFormat = "console" // Human-readable for development
cfg.FileEnabled = true
cfg.FilePath = "/var/log/order-service/app.log"
cfg.ModuleLevels = map[string]string{
    "database": "debug",
    "http":     "info",
    "cache":    "warn",
}

log, err := logger.Init(cfg)
```

#### 2. From YAML File

```go
log, err := logger.InitFromConfig("/etc/myapp/config.yaml")
if err != nil {
    panic(err)
}
```

#### 3. From Environment Variables

```go
log, err := logger.InitFromEnv()
```

#### 4. From Unified Config

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"

cfg, err := config.Load("/etc/myapp/config.yaml")
if err != nil {
    panic(err)
}

log, err := logger.InitFromUnifiedConfig(cfg)
```

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `LOG_LEVEL` | Minimum log level | `debug` |
| `LOG_SERVICE_NAME` | Service name | `my-service` |
| `LOG_SERVICE_VERSION` | Service version | `1.0.0` |
| `LOG_ENVIRONMENT` | Environment | `production` |
| `LOG_CONSOLE_ENABLED` | Enable console output | `true` |
| `LOG_CONSOLE_FORMAT` | Console format | `json` |
| `LOG_OTEL_ENABLED` | Enable OTEL export | `true` |
| `LOG_OTEL_ENDPOINT` | OTEL collector endpoint | `otel-collector:4317` |
| `LOG_OTEL_INSECURE` | Insecure gRPC | `false` |
| `LOG_FILE_ENABLED` | Enable file logging | `true` |
| `LOG_FILE_PATH` | Log file path | `/var/log/app.log` |
| `LOG_FILE_MAX_SIZE_MB` | Max file size | `100` |
| `LOG_FILE_MAX_BACKUPS` | Backup count | `5` |
| `LOG_FILE_MAX_AGE_DAYS` | Retention days | `7` |
| `LOG_FILE_COMPRESS` | Compress old logs | `true` |

### YAML Configuration

```yaml
service:
  name: order-service
  version: 1.0.0
  environment: production

logger:
  level: info
  console_enabled: true
  console_format: json

  otel_enabled: true
  otel_endpoint: otel-collector:4317
  otel_insecure: false

  file_enabled: true
  file_path: /var/log/order-service/app.log
  file_max_size_mb: 100
  file_max_backups: 5
  file_max_age_days: 7
  file_compress: true

  add_caller: true

  module_levels:
    database: debug
    http: info
    cache: warn
```

---

## Usage

### Basic Logging

```go
ctx := context.Background()

// Different log levels
logger.Debug(ctx, "Debugging information")
logger.Info(ctx, "Application event occurred")
logger.Warn(ctx, "Warning: resource usage high")
logger.Error(ctx, "Error processing request")
// logger.Fatal(ctx, "Fatal error") // Logs and calls os.Exit(1)
```

### Structured Fields

Use typed field constructors for optimal performance:

```go
ctx := context.Background()

logger.Info(ctx, "User logged in",
    logger.String("user_id", "usr_12345"),
    logger.String("email", "user@example.com"),
    logger.String("ip_address", "192.168.1.100"),
    logger.Duration("session_duration", 30*time.Minute),
)

logger.Info(ctx, "Order processed",
    logger.String("order_id", "ord_67890"),
    logger.Int("items", 5),
    logger.Float64("total", 149.99),
    logger.Bool("express_shipping", true),
    logger.Time("created_at", time.Now()),
)

// Error with stack trace
err := processOrder()
if err != nil {
    logger.Error(ctx, "Order processing failed",
        logger.String("order_id", "ord_67890"),
        logger.Err(err),
    )
}
```

**Available Field Constructors:**

| Function | Type | Example |
|----------|------|---------|
| `String(key, value)` | string | `logger.String("user_id", "123")` |
| `Strings(key, values)` | []string | `logger.Strings("tags", []string{"a", "b"})` |
| `Int(key, value)` | int | `logger.Int("count", 42)` |
| `Int64(key, value)` | int64 | `logger.Int64("bytes", 1024)` |
| `Float64(key, value)` | float64 | `logger.Float64("price", 19.99)` |
| `Bool(key, value)` | bool | `logger.Bool("active", true)` |
| `Time(key, value)` | time.Time | `logger.Time("created", time.Now())` |
| `Duration(key, value)` | time.Duration | `logger.Duration("elapsed", 5*time.Second)` |
| `Err(error)` | error | `logger.Err(err)` |
| `Any(key, value)` | interface{} | `logger.Any("data", struct{}{})` |
| `Binary(key, value)` | []byte | `logger.Binary("payload", data)` |

### Context-Aware Logging

The logger automatically extracts context values and OpenTelemetry trace information:

```go
// Add context values
ctx := context.Background()
ctx = logger.WithTenantID(ctx, "tenant-123")
ctx = logger.WithRequestID(ctx, "req-456")
ctx = logger.WithUserID(ctx, "user-789")
ctx = logger.WithCorrelationID(ctx, "corr-abc")

// These values are automatically included in all logs
logger.Info(ctx, "Processing request",
    logger.String("action", "create_order"),
)
// Output includes: tenant_id, request_id, user_id, correlation_id
```

**Context Functions:**

| Function | Description |
|----------|-------------|
| `WithTenantID(ctx, id)` | Add tenant identifier |
| `WithRequestID(ctx, id)` | Add request identifier |
| `WithUserID(ctx, id)` | Add user identifier |
| `WithCorrelationID(ctx, id)` | Add correlation identifier |

When used with OpenTelemetry tracing, `trace_id` and `span_id` are automatically extracted from the context.

### Named Loggers

Create named sub-loggers for different components:

```go
// Create named loggers
dbLogger := logger.Named("database")
httpLogger := logger.Named("http")
cacheLogger := logger.Named("cache")

ctx := context.Background()

dbLogger.Info(ctx, "Query executed",
    logger.String("query", "SELECT * FROM users"),
    logger.Duration("duration", 15*time.Millisecond),
)

httpLogger.Info(ctx, "Request received",
    logger.String("method", "POST"),
    logger.String("path", "/api/orders"),
)
```

### Module Loggers

Module loggers provide per-module log level control with dynamic updates:

```go
// Create module loggers
dbLogger := logger.Module("database")
httpLogger := logger.Module("http")

ctx := context.Background()

// Log at different levels
dbLogger.Debug(ctx, "Executing query") // Only logs if database module is at debug level
dbLogger.Info(ctx, "Connection established")

httpLogger.Info(ctx, "Serving request")
httpLogger.Warn(ctx, "Slow response time")

// Dynamic level changes at runtime
logger.SetModuleLevel("database", logger.DebugLevel)
logger.SetModuleLevel("http", logger.WarnLevel)

// Check current module level
currentLevel := logger.ModuleLevel("database")

// Set multiple module levels at once
logger.SetModuleLevels(map[string]logger.Level{
    "database": logger.DebugLevel,
    "http":     logger.InfoLevel,
    "cache":    logger.WarnLevel,
})

// List all configured module levels
levels := logger.ListModuleLevels()
for module, level := range levels {
    fmt.Printf("Module %s: %s\n", module, level)
}

// Reset all module levels to default
logger.ResetModuleLevels()
```

**Module Logger Methods:**

```go
moduleLog := logger.Module("mymodule")

// All standard log methods
moduleLog.Debug(ctx, "message", fields...)
moduleLog.Info(ctx, "message", fields...)
moduleLog.Warn(ctx, "message", fields...)
moduleLog.Error(ctx, "message", fields...)
moduleLog.Fatal(ctx, "message", fields...)

// Add preset fields
enrichedLog := moduleLog.With(logger.String("component", "worker"))

// Get/Set level for this specific module
moduleLog.SetLevel(logger.DebugLevel)
currentLevel := moduleLog.Level()
moduleName := moduleLog.ModuleName()
```

### Dynamic Level Changes

Change log levels at runtime without restarting:

```go
// Change global log level
logger.SetGlobalLevel(logger.DebugLevel)

// Change level on a specific logger instance
log := logger.L()
log.SetLevel(logger.WarnLevel)

// Get current level
currentLevel := log.Level()
fmt.Println("Current level:", currentLevel.String())
```

---

## Global Logger

The package provides a global logger instance for convenience:

```go
// Initialize global logger
logger.Init(cfg)

// Access global logger
log := logger.L()

// Package-level functions use the global logger
logger.Info(ctx, "message")
logger.Debug(ctx, "message")
logger.Warn(ctx, "message")
logger.Error(ctx, "message")

// Create derived loggers from global
namedLog := logger.Named("component")
moduleLog := logger.Module("component")
enrichedLog := logger.With(logger.String("key", "value"))

// Replace global logger (useful for testing)
restore := logger.ReplaceGlobal(newLogger)
defer restore() // Restore original logger

// Sync and close
logger.Sync() // Flush buffered logs
logger.Close() // Shutdown and cleanup
```

---

## OpenTelemetry Integration

Enable OTEL export to send logs to an OpenTelemetry Collector:

```go
cfg := logger.DefaultConfig()
cfg.ServiceName = "my-service"
cfg.OTELEnabled = true
cfg.OTELEndpoint = "otel-collector:4317"
cfg.OTELInsecure = false // Use TLS in production

log, err := logger.Init(cfg)
if err != nil {
    panic(err)
}
defer logger.Close() // Important: ensures logs are flushed

// Logs are automatically exported to the OTEL collector
ctx := context.Background()
logger.Info(ctx, "This log is sent to OTEL collector")
```

**Trace Correlation:**

When using OpenTelemetry tracing, trace and span IDs are automatically extracted:

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

// Start a trace
tracer := otel.Tracer("my-service")
ctx, span := tracer.Start(context.Background(), "operation-name")
defer span.End()

// Logs automatically include trace_id and span_id
logger.Info(ctx, "Processing within trace")
// Output: {"trace_id": "abc123...", "span_id": "def456...", "message": "Processing within trace"}
```

---

## File Logging

Enable rotating file logs:

```go
cfg := logger.DefaultConfig()
cfg.FileEnabled = true
cfg.FilePath = "/var/log/myapp/app.log"
cfg.FileMaxSizeMB = 100     // Rotate when file reaches 100MB
cfg.FileMaxBackups = 5      // Keep 5 old log files
cfg.FileMaxAgeDays = 7      // Delete files older than 7 days
cfg.FileCompress = true     // Compress rotated files with gzip

log, err := logger.Init(cfg)
```

**Log Rotation Behavior:**
- Files are rotated when they reach `FileMaxSizeMB`
- Rotated files are named: `app-2024-01-15T10-30-00.log`
- Compressed files: `app-2024-01-15T10-30-00.log.gz`
- Old files are deleted based on `FileMaxBackups` and `FileMaxAgeDays`

---

## Use Cases

### 1. HTTP Request Logging

```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        ctx = logger.WithRequestID(ctx, r.Header.Get("X-Request-ID"))
        ctx = logger.WithTenantID(ctx, r.Header.Get("X-Tenant-ID"))

        start := time.Now()

        logger.Info(ctx, "Request started",
            logger.String("method", r.Method),
            logger.String("path", r.URL.Path),
            logger.String("remote_addr", r.RemoteAddr),
        )

        // Wrap response writer to capture status
        wrapped := &responseWriter{ResponseWriter: w, status: 200}
        next.ServeHTTP(wrapped, r.WithContext(ctx))

        logger.Info(ctx, "Request completed",
            logger.Int("status", wrapped.status),
            logger.Duration("duration", time.Since(start)),
        )
    })
}
```

### 2. Database Query Logging

```go
var dbLog = logger.Module("database")

func QueryUsers(ctx context.Context, filters map[string]string) ([]User, error) {
    start := time.Now()

    dbLog.Debug(ctx, "Executing query",
        logger.String("table", "users"),
        logger.Any("filters", filters),
    )

    users, err := db.Query(ctx, "SELECT * FROM users WHERE ...")

    if err != nil {
        dbLog.Error(ctx, "Query failed",
            logger.String("table", "users"),
            logger.Duration("duration", time.Since(start)),
            logger.Err(err),
        )
        return nil, err
    }

    dbLog.Debug(ctx, "Query succeeded",
        logger.String("table", "users"),
        logger.Int("rows", len(users)),
        logger.Duration("duration", time.Since(start)),
    )

    return users, nil
}
```

### 3. Background Job Logging

```go
var jobLog = logger.Module("jobs")

func ProcessJob(ctx context.Context, job Job) error {
    jobCtx := logger.WithCorrelationID(ctx, job.ID)

    jobLog.Info(jobCtx, "Job started",
        logger.String("job_type", job.Type),
        logger.String("job_id", job.ID),
    )

    for i, step := range job.Steps {
        jobLog.Debug(jobCtx, "Processing step",
            logger.Int("step", i+1),
            logger.Int("total", len(job.Steps)),
        )

        if err := processStep(jobCtx, step); err != nil {
            jobLog.Error(jobCtx, "Step failed",
                logger.Int("step", i+1),
                logger.Err(err),
            )
            return err
        }
    }

    jobLog.Info(jobCtx, "Job completed",
        logger.String("job_id", job.ID),
    )
    return nil
}
```

### 4. Error Handling with Stack Traces

```go
func ProcessPayment(ctx context.Context, payment Payment) error {
    if err := validatePayment(payment); err != nil {
        logger.Warn(ctx, "Payment validation failed",
            logger.String("payment_id", payment.ID),
            logger.Err(err),
        )
        return fmt.Errorf("validation failed: %w", err)
    }

    if err := chargeCard(ctx, payment); err != nil {
        logger.Error(ctx, "Payment processing failed",
            logger.String("payment_id", payment.ID),
            logger.Float64("amount", payment.Amount),
            logger.String("currency", payment.Currency),
            logger.Err(err),
        )
        return fmt.Errorf("charge failed: %w", err)
    }

    logger.Info(ctx, "Payment processed successfully",
        logger.String("payment_id", payment.ID),
        logger.Float64("amount", payment.Amount),
    )
    return nil
}
```

### 5. Multi-Tenant Application Logging

```go
func HandleTenantRequest(ctx context.Context, tenantID string, req Request) {
    // Add tenant context for all subsequent logs
    ctx = logger.WithTenantID(ctx, tenantID)
    ctx = logger.WithRequestID(ctx, req.ID)

    // All logs automatically include tenant_id and request_id
    logger.Info(ctx, "Processing tenant request")

    // Pass context to all downstream functions
    processOrder(ctx, req.Order)
    updateInventory(ctx, req.Items)
    sendNotification(ctx, req.UserID)
}
```

### 6. Debugging with Dynamic Log Levels

```go
// In production, enable debug logging for specific modules via API
func enableDebugLogging(module string) {
    logger.SetModuleLevel(module, logger.DebugLevel)
    logger.Info(context.Background(), "Debug logging enabled",
        logger.String("module", module),
    )
}

// Disable when done
func disableDebugLogging(module string) {
    logger.SetModuleLevel(module, logger.InfoLevel)
}

// Use in admin endpoint
func handleDebugEndpoint(w http.ResponseWriter, r *http.Request) {
    module := r.URL.Query().Get("module")
    enable := r.URL.Query().Get("enable") == "true"

    if enable {
        enableDebugLogging(module)
    } else {
        disableDebugLogging(module)
    }
}
```

---

## Good Practices vs Bad Practices

### Structured Logging

```go
// GOOD - Structured fields are searchable and parseable
logger.Info(ctx, "User logged in",
    logger.String("user_id", userID),
    logger.String("email", email),
    logger.String("ip", ipAddress),
    logger.String("user_agent", userAgent),
)

// BAD - String interpolation loses structure
logger.Info(ctx, fmt.Sprintf("User %s (%s) logged in from %s using %s",
    userID, email, ipAddress, userAgent))
```

### Context Propagation

```go
// GOOD - Pass context through the call chain
func HandleRequest(ctx context.Context, req *Request) {
    ctx = logger.WithRequestID(ctx, req.ID)
    processOrder(ctx, req.Order)
}

func processOrder(ctx context.Context, order *Order) {
    logger.Info(ctx, "Processing order")  // request_id is included
    saveToDatabase(ctx, order)
}

// BAD - Creating new context loses correlation
func processOrder(ctx context.Context, order *Order) {
    logger.Info(context.Background(), "Processing order")  // Lost request_id!
}
```

### Log Levels

```go
// GOOD - Appropriate levels for different scenarios
logger.Debug(ctx, "Cache lookup", logger.String("key", key))           // Development only
logger.Info(ctx, "Order created", logger.String("order_id", orderID))  // Normal operations
logger.Warn(ctx, "Retry attempt", logger.Int("attempt", 3))            // Concerning but handled
logger.Error(ctx, "Payment failed", logger.Err(err))                   // Errors requiring attention

// BAD - Everything at INFO level
logger.Info(ctx, "Debug: cache lookup")     // Should be Debug
logger.Info(ctx, "ERROR: payment failed")   // Should be Error
```

### Error Context

```go
// GOOD - Include all relevant context for debugging
if err != nil {
    logger.Error(ctx, "Failed to process order",
        logger.String("order_id", orderID),
        logger.String("customer_id", customerID),
        logger.Float64("amount", amount),
        logger.String("payment_method", method),
        logger.Err(err),
    )
}

// BAD - Minimal context makes debugging difficult
if err != nil {
    logger.Error(ctx, "Error occurred", logger.Err(err))
}
```

### Module Loggers

```go
// GOOD - Create module loggers at package level
var (
    dbLog    = logger.Module("database")
    cacheLog = logger.Module("cache")
    httpLog  = logger.Module("http")
)

func queryDatabase(ctx context.Context) {
    dbLog.Debug(ctx, "Executing query")  // Can enable/disable per module
}

// BAD - Creating logger on every call
func queryDatabase(ctx context.Context) {
    log := logger.Module("database")  // Creates new instance every time
    log.Debug(ctx, "Executing query")
}
```

### Field Reuse with With()

```go
// GOOD - Create logger with common fields once
orderLogger := logger.L().With(
    logger.String("component", "order-processor"),
    logger.String("version", "2.0"),
)

func processOrder(ctx context.Context, order *Order) {
    // These fields are included in every log
    orderLogger.Info(ctx, "Processing started", logger.String("order_id", order.ID))
    orderLogger.Info(ctx, "Validation complete")
    orderLogger.Info(ctx, "Order saved")
}

// BAD - Repeating fields in every log call
func processOrder(ctx context.Context, order *Order) {
    logger.Info(ctx, "Processing started",
        logger.String("component", "order-processor"),
        logger.String("version", "2.0"),
        logger.String("order_id", order.ID))
    logger.Info(ctx, "Validation complete",
        logger.String("component", "order-processor"),  // Repetitive!
        logger.String("version", "2.0"))
}
```

---

## Performance Pitfalls

### 1. Expensive Field Construction

```go
// BAD - JSON marshaling happens even if debug is disabled
logger.Debug(ctx, "Request body",
    logger.Any("body", expensiveJSONMarshal(largeStruct)),
)

// GOOD - Check level before expensive operations
if logger.L().Level() <= logger.DebugLevel {
    logger.Debug(ctx, "Request body",
        logger.Any("body", expensiveJSONMarshal(largeStruct)),
    )
}

// BETTER - Use lazy evaluation with a custom field
logger.Debug(ctx, "Request body",
    logger.String("body_preview", largeStruct.Summary()),  // Cheap operation
)
```

### 2. High-Frequency Logging

```go
// BAD - Logging every iteration in a hot loop
for i := 0; i < 1000000; i++ {
    logger.Debug(ctx, "Processing item", logger.Int("index", i))
    processItem(items[i])
}

// GOOD - Log summaries or sample
for i := 0; i < 1000000; i++ {
    processItem(items[i])
}
logger.Info(ctx, "Batch processed", logger.Int("count", 1000000))

// OR - Sample logging for debugging
for i := 0; i < 1000000; i++ {
    if i % 10000 == 0 {
        logger.Debug(ctx, "Processing progress", logger.Int("index", i))
    }
    processItem(items[i])
}
```

### 3. Synchronous File I/O

```go
// Configuration concern - File logging can block
cfg := logger.DefaultConfig()
cfg.FileEnabled = true
cfg.FilePath = "/slow-nfs-mount/app.log"  // Network filesystem = slow

// GOOD - Use local fast storage for logs
cfg.FilePath = "/var/log/app/app.log"  // Local SSD

// OR - Disable file logging and use OTEL collector
cfg.FileEnabled = false
cfg.OTELEnabled = true
```

### 4. Large Log Messages

```go
// BAD - Logging large payloads
logger.Debug(ctx, "Received response",
    logger.String("body", string(largeResponseBody)),  // Could be MB of data
)

// GOOD - Log only what's needed for debugging
logger.Debug(ctx, "Received response",
    logger.Int("size_bytes", len(largeResponseBody)),
    logger.String("content_type", resp.Header.Get("Content-Type")),
    logger.String("body_preview", truncate(string(largeResponseBody), 500)),
)

func truncate(s string, maxLen int) string {
    if len(s) <= maxLen {
        return s
    }
    return s[:maxLen] + "...[truncated]"
}
```

### 5. Missing Close/Sync

```go
// BAD - Logs may be lost on crash
func main() {
    logger.Init(cfg)
    // Application runs...
    // Program exits without flushing
}

// GOOD - Always defer Close
func main() {
    logger.Init(cfg)
    defer logger.Close()  // Flushes all buffered logs

    // Handle graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        logger.Info(context.Background(), "Shutting down")
        logger.Sync()  // Explicit sync before shutdown
    }()

    // Application runs...
}
```

### 6. Nil Context Panic

```go
// BAD - Can panic with nil context
logger.Info(nil, "message")  // PANIC!

// GOOD - Always use valid context
logger.Info(context.Background(), "message")

// OR - Use context.TODO() when context will be added later
logger.Info(context.TODO(), "message")
```

---

## Best Practices

### 1. Initialize Early, Close Properly

```go
func main() {
    log, err := logger.Init(cfg)
    if err != nil {
        panic(err)
    }
    defer logger.Close() // Ensures all logs are flushed

    // Application code...
}
```

### 2. Use Structured Fields Instead of String Formatting

```go
// Good - structured and searchable
logger.Info(ctx, "User created",
    logger.String("user_id", userID),
    logger.String("email", email),
)

// Avoid - harder to search and parse
logger.Info(ctx, fmt.Sprintf("User %s created with email %s", userID, email))
```

### 3. Include Context in All Logs

```go
func handleRequest(ctx context.Context, req *Request) {
    // Add request-scoped values to context
    ctx = logger.WithRequestID(ctx, req.ID)
    ctx = logger.WithTenantID(ctx, req.TenantID)

    // All logs in this request will include these values
    logger.Info(ctx, "Processing request")

    processOrder(ctx, req.Order)
}

func processOrder(ctx context.Context, order *Order) {
    // Context values are preserved
    logger.Info(ctx, "Processing order",
        logger.String("order_id", order.ID),
    )
}
```

### 4. Use Module Loggers for Component-Level Control

```go
var (
    dbLog   = logger.Module("database")
    httpLog = logger.Module("http")
    cacheLog = logger.Module("cache")
)

// In production, you can dynamically enable debug logging
// for specific components without restarting
logger.SetModuleLevel("database", logger.DebugLevel)
```

### 5. Use Appropriate Log Levels

| Level | Use Case |
|-------|----------|
| `Debug` | Detailed debugging information, disabled in production |
| `Info` | General operational events, request processing |
| `Warn` | Unexpected situations that don't cause failures |
| `Error` | Errors that need attention but don't crash the app |
| `Fatal` | Critical errors that require immediate shutdown |

### 6. Add Error Context

```go
if err != nil {
    logger.Error(ctx, "Failed to process payment",
        logger.String("order_id", orderID),
        logger.String("payment_method", method),
        logger.Float64("amount", amount),
        logger.Err(err), // Includes error message and stack trace
    )
}
```

### 7. Use With() for Repeated Fields

```go
// Create a logger with preset fields for a component
orderLogger := logger.L().With(
    logger.String("component", "order-processor"),
    logger.String("version", "2.0"),
)

// All logs from this logger include the preset fields
orderLogger.Info(ctx, "Processing started")
orderLogger.Info(ctx, "Validation complete")
orderLogger.Info(ctx, "Order saved")
```

---

## API Reference

### Initialization Functions

| Function | Description |
|----------|-------------|
| `New(cfg Config) (*Logger, error)` | Create a new logger instance |
| `Init(cfg Config) (*Logger, error)` | Initialize global logger |
| `InitFromConfig(path string) (*Logger, error)` | Initialize from YAML file |
| `InitFromEnv() (*Logger, error)` | Initialize from environment variables |
| `InitFromUnifiedConfig(cfg config.Config) (*Logger, error)` | Initialize from unified config |

### Global Logger Functions

| Function | Description |
|----------|-------------|
| `L() *Logger` | Get global logger instance |
| `ReplaceGlobal(log *Logger) func()` | Replace global logger, returns restore function |
| `Sync() error` | Flush buffered logs |
| `Close() error` | Shutdown and cleanup |

### Logging Functions (Instance and Package-level)

| Function | Description |
|----------|-------------|
| `Debug(ctx, msg, fields...)` | Log at debug level |
| `Info(ctx, msg, fields...)` | Log at info level |
| `Warn(ctx, msg, fields...)` | Log at warn level |
| `Error(ctx, msg, fields...)` | Log at error level |
| `Fatal(ctx, msg, fields...)` | Log at fatal level and exit |

### Logger Methods

| Method | Description |
|--------|-------------|
| `With(fields...) *Logger` | Create logger with preset fields |
| `Named(name string) *Logger` | Create named sub-logger |
| `Module(name string) *ModuleLogger` | Create module logger |
| `SetLevel(level Level)` | Change log level |
| `Level() Level` | Get current log level |
| `Sync() error` | Flush buffered logs |
| `Close() error` | Shutdown and cleanup |

### Module Functions

| Function | Description |
|----------|-------------|
| `Module(name string) *ModuleLogger` | Create module logger from global |
| `SetModuleLevel(module string, level Level)` | Set module log level |
| `ModuleLevel(module string) Level` | Get module log level |
| `SetModuleLevels(map[string]Level)` | Set multiple module levels |
| `ListModuleLevels() map[string]Level` | List all module levels |
| `ResetModuleLevels()` | Reset all module levels |

### Context Functions

| Function | Description |
|----------|-------------|
| `WithTenantID(ctx, id string) context.Context` | Add tenant ID to context |
| `WithRequestID(ctx, id string) context.Context` | Add request ID to context |
| `WithUserID(ctx, id string) context.Context` | Add user ID to context |
| `WithCorrelationID(ctx, id string) context.Context` | Add correlation ID to context |

### Level Type

| Constant | Value | Description |
|----------|-------|-------------|
| `DebugLevel` | -1 | Debug level |
| `InfoLevel` | 0 | Info level |
| `WarnLevel` | 1 | Warning level |
| `ErrorLevel` | 2 | Error level |
| `FatalLevel` | 3 | Fatal level |

| Function | Description |
|----------|-------------|
| `ParseLevel(s string) (Level, error)` | Parse level from string |
| `(l Level) String() string` | Convert level to string |

---

## Example: Complete Application Setup

```go
package main

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

func main() {
    // Initialize logger
    cfg := logger.DefaultConfig()
    cfg.ServiceName = "order-service"
    cfg.ServiceVersion = "1.0.0"
    cfg.Environment = os.Getenv("ENVIRONMENT")
    cfg.Level = os.Getenv("LOG_LEVEL")
    if cfg.Level == "" {
        cfg.Level = "info"
    }
    cfg.OTELEnabled = os.Getenv("OTEL_ENABLED") == "true"
    cfg.OTELEndpoint = os.Getenv("OTEL_ENDPOINT")
    cfg.ModuleLevels = map[string]string{
        "database": "debug",
        "http":     "info",
    }

    log, err := logger.Init(cfg)
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    ctx := context.Background()
    logger.Info(ctx, "Starting application",
        logger.String("environment", cfg.Environment),
    )

    // Create component loggers
    httpLog := logger.Module("http")
    dbLog := logger.Module("database")

    // Setup HTTP server
    http.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
        reqCtx := r.Context()
        reqCtx = logger.WithRequestID(reqCtx, r.Header.Get("X-Request-ID"))
        reqCtx = logger.WithTenantID(reqCtx, r.Header.Get("X-Tenant-ID"))

        httpLog.Info(reqCtx, "Request received",
            logger.String("method", r.Method),
            logger.String("path", r.URL.Path),
        )

        // Process order...
        dbLog.Debug(reqCtx, "Querying database")

        w.WriteHeader(http.StatusOK)
    })

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        logger.Info(ctx, "Shutting down gracefully")
        // Cleanup...
    }()

    logger.Info(ctx, "Server listening", logger.String("addr", ":8080"))
    http.ListenAndServe(":8080", nil)
}
```

---

## License

This package is part of the Motadata Go SDK.
