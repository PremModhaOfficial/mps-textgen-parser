package events

import (
	"context"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Default publisher constants control timeouts for publish and request operations.
const (
	// defaultRequestTimeout is the maximum time to wait for a reply in request-reply
	// operations before returning ErrRequestTimeout. Can be overridden per-call
	// via context deadline or RequestWithTimeout.
	defaultRequestTimeout = 30 * time.Second

	// defaultFlushTimeout is the maximum time to wait for pending messages to be
	// flushed to the NATS server during publisher Close. Respects context deadline
	// if shorter than this default.
	defaultFlushTimeout = 5 * time.Second
)

// Handler is the function signature for publish operations. Every publish
// call ultimately resolves to a Handler invocation. Middleware wraps this
// type to intercept or augment publish behaviour (e.g., tracing, retry).
type Handler func(ctx context.Context, msg *nats.Msg) error

// Middleware intercepts and wraps publish handlers using the decorator pattern.
// Middleware is composed via chainMiddleware, forming an onion-style call chain
// where the outermost middleware executes first and the innermost publishes the
// message to NATS.
type Middleware func(Handler) Handler

// PubAckFuture represents a pending publish acknowledgment for JetStream
// async publishes. Callers can either read from the Ok/Err channels or
// use Wait for a blocking call that respects context cancellation.
type PubAckFuture interface {
	// Ok returns a receive-only channel that delivers the PubAck on success.
	Ok() <-chan *PubAck

	// Err returns a receive-only channel that delivers an error on failure.
	Err() <-chan error

	// Wait blocks until the publish completes and returns the ack, or until
	// the context is cancelled. Exactly one of PubAck or error will be non-nil.
	Wait(ctx context.Context) (*PubAck, error)
}

// PubAck represents a publish acknowledgment from JetStream, confirming
// that a message has been persisted to the specified stream at the given
// sequence number.
type PubAck struct {
	Stream   string // Stream is the name of the JetStream stream that accepted the message.
	Sequence uint64 // Sequence is the stream-level sequence number assigned to the message.
	Domain   string // Domain is the JetStream domain (empty for the default domain).
}

// Publisher sends messages to NATS subjects with optional middleware support.
// It wraps a Connection and provides Core NATS Publish, JetStream PublishAsync,
// and request-reply operations. All methods are safe for concurrent use.
//
// The middleware chain is applied to every publish and request call, enabling
// cross-cutting concerns like tracing, retry, metrics, and circuit breaking
// to be added without modifying application code.
type Publisher struct {
	conn *Connection

	mu     sync.RWMutex // mu protects closed and middleware from concurrent access.
	closed bool         // closed prevents new publishes after Close is called.

	// middleware is the composed middleware chain applied to all publish operations.
	middleware Middleware
}

// NewPublisher creates a new Publisher bound to the given Connection.
// The publisher starts with no middleware; use UseMiddleware to add tracing,
// retry, metrics, or other cross-cutting concerns.
func NewPublisher(conn *Connection) *Publisher {
	return &Publisher{
		conn: conn,
	}
}

// UseMiddleware adds middleware to the publisher's chain. Multiple calls
// compose middleware in order: the first middleware added wraps the outermost
// layer. Middleware is applied to Publish, PublishAsync, and Request calls.
// Must not be called concurrently with publish operations after initial setup.
func (p *Publisher) UseMiddleware(mw Middleware) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.middleware == nil {
		p.middleware = mw
	} else {
		existing := p.middleware
		p.middleware = chainMiddleware(existing, mw)
	}
}

// chainMiddleware chains two middleware functions into one. The resulting
// middleware calls first(second(next)), so first is the outermost wrapper.
func chainMiddleware(first, second Middleware) Middleware {
	return func(next Handler) Handler {
		return first(second(next))
	}
}

// Publish sends a message to the specified subject via Core NATS. The message
// subject field is overwritten with the provided subject. Middleware (tracing,
// retry, etc.) is applied before the actual NATS PublishMsg call. Returns
// ErrConnectionClosed if the publisher has been closed.
func (p *Publisher) Publish(ctx context.Context, subject string, msg *nats.Msg) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return utils.ErrConnectionClosed
	}
	mw := p.middleware
	p.mu.RUnlock()

	handler := p.doPublish(subject)
	if mw != nil {
		handler = mw(handler)
	}

	return handler(ctx, msg)
}

