# JetStream Package

The `jetstream` package provides convenience wrappers for NATS JetStream operations including stream management, consumer management, key-value stores, and object stores.

## Overview

This package contains:

- **StreamManager**: Stream and consumer lifecycle management
- **Consumer**: Pull-based message consumption with context support
- **OrderedConsumer**: Strict ordering guarantees for replay scenarios
- **KVStore**: Key-value store for configuration and state
- **ObjectStore**: Large object storage and retrieval

## Package Structure

```
jetstream/
├── stream.go      # StreamManager for stream and consumer operations
├── consumer.go    # Consumer and OrderedConsumer wrappers
├── kv.go          # Key-Value store utilities
├── objectstore.go # Object store utilities
└── README.md      # This file
```

## Stream Management

### Creating a StreamManager

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/jetstream"
    natsjs "github.com/nats-io/nats.go/jetstream"
)

// Get JetStream context from connection
js := conn.JetStream()

// Create stream manager
sm := jetstream.NewStreamManager(js)
```

### Stream Operations

```go
ctx := context.Background()

// Create a stream
stream, err := sm.CreateStream(ctx, natsjs.StreamConfig{
    Name:      "ORDERS",
    Subjects:  []string{"orders.>"},
    Storage:   natsjs.FileStorage,
    Retention: natsjs.LimitsPolicy,
    MaxAge:    24 * time.Hour,
    Replicas:  3,
})

// Update stream configuration
stream, err = sm.UpdateStream(ctx, natsjs.StreamConfig{
    Name:      "ORDERS",
    Subjects:  []string{"orders.>", "refunds.>"},
    Storage:   natsjs.FileStorage,
    Retention: natsjs.LimitsPolicy,
    MaxAge:    48 * time.Hour,
    Replicas:  3,
})

// Get existing stream
stream, err = sm.Stream(ctx, "ORDERS")

// List all stream names
names, err := sm.StreamNames(ctx)

// Delete a stream
err = sm.DeleteStream(ctx, "ORDERS")
```

### Consumer Operations via StreamManager

```go
// Create or update a consumer on a stream
consumer, err := sm.CreateOrUpdateConsumer(ctx, "ORDERS", natsjs.ConsumerConfig{
    Name:          "order-processor",
    Durable:       "order-processor",
    FilterSubject: "orders.created",
    AckPolicy:     natsjs.AckExplicitPolicy,
    MaxDeliver:    5,
    AckWait:       30 * time.Second,
})

// Access underlying JetStream instance
js := sm.JetStream()
```

## Consumer

### Creating Consumers

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/jetstream"

// Create or update a durable consumer
consumer, err := jetstream.CreateOrUpdateConsumer(ctx, js, "ORDERS", natsjs.ConsumerConfig{
    Name:          "order-processor",
    Durable:       "order-processor",
    FilterSubject: "orders.created",
    AckPolicy:     natsjs.AckExplicitPolicy,
    MaxDeliver:    5,
    AckWait:       30 * time.Second,
})

// Get existing consumer
consumer, err = jetstream.GetConsumer(ctx, js, "ORDERS", "order-processor")

// Delete a consumer
err = jetstream.DeleteConsumer(ctx, js, "ORDERS", "order-processor")
```

### Consuming Messages

#### Continuous Consumption (Push-based)

```go
// Start consuming with callback
consumeCtx, err := consumer.Consume(func(msg natsjs.Msg) {
    // Process message
    fmt.Printf("Received: %s\n", string(msg.Data()))

    // Acknowledge
    msg.Ack()
})
defer consumeCtx.Stop()

// Or with explicit parent context for cancellation
consumeCtx, err = consumer.ConsumeWithContext(ctx, func(msg natsjs.Msg) {
    // Process message
    msg.Ack()
})
```

#### Pull-based Consumption

```go
// Get message iterator
iter, err := consumer.Messages()
defer iter.Stop()

for {
    msg, err := iter.Next()
    if err != nil {
        break
    }
    // Process message
    msg.Ack()
}
```

#### Batch Fetch

```go
// Fetch a batch of messages (blocks until batch is full or timeout)
msgs, err := consumer.Fetch(10)
for msg := range msgs.Messages() {
    // Process message
    msg.Ack()
}

// Fetch immediately available messages (non-blocking)
msgs, err = consumer.FetchNoWait(10)

// Fetch a single message
msg, err := consumer.Next()
msg.Ack()
```

#### Consumer Information

```go
// Get current consumer info (makes server request)
info, err := consumer.Info(ctx)
fmt.Printf("Pending: %d, Delivered: %d\n",
    info.NumPending, info.Delivered.Consumer)

// Get cached info (no server round-trip)
cachedInfo := consumer.CachedInfo()

// Get stream and consumer name
stream := consumer.Stream()
name := consumer.Name()

// Access underlying jetstream.Consumer
raw := consumer.Underlying()
```

### Ordered Consumer

For scenarios requiring strict message ordering with automatic recovery:

