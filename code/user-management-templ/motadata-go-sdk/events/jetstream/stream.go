// Package jetstream provides high-level wrappers around NATS JetStream
// primitives -- streams, consumers, key-value stores, and object stores.
// It simplifies common operations while still exposing the underlying
// JetStream types for advanced use cases.
package jetstream

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// StreamManager provides JetStream stream and consumer management.
//
// A stream is the durable, append-only message log in JetStream. Consumers
// are the way clients read messages from a stream -- each stream can have
// many independent consumers with different delivery semantics (at-most-once,
// at-least-once, exactly-once). StreamManager exposes both stream and
// consumer CRUD so that callers can manage the full lifecycle from a single
// entry point.
type StreamManager struct {
	// js is the underlying JetStream context used for all server operations.
	js jetstream.JetStream
}

// NewStreamManager creates a StreamManager from a JetStream instance.
// It returns ErrJetStreamNotEnabled when js is nil, which typically means
// the NATS connection was established without JetStream capabilities.
func NewStreamManager(js jetstream.JetStream) (*StreamManager, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}
	return &StreamManager{js: js}, nil
}

// CreateStream creates a new stream with the given configuration.
// The StreamConfig controls retention policy, storage backend, subject
// filters, replication factor, and other stream-level settings. An error
// is returned if a stream with the same name already exists.
func (sm *StreamManager) CreateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	return sm.js.CreateStream(ctx, cfg)
}

// UpdateStream updates an existing stream with the given configuration.
// Only certain fields of the StreamConfig can be changed after creation
// (e.g., max messages, max bytes, subjects). Attempting to change immutable
// fields such as the storage type will result in a server-side error.
func (sm *StreamManager) UpdateStream(ctx context.Context, cfg jetstream.StreamConfig) (jetstream.Stream, error) {
	return sm.js.UpdateStream(ctx, cfg)
}

// DeleteStream deletes the named stream and all of its consumers, messages,
// and associated state. This operation is irreversible.
func (sm *StreamManager) DeleteStream(ctx context.Context, name string) error {
	return sm.js.DeleteStream(ctx, name)
}

// Stream returns a handle to the named stream, allowing further operations
// such as message purging, info retrieval, or consumer creation directly
// on the stream object. Returns an error if the stream does not exist.
func (sm *StreamManager) Stream(ctx context.Context, name string) (jetstream.Stream, error) {
	return sm.js.Stream(ctx, name)
}

// StreamNames returns a lister that lazily iterates over all stream names
// on the server. Use the lister's Name() channel and Error() method to
// consume results.
func (sm *StreamManager) StreamNames(ctx context.Context) jetstream.StreamNameLister {
	return sm.js.StreamNames(ctx)
}

// CreateOrUpdateConsumer creates or updates a durable consumer on the given
// stream. This is a convenience method on StreamManager that delegates to the
// underlying JetStream API. For richer consumer wrappers with automatic ack/nak
// handling, see the Consumer type in this package.
func (sm *StreamManager) CreateOrUpdateConsumer(ctx context.Context, stream string, cfg jetstream.ConsumerConfig) (jetstream.Consumer, error) {
	return sm.js.CreateOrUpdateConsumer(ctx, stream, cfg)
}

// JetStream returns the underlying JetStream instance for advanced
// operations not exposed by StreamManager (e.g., direct message access,
// account info queries).
func (sm *StreamManager) JetStream() jetstream.JetStream {
	return sm.js
}