// doPublish returns the terminal publish handler for the given subject.
// This handler creates a producer span, verifies the connection, sets the
// message subject, publishes via NATS, and records the outcome on the span.
func (p *Publisher) doPublish(subject string) Handler {
	return func(ctx context.Context, msg *nats.Msg) error {
		ctx, span := tracer.StartProducer(ctx, "nats.publish",
			tracer.WithAttributes(
				tracer.StringAttr("messaging.destination", subject),
				tracer.StringAttr("message.id", msg.Header.Get(core.HeaderMessageID)),
			),
		)
		defer span.End()

		if !p.conn.IsConnected() {
			span.SetError(utils.ErrNotConnected)
			return utils.ErrNotConnected
		}

		nc := p.conn.Conn()
		if nc == nil {
			span.SetError(utils.ErrNotConnected)
			return utils.ErrNotConnected
		}

		// Ensure the subject is set on the message
		msg.Subject = subject

		p.conn.Touch()

		if err := nc.PublishMsg(msg); err != nil {
			span.SetError(err)
			logger.Error(ctx, "failed to publish NATS message",
				logger.String("subject", subject),
				logger.Err(err),
			)
			return err
		}

		span.SetOK()
		logger.Debug(ctx, "NATS message published",
			logger.String("subject", subject),
			logger.Int("data_size", len(msg.Data)),
		)
		return nil
	}
}

// PublishAsync sends a message asynchronously via JetStream and returns a PubAckFuture.
// The publish executes in a background goroutine; callers should use the future's
// Ok/Err channels or Wait method to observe the result. Returns an error future
// (not nil) if JetStream is unavailable or the publisher is closed, so callers
// can always safely read from the returned future.
func (p *Publisher) PublishAsync(ctx context.Context, subject string, msg *nats.Msg) PubAckFuture {
	future := &pubAckFuture{
		ok:  make(chan *PubAck, 1),
		err: make(chan error, 1),
	}

	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		future.err <- utils.ErrConnectionClosed
		return future
	}
	mw := p.middleware
	p.mu.RUnlock()

	go func() {
		js := p.conn.JetStream()
		if js == nil {
			future.err <- utils.ErrJetStreamNotEnabled
			return
		}

		// Apply middleware for JetStream publish
		handler := p.doJetStreamPublish(js, subject, future)
		if mw != nil {
			handler = mw(handler)
		}

		if err := handler(ctx, msg); err != nil {
			logger.Error(ctx, "async JetStream publish handler failed",
				logger.String("subject", subject),
				logger.Err(err),
			)
		}
	}()

	return future
}

// doJetStreamPublish returns the terminal handler for JetStream publishing.
// It uses jetstream.PublishMsg (not Publish) to preserve NATS headers such as
// trace context, tenant ID, and message ID. If the message has a HeaderMessageID,
// it is forwarded as a JetStream MsgID for server-side deduplication.
func (p *Publisher) doJetStreamPublish(js jetstream.JetStream, subject string, future *pubAckFuture) Handler {
	return func(ctx context.Context, msg *nats.Msg) error {
		if !p.conn.IsConnected() {
			future.err <- utils.ErrNotConnected
			return utils.ErrNotConnected
		}

		p.conn.Touch()

		// Ensure subject is set on the message
		msg.Subject = subject

		// Use jetstream PublishMsg with options to preserve headers
		var opts []jetstream.PublishOpt
		if msg.Header != nil {
			if msgID := msg.Header.Get(core.HeaderMessageID); msgID != "" {
				opts = append(opts, jetstream.WithMsgID(msgID))
			}
		}

		ack, err := js.PublishMsg(ctx, msg, opts...)
		if err != nil {
			future.err <- err
			return err
		}

		future.ok <- &PubAck{
			Stream:   ack.Stream,
			Sequence: ack.Sequence,
			Domain:   ack.Domain,
		}
		return nil
	}
}

