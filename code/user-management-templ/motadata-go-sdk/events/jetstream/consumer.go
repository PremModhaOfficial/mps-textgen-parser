package jetstream

import (
	"context"
	"errors"

	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// MessageHandler processes a JetStream message with context.
// Returning a non-nil error signals that the message should be negatively
// acknowledged (Nak) so the server can redeliver it; returning nil causes
// the message to be acknowledged (Ack) and removed from the pending set.
type MessageHandler func(ctx context.Context, msg jetstream.Msg) error

// Consumer wraps a JetStream pull-based consumer with automatic ack/nak
// handling and convenience methods for different consumption patterns.
//
// JetStream supports two consumer flavours:
//
//   - Consumer (this type): A durable or ephemeral consumer with explicit
//     acknowledgement. It persists its delivery position on the server so
//     that it can resume after restarts, and it supports load-balanced
//     delivery across multiple instances (queue-group semantics). Callers
//     must acknowledge or negatively-acknowledge every message; the Consume
//     and ConsumeWithContext helpers automate this based on the handler's
//     return value.
//
//   - OrderedConsumer: An ephemeral, single-instance consumer that
//     guarantees strict ordering and is automatically recreated on
//     failures. See the OrderedConsumer type below.
//
// Pull vs Push: All consumers created by this package use the pull model.
// In pull mode the client explicitly requests messages from the server,
// which provides natural back-pressure and works well with horizontal
// scaling. The legacy push model (where the server pushes to a delivery
// subject) is not exposed here; the newer pull-based API is recommended
// for all new applications.
type Consumer struct {
	// consumer is the underlying JetStream consumer handle.
	consumer jetstream.Consumer
	// stream is the name of the stream this consumer reads from.
	stream string
	// name is the durable name of the consumer (empty for ephemeral consumers).
	name string
}

// CreateOrUpdateConsumer creates a new consumer or updates the configuration
// of an existing one on the specified stream. The ConsumerConfig controls
// delivery policy (all, new, by start sequence/time), ack policy, max
// redelivery attempts, filter subjects, and other consumer-level settings.
// Returns ErrJetStreamNotEnabled when js is nil.
func CreateOrUpdateConsumer(ctx context.Context, js jetstream.JetStream, stream string, cfg jetstream.ConsumerConfig) (*Consumer, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	c, err := js.CreateOrUpdateConsumer(ctx, stream, cfg)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		stream:   stream,
		name:     cfg.Name,
	}, nil
}

// GetConsumer retrieves an existing consumer by its durable name on the
// given stream. This is useful when a consumer was already created (e.g.,
// by infrastructure tooling) and the application only needs to attach to
// it for message consumption.
func GetConsumer(ctx context.Context, js jetstream.JetStream, stream, name string) (*Consumer, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	c, err := js.Consumer(ctx, stream, name)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		stream:   stream,
		name:     name,
	}, nil
}

// DeleteConsumer deletes a consumer from the stream, removing all of its
// server-side state including pending acks and delivery tracking. Any
// active Consume or Messages sessions for this consumer will be terminated.
func DeleteConsumer(ctx context.Context, js jetstream.JetStream, stream, name string) error {
	if js == nil {
		return utils.ErrJetStreamNotEnabled
	}

	return js.DeleteConsumer(ctx, stream, name)
}

// Consume starts continuous pull-based message consumption with a callback.
// Each message is dispatched to handler in its own goroutine (managed by the
// NATS client). If the handler returns nil the message is Ack'd; if it
// returns an error the message is Nak'd so the server can redeliver it to
// this or another consumer instance. A background context is injected into
// each handler invocation; use ConsumeWithContext when cancellation
// propagation is needed.
//
// The returned ConsumeContext can be stopped by calling Stop(), which
// gracefully drains in-flight messages.
func (c *Consumer) Consume(handler MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	return c.consumer.Consume(func(msg jetstream.Msg) {
		// createBackgroundCtx - use a fresh background context per message.
		ctx := context.Background()

		if err := handler(ctx, msg); err != nil {
			// Nak tells the server to redeliver after the configured backoff.
			_ = msg.Nak()
			return
		}

		// Ack removes the message from the consumer's pending set.
		_ = msg.Ack()
	}, opts...)
}

// ConsumeWithContext starts continuous pull-based message consumption with
// an explicit parent context injected into each handler invocation. Before
// invoking the handler, it checks whether ctx has been cancelled; if so,
// the message is Nak'd so another consumer instance can process it. This
// is useful for graceful shutdown scenarios where the application needs to
// stop processing before the ConsumeContext is fully drained.
func (c *Consumer) ConsumeWithContext(ctx context.Context, handler MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	return c.consumer.Consume(func(msg jetstream.Msg) {
		// earlyExitOnCancel - if the parent context is already done, Nak the
		// message immediately so it can be picked up by another instance.
		select {
		case <-ctx.Done():
			_ = msg.Nak()
			return
		default:
		}

		if err := handler(ctx, msg); err != nil {
			_ = msg.Nak()
			return
		}

		_ = msg.Ack()
	}, opts...)
}

// Messages returns an iterator for pulling messages one at a time via
// the Next() method on the returned MessagesContext. This is the
// recommended approach when the caller wants fine-grained control over
// message processing timing (e.g., processing in a select loop alongside
// other channels). Call Stop() on the returned MessagesContext when done
// to release server-side resources.
func (c *Consumer) Messages(opts ...jetstream.PullMessagesOpt) (jetstream.MessagesContext, error) {
	return c.consumer.Messages(opts...)
}

