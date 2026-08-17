# Events Package - Usage Guide

This guide provides detailed examples for using the events package and all its submodules.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Configuration](#configuration)
3. [Core Module](#core-module)
   - [Messages](#messages)
   - [Headers](#headers)
   - [Connection Pool](#connection-pool)
4. [Transport Layer (NATS)](#transport-layer-nats)
   - [Basic Connection](#basic-connection)
   - [Publishing Messages](#publishing-messages)
   - [Subscribing to Messages](#subscribing-to-messages)
   - [Request-Reply Pattern](#request-reply-pattern)
   - [Batch Publishing](#batch-publishing)
5. [JetStream](#jetstream)
   - [Stream Management](#stream-management)
   - [Pull Consumers](#pull-consumers)
   - [Message Acknowledgment](#message-acknowledgment)
6. [Tenant Management](#tenant-management)
   - [Multi-Tenant Connections](#multi-tenant-connections)
   - [Credential Management](#credential-management)
7. [Microservice Registry](#microservice-registry)
   - [Service Registration](#service-registration)
   - [Instance Management](#instance-management)
   - [Service Discovery](#service-discovery)
   - [Health Checking](#health-checking)
8. [Endpoint Registry](#endpoint-registry)
   - [Endpoint Registration](#endpoint-registration)
   - [Endpoint Discovery](#endpoint-discovery)
9. [Middleware](#middleware)
   - [Retry Middleware](#retry-middleware)
   - [Circuit Breaker](#circuit-breaker)
   - [Tracing Middleware](#tracing-middleware)
   - [Metrics Middleware](#metrics-middleware)
   - [Custom Middleware](#custom-middleware)
10. [Authentication](#authentication)
    - [User/Password](#userpassword)
    - [Token Authentication](#token-authentication)
    - [JWT Authentication](#jwt-authentication)
    - [NKey Authentication](#nkey-authentication)
11. [Encoding](#encoding)
    - [JSON Encoding](#json-encoding)
    - [Protobuf Encoding](#protobuf-encoding)
12. [Utils Package](#utils-package)
13. [Error Handling](#error-handling)
14. [Best Practices](#best-practices)

---

## Quick Start

```go
package main

import (
    "context"
    "log"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

func main() {
    ctx := context.Background()

    // Create configuration
    cfg := config.DefaultConfig()
    cfg.Servers = []string{"nats://localhost:4222"}

    // Create credentials
    creds := events.UserPassCredentials("user", "password")

    // Create client
    client, err := events.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Shutdown(ctx)

    // Connect a tenant
    conn, err := client.Connect(ctx, "tenant-1", creds)
    if err != nil {
        log.Fatal(err)
    }

    // Publish a message
    msg := core.NewMessage([]byte(`{"event": "user.created"}`))
    err = conn.Publisher().Publish(ctx, "events.user", msg)
    if err != nil {
        log.Fatal(err)
    }

    // Subscribe to messages
    sub, err := conn.Subscriber().Subscribe(ctx, "events.>", func(ctx context.Context, msg *core.Message) error {
        log.Printf("Received: %s", string(msg.Data))
        return nil
    })
    if err != nil {
        log.Fatal(err)
    }
    defer sub.Unsubscribe()

    // Keep running
    time.Sleep(10 * time.Second)
}
```

---

## Configuration

### Default Configuration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"

// Get default configuration
cfg := config.DefaultConfig()

// Defaults:
// - Servers: ["nats://localhost:4222"]
// - ConnectTimeout: 10s
// - DrainTimeout: 30s
// - IdleTimeout: 30m
// - CleanupInterval: 1m
// - Stream.Enabled: true
// - Reconnect: infinite with exponential backoff
```

### Custom Configuration

```go
cfg := &config.Config{
    // Server URLs
    Servers: []string{
        "nats://server1:4222",
        "nats://server2:4222",
    },
    
    // Client identifier
    Name: "my-service",
    
    // Timeouts
    ConnectTimeout:  15 * time.Second,
    DrainTimeout:    60 * time.Second,
    IdleTimeout:     1 * time.Hour,
    CleanupInterval: 5 * time.Minute,
    
    // TLS Configuration
    TLS: &config.TLSConfig{
        Enabled:    true,
        CertFile:   "/path/to/cert.pem",
        KeyFile:    "/path/to/key.pem",
        CAFile:     "/path/to/ca.pem",
        SkipVerify: false,
    },
    
    // Reconnection settings
    Reconnect: config.ReconnectConfig{
        MaxAttempts:     -1,                    // Infinite
        InitialInterval: 100 * time.Millisecond,
        MaxInterval:     30 * time.Second,
        Multiplier:      2.0,
        Jitter:          0.1,
    },
    
    // JetStream settings
    Stream: config.StreamConfig{
        Enabled: true,
        Domain:  "hub",
        Prefix:  "",
    },
    
    // Publishing settings
    Publish: config.PublishConfig{
        Retry: config.RetryConfig{
            MaxAttempts:     3,
            InitialInterval: 100 * time.Millisecond,
            MaxInterval:     10 * time.Second,
            Multiplier:      2.0,
        },
        AckTimeout:          5 * time.Second,
        EnableDeduplication: true,
        DeduplicationWindow: 2 * time.Minute,
    },
    
    // Subscription settings
    Subscribe: config.SubscribeConfig{
        QueueGroup:    "my-service-group",
        MaxConcurrent: 10,
        AckWait:       30 * time.Second,
        BatchSize:     100,
        BatchWait:     100 * time.Millisecond,
    },
}

// Validate configuration
if err := config.Validate(cfg); err != nil {
    log.Fatal(err)
}
```

### Loading from Environment

```go
// Load from environment variables
cfg, err := config.LoadFromEnv()
if err != nil {
    log.Fatal(err)
}

// Environment variables:
// NATS_SERVERS=nats://host1:4222,nats://host2:4222
// NATS_NAME=my-service
// NATS_CONNECT_TIMEOUT=15s
// NATS_TLS_ENABLED=true
// NATS_TLS_CERT_FILE=/path/to/cert.pem
// etc.
```

---

## Core Module

### Messages

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Create a new message
msg := core.NewMessage([]byte(`{"key": "value"}`))

// Set subject
msg = msg.WithSubject("events.user.created")

// Set message ID (for deduplication)
msg = msg.WithID("unique-message-id-123")

// Set reply subject (for request-reply)
msg = msg.WithReply("_INBOX.xyz")

// Add headers
msg = msg.WithHeader("Content-Type", "application/json")
msg = msg.WithHeader("X-Tenant-ID", "tenant-123")

// Add multiple headers
headers := core.Headers{
    "X-Trace-ID": []string{"trace-abc"},
    "X-Span-ID":  []string{"span-def"},
}
msg = msg.WithHeaders(headers)

// Access message properties
subject := msg.Subject
data := msg.Data
msgID := msg.ID
timestamp := msg.Timestamp
```

### Headers

```go
// Create headers
headers := make(core.Headers)

// Set a header (replaces existing)
headers.Set("Content-Type", "application/json")

// Add a header value (supports multiple values)
headers.Add("Accept-Encoding", "gzip")
headers.Add("Accept-Encoding", "deflate")

// Get first value
contentType := headers.Get("Content-Type")

// Get all values
encodings := headers.Values("Accept-Encoding")

// Check if header exists
if headers.Has("X-Custom") {
    // ...
}

// Delete header
headers.Del("X-Temp")

// Clone headers
clonedHeaders := headers.Clone()
```

### Message Pooling (Performance)

```go
// Acquire message from pool (reuses memory)
msg := core.AcquireMessage()
msg.Subject = "events.user.created"
msg.Data = []byte(`{"user": "john"}`)

// Use the message...
publisher.Publish(ctx, msg.Subject, msg)

// Release back to pool
core.ReleaseMessage(msg)

// Acquire with data
msg2 := core.AcquireMessageWithData([]byte(`{"key": "value"}`))
defer core.ReleaseMessage(msg2)

// Headers pooling
headers := core.AcquireHeaders()
headers.Set("Content-Type", "application/json")
defer core.ReleaseHeaders(headers)
```

### Connection Pool

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Create pool configuration
poolConfig := core.PoolConfig{
    MinSize:             2,
    MaxSize:             10,
    MaxIdleTime:         30 * time.Minute,
    MaxLifetime:         2 * time.Hour,
    AcquireTimeout:      30 * time.Second,
    HealthCheckInterval: 30 * time.Second,
    HealthCheckTimeout:  5 * time.Second,
    WaitForConnection:   true,
    PreWarm:             true,
    
    // Connection factory
    Factory: func(ctx context.Context) (core.Connection, error) {
        return nats.NewConnection("", cfg, creds), nil
    },
    
    // Callbacks
    OnAcquire: func(conn core.PooledConnection) {
        log.Printf("Connection %s acquired", conn.ID())
    },
    OnRelease: func(conn core.PooledConnection) {
        log.Printf("Connection %s released", conn.ID())
    },
    OnClose: func(conn core.PooledConnection) {
        log.Printf("Connection %s closed", conn.ID())
    },
}

// Create pool
pool, err := core.NewPool(poolConfig)
if err != nil {
    log.Fatal(err)
}
defer pool.Close(ctx)

// Get connection from pool
conn, err := pool.Get(ctx)
if err != nil {
    log.Fatal(err)
}

// Use connection...

// Return to pool
pool.Put(conn)

// Get pool statistics
stats := pool.Stats()
log.Printf("Pool size: %d, Available: %d, In use: %d",
    stats.Size, stats.Available, stats.InUse)
log.Printf("Total acquired: %d, Average acquire time: %v",
    stats.TotalAcquired, stats.AverageAcquireTime)

// Custom health checker
pool.SetHealthChecker(func(ctx context.Context, conn core.PooledConnection) bool {
    // Custom health check logic
    return conn.IsConnected() && checkCustomHealth(conn)
})
```

---

## Transport Layer (NATS)

### Basic Connection

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
)

// Create connection
conn := nats.NewConnection("tenant-1", cfg, creds)

// Connect
if err := conn.Connect(ctx); err != nil {
    log.Fatal(err)
}
defer conn.Close(ctx)

// Check connection state
if conn.IsConnected() {
    log.Println("Connected!")
}

// Get connection state
state := conn.State() // StateConnected, StateReconnecting, etc.

// Get health status
health := conn.Health()
log.Printf("Healthy: %v, State: %s, Message: %s",
    health.Healthy, health.State, health.Message)

// Set error callback
conn.OnError(func(err error) {
    log.Printf("Connection error: %v", err)
})
```

### Publishing Messages

```go
// Create publisher
publisher := nats.NewPublisher(conn)
defer publisher.Close(ctx)

// Simple publish
msg := core.NewMessage([]byte(`{"event": "created"}`))
err := publisher.Publish(ctx, "events.user.created", msg)
if err != nil {
    log.Fatal(err)
}

// Publish with headers
msg = msg.WithHeader("X-Event-Type", "user.created")
msg = msg.WithHeader("X-Tenant-ID", "tenant-123")
err = publisher.Publish(ctx, "events.user.created", msg)

// Async publish (JetStream)
future := publisher.PublishAsync(ctx, "events.user.created", msg)

// Wait for acknowledgment
ack, err := future.Wait(ctx)
if err != nil {
    log.Fatal(err)
}
log.Printf("Published to stream: %s, sequence: %d", ack.Stream, ack.Sequence)

// Or use channels
select {
case ack := <-future.Ok():
    log.Printf("Published: %d", ack.Sequence)
case err := <-future.Err():
    log.Printf("Failed: %v", err)
case <-ctx.Done():
    log.Printf("Timeout")
}
```

### Subscribing to Messages

```go
// Create subscriber
subscriber := nats.NewSubscriber(conn)
defer subscriber.Close(ctx)

// Simple subscription
sub, err := subscriber.Subscribe(ctx, "events.>", func(ctx context.Context, msg *core.Message) error {
    log.Printf("Received on %s: %s", msg.Subject, string(msg.Data))
    return nil
})
if err != nil {
    log.Fatal(err)
}
defer sub.Unsubscribe()

// Queue subscription (load balancing)
sub, err = subscriber.QueueSubscribe(ctx, "events.orders.>", "order-processors", 
    func(ctx context.Context, msg *core.Message) error {
        // Process order - only one subscriber in the queue group receives each message
        return processOrder(msg)
    })

// Wildcard subscriptions
// events.* - matches events.user, events.order (single token)
// events.> - matches events.user.created, events.order.placed (multiple tokens)

// Check subscription status
if sub.IsValid() {
    log.Printf("Subscribed to: %s", sub.Subject())
}

// Drain subscription (process pending, then unsubscribe)
if err := sub.Drain(); err != nil {
    log.Printf("Drain error: %v", err)
}
```

### Request-Reply Pattern

```go
// Send request and wait for reply
request := core.NewMessage([]byte(`{"user_id": "123"}`))
reply, err := publisher.Request(ctx, "api.users.get", request)
if err != nil {
    if errors.Is(err, core.ErrRequestTimeout) {
        log.Println("Request timed out")
    } else if errors.Is(err, core.ErrNoReply) {
        log.Println("No responders available")
    }
    log.Fatal(err)
}
log.Printf("Response: %s", string(reply.Data))

// Request with explicit timeout
reply, err = publisher.RequestWithTimeout(ctx, "api.users.get", request, 5*time.Second)

// Setting up a responder
subscriber.Subscribe(ctx, "api.users.get", func(ctx context.Context, msg *core.Message) error {
    // Parse request
    var req struct {
        UserID string `json:"user_id"`
    }
    json.Unmarshal(msg.Data, &req)
    
    // Get user
    user := getUserByID(req.UserID)
    
    // Send reply
    response, _ := json.Marshal(user)
    replyMsg := core.NewMessage(response)
    return publisher.Publish(ctx, msg.Reply, replyMsg)
})
```

### Batch Publishing

```go
// Create batch publisher
batch := nats.NewBatchPublisher(publisher,
    nats.WithMaxBatchSize(100),          // Auto-flush at 100 messages
    nats.WithFlushInterval(time.Second), // Auto-flush every second
    nats.WithConcurrentFlush(true),      // Publish concurrently
)
defer batch.Close(ctx)

// Add messages to batch
for i := 0; i < 1000; i++ {
    msg := core.NewMessage([]byte(fmt.Sprintf(`{"index": %d}`, i)))
    if err := batch.Add("events.batch", msg); err != nil {
        log.Printf("Add error: %v", err)
    }
}

// Manual flush
if err := batch.Flush(ctx); err != nil {
    log.Printf("Flush error: %v", err)
}

// Check pending count
log.Printf("Pending messages: %d", batch.Count())
```

---

## JetStream

### Stream Management

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Stream configuration
streamCfg := core.StreamConfig{
    Name:        "EVENTS",
    Description: "Application events stream",
    Subjects:    []string{"events.>"},
    
    // Retention
    Retention: core.LimitsPolicy, // or InterestPolicy, WorkQueuePolicy
    MaxMsgs:   1000000,
    MaxBytes:  1024 * 1024 * 1024, // 1GB
    MaxAge:    24 * time.Hour,
    
    // Behavior
    Storage:   core.FileStorage, // or MemoryStorage
    Replicas:  3,
    Discard:   core.DiscardOld, // or DiscardNew
    
    // Deduplication
    DuplicateWindow: 2 * time.Minute,
    
    // Features
    AllowRollup:  false,
    DenyDelete:   false,
    DenyPurge:    false,
    AllowDirect:  true,
    
    // Metadata
    Metadata: map[string]string{
        "environment": "production",
        "team":        "platform",
    },
}

// Create stream (via stream manager interface)
stream, err := streamManager.CreateStream(ctx, streamCfg)
if err != nil {
    log.Fatal(err)
}

// Get stream info
info, err := stream.Info(ctx)
if err != nil {
    log.Fatal(err)
}
log.Printf("Stream %s has %d messages (%d bytes)",
    info.Config.Name, info.State.Messages, info.State.Bytes)

// Purge stream
err = stream.Purge(ctx)

// Purge by subject
err = stream.Purge(ctx, core.WithPurgeSubject("events.user.*"))

// Purge keeping last N messages
err = stream.Purge(ctx, core.WithPurgeKeep(1000))

// Delete specific message
err = stream.DeleteMessage(ctx, 12345)

// Get message by sequence
rawMsg, err := stream.GetMessage(ctx, 12345)

// Get last message for subject
rawMsg, err = stream.GetLastMessageForSubject(ctx, "events.user.created")
```

### Pull Consumers

```go
// Consumer configuration
consumerCfg := core.ConsumerConfig{
    Name:        "event-processor",
    Description: "Processes application events",
    
    // Delivery
    DeliverPolicy: core.DeliverAll, // or DeliverLast, DeliverNew, etc.
    AckPolicy:     core.AckExplicit,
    AckWait:       30 * time.Second,
    
    // Filtering
    FilterSubject: "events.user.>",
    
    // Limits
    MaxDeliver:    5,
    MaxAckPending: 1000,
    MaxWaiting:    512,
    
    // Backoff for redeliveries
    BackOff: []time.Duration{
        1 * time.Second,
        5 * time.Second,
        30 * time.Second,
        1 * time.Minute,
    },
    
    // Metadata
    Metadata: map[string]string{
        "processor": "main",
    },
}

// Create consumer
consumer, err := stream.CreateConsumer(ctx, consumerCfg)
if err != nil {
    log.Fatal(err)
}

// Fetch messages (blocking)
batch, err := consumer.Fetch(100, 
    core.WithFetchMaxWait(5*time.Second),
)
if err != nil {
    log.Fatal(err)
}

for msg := range batch.Messages() {
    // Process message
    log.Printf("Processing: %s", string(msg.Data()))
    
    // Acknowledge
    if err := msg.Ack(); err != nil {
        log.Printf("Ack error: %v", err)
    }
}

// Fetch by bytes
batch, err = consumer.FetchBytes(1024*1024) // 1MB max

// Fetch without waiting (returns immediately available)
batch, err = consumer.FetchNoWait(100)

// Get single message
msg, err := consumer.Next(core.WithFetchMaxWait(time.Second))
if err != nil {
    if errors.Is(err, core.ErrNoMessages) {
        log.Println("No messages available")
    }
}

// Continuous consumption with handler
consumeCtx, err := consumer.Consume(func(msg core.Msg) {
    defer msg.Ack()
    
    // Process message
    processEvent(msg.Data())
}, 
    core.WithConsumeMaxMessages(100),
    core.WithConsumeExpiry(30*time.Second),
    core.WithConsumeErrorHandler(func(cc core.ConsumeContext, err error) {
        log.Printf("Consume error: %v", err)
    }),
)
if err != nil {
    log.Fatal(err)
}

// Stop consumption
consumeCtx.Stop()

// Or drain (finish pending, then stop)
consumeCtx.Drain()

// Iterator pattern
iter, err := consumer.Messages(
    core.WithStopAfter(1000),
    core.WithIteratorHeartbeat(5*time.Second),
)
if err != nil {
    log.Fatal(err)
}
defer iter.Stop()

for {
    msg, err := iter.Next()
    if err != nil {
        break
    }
    processMessage(msg)
    msg.Ack()
}
```

### Message Acknowledgment

```go
// Positive acknowledgment
err := msg.Ack()

// Double ack (wait for server confirmation)
err := msg.DoubleAck(ctx)

// Negative ack (redeliver immediately)
err := msg.Nak()

// Negative ack with delay
err := msg.NakWithDelay(5 * time.Second)

// Terminate (no more redeliveries)
err := msg.Term()

// Terminate with reason
err := msg.TermWithReason("invalid message format")

// Signal in progress (extend ack deadline)
err := msg.InProgress()

// Get message metadata
meta, err := msg.Metadata()
if err == nil {
    log.Printf("Stream: %s, Consumer: %s, Sequence: %d, Delivered: %d",
        meta.Stream, meta.Consumer, meta.Sequence.Stream, meta.NumDelivered)
}
```

---

## Tenant Management

### Multi-Tenant Connections

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
)

// Create tenant manager
managerCfg := tenant.ManagerConfig{
    Config:          cfg,
    Factory:         &nats.Factory{},
    ErrorBufferSize: 100,
}

manager, err := tenant.NewManager(managerCfg)
if err != nil {
    log.Fatal(err)
}
defer manager.Shutdown(ctx)

// Connect tenant
conn, err := manager.Connect(ctx, "tenant-1", creds1)
if err != nil {
    log.Fatal(err)
}

// Get existing connection
conn, ok := manager.Get("tenant-1")
if !ok {
    log.Println("Tenant not connected")
}

// Disconnect tenant
err = manager.Disconnect(ctx, "tenant-1")

// Get all connected tenant IDs
tenantIDs := manager.TenantIDs()

// Get active connection count
count := manager.ActiveConnections()

// Health status for all tenants
health := manager.Health()
for tenantID, status := range health {
    log.Printf("Tenant %s: healthy=%v, state=%s",
        tenantID, status.Healthy, status.State)
}

// Health for specific tenant
status, ok := manager.TenantHealth("tenant-1")

// Handle tenant errors
go func() {
    for err := range manager.Errors() {
        log.Printf("Tenant %s error: %v", err.TenantID, err.Err)
    }
}()

// Or use callback
manager.OnError(func(err core.TenantError) {
    log.Printf("Tenant %s: %s: %v", err.TenantID, err.Op, err.Err)
    
    // Auto-reconnect on error
    go func() {
        time.Sleep(5 * time.Second)
        manager.Connect(ctx, err.TenantID, nil)
    }()
})
```

### Credential Management

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

// Create credential manager
credManager := auth.NewCredentialManager()

// Store credentials
creds := auth.JWT("jwt-token", "nkey-seed")
credManager.StoreCredentials("tenant-1", creds)

// Get credentials
creds, err := credManager.GetCredentials("tenant-1")

// Register from payload
payload := auth.RegistrationPayload{
    TenantID: "tenant-2",
    JWT:      "jwt-token-for-tenant-2",
    Seed:     "SUAG...",
}
if err := payload.Validate(); err != nil {
    log.Fatal(err)
}
credManager.RegisterFromPayload(payload)

// Connect using credential manager
manager.SetCredentialManager(credManager)
conn, err := manager.ConnectWithCredentialManager(ctx, "tenant-1")

// Or connect with payload directly
conn, err := manager.ConnectWithPayload(ctx, payload)

// Remove credentials
credManager.RemoveCredentials("tenant-1")

// Check if credentials exist
if credManager.HasCredentials("tenant-1") {
    // ...
}
```

---

## Microservice Registry

### Service Registration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"

// Create registry
registry := microservice.NewRegistry()

// Create manager
mgr := microservice.NewManager(microservice.ManagerConfig{
    Registry:            registry,
    HealthCheckInterval: 30 * time.Second,
    HealthCheckTimeout:  5 * time.Second,
})
defer mgr.Shutdown(ctx)

// Define service
service := &microservice.Info{
    ID:      "user-service-1",
    Name:    "user-service",
    Version: "1.0.0",
    Type:    microservice.ServiceTypeAPI,
    Status:  microservice.StatusRunning,
    
    Tags:         []string{"production", "us-east"},
    Capabilities: []string{"oauth2", "rbac"},
    Dependencies: []string{"database-service", "cache-service"},
    
    Metadata: map[string]string{
        "team":     "platform",
        "oncall":   "platform@example.com",
    },
}

// Register service
req := &microservice.RegistrationRequest{
    Service:  service,
    Override: false,
}
if err := req.Validate(); err != nil {
    log.Fatal(err)
}

err := mgr.RegisterService(ctx, req)
if err != nil {
    log.Fatal(err)
}

// Update service
service.Version = "1.1.0"
err = mgr.UpdateService(ctx, service)

// Deregister service
err = mgr.DeregisterService(ctx, &microservice.DeregistrationRequest{
    ServiceID: "user-service-1",
    Reason:    "Scaling down",
})
```

### Instance Management

```go
// Define instance
instance := &microservice.Instance{
    ID:        "user-service-1-instance-a",
    ServiceID: "user-service-1",
    Host:      "10.0.1.100",
    Port:      8080,
    Status:    microservice.StatusRunning,
    
    HealthCheckPath: "/health",
    Weight:          100,
    
    // Resource metrics
    CPUUsage:    45.5,
    MemoryUsage: 62.3,
    
    Metadata: map[string]string{
        "pod":       "user-service-abc123",
        "node":      "node-1",
        "zone":      "us-east-1a",
    },
}

// Register instance
req := &microservice.InstanceRegistrationRequest{
    Instance: instance,
    TTL:      5 * time.Minute, // Auto-expire if not renewed
}
err := mgr.RegisterInstance(ctx, req)

// Update instance status
instance.Status = microservice.StatusDegraded
instance.CPUUsage = 95.0
err = mgr.UpdateInstance(ctx, instance)

// Deregister instance
err = mgr.DeregisterInstance(ctx, &microservice.InstanceDeregistrationRequest{
    InstanceID: instance.ID,
    ServiceID:  service.ID,
    Reason:     "Graceful shutdown",
})
```

### Service Discovery

```go
// Discover services by name
services := mgr.Discover("user-service")

// Discover with options
services = mgr.Discover("user-service",
    microservice.WithType(microservice.ServiceTypeAPI),
    microservice.WithCapabilities("oauth2"),
    microservice.WithTags("production"),
    microservice.OnlyHealthy(),
    microservice.WithMinInstances(2),
)

// Get all services
allServices := mgr.ListServices()

// Get service by ID
service, err := mgr.GetService(ctx, "user-service-1")

// Get instances of a service
instances, err := mgr.GetInstances(ctx, "user-service-1")

// Query services
query := &microservice.Query{
    Type:        microservice.ServiceTypeAPI,
    Tags:        []string{"production"},
    OnlyHealthy: true,
}
services, err = mgr.Query(ctx, query)
```

### Health Checking

```go
// Set up health check callbacks
mgr.OnServiceHealthy(func(event microservice.Event) {
    log.Printf("Service %s is now healthy", event.Service.ID)
})

mgr.OnServiceUnhealthy(func(event microservice.Event) {
    log.Printf("Service %s is unhealthy: %s", 
        event.Service.ID, event.Reason)
    
    // Alert ops team
    sendAlert(event)
})

mgr.OnInstanceHealthy(func(event microservice.Event) {
    log.Printf("Instance %s of %s is healthy",
        event.Instance.ID, event.Service.ID)
})

mgr.OnInstanceUnhealthy(func(event microservice.Event) {
    log.Printf("Instance %s is unhealthy", event.Instance.ID)
})

// Manual health check
result := mgr.CheckHealth(ctx, "user-service-1")
if !result.Healthy {
    log.Printf("Health check failed: %s", result.Message)
}

// Get health summary
summary := mgr.HealthSummary()
log.Printf("Total: %d, Healthy: %d, Unhealthy: %d",
    summary.Total, summary.Healthy, summary.Unhealthy)
```

### Service Lifecycle Events

```go
// Register callbacks for all lifecycle events
mgr.OnServiceRegistered(func(event microservice.Event) {
    log.Printf("Service registered: %s", event.Service.Name)
})

mgr.OnServiceDeregistered(func(event microservice.Event) {
    log.Printf("Service deregistered: %s (reason: %s)", 
        event.Service.Name, event.Reason)
})

mgr.OnServiceUpdated(func(event microservice.Event) {
    log.Printf("Service updated: %s", event.Service.Name)
})

mgr.OnServiceStatusChange(func(event microservice.Event) {
    log.Printf("Service %s status: %s -> %s",
        event.Service.Name, event.OldStatus, event.NewStatus)
})

mgr.OnInstanceRegistered(func(event microservice.Event) {
    log.Printf("Instance registered: %s (%s:%d)",
        event.Instance.ID, event.Instance.Host, event.Instance.Port)
})

mgr.OnInstanceDeregistered(func(event microservice.Event) {
    log.Printf("Instance deregistered: %s", event.Instance.ID)
})

mgr.OnInstanceStatusChange(func(event microservice.Event) {
    log.Printf("Instance %s status: %s -> %s",
        event.Instance.ID, event.OldStatus, event.NewStatus)
})
```

---

## Endpoint Registry

### Endpoint Registration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"

// Create registry and manager
registry := endpoint.NewRegistry()
mgr := endpoint.NewManager(endpoint.ManagerConfig{
    Registry:            registry,
    HealthCheckInterval: 30 * time.Second,
})
defer mgr.Shutdown(ctx)

// Define endpoint
ep := &endpoint.Info{
    ID:         "user-api-get-users",
    ServiceID:  "user-service-1",
    InstanceID: "user-service-1-instance-a",
    
    Host:     "10.0.1.100",
    Port:     8080,
    Path:     "/api/v1/users",
    Protocol: endpoint.ProtocolHTTPS,
    Methods:  []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST},
    
    Status:          endpoint.StatusHealthy,
    HealthCheckPath: "/health",
    
    // Load balancing
    Weight:   100,
    Priority: 1,
    
    // Geographic
    Region: "us-east-1",
    Zone:   "us-east-1a",
    
    // Versioning
    Version:    "1.0.0",
    APIVersion: "v1",
    
    Tags: []string{"public", "rest"},
    Metadata: map[string]string{
        "rateLimit": "1000/min",
        "auth":      "jwt",
    },
}

// Register endpoint
req := &endpoint.RegistrationRequest{
    Endpoint: ep,
    TTL:      5 * time.Minute,
}
err := mgr.RegisterEndpoint(ctx, req)

// Update endpoint
ep.Weight = 50 // Reduce traffic
err = mgr.UpdateEndpoint(ctx, ep)

// Deregister endpoint
err = mgr.DeregisterEndpoint(ctx, &endpoint.DeregistrationRequest{
    EndpointID: ep.ID,
    ServiceID:  ep.ServiceID,
    Reason:     "Maintenance",
})
```

### Endpoint Discovery

```go
// Discover endpoints for a service
endpoints := mgr.Discover("user-service-1")

// Discover with options
endpoints = mgr.Discover("user-service-1",
    endpoint.WithProtocol(endpoint.ProtocolGRPC),
    endpoint.WithRegion("us-east-1"),
    endpoint.OnlyHealthy(),
    endpoint.WithMethod(endpoint.MethodGET),
)

// Query endpoints
query := &endpoint.Query{
    ServiceIDs:    []string{"user-service", "auth-service"},
    Protocol:      endpoint.ProtocolHTTPS,
    OnlyHealthy:   true,
    OnlyAvailable: true,
    Tags:          []string{"public"},
    Region:        "us-east-1",
}
endpoints, err := mgr.Query(ctx, query)

// Get endpoint by ID
ep, err := mgr.GetEndpoint(ctx, "user-api-get-users")

// Get endpoints by service
endpoints, err = mgr.GetEndpointsByService(ctx, "user-service-1")

// Load balancing: get next endpoint
ep = mgr.NextEndpoint("user-service-1") // Round-robin
ep = mgr.WeightedEndpoint("user-service-1") // Weighted random
```

---

## Middleware

### Retry Middleware

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Create retry middleware
retryMW := middleware.NewRetryMiddleware(middleware.RetryConfig{
    MaxAttempts:     3,
    InitialInterval: 100 * time.Millisecond,
    MaxInterval:     10 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.1,
    
    // Custom retry condition
    ShouldRetry: func(err error) bool {
        // Don't retry on validation errors
        if errors.Is(err, core.ErrInvalidMessage) {
            return false
        }
        return core.IsRetryable(err)
    },
    
    // Callback on retry
    OnRetry: func(attempt int, err error) {
        log.Printf("Retry attempt %d: %v", attempt, err)
    },
})

// Apply to publisher
publisher.UseMiddleware(retryMW.Publish())

// Apply to subscriber
subscriber.UseMiddleware(retryMW.Subscribe())
```

### Circuit Breaker

```go
// Create circuit breaker
cb := middleware.NewCircuitBreaker(middleware.CircuitBreakerConfig{
    Name:              "nats-publisher",
    FailureThreshold:  5,
    SuccessThreshold:  3,
    Timeout:           30 * time.Second,
    HalfOpenMaxCalls:  3,
    
    // Optional: classify errors
    IsFailure: func(err error) bool {
        // Don't count client errors as failures
        if errors.Is(err, core.ErrInvalidMessage) {
            return false
        }
        return err != nil
    },
    
    // State change callback
    OnStateChange: func(from, to middleware.CircuitState) {
        log.Printf("Circuit breaker: %s -> %s", from, to)
        if to == middleware.CircuitOpen {
            sendAlert("Circuit breaker opened!")
        }
    },
})

// Get circuit breaker middleware
cbMW := middleware.NewCircuitBreakerMiddleware(cb)
publisher.UseMiddleware(cbMW.Publish())

// Check circuit state
state := cb.State() // CircuitClosed, CircuitOpen, CircuitHalfOpen

// Manual control
cb.Reset()
cb.Trip() // Force open
```

### Tracing Middleware

```go
// Create tracing middleware (uses OTEL)
tracingMW := middleware.NewTracingMiddleware(middleware.TracingConfig{
    ServiceName: "my-service",
    
    // Custom span naming
    SpanNameFunc: func(subject string) string {
        return fmt.Sprintf("NATS %s", subject)
    },
    
    // Add custom attributes
    AttributeFunc: func(msg *core.Message) []attribute.KeyValue {
        return []attribute.KeyValue{
            attribute.String("tenant.id", msg.Headers.Get("X-Tenant-ID")),
        }
    },
})

publisher.UseMiddleware(tracingMW.Publish())
subscriber.UseMiddleware(tracingMW.Subscribe())
```

### Metrics Middleware

```go
// Create metrics middleware
metricsMW := middleware.NewMetricsMiddleware(middleware.MetricsConfig{
    Namespace: "myapp",
    Subsystem: "nats",
    
    // Histogram buckets for latency
    LatencyBuckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
    
    // Custom labels
    LabelFunc: func(msg *core.Message) map[string]string {
        return map[string]string{
            "tenant": msg.Headers.Get("X-Tenant-ID"),
        }
    },
})

publisher.UseMiddleware(metricsMW.Publish())
subscriber.UseMiddleware(metricsMW.Subscribe())

// Access metrics
publishCount := metricsMW.PublishCount("events.user.created")
publishLatency := metricsMW.PublishLatency("events.user.created")
errorCount := metricsMW.ErrorCount("events.user.created")
```

### Logging Middleware

```go
// Create logging middleware (uses OTEL logger)
loggingMW := middleware.NewLoggingMiddleware(middleware.LoggingConfig{
    Level: "debug",
    
    // Log fields
    IncludeHeaders: true,
    IncludeData:    false, // Don't log message data
    MaxDataSize:    1024,  // Truncate if logging data
    
    // Custom fields
    FieldsFunc: func(msg *core.Message) map[string]any {
        return map[string]any{
            "tenant_id": msg.Headers.Get("X-Tenant-ID"),
        }
    },
})

publisher.UseMiddleware(loggingMW.Publish())
subscriber.UseMiddleware(loggingMW.Subscribe())
```

### Custom Middleware

```go
// Publish middleware
func myPublishMiddleware(next middleware.PublishHandler) middleware.PublishHandler {
    return func(ctx context.Context, subject string, msg *core.Message) error {
        // Before publishing
        start := time.Now()
        log.Printf("Publishing to %s", subject)
        
        // Add custom header
        msg.Headers.Set("X-Published-At", time.Now().Format(time.RFC3339))
        
        // Call next middleware/handler
        err := next(ctx, subject, msg)
        
        // After publishing
        duration := time.Since(start)
        if err != nil {
            log.Printf("Publish to %s failed after %v: %v", subject, duration, err)
        } else {
            log.Printf("Published to %s in %v", subject, duration)
        }
        
        return err
    }
}

// Subscribe middleware
func mySubscribeMiddleware(next middleware.SubscribeHandler) middleware.SubscribeHandler {
    return func(ctx context.Context, msg *core.Message) error {
        // Before processing
        log.Printf("Received from %s", msg.Subject)
        
        // Add to context
        ctx = context.WithValue(ctx, "received_at", time.Now())
        
        // Call next middleware/handler
        err := next(ctx, msg)
        
        // After processing
        if err != nil {
            log.Printf("Processing %s failed: %v", msg.Subject, err)
        }
        
        return err
    }
}

// Apply middleware
publisher.UseMiddleware(myPublishMiddleware)
subscriber.UseMiddleware(mySubscribeMiddleware)

// Chain multiple middleware
publisher.UseMiddleware(
    middleware.Chain(
        loggingMW.Publish(),
        tracingMW.Publish(),
        retryMW.Publish(),
        cbMW.Publish(),
    ),
)
```

---

## Authentication

### User/Password

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

creds := auth.UserPass("myuser", "mypassword")

// Or use events package helper
creds = events.UserPassCredentials("myuser", "mypassword")
```

### Token Authentication

```go
creds := auth.Token("my-secret-token")

// Or use events package helper
creds = events.TokenCredentials("my-secret-token")
```

### JWT Authentication

```go
// JWT with NKey seed
creds := auth.JWT(
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "SUAG4PHKQ...",
)

// Validate credentials
if creds.IsEmpty() {
    log.Fatal("Invalid credentials")
}
```

### NKey Authentication

```go
// NKey from seed file
creds := auth.NKey("/path/to/nkey.seed")

// Or from seed string
creds = &auth.Credentials{
    Type: auth.TypeNKey,
    Seed: "SUAG...",
}
```

### Credentials File

```go
// Use NATS credentials file
creds := auth.CredentialsFile("/path/to/user.creds")
```

---

## Encoding

### JSON Encoding

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"

// Create JSON encoder
encoder := encoding.NewJSONEncoder()

// Encode
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
user := &User{ID: "123", Name: "John"}
data, err := encoder.Encode(user)

// Decode
var decoded User
err = encoder.Decode(data, &decoded)

// Create message with encoded data
msg := core.NewMessage(data)
msg = msg.WithHeader("Content-Type", "application/json")
```

### Protobuf Encoding

```go
// Create Protobuf encoder
encoder := encoding.NewProtobufEncoder()

// Encode (requires proto.Message)
protoUser := &pb.User{Id: "123", Name: "John"}
data, err := encoder.Encode(protoUser)

// Decode
var decoded pb.User
err = encoder.Decode(data, &decoded)

// Create message
msg := core.NewMessage(data)
msg = msg.WithHeader("Content-Type", "application/protobuf")
```

### Auto-Detection

```go
// Get encoder by content type
encoder := encoding.GetEncoder("application/json")
encoder = encoding.GetEncoder("application/protobuf")
encoder = encoding.GetEncoder("application/msgpack")

// Register custom encoder
encoding.RegisterEncoder("application/custom", &MyCustomEncoder{})
```

---

## Utils Package

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Subject utilities
valid := utils.ValidateSubject("events.user.created") // true
sanitized := utils.SanitizeSubject("events user@created") // "eventsusercreated"
subject := utils.BuildSubject("events", "user", "created") // "events.user.created"
tokens := utils.ParseSubject("events.user.created") // ["events", "user", "created"]

// ID generation
msgID := utils.GenerateMessageID()     // 32-char hex
corrID := utils.GenerateCorrelationID() // 24-char hex
traceID := utils.GenerateTraceID()     // 32-char hex (W3C compatible)
spanID := utils.GenerateSpanID()       // 16-char hex (W3C compatible)

// String utilities
truncated := utils.TruncateString("hello world", 8) // "hello..."
value := utils.CoalesceString("", "", "default") // "default"
value = utils.StringOrDefault("", "fallback") // "fallback"

// Type checking
if utils.IsZeroValue(myValue) {
    // value is zero
}
if utils.IsNil(myInterface) {
    // interface or pointer is nil
}

// Time utilities
millis := utils.UnixMillis()
ts := utils.TimestampString() // RFC3339 format

// Collection utilities
if utils.Contains(slice, element) {
    // element in slice
}
unique := utils.Unique([]int{1, 2, 2, 3}) // [1, 2, 3]
filtered := utils.Filter(slice, func(x int) bool { return x > 0 })
mapped := utils.Map(slice, func(x int) string { return strconv.Itoa(x) })

// Map utilities
merged := utils.MergeMaps(map1, map2)
copied := utils.CopyMap(original)
keys := utils.Keys(myMap)
values := utils.Values(myMap)

// Error utilities
wrapped := utils.WrapError(err, "context")
wrapped = utils.WrapErrorf(err, "failed to process %s", id)

collector := utils.NewErrorCollector()
collector.Add(err1)
collector.Add(err2)
if collector.HasErrors() {
    return collector.Error()
}

root := utils.RootCause(wrappedError)

// Thread-safe collections
safeMap := utils.NewSafeMap[string, int]()
safeMap.Set("key", 42)
value, ok := safeMap.Get("key")
safeMap.Range(func(k string, v int) bool {
    // iterate
    return true
})

safeSlice := utils.NewSafeSlice[string]()
safeSlice.Append("a", "b", "c")
val, ok := safeSlice.Get(0)

// Optional type
opt := utils.Some(42)
if opt.IsPresent() {
    val := opt.Value()
}
val := opt.ValueOr(0)

// Result type
result := utils.Ok(data)
result = utils.Err[Data](err)
if result.IsOk() {
    data := result.Value()
}
```

---

## Error Handling

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Check specific errors
if errors.Is(err, core.ErrNotConnected) {
    // Handle not connected
}
if errors.Is(err, core.ErrPublishTimeout) {
    // Handle timeout
}
if errors.Is(err, core.ErrNoReply) {
    // No responders available
}

// Check error categories
if core.IsRetryable(err) {
    // Safe to retry
}
if core.IsTemporary(err) {
    // Temporary error, will likely resolve
}

// Error types
var tenantErr core.TenantError
if errors.As(err, &tenantErr) {
    log.Printf("Tenant %s: %v", tenantErr.TenantID, tenantErr.Err)
}

var serErr core.SerializationError
if errors.As(err, &serErr) {
    log.Printf("Serialization failed for %s: %v", serErr.Type, serErr.Err)
}

// Create structured errors
err := core.NewError("publish", "send_message", originalErr)
err = err.WithDetails(map[string]any{
    "subject": "events.user.created",
    "size":    len(data),
})

// Multi-error handling
var multiErr *core.MultiError
if errors.As(err, &multiErr) {
    for _, e := range multiErr.All() {
        log.Printf("Error: %v", e)
    }
}

// Common error variables
core.ErrNotConnected       // Not connected to NATS
core.ErrAlreadyConnected   // Already connected
core.ErrConnectionClosed   // Connection was closed
core.ErrConnectionTimeout  // Connection timed out
core.ErrPublishFailed      // Publish operation failed
core.ErrPublishTimeout     // Publish timed out
core.ErrNoAck              // No acknowledgment received
core.ErrInvalidSubject     // Invalid subject name
core.ErrInvalidMessage     // Invalid message
core.ErrMessageTooLarge    // Message exceeds size limit
core.ErrRequestTimeout     // Request-reply timed out
core.ErrNoReply            // No responders available
core.ErrSubscriptionClosed // Subscription was closed
core.ErrAuthFailed         // Authentication failed
core.ErrPermissionDenied   // Permission denied
core.ErrTenantNotFound     // Tenant not found
core.ErrShutdownInProgress // Shutdown in progress
```

---

## Best Practices

### 1. Connection Management

```go
// DO: Use connection pooling for high-throughput scenarios
pool, _ := core.NewPool(core.PoolConfig{
    MinSize: 2,
    MaxSize: 10,
    Factory: factory,
})

// DO: Properly handle shutdown
defer func() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    client.Shutdown(ctx)
}()

// DO: Monitor connection health
health := conn.Health()
if !health.Healthy {
    log.Printf("Unhealthy: %s", health.Message)
}
```

### 2. Message Handling

```go
// DO: Use message pooling for high throughput
msg := core.AcquireMessage()
defer core.ReleaseMessage(msg)

// DO: Set appropriate message IDs for deduplication
msg = msg.WithID(uuid.New().String())

// DO: Include correlation IDs for tracing
msg = msg.WithHeader(core.HeaderCorrelationID, correlationID)

// DON'T: Store large data in headers
// DO: Keep headers small, put data in message body
```

### 3. Error Handling

```go
// DO: Check if errors are retryable
if !core.IsRetryable(err) {
    return err // Don't retry
}

// DO: Use circuit breakers for external calls
cb := middleware.NewCircuitBreaker(cfg)
publisher.UseMiddleware(middleware.NewCircuitBreakerMiddleware(cb).Publish())

// DO: Log errors with context
logger.Error(ctx, "publish failed",
    logger.String("subject", subject),
    logger.Err(err),
)
```

### 4. Subscription Best Practices

```go
// DO: Use queue groups for load balancing
sub, _ := subscriber.QueueSubscribe(ctx, "events.>", "processor-group", handler)

// DO: Handle errors in message handlers
handler := func(ctx context.Context, msg *core.Message) error {
    if err := process(msg); err != nil {
        logger.Error(ctx, "processing failed", logger.Err(err))
        return err // Will trigger retry if configured
    }
    return nil
}

// DO: Use drain for graceful shutdown
if err := sub.Drain(); err != nil {
    log.Printf("Drain error: %v", err)
}
```

### 5. JetStream Best Practices

```go
// DO: Use pull consumers (recommended by NATS)
batch, _ := consumer.Fetch(100)
for msg := range batch.Messages() {
    process(msg)
    msg.Ack()
}

// DO: Use backoff for redeliveries
consumerCfg.BackOff = []time.Duration{
    1 * time.Second,
    5 * time.Second,
    30 * time.Second,
}

// DO: Set appropriate ack wait times
consumerCfg.AckWait = 30 * time.Second

// DO: Use InProgress for long-running processing
go func() {
    for !done {
        msg.InProgress()
        time.Sleep(10 * time.Second)
    }
}()
```

### 6. Multi-Tenant Best Practices

```go
// DO: Isolate tenant connections
conn, _ := manager.Connect(ctx, tenantID, creds)

// DO: Handle tenant errors
manager.OnError(func(err core.TenantError) {
    metrics.IncrCounter("tenant_errors", err.TenantID)
    log.Printf("Tenant %s error: %v", err.TenantID, err.Err)
})

// DO: Monitor tenant health
for id, status := range manager.Health() {
    if !status.Healthy {
        alertOnCall(id, status.Message)
    }
}
```

### 7. Observability

```go
// DO: Use all observability middleware
publisher.UseMiddleware(middleware.Chain(
    loggingMW.Publish(),
    tracingMW.Publish(),
    metricsMW.Publish(),
))

// DO: Include trace context in messages
ctx = core.WithTraceContext(ctx, traceCtx)
headers := core.ExtractHeaders(ctx, nil)
msg = msg.WithHeaders(headers)

// DO: Extract trace context from received messages
ctx = core.InjectContext(ctx, msg.Headers)
```

---

## Complete Example: Event-Driven Microservice

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
)

type UserCreatedEvent struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Configuration
    cfg := config.DefaultConfig()
    cfg.Servers = []string{os.Getenv("NATS_URL")}
    cfg.Name = "user-service"

    // Create client
    client, err := events.NewClient(cfg)
    if err != nil {
        log.Fatal(err)
    }

    // Connect
    creds := events.UserPassCredentials(
        os.Getenv("NATS_USER"),
        os.Getenv("NATS_PASS"),
    )
    conn, err := client.Connect(ctx, "main", creds)
    if err != nil {
        log.Fatal(err)
    }

    // Set up middleware
    retryMW := middleware.NewRetryMiddleware(middleware.RetryConfig{
        MaxAttempts: 3,
    })
    tracingMW := middleware.NewTracingMiddleware(middleware.TracingConfig{
        ServiceName: "user-service",
    })

    publisher := conn.Publisher()
    publisher.UseMiddleware(middleware.Chain(
        tracingMW.Publish(),
        retryMW.Publish(),
    ))

    subscriber := conn.Subscriber()
    subscriber.UseMiddleware(tracingMW.Subscribe())

    // Register service
    registry := microservice.NewRegistry()
    mgr := microservice.NewManager(microservice.ManagerConfig{
        Registry: registry,
    })

    service := &microservice.Info{
        ID:      "user-service-1",
        Name:    "user-service",
        Version: "1.0.0",
        Type:    microservice.ServiceTypeAPI,
        Status:  microservice.StatusRunning,
    }
    mgr.RegisterService(ctx, &microservice.RegistrationRequest{Service: service})

    // Subscribe to events
    sub, err := subscriber.QueueSubscribe(ctx, "events.user.>", "user-service",
        func(ctx context.Context, msg *core.Message) error {
            var event UserCreatedEvent
            if err := json.Unmarshal(msg.Data, &event); err != nil {
                log.Printf("Invalid event: %v", err)
                return nil // Don't retry invalid messages
            }

            log.Printf("User created: %s (%s)", event.UserID, event.Email)

            // Process event...
            
            return nil
        })
    if err != nil {
        log.Fatal(err)
    }

    // Publish sample event
    event := UserCreatedEvent{
        UserID:    "user-123",
        Email:     "user@example.com",
        CreatedAt: time.Now(),
    }
    data, _ := json.Marshal(event)
    msg := core.NewMessage(data)
    msg = msg.WithID("event-" + event.UserID)
    
    if err := publisher.Publish(ctx, "events.user.created", msg); err != nil {
        log.Printf("Publish error: %v", err)
    }

    // Wait for shutdown signal
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh

    // Graceful shutdown
    log.Println("Shutting down...")
    
    sub.Drain()
    mgr.Shutdown(ctx)
    
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    
    if err := client.Shutdown(shutdownCtx); err != nil {
        log.Printf("Shutdown error: %v", err)
    }
    
    log.Println("Shutdown complete")
}
```

---

For more details, see the [Technical Architecture Documentation](README.md) and [API Reference](https://pkg.go.dev).