```go
// Create ordered consumer
ordered, err := jetstream.CreateOrderedConsumer(ctx, js, "ORDERS", natsjs.OrderedConsumerConfig{
    FilterSubjects: []string{"orders.>"},
})

// Consume in order
consumeCtx, err := ordered.Consume(func(msg natsjs.Msg) {
    // Messages arrive in strict sequence order
    msg.Ack()
})

// Or use iterator
iter, err := ordered.Messages()

// Or batch fetch
msgs, err := ordered.Fetch(10)

// Get stream name
stream := ordered.Stream()
```

## Key-Value Store

### Creating and Managing KV Stores

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/jetstream"

ctx := context.Background()

// Create a new KV store
kv, err := jetstream.CreateKVStore(ctx, js, natsjs.KeyValueConfig{
    Bucket:      "config",
    Description: "Application configuration",
    History:     5,
    TTL:         24 * time.Hour,
    Storage:     natsjs.FileStorage,
    Replicas:    3,
})

// Get existing KV store
kv, err = jetstream.GetKVStore(ctx, js, "config")

// List all KV store bucket names
names, err := jetstream.KVStoreNames(ctx, js)

// Delete a KV store
err = jetstream.DeleteKVStore(ctx, js, "config")
```

### KV Operations

```go
// Put a value
revision, err := kv.Put(ctx, "app.setting", []byte("value"))

// Create (fails if key exists)
revision, err = kv.Create(ctx, "app.new-key", []byte("initial"))

// Update (optimistic concurrency with revision)
revision, err = kv.Update(ctx, "app.setting", []byte("new-value"), revision)

// Get a value
entry, err := kv.Get(ctx, "app.setting")
fmt.Printf("Key: %s, Value: %s, Revision: %d\n",
    entry.Key(), string(entry.Value()), entry.Revision())

// Delete a key
err = kv.Delete(ctx, "app.setting")
```

### Key Management

```go
// List all keys
keys, err := kv.Keys(ctx)
for _, key := range keys {
    fmt.Println(key)
}

// Get history of a key
entries, err := kv.History(ctx, "app.setting")
for _, entry := range entries {
    fmt.Printf("Rev %d: %s\n", entry.Revision(), string(entry.Value()))
}

// Purge a key (remove all revisions)
err = kv.Purge(ctx, "app.setting")

// Purge all delete markers
err = kv.PurgeDeletes(ctx)
```

### Watching for Changes

```go
// Watch a specific key
watcher, err := kv.Watch(ctx, "app.setting")
defer watcher.Stop()

for entry := range watcher.Updates() {
    if entry == nil {
        continue // Initial values done
    }
    fmt.Printf("Updated: %s = %s\n", entry.Key(), string(entry.Value()))
}

// Watch all keys
watcher, err = kv.WatchAll(ctx)
```

### KV Store Status

```go
status, err := kv.Status(ctx)
fmt.Printf("Bucket: %s, Keys: %d, Bytes: %d\n",
    status.Bucket(), status.Values(), status.Bytes())

bucketName := kv.BucketName()
```

## Object Store

### Creating and Managing Object Stores

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/jetstream"

ctx := context.Background()

// Create object store
os, err := jetstream.CreateObjectStore(ctx, js, natsjs.ObjectStoreConfig{
    Bucket:      "assets",
    Description: "Application assets",
    MaxBytes:    1 << 30, // 1 GB
    Storage:     natsjs.FileStorage,
    Replicas:    3,
})

// Get existing object store
os, err = jetstream.GetObjectStore(ctx, js, "assets")

// List all object store names
names, err := jetstream.ObjectStoreNames(ctx, js)

// Delete object store
err = jetstream.DeleteObjectStore(ctx, js, "assets")
```

### Object Operations

```go
// Put an object (from io.Reader)
file, _ := os.Open("image.png")
defer file.Close()

info, err := os.Put(ctx, natsjs.ObjectMeta{
    Name:        "images/logo.png",
    Description: "Company logo",
}, file)

// Put bytes directly
info, err = os.PutBytes(ctx, "config/app.json", []byte(`{"key": "value"}`))

// Get an object (returns io.ReadCloser)
result, err := os.Get(ctx, "images/logo.png")
defer result.Close()
data, _ := io.ReadAll(result)

// Get object info without downloading
info, err = os.GetInfo(ctx, "images/logo.png")
fmt.Printf("Name: %s, Size: %d\n", info.Name, info.Size)

// Delete an object
err = os.Delete(ctx, "images/logo.png")
```

### Listing and Watching

```go
// List all objects
objects, err := os.List(ctx)
for _, info := range objects {
    fmt.Printf("%s (%d bytes)\n", info.Name, info.Size)
}

// Watch for changes
watcher, err := os.Watch(ctx)
defer watcher.Stop()

for info := range watcher.Updates() {
    if info == nil {
        continue
    }
    fmt.Printf("Object changed: %s\n", info.Name)
}
```

### Sealing and Status

```go
// Seal the object store (no more writes)
err = os.Seal(ctx)

// Get object store status
status, err := os.Status(ctx)
fmt.Printf("Bucket: %s, Objects: %d\n",
    status.Bucket(), status.Size())

bucketName := os.BucketName()
```