// Fetch retrieves up to batch messages synchronously, blocking until the
// batch is filled or the fetch timeout expires. The returned MessageBatch
// provides a Messages() channel for iteration and an Error() method to
// check for fetch-level errors.
func (c *Consumer) Fetch(batch int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	return c.consumer.Fetch(batch, opts...)
}

// FetchNoWait retrieves only the messages that are immediately available
// on the server, up to batch size, without blocking. This is useful for
// polling patterns or batch-processing jobs that should not wait for new
// messages to arrive.
func (c *Consumer) FetchNoWait(batch int) (jetstream.MessageBatch, error) {
	return c.consumer.FetchNoWait(batch)
}

// Next retrieves the next single message, blocking until one is available
// or the fetch timeout expires. Internally it issues a Fetch(1) call and
// returns the first (and only) message from the batch. If no message is
// available within the timeout, it returns an error. This is a convenience
// for request-reply or sequential-processing patterns where only one
// message at a time is needed.
func (c *Consumer) Next(opts ...jetstream.FetchOpt) (jetstream.Msg, error) {
	// fetchSingle - fetch exactly one message from the server.
	batch, err := c.consumer.Fetch(1, opts...)
	if err != nil {
		return nil, err
	}

	// Return the first message from the batch channel.
	for msg := range batch.Messages() {
		return msg, nil
	}

	if batch.Error() != nil {
		return nil, batch.Error()
	}

	return nil, errors.New("no messages available")
}

// Info returns up-to-date consumer information and status by issuing a
// round-trip to the NATS server. The returned ConsumerInfo includes
// delivery counts, ack floor, number of pending messages, and the
// consumer's full configuration.
func (c *Consumer) Info(ctx context.Context) (*jetstream.ConsumerInfo, error) {
	return c.consumer.Info(ctx)
}

// CachedInfo returns the last known consumer info without a server
// round-trip. The data may be stale; use Info for authoritative state.
func (c *Consumer) CachedInfo() *jetstream.ConsumerInfo {
	return c.consumer.CachedInfo()
}

// Stream returns the name of the stream this consumer reads from.
// A consumer is always bound to exactly one stream.
func (c *Consumer) Stream() string {
	return c.stream
}

// Name returns the durable consumer name. For ephemeral consumers created
// without an explicit name, this may be empty.
func (c *Consumer) Name() string {
	return c.name
}

// Underlying returns the raw jetstream.Consumer for advanced use cases
// that are not covered by this wrapper (e.g., custom ack policies or
// direct access to the consumer's ordered delivery internals).
func (c *Consumer) Underlying() jetstream.Consumer {
	return c.consumer
}

// OrderedConsumer wraps a JetStream ordered consumer for replay and
// stream-replay scenarios.
//
// Unlike a regular Consumer, an OrderedConsumer:
//   - Is always ephemeral -- it has no durable name and cannot be shared.
//   - Guarantees strict, gap-free message ordering.
//   - Is automatically recreated by the NATS client if the server reports a
//     sequence mismatch or the consumer is otherwise disrupted.
//   - Does not require explicit Ack/Nak because the server tracks delivery
//     by heartbeats and flow control rather than per-message acknowledgement.
//
// Use an OrderedConsumer when you need to replay all messages in a stream
// in order (e.g., materializing a read model) and do not need load-balanced
// delivery across multiple instances.
type OrderedConsumer struct {
	// consumer is the underlying JetStream ordered consumer handle.
	consumer jetstream.Consumer
	// stream is the name of the stream this ordered consumer reads from.
	stream string
}

// CreateOrderedConsumer creates an ordered consumer for the given stream.
// Ordered consumers are ephemeral, automatically recreated on failures,
// and guarantee strict ordered delivery. The OrderedConsumerConfig can
// specify filter subjects and an optional starting sequence or time to
// begin replay from a particular point in the stream.
// Returns ErrJetStreamNotEnabled when js is nil.
func CreateOrderedConsumer(ctx context.Context, js jetstream.JetStream, stream string, cfg jetstream.OrderedConsumerConfig) (*OrderedConsumer, error) {
	if js == nil {
		return nil, utils.ErrJetStreamNotEnabled
	}

	c, err := js.OrderedConsumer(ctx, stream, cfg)
	if err != nil {
		return nil, err
	}

	return &OrderedConsumer{
		consumer: c,
		stream:   stream,
	}, nil
}

// Consume starts continuous ordered message consumption. Because ordered
// consumers do not use explicit acknowledgements, the handler's return
// value is ignored -- messages are delivered exactly once in sequence
// regardless of handler errors. The NATS client handles flow control
// and heartbeat-based tracking internally.
func (oc *OrderedConsumer) Consume(handler MessageHandler, opts ...jetstream.PullConsumeOpt) (jetstream.ConsumeContext, error) {
	return oc.consumer.Consume(func(msg jetstream.Msg) {
		ctx := context.Background()
		// No ack/nak needed -- ordered consumers track progress automatically.
		_ = handler(ctx, msg)
	}, opts...)
}

// Messages returns an iterator for pulling ordered messages one at a time.
// The iterator preserves strict ordering guarantees. Call Stop() on the
// returned MessagesContext when done.
func (oc *OrderedConsumer) Messages(opts ...jetstream.PullMessagesOpt) (jetstream.MessagesContext, error) {
	return oc.consumer.Messages(opts...)
}

// Fetch retrieves up to batch messages synchronously from the ordered
// consumer. Messages are guaranteed to be in stream-sequence order.
func (oc *OrderedConsumer) Fetch(batch int, opts ...jetstream.FetchOpt) (jetstream.MessageBatch, error) {
	return oc.consumer.Fetch(batch, opts...)
}

// Stream returns the name of the stream this ordered consumer reads from.
func (oc *OrderedConsumer) Stream() string {
	return oc.stream
}
