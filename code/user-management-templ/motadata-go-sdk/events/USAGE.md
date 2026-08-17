# Events Package - Step-by-Step Usage Guide

This guide provides detailed, step-by-step examples for using the events package and all its submodules.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Configuration](#configuration)
   - [Default Configuration](#default-configuration)
   - [Environment Variables](#environment-variables)
   - [Functional Options](#functional-options)
3. [Connection Management](#connection-management)
   - [Basic Connection](#basic-connection)
   - [Connection with TLS](#connection-with-tls)
   - [Connection Health Monitoring](#connection-health-monitoring)
4. [Publishing Messages](#publishing-messages)
   - [Synchronous Publish](#synchronous-publish)
   - [Asynchronous Publish (JetStream)](#asynchronous-publish-jetstream)
   - [Request-Reply Pattern](#request-reply-pattern)
   - [Batch Publishing](#batch-publishing)
5. [Subscribing to Messages](#subscribing-to-messages)
   - [Basic Subscription](#basic-subscription)
   - [Queue Subscriptions](#queue-subscriptions)
   - [Wildcard Subscriptions](#wildcard-subscriptions)
6. [JetStream](#jetstream)
   - [Stream Management](#stream-management)
   - [Consumer Management](#consumer-management)
   - [Key-Value Store](#key-value-store)
   - [Object Store](#object-store)
7. [Tenant Management](#tenant-management)
   - [Multi-Tenant Connections](#multi-tenant-connections)
   - [Credential Management](#credential-management)
   - [Idle Connection Cleanup](#idle-connection-cleanup)
8. [Authentication](#authentication)
   - [User/Password](#userpassword)
   - [Token Authentication](#token-authentication)
   - [JWT Authentication](#jwt-authentication)
   - [NKey Authentication](#nkey-authentication)
   - [Credential Manager](#credential-manager)
9. [Middleware](#middleware)
   - [Middleware Stack](#middleware-stack)
   - [Tracing Middleware](#tracing-middleware)
   - [Retry Middleware](#retry-middleware)
   - [Circuit Breaker](#circuit-breaker)
   - [Rate Limiter](#rate-limiter)
   - [Metrics Middleware](#metrics-middleware)
   - [Logging Middleware](#logging-middleware)
   - [Custom Middleware](#custom-middleware)
10. [Microservice Registry](#microservice-registry)
    - [Service Registration](#service-registration)
    - [Instance Management](#instance-management)
    - [Service Discovery](#service-discovery)
    - [Health Checking](#health-checking)
11. [Endpoint Registry](#endpoint-registry)
    - [Endpoint Registration](#endpoint-registration)
    - [Endpoint Discovery](#endpoint-discovery)
    - [Endpoint Health Checking](#endpoint-health-checking)
12. [Context Propagation](#context-propagation)
13. [Error Handling](#error-handling)
14. [Utilities](#utilities)
15. [Best Practices](#best-practices)

---

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
)

func main() {
    ctx := context.Background()

    // Step 1: Create client
    client, err := events.NewClient(
        events.WithServers("nats://localhost:4222"),
        events.WithName("my-service"),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown(ctx)

    // Step 2: Connect a tenant
    conn, err := client.ConnectWithJWT(ctx, "tenant-1", jwtToken, seed)
    if err != nil {
        log.Fatal(err)
    }

    // Step 3: Create publisher and subscriber
    pub := events.NewPublisher(conn)
    sub := events.NewSubscriber(conn)

    // Step 4: Subscribe to messages
    subscription, err := sub.Subscribe(ctx, "events.>",
        func(ctx context.Context, msg *nats.Msg) error {
            fmt.Printf("Received on %s: %s\n", msg.Subject, string(msg.Data))
            return nil
        },
    )
    if err != nil {
        log.Fatal(err)
    }
    defer subscription.Unsubscribe()

    // Step 5: Publish a message
    msg := &nats.Msg{
        Subject: "events.user.created",
        Data:    []byte(`{"user_id": "123", "name": "John"}`),
    }
    err = pub.Publish(ctx, "events.user.created", msg)
    if err != nil {
        log.Fatal(err)
    }

    time.Sleep(time.Second) // Wait for message delivery
}
```

---

## Configuration

### Default Configuration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"

// Create default configuration
cfg := config.DefaultEventsConfig()
// Servers: ["nats://localhost:4222"]
// ConnectTimeout: 10s
// DrainTimeout: 30s
// IdleTimeout: 30m
// CleanupInterval: 1m
// Stream.Enabled: true
// Reconnect.MaxAttempts: -1 (infinite)

// Clone and modify
myCfg := cfg.Clone()
myCfg.Servers = []string{"nats://server1:4222", "nats://server2:4222"}
```

### Environment Variables

```bash
# Set environment variables before starting your application
export EVENTS_SERVERS="nats://server1:4222,nats://server2:4222"
export EVENTS_NAME="my-service"
export EVENTS_CONNECT_TIMEOUT="10s"
export EVENTS_DRAIN_TIMEOUT="30s"
export EVENTS_IDLE_TIMEOUT="30m"
export EVENTS_TLS_ENABLED="true"
export EVENTS_TLS_CERT="/path/to/cert.pem"
export EVENTS_TLS_KEY="/path/to/key.pem"
export EVENTS_TLS_CA="/path/to/ca.pem"
export EVENTS_STREAM_ENABLED="true"
export EVENTS_STREAM_DOMAIN="hub"
```

```go
// Load configuration from environment
cfg, err := config.LoadEventsFromEnv()
if err != nil {
    log.Fatal(err)
}

// Or with custom prefix
cfg, err = config.LoadEventsFromEnvWithPrefix("MYAPP")
// Uses MYAPP_SERVERS, MYAPP_NAME, etc.
```

### Functional Options

```go
// Create client with functional options
client, err := events.NewClient(
    // Server configuration
    events.WithServers("nats://server1:4222", "nats://server2:4222"),
    events.WithName("order-service"),
    events.WithConnectTimeout(10 * time.Second),
    events.WithDrainTimeout(30 * time.Second),

    // TLS configuration
    events.WithTLS(&config.TLSConfig{
        Enabled:  true,
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
        CAFile:   "/path/to/ca.pem",
    }),

    // Or shorthand TLS from files
    events.WithTLSFiles("/path/to/cert.pem", "/path/to/key.pem", "/path/to/ca.pem"),

    // Reconnection settings
    events.WithReconnect(config.ReconnectConfig{
        MaxAttempts:     -1,     // Infinite
        InitialInterval: 100 * time.Millisecond,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          0.1,
    }),

    // JetStream settings
    events.WithStream(config.StreamConfig{
        Enabled: true,
        Domain:  "hub",
    }),

    // Idle connection cleanup
    events.WithIdleTimeout(30 * time.Minute),
    events.WithCleanupInterval(time.Minute),
)
```

---

## Connection Management

### Basic Connection

```go
// Step 1: Create client
client, err := events.NewClient(
    events.WithServers("nats://localhost:4222"),
)
defer client.Shutdown(context.Background())

// Step 2: Connect with JWT (recommended for production)
conn, err := client.ConnectWithJWT(ctx, "tenant-1", jwtToken, nkeySeed)
if err != nil {
    log.Fatal(err)
}

// Step 3: Check connection state
fmt.Println("Connected:", conn.IsConnected())
fmt.Println("State:", conn.State())
fmt.Println("Tenant:", conn.TenantID())

// Step 4: Get raw NATS connection (for advanced usage)
natsConn := conn.Conn()

// Step 5: Get JetStream context
js := conn.JetStream()
```

### Connection with TLS

```go
client, err := events.NewClient(
    events.WithServers("nats://secure-server:4222"),
    events.WithTLS(&config.TLSConfig{
        Enabled:  true,
        CertFile: "/etc/certs/client.pem",
        KeyFile:  "/etc/certs/client-key.pem",
        CAFile:   "/etc/certs/ca.pem",
    }),
)

// Or skip verification (development only!)
client, err = events.NewClient(
    events.WithServers("nats://localhost:4222"),
    events.WithTLSSkipVerify(true),
)
```

### Connection Health Monitoring

```go
// Step 1: Get health status
health := conn.Health()
fmt.Printf("Healthy: %v, State: %s, Message: %s\n",
    health.Healthy, health.State, health.Message)

// Step 2: Set up health monitor
monitor := events.NewHealthMonitor(conn, events.DefaultHealthMonitorConfig())
monitor.Start()
defer monitor.Stop()

// Step 3: Set health change callback
monitor.OnHealthChange(func(healthy bool) {
    if !healthy {
        log.Println("Connection became unhealthy!")
    }
})

// Step 4: Get connection statistics
stats := monitor.Stats()
fmt.Printf("Published: %d, Received: %d, Errors: %d\n",
    stats.MessagesPublished, stats.MessagesReceived, stats.Errors)
fmt.Printf("Avg RTT: %v, Uptime: %.0fs\n",
    stats.AvgRTT, stats.UptimeSeconds)

// Step 5: Record operations manually
monitor.RecordPublish(len(data))
monitor.RecordReceive(len(data))
monitor.RecordError("publish")
```

---

## Publishing Messages

### Synchronous Publish

```go
// Step 1: Create publisher
pub := events.NewPublisher(conn)
defer pub.Close(ctx)

// Step 2: Create a NATS message
msg := &nats.Msg{
    Subject: "orders.created",
    Data:    []byte(`{"order_id": "ORD-123", "amount": 99.99}`),
    Header:  nats.Header{},
}

// Step 3: Add headers
msg.Header.Set("Content-Type", "application/json")
msg.Header.Set("X-Tenant-ID", "tenant-1")
msg.Header.Set("Nats-Msg-Id", "unique-msg-id") // For deduplication

// Step 4: Publish
err := pub.Publish(ctx, "orders.created", msg)
if err != nil {
    log.Printf("Publish failed: %v", err)
}
```

### Asynchronous Publish (JetStream)

```go
// Step 1: Create publisher
pub := events.NewPublisher(conn)

// Step 2: Publish asynchronously (returns future)
msg := &nats.Msg{
    Subject: "orders.created",
    Data:    orderJSON,
}
future := pub.PublishAsync(ctx, "orders.created", msg)

// Step 3a: Non-blocking check
select {
case ack := <-future.Ok():
    fmt.Printf("Published to stream %s, sequence %d\n", ack.Stream, ack.Sequence)
case err := <-future.Err():
    fmt.Printf("Publish failed: %v\n", err)
case <-time.After(5 * time.Second):
    fmt.Println("Publish timed out")
}

// Step 3b: Or blocking wait
ack, err := future.Wait(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Sequence: %d\n", ack.Sequence)
```

### Request-Reply Pattern

```go
// --- Requester Side ---
pub := events.NewPublisher(conn)

// Step 1: Send request and wait for reply
request := &nats.Msg{
    Subject: "users.get",
    Data:    []byte(`{"user_id": "123"}`),
}
response, err := pub.Request(ctx, "users.get", request)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User: %s\n", string(response.Data))

// Step 2: Request with explicit timeout
response, err = pub.RequestWithTimeout(ctx, "users.get", request, 5*time.Second)

// --- Responder Side ---
sub := events.NewSubscriber(conn)
sub.Subscribe(ctx, "users.get", func(ctx context.Context, msg *nats.Msg) error {
    // Process request
    user := getUserByID(string(msg.Data))

    // Send reply
    reply := &nats.Msg{
        Subject: msg.Reply,
        Data:    userJSON,
    }
    return pub.Publish(ctx, msg.Reply, reply)
})
```

### Batch Publishing

```go
// Step 1: Create batch publisher
batch := events.NewBatchPublisher(pub,
    events.WithMaxBatchSize(100),           // Flush at 100 messages
    events.WithFlushInterval(time.Second),  // Or every second
    events.WithConcurrentFlush(true),       // Parallel flush
    events.WithMaxFlushWorkers(4),          // 4 flush workers
    events.WithFlushErrorCallback(func(err error) {
        log.Printf("Batch flush error: %v", err)
    }),
)
defer batch.Close(ctx)

// Step 2: Add messages to batch
for _, order := range orders {
    msg := &nats.Msg{
        Subject: "orders.created",
        Data:    orderJSON,
    }
    batch.Add("orders.created", msg)
}

// Step 3: Check pending count
fmt.Printf("Pending messages: %d\n", batch.Count())

// Step 4: Manual flush (auto-flush also runs on interval/size)
err := batch.Flush(ctx)
if err != nil {
    log.Printf("Flush failed: %v", err)
}

// Step 5: Add multiple at once
messages := []events.BatchMessage{
    {Subject: "orders.created", Msg: msg1},
    {Subject: "orders.updated", Msg: msg2},
}
batch.AddMultiple(messages)
```

---

## Subscribing to Messages

### Basic Subscription

```go
// Step 1: Create subscriber
sub := events.NewSubscriber(conn)
defer sub.Close(ctx)

// Step 2: Subscribe with message handler
subscription, err := sub.Subscribe(ctx, "orders.created",
    func(ctx context.Context, msg *nats.Msg) error {
        fmt.Printf("Order: %s\n", string(msg.Data))

        // Extract context values (set by middleware)
        tenantID, _ := events.TenantIDFromContext(ctx)
        corrID, _ := events.CorrelationIDFromContext(ctx)

        log.Printf("Processing order for tenant %s (correlation: %s)",
            tenantID, corrID)

        return nil
    },
)
if err != nil {
    log.Fatal(err)
}

// Step 3: Check subscription status
fmt.Println("Subject:", subscription.Subject())
fmt.Println("Valid:", subscription.IsValid())

// Step 4: Unsubscribe when done
subscription.Unsubscribe()

// Or drain (process pending then unsubscribe)
subscription.Drain()
```

### Queue Subscriptions

Load balancing across multiple instances:

```go
// Step 1: Create subscriber
sub := events.NewSubscriber(conn)

// Step 2: Subscribe with queue group
// Only ONE subscriber in the group receives each message
subscription, err := sub.QueueSubscribe(ctx, "orders.>", "order-processors",
    func(ctx context.Context, msg *nats.Msg) error {
        fmt.Printf("Processing: %s\n", string(msg.Data))
        return nil
    },
)

// Instance 1: QueueSubscribe("orders.>", "order-processors", handler)
// Instance 2: QueueSubscribe("orders.>", "order-processors", handler)
// Instance 3: QueueSubscribe("orders.>", "order-processors", handler)
// Each message goes to exactly ONE of the three instances
```

### Wildcard Subscriptions

```go
// Single-token wildcard: * matches one token
sub.Subscribe(ctx, "orders.*", handler)
// Matches: orders.created, orders.updated, orders.deleted
// Not: orders.item.created

// Multi-token wildcard: > matches one or more tokens
sub.Subscribe(ctx, "orders.>", handler)
// Matches: orders.created, orders.item.created, orders.us.east.created

// Specific subject
sub.Subscribe(ctx, "orders.created", handler)
// Matches only: orders.created
```

---

## JetStream

### Stream Management

```go
import natsjs "github.com/nats-io/nats.go/jetstream"

// Step 1: Get JetStream context
js := conn.JetStream()

// Step 2: Create stream manager
sm := events.NewStreamManager(js)

// Step 3: Create a stream
stream, err := sm.CreateStream(ctx, natsjs.StreamConfig{
    Name:      "ORDERS",
    Subjects:  []string{"orders.>"},
    Storage:   natsjs.FileStorage,
    Retention: natsjs.LimitsPolicy,
    MaxAge:    7 * 24 * time.Hour, // 7 days
    MaxBytes:  1 << 30,            // 1 GB
    Replicas:  3,
})

// Step 4: List streams
names, err := sm.StreamNames(ctx)

// Step 5: Delete stream
err = sm.DeleteStream(ctx, "ORDERS")
```

### Consumer Management

```go
// Step 1: Create durable consumer
consumer, err := events.CreateOrUpdateConsumer(ctx, js, "ORDERS", natsjs.ConsumerConfig{
    Name:          "order-processor",
    Durable:       "order-processor",
    FilterSubject: "orders.created",
    AckPolicy:     natsjs.AckExplicitPolicy,
    MaxDeliver:    5,
    AckWait:       30 * time.Second,
})

// Step 2: Consume messages continuously
consumeCtx, err := consumer.Consume(func(msg natsjs.Msg) {
    fmt.Printf("Processing: %s\n", string(msg.Data()))

    // Acknowledge message
    msg.Ack()
})
defer consumeCtx.Stop()

// Step 3: Or fetch batches
msgs, err := consumer.Fetch(10)
for msg := range msgs.Messages() {
    process(msg)
    msg.Ack()
}

// Step 4: Get consumer info
info, err := consumer.Info(ctx)
fmt.Printf("Pending: %d\n", info.NumPending)
```

### Key-Value Store

```go
// Step 1: Create KV store
kv, err := events.CreateKVStore(ctx, js, natsjs.KeyValueConfig{
    Bucket:  "app-config",
    History: 5,
    TTL:     24 * time.Hour,
})

// Step 2: Put values
revision, err := kv.Put(ctx, "feature.dark-mode", []byte("true"))

// Step 3: Get values
entry, err := kv.Get(ctx, "feature.dark-mode")
fmt.Printf("Value: %s (rev: %d)\n", string(entry.Value()), entry.Revision())

// Step 4: Watch for changes
watcher, err := kv.Watch(ctx, "feature.*")
defer watcher.Stop()
go func() {
    for entry := range watcher.Updates() {
        if entry != nil {
            fmt.Printf("Config changed: %s = %s\n",
                entry.Key(), string(entry.Value()))
        }
    }
}()

// Step 5: Delete
err = kv.Delete(ctx, "feature.dark-mode")

// Step 6: List keys
keys, err := kv.Keys(ctx)
```

### Object Store

```go
// Step 1: Create object store
objStore, err := events.CreateObjectStore(ctx, js, natsjs.ObjectStoreConfig{
    Bucket:   "reports",
    MaxBytes: 1 << 30, // 1 GB
})

// Step 2: Store an object
info, err := objStore.PutBytes(ctx, "monthly/2024-01.pdf", reportData)

// Step 3: Retrieve an object
result, err := objStore.Get(ctx, "monthly/2024-01.pdf")
defer result.Close()
data, _ := io.ReadAll(result)

// Step 4: List objects
objects, err := objStore.List(ctx)
for _, info := range objects {
    fmt.Printf("%s (%d bytes)\n", info.Name, info.Size)
}

// Step 5: Delete
err = objStore.Delete(ctx, "monthly/2024-01.pdf")
```

---

## Tenant Management

### Multi-Tenant Connections

```go
// Step 1: Create client
client, err := events.NewClient(
    events.WithServers("nats://localhost:4222"),
    events.WithIdleTimeout(30 * time.Minute),
)
defer client.Shutdown(ctx)

// Step 2: Connect multiple tenants (each gets isolated connection)
conn1, err := client.ConnectWithJWT(ctx, "tenant-1", jwt1, seed1)
conn2, err := client.ConnectWithJWT(ctx, "tenant-2", jwt2, seed2)

// Step 3: Each tenant has its own publisher/subscriber
pub1 := events.NewPublisher(conn1)
pub2 := events.NewPublisher(conn2)

// Step 4: Get existing connection
existingConn, err := client.Get("tenant-1")

// Step 5: Disconnect specific tenant
client.Disconnect(ctx, "tenant-2")

// Step 6: Check all tenant health
manager := client.Manager()
health := manager.Health()
for tenantID, status := range health {
    fmt.Printf("Tenant %s: healthy=%v\n", tenantID, status.IsHealthy())
}

// Step 7: List connected tenants
tenantIDs := manager.TenantIDs()
fmt.Printf("Active tenants: %v\n", tenantIDs)
```

### Credential Management

```go
// Step 1: Create credential manager
credManager := events.NewCredentialManager(auth.DefaultCredentialManagerConfig())

// Step 2: Register tenant credentials
credManager.Register("tenant-1", jwt1, seed1)
credManager.Register("tenant-2", jwt2, seed2)

// Step 3: Or register from HTTP payload
payload := events.RegistrationPayload{
    TenantID: "tenant-3",
    JWT:      jwt3,
    Seed:     seed3,
}
credManager.RegisterFromPayload(payload)

// Step 4: Connect using registered credentials
conn, err := client.Connect(ctx, "tenant-1", nil) // Uses credential manager

// Step 5: Handle credential refresh
credManager.OnRefreshNeeded(func(tenantID string) {
    newJWT, newSeed := fetchNewCredentials(tenantID)
    credManager.Register(tenantID, newJWT, newSeed)
})

// Step 6: Start credential manager lifecycle
credManager.Start()
defer credManager.Stop()
```

### Idle Connection Cleanup

```go
// Connections idle for longer than IdleTimeout are automatically closed
client, err := events.NewClient(
    events.WithIdleTimeout(15 * time.Minute),    // Close after 15 min idle
    events.WithCleanupInterval(30 * time.Second), // Check every 30s
)

// Activity is tracked automatically:
// - Publishing a message calls Touch()
// - Receiving a message calls Touch()
// - Getting a connection calls Touch()
```

---

## Authentication

### User/Password

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

creds := auth.UserPass("admin", "secret-password")
conn, err := client.Connect(ctx, "tenant-1", creds)
```

### Token Authentication

```go
creds := auth.Token("my-static-auth-token")
conn, err := client.Connect(ctx, "tenant-1", creds)
```

### JWT Authentication

```go
// Recommended for production multi-tenant deployments
creds := auth.JWT(jwtToken, nkeySeed)
conn, err := client.Connect(ctx, "tenant-1", creds)

// Or use the shorthand
conn, err = client.ConnectWithJWT(ctx, "tenant-1", jwtToken, nkeySeed)
```

### NKey Authentication

```go
// From file
creds := auth.NKey("/path/to/user.nkey")

// From seed string
creds = auth.NKeySeed("SUAM...")

conn, err := client.Connect(ctx, "tenant-1", creds)
```

### Credential Manager

```go
// Step 1: Create and start
cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
cm.Start()
defer cm.Stop()

// Step 2: Register credentials
cm.Register("tenant-1", jwt1, seed1)

// Step 3: Check if valid
if cm.HasValid("tenant-1") {
    creds, _ := cm.GetCredentials("tenant-1")
    fmt.Printf("Auth type: %v\n", creds.Type)
}

// Step 4: Load from file
cm.LoadFromFile("/path/to/credentials.json")

// Step 5: List tenants
tenants := cm.List()

// Step 6: Remove credentials
cm.Remove("tenant-1")
```

---

## Middleware

### Middleware Stack

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Step 1: Create stack
stack := middleware.NewStack()

// Step 2: Add middleware (order matters - first added is outermost)
stack.UseInterceptor(middleware.NewTracingMiddleware())
stack.UsePublish(middleware.CircuitBreakerMiddleware(cb))
stack.UsePublish(middleware.RateLimitMiddleware(rl))
stack.UseInterceptor(middleware.NewMetricsCollector())
stack.UsePublish(middleware.Retry(retryCfg))
stack.UseInterceptor(middleware.NewLoggingMiddleware())

// Step 3: Wrap handlers
publish := stack.WrapPublish(func(ctx context.Context, subject string, msg *nats.Msg) error {
    return publisher.Publish(ctx, subject, msg)
})

subscribe := stack.WrapSubscribe(func(ctx context.Context, msg *nats.Msg) error {
    return processMessage(ctx, msg)
})

// Step 4: Use wrapped handlers
err := publish(ctx, "orders.created", orderMsg)
```

### Tracing Middleware

```go
// Basic tracing
tracing := middleware.NewTracingMiddleware()
stack.UseInterceptor(tracing)

// With OTEL enabled
tracing = middleware.NewTracingMiddlewareWithOTEL()
stack.UseInterceptor(tracing)

// Trace context is automatically:
// - Injected into message headers on publish
// - Extracted from message headers on subscribe
// Supports: W3C traceparent, B3, custom X-Trace-ID
```

### Retry Middleware

```go
// Default retry (3 attempts, exponential backoff)
retry := middleware.NewRetryMiddleware(middleware.DefaultRetryConfig())
stack.UsePublish(retry.InterceptPublish())

// Custom retry
retry = middleware.NewRetryMiddleware(middleware.RetryConfig{
    MaxAttempts:     5,
    InitialInterval: 200 * time.Millisecond,
    MaxInterval:     10 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.2,
})

// Only retry specific errors
retry.WithShouldRetry(func(err error) bool {
    return errors.Is(err, utils.ErrPublishTimeout) ||
           errors.Is(err, utils.ErrNotConnected)
})

stack.UsePublish(retry.InterceptPublish())
stack.UseSubscribe(retry.InterceptSubscribe())
```

### Circuit Breaker

```go
// Step 1: Create circuit breaker
cb := middleware.NewCircuitBreaker(middleware.CircuitBreakerConfig{
    MaxFailures:      5,             // Open after 5 failures
    ResetTimeout:     30 * time.Second, // Wait before half-open
    HalfOpenRequests: 3,             // Allow 3 test requests
    SuccessThreshold: 2,             // Close after 2 successes
})

// Step 2: Use as middleware
stack.UsePublish(middleware.CircuitBreakerMiddleware(cb))

// Step 3: Monitor state
fmt.Printf("State: %v, Failures: %d\n", cb.State(), cb.Failures())

// Step 4: Per-subject circuit breakers (recommended)
mcb := middleware.NewMultiCircuitBreaker(middleware.DefaultCircuitBreakerConfig())
stack.UsePublish(middleware.MultiCircuitBreakerMiddleware(mcb))
// "orders.create" and "users.get" have independent circuit breakers
```

### Rate Limiter

```go
// Step 1: Token bucket rate limiter
rl := middleware.NewRateLimiter(middleware.RateLimiterConfig{
    Rate:  100,  // 100 tokens per second
    Burst: 200,  // Allow burst of 200
})

// Step 2: Use as middleware (reject when exceeded)
stack.UsePublish(middleware.RateLimitMiddleware(rl))

// Step 3: Or use blocking middleware (wait for token)
stack.UsePublish(middleware.RateLimitWaitMiddleware(rl))

// Step 4: Per-subject rate limiting
psrl := middleware.NewPerSubjectRateLimiter(middleware.RateLimiterConfig{
    Rate:  50,
    Burst: 100,
})
stack.UsePublish(middleware.PerSubjectRateLimitMiddleware(psrl))

// Step 5: Sliding window rate limiter
swl := middleware.NewSlidingWindowLimiter(1000, time.Minute)
stack.UsePublish(middleware.SlidingWindowMiddleware(swl))
```

### Metrics Middleware

```go
// Basic metrics
mc := middleware.NewMetricsCollector()
stack.UseInterceptor(mc)

// OTEL metrics
otelMetrics := middleware.NewOTELMetricsMiddleware("my-service")
stack.UseInterceptor(otelMetrics)

// Per-subject metrics
psm := middleware.NewPerSubjectMetrics()
stack.UseInterceptor(psm)

// Get collected metrics
metrics := mc.Collect()
fmt.Printf("Published: %d, Errors: %d\n",
    metrics.PublishCount, metrics.PublishErrors)
```

### Logging Middleware

```go
// Basic logging
logging := middleware.NewLoggingMiddleware()
stack.UseInterceptor(logging)

// With custom level and payload logging
logging.WithLevel(middleware.LogLevelDebug)
logging.WithPayload(true) // Log message payloads

// OTEL logging
otelLogger := middleware.NewOTELLogger("my-service")
stack.UseInterceptor(middleware.OTELLoggingMiddleware("my-service"))
```

### Custom Middleware

```go
// Simple timing middleware
func TimingMiddleware() middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *nats.Msg) error {
            start := time.Now()
            err := next(ctx, subject, msg)
            duration := time.Since(start)

            if duration > time.Second {
                log.Printf("SLOW: Publish to %s took %v", subject, duration)
            }
            return err
        }
    }
}

stack.UsePublish(TimingMiddleware())
```

---

## Microservice Registry

### Service Registration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"

// Step 1: Create manager
manager, err := microservice.NewManager(microservice.ManagerConfig{
    Publisher: pub, // Optional: for event publishing
    HealthCheck: microservice.HealthCheckConfig{
        Enabled:            true,
        Interval:           30 * time.Second,
        Timeout:            5 * time.Second,
        HealthyThreshold:   2,
        UnhealthyThreshold: 3,
    },
})
defer manager.Shutdown(ctx)

// Step 2: Register service
svc := &microservice.Info{
    ID:           "order-service",
    Name:         "order-service",
    DisplayName:  "Order Service",
    Type:         microservice.ServiceTypeAPI,
    Version:      "2.0.0",
    APIVersion:   "v2",
    Capabilities: []string{"orders", "checkout", "refunds"},
    Dependencies: []string{"user-service", "inventory-service"},
    Tags:         []string{"api", "commerce"},
}
err = manager.RegisterServiceDirect(ctx, svc)

// Step 3: Register distributed handler
handler := microservice.NewHandler(manager)
handler.RegisterQueueSubscriptions(ctx, sub, "service-registry")
```

### Instance Management

```go
// Step 1: Register instance
inst := &microservice.Instance{
    ID:              "order-service-1",
    ServiceID:       "order-service",
    Host:            "10.0.1.50",
    Port:            8080,
    Protocol:        "http",
    Status:          microservice.StatusRunning,
    HealthCheckPath: "/health",
    Weight:          100,
    Region:          "us-east-1",
    Zone:            "us-east-1a",
    PodName:         os.Getenv("POD_NAME"),
    NodeName:        os.Getenv("NODE_NAME"),
}
err := manager.RegisterInstanceDirect(ctx, inst)

// Step 2: Set up heartbeat
go func() {
    ticker := time.NewTicker(10 * time.Second) // 1/3 of TTL
    for range ticker.C {
        manager.RenewInstance(inst.ID, 30*time.Second)
    }
}()

// Step 3: Update status when needed
manager.UpdateInstanceStatus(inst.ID, microservice.StatusDegraded)

// Step 4: Deregister on shutdown
manager.DeregisterInstance(ctx, inst.ID)
```

### Service Discovery

```go
// Discover services by name
services := manager.Discover("user-service",
    microservice.OnlyHealthy(),
)

// Discover by capability
services = manager.Discover("",
    microservice.WithCapabilities("payments"),
    microservice.OnlyHealthy(),
)

// Discover by type
services = manager.Discover("",
    microservice.WithType(microservice.ServiceTypeAPI),
    microservice.WithVersion("2.0.0"),
)

// Get healthy instances for load balancing
instances := manager.DiscoverInstances("user-service", true)
for _, inst := range instances {
    fmt.Printf("%s:%d (weight: %d)\n", inst.Host, inst.Port, inst.Weight)
}

// Check dependencies at startup
for _, dep := range svc.Dependencies {
    if services := manager.Discover(dep, microservice.OnlyHealthy()); len(services) == 0 {
        log.Printf("WARNING: Dependency %s not available", dep)
    }
}
```

### Health Checking

```go
// Health checks run automatically when configured
// Implement /health endpoint in your service:
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":  "healthy",
        "version": "2.0.0",
    })
})

// Monitor health changes
registry := manager.Registry()
registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
    if new == microservice.StatusFailed {
        log.Printf("ALERT: Instance %s failed!", inst.ID)
    }
})
```

---

## Endpoint Registry

### Endpoint Registration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"

// Step 1: Create endpoint manager
epManager, err := endpoint.NewManager(endpoint.ManagerConfig{
    Publisher: pub,
    HealthCheck: endpoint.HealthCheckConfig{
        Enabled:            true,
        Interval:           30 * time.Second,
        Timeout:            5 * time.Second,
        HealthyThreshold:   2,
        UnhealthyThreshold: 3,
    },
})
defer epManager.Shutdown(ctx)

// Step 2: Register endpoint
ep := &endpoint.Info{
    ID:              "user-api-v1",
    Name:            "User API v1",
    ServiceID:       "user-service",
    Host:            "localhost",
    Port:            8080,
    Path:            "/api/v1/users",
    Protocol:        endpoint.ProtocolHTTP,
    Methods:         []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST},
    Status:          endpoint.StatusHealthy,
    HealthCheckPath: "/health",
    Tags:            []string{"api", "users", "v1"},
    Version:         "1.0.0",
}
err = epManager.RegisterEndpoint(ctx, ep)

// Step 3: Set up distributed handler
epHandler := endpoint.NewHandler(epManager)
epHandler.RegisterQueueSubscriptions(ctx, sub, "endpoint-workers")

// Step 4: Set up heartbeat
go func() {
    ticker := time.NewTicker(2 * time.Minute)
    for range ticker.C {
        epClient := endpoint.NewClient(pub)
        epClient.Renew(ctx, "user-api-v1", 5*time.Minute)
    }
}()
```

### Endpoint Discovery

```go
// Discover by service
endpoints := epManager.Discover("user-service",
    endpoint.OnlyHealthy(),
    endpoint.WithTags("api"),
)

// Discover with filters
endpoints = epManager.Discover("user-service",
    endpoint.WithProtocol(endpoint.ProtocolHTTP),
    endpoint.WithVersion("1.0.0"),
    endpoint.OnlyHealthy(),
)

// Use discovered endpoint
for _, ep := range endpoints {
    url := fmt.Sprintf("%s://%s:%d%s", ep.Protocol, ep.Host, ep.Port, ep.Path)
    fmt.Printf("Endpoint: %s [%v]\n", url, ep.Methods)
}

// Remote discovery via NATS
epClient := endpoint.NewClient(pub)
resp, err := epClient.Discover(ctx, "user-service", endpoint.OnlyHealthy())
```

### Endpoint Health Checking

```go
// Manual health update
epManager.UpdateHealth(ctx, "user-api-v1", endpoint.HealthUpdate{
    Status:  endpoint.StatusMaintenance,
    Message: "Scheduled maintenance window",
})

// Monitor health changes
registry := epManager.Registry()
registry.OnStatusChange(func(ep *endpoint.Info, old, new endpoint.Status) {
    log.Printf("Endpoint %s: %s -> %s", ep.ID, old, new)

    if new == endpoint.StatusUnhealthy {
        alerting.Warn("Endpoint unhealthy", ep.ID)
    }
})
```

---

## Context Propagation

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Step 1: Add context values in your HTTP handler
func handleHTTPRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    ctx = core.WithTenantID(ctx, r.Header.Get("X-Tenant-ID"))
    ctx = core.WithCorrelationID(ctx, r.Header.Get("X-Correlation-ID"))
    ctx = core.WithMessageID(ctx, uuid.New().String())

    // Step 2: Publish with enriched context
    // Middleware automatically extracts context into headers
    pub.Publish(ctx, "events.request", msg)
}

// Step 3: In subscriber, context is automatically populated
sub.Subscribe(ctx, "events.>", func(ctx context.Context, msg *nats.Msg) error {
    tenantID, _ := core.TenantIDFromContext(ctx)
    corrID, _ := core.CorrelationIDFromContext(ctx)
    traceCtx, _ := core.TraceContextFromContext(ctx)

    log.Printf("Tenant: %s, Correlation: %s, Trace: %s",
        tenantID, corrID, traceCtx.TraceID)
    return nil
})

// Step 4: Manual header extraction/injection
headers := make(nats.Header)
core.ExtractHeaders(ctx, headers) // Context -> Headers

ctx = core.InjectContext(ctx, headers) // Headers -> Context
```

---

## Error Handling

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Connection errors
conn, err := client.Connect(ctx, tenantID, creds)
if err != nil {
    switch {
    case errors.Is(err, utils.ErrConnectionTimeout):
        log.Println("Connection timed out, retrying...")
    case errors.Is(err, utils.ErrAuthFailed):
        log.Println("Authentication failed, check credentials")
    case errors.Is(err, utils.ErrAlreadyConnected):
        log.Println("Already connected")
    }
}

// Publishing errors
err = pub.Publish(ctx, subject, msg)
if err != nil {
    if utils.IsRetryable(err) {
        // Safe to retry
        time.Sleep(time.Second)
        err = pub.Publish(ctx, subject, msg)
    }
}

// Middleware errors
if errors.Is(err, middleware.ErrCircuitOpen) {
    log.Println("Circuit breaker is open, service is degraded")
}
if errors.Is(err, middleware.ErrRateLimitExceeded) {
    log.Println("Rate limit exceeded, slow down")
}

// Error collection for shutdown
collector := utils.NewErrorCollector()
collector.Add(pub.Close(ctx))
collector.Add(sub.Close(ctx))
collector.Add(client.Shutdown(ctx))
if collector.HasErrors() {
    log.Printf("Shutdown errors: %v", collector.Error())
}
```

---

## Utilities

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Thread-safe map
connections := utils.NewSafeMap[string, *events.Connection]()
connections.Set("tenant-1", conn)
conn, ok := connections.Get("tenant-1")

// Optional values
func findUser(id string) utils.Optional[*User] {
    user, ok := db.Get(id)
    if !ok {
        return utils.None[*User]()
    }
    return utils.Some(user)
}

result := findUser("123")
user := result.ValueOr(defaultUser)

// Safe goroutines
utils.SafeGo(func() {
    // Won't crash the app if it panics
    processMessages()
})

// Error classification
if utils.IsRetryable(err) { /* retry */ }
if utils.IsTemporary(err) { /* wait and retry */ }
```

---

## Best Practices

### 1. Always Use Graceful Shutdown

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    <-sigChan
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Shutdown in reverse order of creation
    epManager.Shutdown(ctx)
    svcManager.Shutdown(ctx)
    batch.Close(ctx)
    pub.Close(ctx)
    sub.Close(ctx)
    client.Shutdown(ctx)
}()
```

### 2. Use Middleware for Production

```go
stack := middleware.NewStack()
stack.UseInterceptor(middleware.NewTracingMiddlewareWithOTEL())
stack.UsePublish(middleware.CircuitBreakerMiddleware(cb))
stack.UsePublish(middleware.RateLimitMiddleware(rl))
stack.UseInterceptor(middleware.NewOTELMetricsMiddleware("svc"))
stack.UsePublish(middleware.Retry(middleware.RetryConfig{MaxAttempts: 3}))
stack.UseInterceptor(middleware.NewLoggingMiddleware())
```

### 3. Use Queue Groups for Scalability

```go
// Every instance of your service should use the same queue group
sub.QueueSubscribe(ctx, "orders.>", "order-service", handler)
// Messages are load-balanced across all instances
```

### 4. Implement Health Endpoints

```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    health := client.Manager().Health()
    allHealthy := true
    for _, status := range health {
        if !status.IsHealthy() {
            allHealthy = false
        }
    }
    if allHealthy {
        w.WriteHeader(http.StatusOK)
    } else {
        w.WriteHeader(http.StatusServiceUnavailable)
    }
})
```

### 5. Use Batch Publishing for High Throughput

```go
// Instead of publishing one at a time:
for _, msg := range messages {
    pub.Publish(ctx, subject, msg) // Slow
}

// Use batch publisher:
batch := events.NewBatchPublisher(pub,
    events.WithMaxBatchSize(100),
    events.WithFlushInterval(time.Second),
)
for _, msg := range messages {
    batch.Add(subject, msg) // Fast, batched
}
batch.Flush(ctx)
```

### 6. Use Per-Subject Middleware for Isolation

```go
// Per-subject circuit breakers prevent one failing subject from
// blocking all other subjects
mcb := middleware.NewMultiCircuitBreaker(middleware.DefaultCircuitBreakerConfig())
stack.UsePublish(middleware.MultiCircuitBreakerMiddleware(mcb))
```

### 7. Monitor Error Channels

```go
go func() {
    for err := range client.Manager().Errors() {
        log.Printf("Tenant %s error: %v", err.TenantID, err.Err)
        metrics.IncrementCounter("tenant_errors", err.TenantID)
    }
}()
```

### 8. Use JetStream for Reliable Messaging

```go
// Use JetStream when you need:
// - Message persistence (survives server restarts)
// - At-least-once delivery
// - Message replay
// - Deduplication

// Use Core NATS when you need:
// - Lowest latency
// - Fire-and-forget
// - Simple pub/sub without persistence
```