## Step-by-Step Usage Guide

### Step 1: Set Up JetStream Connection

```go
import (
    "context"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
)

// Create events client with JetStream enabled
client, err := events.NewClient(
    events.WithServers("nats://localhost:4222"),
    events.WithStream(config.StreamConfig{Enabled: true}),
)
defer client.Shutdown(context.Background())

// Connect tenant
conn, err := client.ConnectWithJWT(ctx, "tenant-1", jwt, seed)

// Get JetStream context
js := conn.JetStream()
```

### Step 2: Create Streams

```go
sm := events.NewStreamManager(js)

stream, err := sm.CreateStream(ctx, natsjs.StreamConfig{
    Name:     "EVENTS",
    Subjects: []string{"events.>"},
    Storage:  natsjs.FileStorage,
})
```

### Step 3: Publish to JetStream

```go
pub := events.NewPublisher(conn)

// Async publish with acknowledgment
future := pub.PublishAsync(ctx, "events.order.created", orderMsg)

ack, err := future.Wait(ctx)
fmt.Printf("Published to stream %s, seq %d\n", ack.Stream, ack.Sequence)
```

### Step 4: Create Consumers

```go
consumer, err := events.CreateOrUpdateConsumer(ctx, js, "EVENTS", natsjs.ConsumerConfig{
    Name:          "order-handler",
    Durable:       "order-handler",
    FilterSubject: "events.order.>",
    AckPolicy:     natsjs.AckExplicitPolicy,
})
```

### Step 5: Consume Messages

```go
consumeCtx, err := consumer.Consume(func(msg natsjs.Msg) {
    fmt.Printf("Processing: %s\n", string(msg.Data()))
    msg.Ack()
})
defer consumeCtx.Stop()
```

### Step 6: Use KV Store for Configuration

```go
kv, err := events.CreateKVStore(ctx, js, natsjs.KeyValueConfig{
    Bucket: "service-config",
})

kv.Put(ctx, "feature.enabled", []byte("true"))

entry, _ := kv.Get(ctx, "feature.enabled")
fmt.Println(string(entry.Value())) // "true"
```

### Step 7: Use Object Store for Large Data

```go
objStore, err := events.CreateObjectStore(ctx, js, natsjs.ObjectStoreConfig{
    Bucket: "reports",
})

objStore.PutBytes(ctx, "monthly/2024-01.pdf", reportData)
```

## Error Handling

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// JetStream errors
err := sm.CreateStream(ctx, cfg)
if err != nil {
    if errors.Is(err, utils.ErrJetStreamNotEnabled) {
        // JetStream not enabled on server
    }
    // Handle other errors
}
```

## Best Practices

### 1. Use Durable Consumers for Reliability

```go
consumer, _ := jetstream.CreateOrUpdateConsumer(ctx, js, "ORDERS", natsjs.ConsumerConfig{
    Durable:   "order-processor",  // Survives restarts
    AckPolicy: natsjs.AckExplicitPolicy,
    MaxDeliver: 5,                  // Max redeliveries
    AckWait:   30 * time.Second,   // Ack deadline
})
```

### 2. Configure Appropriate Retention

```go
sm.CreateStream(ctx, natsjs.StreamConfig{
    Name:      "EVENTS",
    Retention: natsjs.LimitsPolicy,  // Keep by limits
    MaxAge:    7 * 24 * time.Hour,   // 7 days
    MaxBytes:  1 << 30,              // 1 GB max
    MaxMsgs:   1_000_000,            // 1M messages max
    Discard:   natsjs.DiscardOld,    // Discard oldest
})
```

### 3. Use Optimistic Concurrency for KV Updates

```go
entry, _ := kv.Get(ctx, "counter")
revision := entry.Revision()

// Update with revision check (fails if modified since read)
_, err := kv.Update(ctx, "counter", newValue, revision)
if err != nil {
    // Conflict - re-read and retry
}
```

### 4. Handle Consumer Redeliveries

```go
consumer.Consume(func(msg natsjs.Msg) {
    meta, _ := msg.Metadata()

    if meta.NumDelivered > 3 {
        // Too many attempts - send to DLQ
        msg.Term()
        return
    }

    if err := process(msg); err != nil {
        msg.Nak() // Request redelivery
        return
    }

    msg.Ack()
})
```

## Thread Safety

All JetStream wrappers are thread-safe:
- `StreamManager` is safe for concurrent stream operations
- `Consumer` is safe for concurrent consumption
- `KVStore` is safe for concurrent reads and writes
- `ObjectStore` is safe for concurrent operations

## Testing

```bash
# Run tests (requires NATS server with JetStream)
go test ./events/jetstream/...

# Run with coverage
go test -cover ./events/jetstream/...
```

## Usage in Other Packages

This package is used by:
- `events` (main) - Re-exports for convenience (NewStreamManager, CreateKVStore, etc.)
- Application code - Direct JetStream operations