// Request sends a message and waits for a reply using the NATS request-reply pattern.
// Middleware is applied to the outgoing message as a pre-send hook (the middleware's
// terminal handler is a no-op; the actual send uses nc.RequestMsg). The timeout is
// derived from the context deadline or falls back to defaultRequestTimeout.
// Returns ErrRequestTimeout on timeout, ErrNoReply when no responders are available.
func (p *Publisher) Request(ctx context.Context, subject string, msg *nats.Msg) (*nats.Msg, error) {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return nil, utils.ErrConnectionClosed
	}
	mw := p.middleware
	p.mu.RUnlock()

	if !p.conn.IsConnected() {
		return nil, utils.ErrNotConnected
	}

	nc := p.conn.Conn()
	if nc == nil {
		return nil, utils.ErrNotConnected
	}

	// Ensure the subject is set on the message
	msg.Subject = subject

	// Apply middleware to the outgoing message (for tracing, metrics, etc.)
	if mw != nil {
		// Run middleware as a pre-send hook — the actual send is request-reply
		mwHandler := mw(func(_ context.Context, _ *nats.Msg) error {
			return nil // no-op terminal handler; actual send happens below
		})
		if err := mwHandler(ctx, msg); err != nil {
			return nil, err
		}
	}

	p.conn.Touch()

	// Calculate timeout from context or use default
	timeout := defaultRequestTimeout
	if deadline, ok := ctx.Deadline(); ok {
		if remaining := time.Until(deadline); remaining > 0 {
			timeout = remaining
		}
	}

	// Send request and wait for reply
	replyMsg, err := nc.RequestMsg(msg, timeout)
	if err != nil {
		if err == nats.ErrTimeout {
			return nil, utils.ErrRequestTimeout
		}
		if err == nats.ErrNoResponders {
			return nil, utils.ErrNoReply
		}
		return nil, err
	}

	return replyMsg, nil
}

// RequestWithTimeout sends a request with an explicit timeout, wrapping the
// provided context with context.WithTimeout. This is a convenience wrapper
// around Request for callers who prefer specifying the timeout directly
// rather than via context.
func (p *Publisher) RequestWithTimeout(ctx context.Context, subject string, msg *nats.Msg, timeout time.Duration) (*nats.Msg, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return p.Request(ctx, subject, msg)
}

// Close marks the publisher as closed and flushes any pending messages to the
// NATS server. After Close returns, all subsequent Publish and PublishAsync calls
// will return ErrConnectionClosed. The flush timeout respects the context deadline
// or falls back to defaultFlushTimeout.
func (p *Publisher) Close(ctx context.Context) error {
	p.mu.Lock()
	p.closed = true
	p.mu.Unlock()

	// Flush any pending messages to ensure they are sent
	if p.conn == nil || !p.conn.IsConnected() {
		return nil
	}

	nc := p.conn.Conn()
	if nc == nil {
		return nil
	}

	return nc.FlushTimeout(timeoutFromContext(ctx, defaultFlushTimeout))
}

// pubAckFuture implements PubAckFuture using buffered channels. Both ok and
// err channels have a buffer of 1, ensuring the producer goroutine never blocks.
// Exactly one of the two channels receives a value for each future.
type pubAckFuture struct {
	ok  chan *PubAck // ok receives the ack when the JetStream publish succeeds.
	err chan error   // err receives the error when the JetStream publish fails.
}

// Ok returns a receive-only channel that delivers the PubAck on success.
func (f *pubAckFuture) Ok() <-chan *PubAck {
	return f.ok
}

// Err returns a receive-only channel that delivers an error on failure.
func (f *pubAckFuture) Err() <-chan error {
	return f.err
}

// Wait blocks until the future resolves or the context is cancelled. It
// selects across the ok channel, err channel, and context.Done, returning
// whichever fires first.
func (f *pubAckFuture) Wait(ctx context.Context) (*PubAck, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ack := <-f.ok:
		return ack, nil
	case err := <-f.err:
		return nil, err
	}
}

// defaultMaxConcurrentFlush limits the number of goroutines spawned during
// concurrent batch flush to prevent unbounded goroutine creation.
const defaultMaxConcurrentFlush = 64

// BatchPublisher accumulates messages and publishes them in batches, either
// on-demand via Flush or automatically when maxBatchSize is reached or
// flushInterval elapses. This amortises per-message overhead and is useful
// for high-throughput producers.
//
// Concurrent flush uses a channel-based semaphore (maxFlushWorkers) to bound
// the number of in-flight goroutines. All methods are safe for concurrent use.
type BatchPublisher struct {
	publisher *Publisher      // publisher is the underlying Publisher used for each message.
	messages  []*batchMessage // messages is the pending batch buffer.
	mu        sync.Mutex      // mu protects messages from concurrent Add/Flush access.

	// Configuration
	maxBatchSize    int           // maxBatchSize triggers auto-flush when reached (0 = no limit).
	flushInterval   time.Duration // flushInterval triggers periodic auto-flush (0 = disabled).
	concurrent      bool          // concurrent enables parallel message publishing during flush.
	maxFlushWorkers int           // maxFlushWorkers caps goroutines in concurrent flush mode.

	// Auto-flush lifecycle
	done   <-chan struct{}     // done is closed when cancel is called, stopping the auto-flush goroutine.
	cancel context.CancelFunc // cancel stops the auto-flush background loop.
	wg     sync.WaitGroup     // wg tracks the auto-flush goroutine for graceful shutdown.

	// onFlushError is an optional callback invoked when auto-flush encounters an error.
	// If nil, errors are logged via the otel logger.
	onFlushError func(error)
}

