# Metrics Package

A production-ready metrics collection package built on [OpenTelemetry](https://opentelemetry.io/) for comprehensive application observability.

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
- [Metric Types](#metric-types)
  - [Counter](#counter)
  - [Gauge](#gauge)
  - [Histogram](#histogram)
  - [Timer](#timer)
- [Usage](#usage)
  - [Basic Metrics](#basic-metrics)
  - [Labels and Dimensions](#labels-and-dimensions)
  - [Namespaced Registries](#namespaced-registries)
  - [Service Metrics](#service-metrics)
  - [Helper Functions](#helper-functions)
- [Global Registry](#global-registry)
- [OpenTelemetry Integration](#opentelemetry-integration)
- [Best Practices](#best-practices)
- [API Reference](#api-reference)

---

## Cross-Module Usage

The metrics package is designed for use across multiple packages in your service. **Initialize once, use everywhere.**

```
myapp/
├── main.go              # Initialize metrics here
├── handlers/
│   └── order.go         # Use metrics.Namespace("http")
├── services/
│   └── order.go         # Use metrics.NewServiceMetrics("orders")
└── repository/
    └── order.go         # Use metrics.Namespace("database")
```

```go
// main.go - Initialize once
func main() {
    metrics.Init(cfg)
    defer metrics.Shutdown(context.Background())
    // ...
}

// handlers/order.go - Use namespaced metrics
package handlers

import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"

var httpMetrics = metrics.Namespace("http")

func HandleOrder(ctx context.Context) {
    httpMetrics.Counter("requests_total", "Total requests").Inc(ctx)

    timer := httpMetrics.Timer(ctx, "duration_ms", "Request duration")
    defer timer.Stop()
    // ...
}

// services/order.go - Use ServiceMetrics for standard patterns
package services

import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"

var orderMetrics = metrics.NewServiceMetrics("order_service")

func CreateOrder(ctx context.Context) error {
    timer := orderMetrics.Timer(ctx, "create")
    defer timer.Stop()

    orderMetrics.Requests("create").Inc(ctx)
    return nil
}

// repository/order.go - Use namespace for database metrics
package repository

import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"

var dbMetrics = metrics.Namespace("database")

func SaveOrder(ctx context.Context) error {
    dbMetrics.Counter("queries_total", "Queries").Inc(ctx,
        metrics.Labels{"table": "orders", "operation": "insert"})
    return nil
}
```

| Approach | Usage | Best For |
|----------|-------|----------|
| **Package functions** | `metrics.NewCounter()` | Quick, simple metrics |
| **Namespace** | `metrics.Namespace("http")` | Organizing by component with prefixed names |
| **ServiceMetrics** | `metrics.NewServiceMetrics("name")` | Standard request/error/duration patterns |

---

## Features

- **OpenTelemetry Native**: Built on OTEL SDK for industry-standard observability
- **Multiple Metric Types**: Counters, Gauges, Histograms, and Timers
- **Labeled Metrics**: Rich dimensional data with flexible label support
- **Namespaced Registries**: Organize metrics by service or component
- **Automatic Deduplication**: Metrics are reused when created with same name
- **No-op Fallback**: Graceful degradation when metrics are disabled
- **Service Metrics Builder**: Pre-built patterns for common service metrics
- **OTLP Export**: Native gRPC export to OpenTelemetry Collectors

---

## Installation

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
```

---

## Quick Start

### Minimal Setup

```go
package main

import (
    "context"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
)

func main() {
    // Initialize with defaults
    cfg := metrics.DefaultConfig()
    cfg.Enabled = true
    cfg.ServiceName = "my-service"
    cfg.ServiceVersion = "1.0.0"
    cfg.Environment = "production"

    if err := metrics.Init(cfg); err != nil {
        panic(err)
    }
    defer metrics.Shutdown(context.Background())

    // Create and use metrics
    ctx := context.Background()

    requestCounter := metrics.NewCounter("http_requests_total", "Total HTTP requests")
    requestCounter.Inc(ctx)

    activeConnections := metrics.NewGauge("active_connections", "Active connections")
    activeConnections.Set(ctx, 42)

    latencyHistogram := metrics.NewHistogram("request_duration_ms", "Request duration")
    latencyHistogram.Observe(ctx, 125.5)
}
```

### Using Environment Variables

```go
package main

import (
    "context"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
)

func main() {
    if err := metrics.InitFromEnv(); err != nil {
        panic(err)
    }
    defer metrics.Shutdown(context.Background())

    ctx := context.Background()
    counter := metrics.NewCounter("events_total", "Total events")
    counter.Inc(ctx)
}
```

---

## Configuration

### Configuration Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `Enabled` | bool | `true` | Enable/disable metrics collection |
| `ServiceName` | string | `"app"` | Service identifier for all metrics |
| `ServiceVersion` | string | `"0.0.0"` | Service version |
| `Environment` | string | `"development"` | Deployment environment |
| `OTELEnabled` | bool | `false` | Enable OpenTelemetry export |
| `OTELEndpoint` | string | `"localhost:4317"` | OTEL collector gRPC endpoint |
| `OTELInsecure` | bool | `true` | Use insecure gRPC connection |
| `OTELProtocol` | string | `"grpc"` | Export protocol: `grpc` or `http` |
| `PrometheusEnabled` | bool | `false` | Enable Prometheus endpoint |
| `PrometheusPort` | int | `9090` | Prometheus metrics port |
| `PrometheusPath` | string | `"/metrics"` | Prometheus metrics path |
| `ExportInterval` | time.Duration | `15s` | Metric export interval |
| `DefaultLabels` | map[string]string | `{}` | Labels applied to all metrics |

### Configuration Methods

#### 1. Programmatic Configuration

```go
cfg := metrics.DefaultConfig()
cfg.Enabled = true
cfg.ServiceName = "order-service"
cfg.ServiceVersion = "2.1.0"
cfg.Environment = "production"
cfg.OTELEnabled = true
cfg.OTELEndpoint = "otel-collector:4317"
cfg.OTELInsecure = false
cfg.ExportInterval = 10 * time.Second
cfg.DefaultLabels = map[string]string{
    "region": "us-east-1",
    "cluster": "prod-01",
}

if err := metrics.Init(cfg); err != nil {
    panic(err)
}
```

#### 2. From Environment Variables

```go
if err := metrics.InitFromEnv(); err != nil {
    panic(err)
}
```

#### 3. From Unified Config

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"

cfg, err := config.Load("/etc/myapp/config.yaml")
if err != nil {
    panic(err)
}

if err := metrics.InitFromUnifiedConfig(cfg); err != nil {
    panic(err)
}
```

#### 4. Must Initialize (Panics on Error)

```go
// Useful for applications that cannot function without metrics
metrics.MustInit(cfg)
```

### Environment Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `METRICS_ENABLED` | Enable metrics collection | `true` |
| `METRICS_SERVICE_NAME` | Service name | `my-service` |
| `METRICS_SERVICE_VERSION` | Service version | `1.0.0` |
| `METRICS_ENVIRONMENT` | Environment | `production` |
| `METRICS_OTEL_ENABLED` | Enable OTEL export | `true` |
| `METRICS_OTEL_ENDPOINT` | OTEL collector endpoint | `otel-collector:4317` |
| `METRICS_OTEL_INSECURE` | Insecure gRPC | `false` |
| `METRICS_OTEL_PROTOCOL` | Export protocol | `grpc` |
| `METRICS_PROMETHEUS_ENABLED` | Enable Prometheus | `true` |
| `METRICS_PROMETHEUS_PORT` | Prometheus port | `9090` |
| `METRICS_PROMETHEUS_PATH` | Prometheus path | `/metrics` |
| `METRICS_EXPORT_INTERVAL` | Export interval | `15s` |

### YAML Configuration

```yaml
service:
  name: order-service
  version: 1.0.0
  environment: production

metrics:
  enabled: true

  otel_enabled: true
  otel_endpoint: otel-collector:4317
  otel_insecure: false
  otel_protocol: grpc

  prometheus_enabled: true
  prometheus_port: 9090
  prometheus_path: /metrics

  export_interval: 15s

  default_labels:
    region: us-east-1
    cluster: prod-01
```

---

## Metric Types

### Counter

A monotonically increasing value. Use for counting events, requests, errors, etc.

```go
ctx := context.Background()

// Create a counter
requestCounter := metrics.NewCounter(
    "http_requests_total",
    "Total number of HTTP requests",
)

// Increment by 1
requestCounter.Inc(ctx)

// Increment with labels
requestCounter.Inc(ctx, metrics.Labels{"method": "GET", "path": "/api/users"})

// Add a specific value (must be >= 0)
requestCounter.Add(ctx, 5)
requestCounter.Add(ctx, 10, metrics.Labels{"status": "200"})

// Get metric name
name := requestCounter.Name() // "http_requests_total"
```

**Use Cases:**
- HTTP request counts
- Error counts
- Events processed
- Messages sent/received
- Cache hits/misses

### Gauge

A value that can go up or down. Use for current state measurements.

```go
ctx := context.Background()

// Create a gauge
activeConnections := metrics.NewGauge(
    "active_connections",
    "Number of active connections",
)

// Set to a specific value
activeConnections.Set(ctx, 42)

// Increment by 1
activeConnections.Inc(ctx)

// Decrement by 1
activeConnections.Dec(ctx)

// Add (can be negative)
activeConnections.Add(ctx, 5)
activeConnections.Add(ctx, -3)

// With labels
activeConnections.Set(ctx, 10, metrics.Labels{"protocol": "http"})
activeConnections.Set(ctx, 25, metrics.Labels{"protocol": "grpc"})
```

**Use Cases:**
- Active connections/sessions
- Queue depth
- Memory usage
- Temperature readings
- Current inventory levels

### Histogram

Records observations in configurable buckets. Use for measuring distributions.

```go
ctx := context.Background()

// Create a histogram with default buckets
latencyHistogram := metrics.NewHistogram(
    "request_duration_ms",
    "Request duration in milliseconds",
)

// Record a value
latencyHistogram.Observe(ctx, 125.5)
latencyHistogram.Observe(ctx, 250.0, metrics.Labels{"endpoint": "/api/orders"})

// Record duration from start time
startTime := time.Now()
// ... do work ...
latencyHistogram.ObserveDuration(ctx, startTime)
latencyHistogram.ObserveDuration(ctx, startTime, metrics.Labels{"operation": "query"})
```

**Custom Buckets:**

```go
// Create histogram with custom buckets via registry
registry := metrics.R()
customHistogram := registry.HistogramWithBuckets(
    "response_size_bytes",
    "Response size in bytes",
    "bytes",
    []float64{100, 500, 1000, 5000, 10000, 50000, 100000},
)
```

**Default Bucket Sets:**

```go
// Default latency buckets (milliseconds)
metrics.DefaultHistogramBuckets = []float64{
    1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000,
}

// Default size buckets (bytes)
metrics.DefaultSizeBuckets = []float64{
    100, 1000, 10000, 100000, 1000000, 10000000, 100000000,
}
```

**Use Cases:**
- Request latencies
- Response sizes
- Processing times
- Queue wait times
- Batch sizes

### Timer

A convenience wrapper for timing operations.

```go
ctx := context.Background()

// Create and start a timer
timer := metrics.NewTimer(ctx, "operation_duration_ms", "Operation duration")

// Do work...
processData()

// Stop the timer (records duration)
timer.Stop()

// Or stop with additional labels
timer.StopWithLabels(metrics.Labels{"status": "success"})
```

**Using Registry Timer:**

```go
registry := metrics.R()

func processRequest(ctx context.Context) {
    timer := registry.Timer(ctx, "request_processing_ms", "Request processing time",
        metrics.Labels{"handler": "orders"})
    defer timer.Stop()

    // Process request...
}
```

---

## Usage

### Basic Metrics

```go
ctx := context.Background()

// Counters for events
loginCounter := metrics.NewCounter("user_logins_total", "Total user logins")
loginCounter.Inc(ctx)

errorCounter := metrics.NewCounter("errors_total", "Total errors")
errorCounter.Inc(ctx, metrics.Labels{"type": "validation"})

// Gauges for current state
queueDepth := metrics.NewGauge("queue_depth", "Current queue depth")
queueDepth.Set(ctx, 150)

// Histograms for distributions
responseTime := metrics.NewHistogram("response_time_ms", "Response time")
responseTime.Observe(ctx, 45.2)
```

### Labels and Dimensions

Labels add dimensions to your metrics for filtering and grouping:

```go
ctx := context.Background()

// Counter with labels
httpRequests := metrics.NewCounter("http_requests_total", "HTTP requests")
httpRequests.Inc(ctx, metrics.Labels{
    "method": "POST",
    "path":   "/api/orders",
    "status": "201",
})

// Gauge with labels
connections := metrics.NewGauge("connections", "Active connections")
connections.Set(ctx, 10, metrics.Labels{"protocol": "http"})
connections.Set(ctx, 25, metrics.Labels{"protocol": "grpc"})
connections.Set(ctx, 5, metrics.Labels{"protocol": "websocket"})

// Histogram with labels
latency := metrics.NewHistogram("latency_ms", "Request latency")
latency.Observe(ctx, 50, metrics.Labels{
    "service": "user-service",
    "method":  "GetUser",
})
```

**Merging Labels:**

```go
baseLabels := metrics.Labels{"service": "orders", "version": "2.0"}
requestLabels := metrics.Labels{"method": "POST", "path": "/create"}

// Merge labels (second set takes precedence)
merged := baseLabels.Merge(requestLabels)
// Result: {"service": "orders", "version": "2.0", "method": "POST", "path": "/create"}
```

### Namespaced Registries

Create registries with namespace prefixes for organizing metrics:

```go
// Create namespaced registries
orderMetrics := metrics.Namespace("orders")
paymentMetrics := metrics.Namespace("payments")
inventoryMetrics := metrics.Namespace("inventory")

ctx := context.Background()

// Metrics are automatically prefixed
orderCounter := orderMetrics.Counter("created_total", "Orders created")
// Actual metric name: "orders_created_total"

paymentGauge := paymentMetrics.Gauge("pending", "Pending payments")
// Actual metric name: "payments_pending"

inventoryHistogram := inventoryMetrics.Histogram("lookup_ms", "Inventory lookup time")
// Actual metric name: "inventory_lookup_ms"
```

### Service Metrics

Pre-built patterns for common service metrics:

```go
// Create service metrics for a component
orderService := metrics.NewServiceMetrics("order_service")

ctx := context.Background()

// Track requests
orderService.Requests("create").Inc(ctx)
orderService.Requests("update").Inc(ctx)
orderService.Requests("delete").Inc(ctx)

// Track errors
orderService.Errors("create").Inc(ctx)

// Track duration
orderService.Duration("create").Observe(ctx, 125.5)

// Track active operations
orderService.Active("processing").Inc(ctx)
defer orderService.Active("processing").Dec(ctx)

// Use timer for operations
timer := orderService.Timer(ctx, "create")
defer timer.Stop()

// Record complete request with duration and error status
startTime := time.Now()
err := processOrder()
duration := time.Since(startTime)
orderService.RecordRequest(ctx, "create", duration, err)
```

**ServiceMetrics Methods:**

| Method | Metric Created | Description |
|--------|----------------|-------------|
| `Requests(op)` | `{namespace}_requests_total` | Counter for request counts |
| `Errors(op)` | `{namespace}_errors_total` | Counter for error counts |
| `Duration(op)` | `{namespace}_duration_ms` | Histogram for durations |
| `Active(op)` | `{namespace}_active` | Gauge for active operations |
| `Timer(ctx, op)` | Uses duration histogram | Timer for operations |
| `RecordRequest(...)` | Uses all above | Record complete request |

### Helper Functions

Convenience functions for common metric patterns:

```go
ctx := context.Background()

// Request counter with standard naming
requestCounter := metrics.RequestCounter("api")
// Creates: "api_requests_total"

// Error counter with standard naming
errorCounter := metrics.ErrorCounter("api")
// Creates: "api_errors_total"

// Duration histogram with default buckets
durationHist := metrics.DurationHistogram("api")
// Creates: "api_duration_ms" with default latency buckets

// Size histogram with default buckets
sizeHist := metrics.SizeHistogram("response")
// Creates: "response_bytes" with default size buckets

// Active gauge with standard naming
activeGauge := metrics.ActiveGauge("connections")
// Creates: "connections_active"
```

**Timing Functions:**

```go
ctx := context.Background()

// Time a function execution
metrics.TimeFunc(ctx, "operation_duration_ms", "Operation duration", func() {
    // Do work...
    processData()
})

// Time with labels
metrics.TimeFuncWithLabels(ctx, "query_duration_ms", "Query duration",
    metrics.Labels{"table": "users"},
    func() {
        // Execute query...
    },
)
```

---

## Global Registry

The package provides a global registry for convenience:

```go
// Initialize global metrics
metrics.Init(cfg)

// Access global registry
registry := metrics.R()

// Package-level functions use the global registry
counter := metrics.NewCounter("name", "description")
gauge := metrics.NewGauge("name", "description")
histogram := metrics.NewHistogram("name", "description")
timer := metrics.NewTimer(ctx, "name", "description")

// Create namespaced registry from global
nsRegistry := metrics.Namespace("myservice")

// Get registry statistics
stats := registry.Stats()
fmt.Printf("Registered: %d counters, %d gauges, %d histograms\n",
    stats.Counters, stats.Gauges, stats.Histograms)

// Shutdown (flushes pending metrics)
metrics.Shutdown(context.Background())
```

---

## OpenTelemetry Integration

### Exporting to OTEL Collector

```go
cfg := metrics.DefaultConfig()
cfg.Enabled = true
cfg.ServiceName = "my-service"
cfg.OTELEnabled = true
cfg.OTELEndpoint = "otel-collector:4317"
cfg.OTELInsecure = false // Use TLS in production
cfg.ExportInterval = 10 * time.Second

if err := metrics.Init(cfg); err != nil {
    panic(err)
}
defer metrics.Shutdown(context.Background()) // Flushes pending metrics
```

### Resource Attributes

Metrics are automatically tagged with:
- `service.name` - Service name from config
- `service.version` - Service version from config
- `deployment.environment` - Environment from config

### Default Labels

Apply labels to all metrics:

```go
cfg.DefaultLabels = map[string]string{
    "region":     "us-east-1",
    "cluster":    "prod-01",
    "datacenter": "dc1",
}
```

---

## Use Cases

### 1. HTTP Server Metrics

```go
var httpMetrics = metrics.Namespace("http")

func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Track active requests
        activeGauge := httpMetrics.Gauge("requests_active", "Active requests")
        activeGauge.Inc(r.Context())
        defer activeGauge.Dec(r.Context())

        // Time the request
        timer := httpMetrics.Timer(r.Context(), "request_duration_ms", "Request duration",
            metrics.Labels{"method": r.Method, "path": r.URL.Path})

        // Wrap response to capture status
        wrapped := &responseWriter{ResponseWriter: w, status: 200}
        next.ServeHTTP(wrapped, r)

        // Record metrics
        timer.StopWithLabels(metrics.Labels{"status": strconv.Itoa(wrapped.status)})

        httpMetrics.Counter("requests_total", "Total requests").Inc(r.Context(),
            metrics.Labels{
                "method": r.Method,
                "path":   r.URL.Path,
                "status": strconv.Itoa(wrapped.status),
            })
    })
}
```

### 2. Database Connection Pool Metrics

```go
var dbMetrics = metrics.Namespace("database")

type MetricsPool struct {
    pool *sql.DB
}

func (p *MetricsPool) Stats(ctx context.Context) {
    stats := p.pool.Stats()

    dbMetrics.Gauge("connections_open", "Open connections").Set(ctx, float64(stats.OpenConnections))
    dbMetrics.Gauge("connections_idle", "Idle connections").Set(ctx, float64(stats.Idle))
    dbMetrics.Gauge("connections_in_use", "In-use connections").Set(ctx, float64(stats.InUse))
    dbMetrics.Gauge("connections_max", "Max open connections").Set(ctx, float64(stats.MaxOpenConnections))

    dbMetrics.Counter("connections_wait_total", "Total wait count").Add(ctx, float64(stats.WaitCount))
    dbMetrics.Histogram("connections_wait_duration_ms", "Wait duration").Observe(ctx,
        float64(stats.WaitDuration.Milliseconds()))
}

func (p *MetricsPool) Query(ctx context.Context, query string) (*sql.Rows, error) {
    timer := dbMetrics.Timer(ctx, "query_duration_ms", "Query duration",
        metrics.Labels{"operation": "query"})
    defer timer.Stop()

    rows, err := p.pool.QueryContext(ctx, query)
    if err != nil {
        dbMetrics.Counter("query_errors_total", "Query errors").Inc(ctx,
            metrics.Labels{"operation": "query"})
    }
    return rows, err
}
```

### 3. Message Queue Metrics

```go
var queueMetrics = metrics.Namespace("queue")

type MetricsQueue struct {
    queue MessageQueue
}

func (q *MetricsQueue) Publish(ctx context.Context, msg Message) error {
    timer := queueMetrics.Timer(ctx, "publish_duration_ms", "Publish duration")
    defer timer.Stop()

    queueMetrics.Counter("messages_published_total", "Published messages").Inc(ctx,
        metrics.Labels{"topic": msg.Topic})
    queueMetrics.Histogram("message_size_bytes", "Message size").Observe(ctx,
        float64(len(msg.Body)), metrics.Labels{"topic": msg.Topic})

    return q.queue.Publish(ctx, msg)
}

func (q *MetricsQueue) Consume(ctx context.Context, handler MessageHandler) {
    for msg := range q.queue.Messages() {
        start := time.Now()

        queueMetrics.Gauge("queue_depth", "Queue depth").Set(ctx, float64(q.queue.Len()),
            metrics.Labels{"topic": msg.Topic})

        err := handler(ctx, msg)

        labels := metrics.Labels{"topic": msg.Topic}
        if err != nil {
            labels["status"] = "error"
            queueMetrics.Counter("messages_failed_total", "Failed messages").Inc(ctx, labels)
        } else {
            labels["status"] = "success"
            queueMetrics.Counter("messages_processed_total", "Processed messages").Inc(ctx, labels)
        }

        queueMetrics.Histogram("processing_duration_ms", "Processing duration").Observe(ctx,
            float64(time.Since(start).Milliseconds()), labels)
    }
}
```

### 4. Cache Metrics

```go
var cacheMetrics = metrics.Namespace("cache")

type MetricsCache struct {
    cache Cache
}

func (c *MetricsCache) Get(ctx context.Context, key string) (interface{}, bool) {
    timer := cacheMetrics.Timer(ctx, "operation_duration_ms", "Cache operation duration",
        metrics.Labels{"operation": "get"})
    defer timer.Stop()

    value, found := c.cache.Get(key)

    if found {
        cacheMetrics.Counter("hits_total", "Cache hits").Inc(ctx)
    } else {
        cacheMetrics.Counter("misses_total", "Cache misses").Inc(ctx)
    }

    return value, found
}

func (c *MetricsCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) {
    timer := cacheMetrics.Timer(ctx, "operation_duration_ms", "Cache operation duration",
        metrics.Labels{"operation": "set"})
    defer timer.Stop()

    c.cache.Set(key, value, ttl)
    cacheMetrics.Counter("sets_total", "Cache sets").Inc(ctx)
}

func (c *MetricsCache) Stats(ctx context.Context) {
    stats := c.cache.Stats()
    cacheMetrics.Gauge("items_count", "Items in cache").Set(ctx, float64(stats.ItemCount))
    cacheMetrics.Gauge("size_bytes", "Cache size in bytes").Set(ctx, float64(stats.SizeBytes))

    hitRate := float64(stats.Hits) / float64(stats.Hits+stats.Misses) * 100
    cacheMetrics.Gauge("hit_rate_percent", "Cache hit rate").Set(ctx, hitRate)
}
```

### 5. Business Metrics

```go
var businessMetrics = metrics.NewServiceMetrics("orders")

func ProcessOrder(ctx context.Context, order Order) error {
    // Track the full operation
    timer := businessMetrics.Timer(ctx, "process")
    defer timer.Stop()

    // Validate
    if err := validate(order); err != nil {
        businessMetrics.Errors("process").Inc(ctx, metrics.Labels{"reason": "validation"})
        return err
    }

    // Process payment
    if err := processPayment(ctx, order); err != nil {
        businessMetrics.Errors("process").Inc(ctx, metrics.Labels{"reason": "payment"})
        return err
    }

    // Record success
    businessMetrics.Requests("process").Inc(ctx, metrics.Labels{"status": "success"})

    // Business-specific metrics
    metrics.NewCounter("orders_revenue_total", "Total revenue").Add(ctx, order.Total,
        metrics.Labels{
            "currency": order.Currency,
            "region":   order.Region,
        })

    metrics.NewHistogram("order_items_count", "Items per order").Observe(ctx,
        float64(len(order.Items)))

    return nil
}
```

### 6. Rate Limiting Metrics

```go
var rateLimitMetrics = metrics.Namespace("ratelimit")

type MetricsRateLimiter struct {
    limiter RateLimiter
}

func (r *MetricsRateLimiter) Allow(ctx context.Context, key string) bool {
    allowed := r.limiter.Allow(key)

    labels := metrics.Labels{"key_prefix": extractPrefix(key)}

    if allowed {
        rateLimitMetrics.Counter("requests_allowed_total", "Allowed requests").Inc(ctx, labels)
    } else {
        rateLimitMetrics.Counter("requests_denied_total", "Denied requests").Inc(ctx, labels)
    }

    // Track current rate
    rateLimitMetrics.Gauge("current_rate", "Current request rate").Set(ctx,
        r.limiter.CurrentRate(key), labels)

    return allowed
}
```

---

## Good Practices vs Bad Practices

### Metric Naming

```go
// GOOD - Clear, consistent naming with units
metrics.NewCounter("http_requests_total", "Total HTTP requests")
metrics.NewHistogram("http_request_duration_seconds", "Request duration in seconds")
metrics.NewGauge("db_connections_active", "Currently active database connections")

// BAD - Unclear, inconsistent naming
metrics.NewCounter("requests", "requests")                    // No unit, vague description
metrics.NewHistogram("duration", "how long things take")      // No context, no unit
metrics.NewGauge("connections", "connections")                // Which connections?
```

### Label Cardinality

```go
// GOOD - Bounded label values
counter.Inc(ctx, metrics.Labels{
    "method":   r.Method,           // ~10 values (GET, POST, PUT, DELETE, etc.)
    "status":   strconv.Itoa(code), // ~50 values (200, 201, 400, 404, 500, etc.)
    "endpoint": "/api/users",       // Known set of endpoints
})

// BAD - Unbounded label values (cardinality explosion!)
counter.Inc(ctx, metrics.Labels{
    "user_id":    userID,           // Millions of unique values!
    "request_id": requestID,        // Unique per request!
    "timestamp":  time.Now().String(), // Always unique!
    "ip_address": clientIP,         // Thousands of unique IPs
})

// If you need user-level metrics, use a separate dimension
// or aggregate by user segments instead
counter.Inc(ctx, metrics.Labels{
    "user_tier": "premium",         // Bounded: free, basic, premium, enterprise
    "region":    "us-east",         // Bounded: known regions
})
```

### Metric Reuse

```go
// GOOD - Define metrics once, reuse everywhere
var (
    requestCounter = metrics.NewCounter("api_requests_total", "Total API requests")
    latencyHist    = metrics.NewHistogram("api_latency_ms", "API latency")
)

func handleRequest(ctx context.Context) {
    requestCounter.Inc(ctx)
    start := time.Now()
    defer func() {
        latencyHist.ObserveDuration(ctx, start)
    }()
    // Process...
}

// ALSO GOOD - Metrics are auto-deduplicated by name
func handler1(ctx context.Context) {
    counter := metrics.NewCounter("api_requests_total", "Total API requests")
    counter.Inc(ctx)  // Same underlying metric as handler2
}

func handler2(ctx context.Context) {
    counter := metrics.NewCounter("api_requests_total", "Total API requests")
    counter.Inc(ctx)  // Same underlying metric as handler1
}
```

### Timer Usage

```go
// GOOD - Use defer for automatic timing
func processRequest(ctx context.Context) error {
    timer := metrics.NewTimer(ctx, "request_duration_ms", "Request duration")
    defer timer.Stop()

    // All paths automatically timed
    if err := validate(); err != nil {
        return err
    }
    return process()
}

// BAD - Manual timing with multiple exit points
func processRequest(ctx context.Context) error {
    start := time.Now()

    if err := validate(); err != nil {
        // Forgot to record timing here!
        return err
    }

    result := process()

    metrics.NewHistogram("request_duration_ms", "").Observe(ctx,
        float64(time.Since(start).Milliseconds()))
    return result
}
```

### ServiceMetrics Pattern

```go
// GOOD - Use ServiceMetrics for consistent patterns
var orderMetrics = metrics.NewServiceMetrics("order_service")

func CreateOrder(ctx context.Context, order Order) error {
    start := time.Now()
    err := doCreateOrder(ctx, order)
    orderMetrics.RecordRequest(ctx, "create", time.Since(start), err)
    return err
}

// BAD - Inconsistent metric naming across operations
func CreateOrder(ctx context.Context, order Order) error {
    metrics.NewCounter("orders_created", "").Inc(ctx)  // Different pattern
    // ...
}

func UpdateOrder(ctx context.Context, order Order) error {
    metrics.NewCounter("order_updates_total", "").Inc(ctx)  // Yet another pattern
    // ...
}
```

---

## Performance Pitfalls

### 1. Cardinality Explosion

```go
// BAD - Creates millions of metric series, crashes Prometheus/OTEL collector
for _, user := range users {
    metrics.NewCounter("user_actions_total", "User actions").Inc(ctx,
        metrics.Labels{"user_id": user.ID})  // Each user = new time series
}

// GOOD - Aggregate by bounded dimensions
actionsByTier := make(map[string]int)
for _, user := range users {
    actionsByTier[user.Tier]++
}
for tier, count := range actionsByTier {
    metrics.NewCounter("user_actions_total", "User actions").Add(ctx, float64(count),
        metrics.Labels{"tier": tier})  // Only a few tiers
}
```

### 2. Frequent Metric Creation in Hot Paths

```go
// BAD - Metric lookup overhead in hot loop
for i := 0; i < 1000000; i++ {
    metrics.NewCounter("iterations_total", "Iterations").Inc(ctx)
    // NewCounter does map lookup each time
}

// GOOD - Create metric once outside loop
counter := metrics.NewCounter("iterations_total", "Iterations")
for i := 0; i < 1000000; i++ {
    counter.Inc(ctx)
}

// EVEN BETTER - Batch the increment
counter := metrics.NewCounter("iterations_total", "Iterations")
// ... process loop ...
counter.Add(ctx, 1000000)  // Single increment
```

### 3. Labels Allocation

```go
// BAD - Allocates new map on every call
for _, req := range requests {
    counter.Inc(ctx, metrics.Labels{
        "method": req.Method,
        "path":   req.Path,
    })  // New map allocation each iteration
}

// GOOD for high-frequency paths - Reuse label maps
labels := metrics.Labels{}
for _, req := range requests {
    labels["method"] = req.Method
    labels["path"] = req.Path
    counter.Inc(ctx, labels)
}
// Note: This works because labels are copied internally
```

### 4. Histogram Bucket Overhead

```go
// BAD - Too many buckets (100+ values)
buckets := make([]float64, 100)
for i := range buckets {
    buckets[i] = float64(i * 10)
}
registry.HistogramWithBuckets("latency_ms", "Latency", "ms", buckets)

// GOOD - Reasonable number of buckets (10-20)
registry.HistogramWithBuckets("latency_ms", "Latency", "ms",
    []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000})
```

### 5. Missing Shutdown

```go
// BAD - Metrics may not be exported on shutdown
func main() {
    metrics.Init(cfg)
    // Application runs...
    // Program exits without flushing
}

// GOOD - Always shutdown with context
func main() {
    metrics.Init(cfg)
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        metrics.Shutdown(ctx)  // Flushes pending metrics
    }()

    // Application runs...
}
```

### 6. Timer Double-Stop

```go
// BAD - Timer stopped twice
func process(ctx context.Context) {
    timer := metrics.NewTimer(ctx, "duration_ms", "Duration")
    defer timer.Stop()

    if err := step1(); err != nil {
        timer.Stop()  // Already stopped by defer!
        return
    }
}

// GOOD - Timer handles double-stop safely (no panic, but wasteful)
// Better to just use defer consistently
func process(ctx context.Context) {
    timer := metrics.NewTimer(ctx, "duration_ms", "Duration")
    defer timer.Stop()  // Only place timer is stopped

    if err := step1(); err != nil {
        // Error handling, timer stopped by defer
        return
    }
}
```

### 7. Gauge vs Counter Confusion

```go
// BAD - Using gauge for cumulative value
for _, sale := range sales {
    // Gauge loses history when scraped at different times
    metrics.NewGauge("sales_total", "Total sales").Add(ctx, sale.Amount)
}

// GOOD - Use counter for cumulative values
for _, sale := range sales {
    metrics.NewCounter("sales_total", "Total sales").Add(ctx, sale.Amount)
}

// BAD - Using counter for current state
metrics.NewCounter("connections_active", "Active connections").Inc(ctx)
// Counter can only increase, can't show connections closing

// GOOD - Use gauge for current state
gauge := metrics.NewGauge("connections_active", "Active connections")
gauge.Inc(ctx)  // Connection opened
// ...
gauge.Dec(ctx)  // Connection closed
```

---

## Best Practices

### 1. Initialize Early, Shutdown Properly

```go
func main() {
    if err := metrics.Init(cfg); err != nil {
        panic(err)
    }
    defer metrics.Shutdown(context.Background()) // Ensures metrics are flushed

    // Application code...
}
```

### 2. Use Consistent Naming Conventions

```go
// Good - clear, consistent naming
metrics.NewCounter("http_requests_total", "Total HTTP requests")
metrics.NewHistogram("http_request_duration_seconds", "HTTP request duration")
metrics.NewGauge("http_connections_active", "Active HTTP connections")

// Naming patterns:
// - Counters: <noun>_total (http_requests_total)
// - Gauges: <noun>_<state> (connections_active, queue_depth)
// - Histograms: <noun>_<unit> (request_duration_seconds, response_size_bytes)
```

### 3. Use Meaningful Labels

```go
// Good - actionable dimensions
counter.Inc(ctx, metrics.Labels{
    "method":   "POST",
    "endpoint": "/api/orders",
    "status":   "201",
})

// Avoid - high cardinality labels
counter.Inc(ctx, metrics.Labels{
    "user_id":    userID,    // Millions of unique values
    "request_id": requestID, // Unique per request
    "timestamp":  time.Now().String(), // Always unique
})
```

### 4. Reuse Metric Instances

```go
// Good - create once, use many times
var (
    requestCounter = metrics.NewCounter("requests_total", "Total requests")
    latencyHist    = metrics.NewHistogram("latency_ms", "Request latency")
)

func handleRequest(ctx context.Context) {
    requestCounter.Inc(ctx)
    start := time.Now()
    defer func() {
        latencyHist.ObserveDuration(ctx, start)
    }()
    // ...
}

// Also good - metrics are automatically deduplicated
func handler1(ctx context.Context) {
    counter := metrics.NewCounter("requests_total", "Total requests")
    counter.Inc(ctx) // Same metric instance as handler2
}

func handler2(ctx context.Context) {
    counter := metrics.NewCounter("requests_total", "Total requests")
    counter.Inc(ctx) // Same metric instance as handler1
}
```

### 5. Use Namespaces for Organization

```go
// Organize metrics by component
var (
    dbMetrics    = metrics.Namespace("database")
    cacheMetrics = metrics.Namespace("cache")
    httpMetrics  = metrics.Namespace("http")
)

func queryDatabase(ctx context.Context) {
    timer := dbMetrics.Timer(ctx, "query_duration_ms", "Query duration")
    defer timer.Stop()
    // ...
}
```

### 6. Handle Disabled Metrics Gracefully

```go
// Metrics package provides no-op implementations when disabled
cfg.Enabled = false
metrics.Init(cfg)

// These calls are safe and do nothing
counter := metrics.NewCounter("name", "desc")
counter.Inc(ctx) // No-op, no panic
```

### 7. Use ServiceMetrics for Standard Patterns

```go
// Instead of creating individual metrics
requestCounter := metrics.NewCounter("orders_requests_total", "...")
errorCounter := metrics.NewCounter("orders_errors_total", "...")
durationHist := metrics.NewHistogram("orders_duration_ms", "...")

// Use ServiceMetrics for consistent patterns
orderMetrics := metrics.NewServiceMetrics("orders")

func processOrder(ctx context.Context) error {
    start := time.Now()
    err := doProcess()
    orderMetrics.RecordRequest(ctx, "process", time.Since(start), err)
    return err
}
```

### 8. Choose Appropriate Histogram Buckets

```go
registry := metrics.R()

// For latencies (milliseconds)
latencyHist := registry.HistogramWithBuckets(
    "api_latency_ms",
    "API latency",
    "ms",
    []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
)

// For sizes (bytes)
sizeHist := registry.HistogramWithBuckets(
    "payload_size_bytes",
    "Payload size",
    "bytes",
    []float64{256, 512, 1024, 4096, 16384, 65536, 262144, 1048576},
)

// For counts
batchHist := registry.HistogramWithBuckets(
    "batch_size",
    "Batch size",
    "items",
    []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
)
```

---

## API Reference

### Initialization Functions

| Function | Description |
|----------|-------------|
| `Init(cfg Config) error` | Initialize global metrics |
| `InitFromEnv() error` | Initialize from environment variables |
| `InitFromUnifiedConfig(cfg config.Config) error` | Initialize from unified config |
| `MustInit(cfg Config)` | Initialize or panic |
| `Shutdown(ctx context.Context) error` | Shutdown and flush metrics |

### Configuration Functions

| Function | Description |
|----------|-------------|
| `DefaultConfig() Config` | Get default configuration |
| `ConfigFromEnv() Config` | Load config from environment |

### Global Registry Functions

| Function | Description |
|----------|-------------|
| `R() *Registry` | Get global registry |
| `Namespace(name string) *Registry` | Create namespaced registry |
| `NewCounter(name, desc string, labels ...Labels) Counter` | Create counter |
| `NewGauge(name, desc string, labels ...Labels) Gauge` | Create gauge |
| `NewHistogram(name, desc string, labels ...Labels) Histogram` | Create histogram |
| `NewTimer(ctx, name, desc string, labels ...Labels) Timer` | Create timer |

### Helper Metric Functions

| Function | Description |
|----------|-------------|
| `RequestCounter(name string, labels ...Labels) Counter` | Create request counter |
| `ErrorCounter(name string, labels ...Labels) Counter` | Create error counter |
| `DurationHistogram(name string, labels ...Labels) Histogram` | Create duration histogram |
| `SizeHistogram(name string, labels ...Labels) Histogram` | Create size histogram |
| `ActiveGauge(name string, labels ...Labels) Gauge` | Create active gauge |
| `TimeFunc(ctx, name, desc string, fn func())` | Time a function |
| `TimeFuncWithLabels(ctx, name, desc string, labels Labels, fn func())` | Time with labels |

### Registry Methods

| Method | Description |
|--------|-------------|
| `Counter(name, desc string, labels ...Labels) Counter` | Create/get counter |
| `CounterWithUnit(name, desc, unit string, labels ...Labels) Counter` | Counter with unit |
| `Gauge(name, desc string, labels ...Labels) Gauge` | Create/get gauge |
| `GaugeWithUnit(name, desc, unit string, labels ...Labels) Gauge` | Gauge with unit |
| `Histogram(name, desc string, labels ...Labels) Histogram` | Create/get histogram |
| `HistogramWithBuckets(name, desc, unit string, buckets []float64, labels ...Labels) Histogram` | Histogram with buckets |
| `Timer(ctx, name, desc string, labels ...Labels) Timer` | Create timer |
| `Stats() RegistryStats` | Get registry statistics |
| `Shutdown(ctx context.Context) error` | Shutdown registry |

### Counter Interface

| Method | Description |
|--------|-------------|
| `Inc(ctx context.Context, labels ...Labels)` | Increment by 1 |
| `Add(ctx context.Context, value float64, labels ...Labels)` | Add value (>= 0) |
| `Name() string` | Get metric name |

### Gauge Interface

| Method | Description |
|--------|-------------|
| `Set(ctx context.Context, value float64, labels ...Labels)` | Set value |
| `Inc(ctx context.Context, labels ...Labels)` | Increment by 1 |
| `Dec(ctx context.Context, labels ...Labels)` | Decrement by 1 |
| `Add(ctx context.Context, value float64, labels ...Labels)` | Add value (can be negative) |
| `Name() string` | Get metric name |

### Histogram Interface

| Method | Description |
|--------|-------------|
| `Observe(ctx context.Context, value float64, labels ...Labels)` | Record value |
| `ObserveDuration(ctx context.Context, start time.Time, labels ...Labels)` | Record duration |
| `Name() string` | Get metric name |

### Timer Interface

| Method | Description |
|--------|-------------|
| `Stop()` | Stop timer and record duration |
| `StopWithLabels(labels Labels)` | Stop with additional labels |

### Labels Type

| Method | Description |
|--------|-------------|
| `Merge(other Labels) Labels` | Merge two label sets |

### ServiceMetrics

| Method | Description |
|--------|-------------|
| `Requests(operation string) Counter` | Get requests counter |
| `Errors(operation string) Counter` | Get errors counter |
| `Duration(operation string) Histogram` | Get duration histogram |
| `Active(operation string) Gauge` | Get active gauge |
| `Timer(ctx context.Context, operation string) Timer` | Create timer |
| `RecordRequest(ctx, operation string, duration time.Duration, err error)` | Record request |

### Default Bucket Values

```go
// Default histogram buckets for latency (milliseconds)
DefaultHistogramBuckets = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

// Default histogram buckets for sizes (bytes)
DefaultSizeBuckets = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000}
```

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

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
)

// Global metrics
var (
    httpMetrics = metrics.Namespace("http")
    dbMetrics   = metrics.Namespace("database")
)

func main() {
    // Initialize metrics
    cfg := metrics.DefaultConfig()
    cfg.Enabled = true
    cfg.ServiceName = "order-service"
    cfg.ServiceVersion = "1.0.0"
    cfg.Environment = os.Getenv("ENVIRONMENT")
    cfg.OTELEnabled = os.Getenv("OTEL_ENABLED") == "true"
    cfg.OTELEndpoint = os.Getenv("OTEL_ENDPOINT")
    cfg.DefaultLabels = map[string]string{
        "region": os.Getenv("REGION"),
    }

    if err := metrics.Init(cfg); err != nil {
        panic(err)
    }
    defer metrics.Shutdown(context.Background())

    // Setup HTTP handler with metrics
    http.HandleFunc("/api/orders", handleOrders)

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
        <-sigChan

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        metrics.Shutdown(ctx)
    }()

    http.ListenAndServe(":8080", nil)
}

func handleOrders(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Track request
    timer := httpMetrics.Timer(ctx, "request_duration_ms", "Request duration",
        metrics.Labels{"method": r.Method, "path": r.URL.Path})
    defer timer.Stop()

    httpMetrics.Counter("requests_total", "Total requests").Inc(ctx,
        metrics.Labels{"method": r.Method, "path": r.URL.Path})

    // Track active requests
    activeGauge := httpMetrics.Gauge("active_requests", "Active requests")
    activeGauge.Inc(ctx)
    defer activeGauge.Dec(ctx)

    // Process request
    if err := processOrder(ctx); err != nil {
        httpMetrics.Counter("errors_total", "Total errors").Inc(ctx,
            metrics.Labels{"method": r.Method, "path": r.URL.Path, "error": "processing"})
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func processOrder(ctx context.Context) error {
    // Database query with metrics
    timer := dbMetrics.Timer(ctx, "query_duration_ms", "Query duration",
        metrics.Labels{"operation": "insert"})
    defer timer.Stop()

    // Simulate database operation
    time.Sleep(50 * time.Millisecond)

    dbMetrics.Counter("queries_total", "Total queries").Inc(ctx,
        metrics.Labels{"operation": "insert", "table": "orders"})

    return nil
}
```

---

## License

This package is part of the Motadata Go SDK.