// batchMessage pairs a subject with the NATS message to be published.
type batchMessage struct {
	subject string    // subject is the NATS subject to publish to.
	msg     *nats.Msg // msg is the NATS message payload including headers.
}

// BatchPublisherOption configures a BatchPublisher using the functional options pattern.
type BatchPublisherOption func(*BatchPublisher)

// WithMaxBatchSize sets the maximum number of messages buffered before an
// automatic flush is triggered. A value of 0 (default) disables size-based
// auto-flush; callers must explicitly call Flush.
func WithMaxBatchSize(size int) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.maxBatchSize = size
	}
}

// WithFlushInterval sets the time-based auto-flush interval. A background
// goroutine will flush any pending messages at this interval. A value of 0
// (default) disables periodic auto-flush. The goroutine is stopped on Close.
func WithFlushInterval(d time.Duration) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.flushInterval = d
	}
}

// WithConcurrentFlush enables or disables concurrent message publishing during
// flush. When enabled, messages are published in parallel goroutines bounded
// by maxFlushWorkers (default 64). When disabled, messages are published
// sequentially. Concurrent mode is faster but uses more resources.
func WithConcurrentFlush(concurrent bool) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.concurrent = concurrent
	}
}

// WithMaxFlushWorkers sets the maximum number of concurrent goroutines used
// during a concurrent flush. The semaphore pattern ensures no more than n
// goroutines are in-flight at any time. Values <= 0 are ignored, preserving
// the default of 64.
func WithMaxFlushWorkers(n int) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		if n > 0 {
			bp.maxFlushWorkers = n
		}
	}
}

// WithFlushErrorCallback sets a callback invoked when an auto-flush (timer or
// size-triggered) encounters an error. Without this callback, errors are logged
// via the otel logger. The callback runs on the auto-flush goroutine, so it
// should not block.
func WithFlushErrorCallback(cb func(error)) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.onFlushError = cb
	}
}

// NewBatchPublisher creates a new BatchPublisher wrapping the given Publisher.
// Options configure batch size, flush interval, concurrency, and error handling.
// If a flush interval is set, a background auto-flush goroutine is started
// immediately and must be stopped by calling Close.
func NewBatchPublisher(publisher *Publisher, opts ...BatchPublisherOption) *BatchPublisher {
	ctx, cancel := context.WithCancel(context.Background())
	bp := &BatchPublisher{
		publisher:       publisher,
		messages:        make([]*batchMessage, 0, 100),
		maxBatchSize:    0, // No limit by default
		flushInterval:   0, // No auto-flush by default
		maxFlushWorkers: defaultMaxConcurrentFlush,
		done:            ctx.Done(),
		cancel:          cancel,
	}

	for _, opt := range opts {
		opt(bp)
	}

	// Start auto-flush timer if interval is set
	if bp.flushInterval > 0 {
		bp.startAutoFlush()
	}

	return bp
}

// startAutoFlush starts the background goroutine that periodically flushes
// pending messages at the configured flushInterval. The goroutine runs until
// the BatchPublisher's context is cancelled via Close.
func (bp *BatchPublisher) startAutoFlush() {
	bp.wg.Add(1)
	go func() {
		defer bp.wg.Done()
		ticker := time.NewTicker(bp.flushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-bp.done:
				return
			case <-ticker.C:
				bp.tryAutoFlush()
			}
		}
	}()
}

// tryAutoFlush flushes pending messages and handles any errors.
func (bp *BatchPublisher) tryAutoFlush() {
	if bp.Count() == 0 {
		return
	}

	err := bp.Flush(context.Background())
	if err == nil {
		return
	}

	if bp.onFlushError != nil {
		bp.onFlushError(err)
	} else {
		logger.Error(context.Background(), "batch auto-flush failed",
			logger.Err(err),
		)
	}
}

// Add appends a message to the pending batch. If maxBatchSize is configured and
// the batch reaches that size, Flush is called automatically. Returns the error
// from the auto-flush, if any; the message is always added to the batch regardless.
func (bp *BatchPublisher) Add(subject string, msg *nats.Msg) error {
	bp.mu.Lock()
	bp.messages = append(bp.messages, &batchMessage{subject: subject, msg: msg})
	shouldFlush := bp.maxBatchSize > 0 && len(bp.messages) >= bp.maxBatchSize
	bp.mu.Unlock()

	if shouldFlush {
		return bp.Flush(context.Background())
	}
	return nil
}

// AddMultiple adds multiple messages to the batch in a single lock acquisition.
// The map keys are subjects and the values are the corresponding NATS messages.
// Triggers auto-flush if the batch reaches maxBatchSize after the additions.
func (bp *BatchPublisher) AddMultiple(messages map[string]*nats.Msg) error {
	bp.mu.Lock()
	for subject, msg := range messages {
		bp.messages = append(bp.messages, &batchMessage{subject: subject, msg: msg})
	}
	shouldFlush := bp.maxBatchSize > 0 && len(bp.messages) >= bp.maxBatchSize
	bp.mu.Unlock()

	if shouldFlush {
		return bp.Flush(context.Background())
	}
	return nil
}

// Flush atomically drains the pending batch and publishes all messages.
// The batch buffer is swapped under the mutex so new messages can be added
// while flush is in progress. Uses sequential or concurrent mode depending
// on the concurrent flag. Returns a MultiError if multiple publishes fail.
func (bp *BatchPublisher) Flush(ctx context.Context) error {
	bp.mu.Lock()
	if len(bp.messages) == 0 {
		bp.mu.Unlock()
		return nil
	}
	messages := bp.messages
	bp.messages = make([]*batchMessage, 0, cap(messages))
	bp.mu.Unlock()

	if bp.concurrent {
		return bp.flushConcurrent(ctx, messages)
	}
	return bp.flushSequential(ctx, messages)
}

// flushSequential publishes messages one at a time in order. Errors are
// collected and returned as a single MultiError; a failure does not abort
// the remaining messages.
func (bp *BatchPublisher) flushSequential(ctx context.Context, messages []*batchMessage) error {
	var errs []error
	for _, msg := range messages {
		if err := bp.publisher.Publish(ctx, msg.subject, msg.msg); err != nil {
			errs = append(errs, err)
		}
	}
	return collectErrors(errs)
}

// flushConcurrent publishes messages in parallel using a channel-based semaphore
// to bound the number of in-flight goroutines to maxFlushWorkers. All goroutines
// must complete before the method returns. Errors are collected via a buffered
// channel and returned as a single MultiError.
func (bp *BatchPublisher) flushConcurrent(ctx context.Context, messages []*batchMessage) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(messages))

	// Use a semaphore to bound concurrency
	sem := make(chan struct{}, bp.maxFlushWorkers)

	for _, msg := range messages {
		wg.Add(1)
		sem <- struct{}{} // Acquire semaphore slot
		go func(m *batchMessage) {
			defer wg.Done()
			defer func() { <-sem }() // Release semaphore slot
			if err := bp.publisher.Publish(ctx, m.subject, m.msg); err != nil {
				errCh <- err
			}
		}(msg)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	return collectErrors(errs)
}

// timeoutFromContext returns the remaining time from the context deadline,
// or the provided default if no deadline is set.
func timeoutFromContext(ctx context.Context, defaultTimeout time.Duration) time.Duration {
	if deadline, ok := ctx.Deadline(); ok {
		if remaining := time.Until(deadline); remaining > 0 {
			return remaining
		}
	}
	return defaultTimeout
}

// collectErrors combines multiple errors into a single error. Returns nil if
// the slice is empty, the single error if there is exactly one, or a MultiError
// wrapping all errors if there are multiple.
func collectErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return utils.NewMultiError(errs)
}

// Count returns the number of messages currently buffered and awaiting flush.
func (bp *BatchPublisher) Count() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.messages)
}

// Close stops the auto-flush goroutine (if running), waits for it to exit,
// and performs a final flush of any remaining messages. This ensures no messages
// are lost on shutdown.
func (bp *BatchPublisher) Close(ctx context.Context) error {
	bp.cancel()
	bp.wg.Wait()

	// Final flush
	return bp.Flush(ctx)
}
