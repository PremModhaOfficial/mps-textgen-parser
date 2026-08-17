package events

import (
	"context"
	"sync"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"

	"github.com/nats-io/nats.go"
)

const attrMessagingDestination = "messaging.destination"

// MessageHandler processes received NATS messages. Returning an error marks
// the consumer span as failed but does not nack or retry the message — that
// responsibility belongs to the middleware layer (e.g., retry middleware).
type MessageHandler func(ctx context.Context, msg *nats.Msg) error

// SubscribeMiddleware intercepts message handling using the decorator pattern,
// analogous to the publisher's Middleware type. Common uses include logging,
// error recovery, and message filtering.
type SubscribeMiddleware func(MessageHandler) MessageHandler

// Subscription represents an active NATS subscription. It abstracts over
// both regular and queue subscriptions, providing lifecycle control via
// Unsubscribe and Drain.
type Subscription interface {
	// Subject returns the NATS subject pattern this subscription listens on.
	Subject() string

	// Unsubscribe immediately removes the subscription. Pending messages
	// may be dropped; prefer Drain for graceful shutdown.
	Unsubscribe() error

	// Drain unsubscribes and waits for all in-flight messages to be processed
	// before the subscription is fully closed. This is the preferred method
	// for graceful shutdown.
	Drain() error

	// IsValid returns true if the underlying NATS subscription is still active
	// and has not been unsubscribed or drained.
	IsValid() bool
}

// Subscriber subscribes to NATS subjects and receives messages with optional
// middleware support. It tracks all active subscriptions for bulk cleanup via
// Close and supports both regular and queue group subscriptions. All methods
// are safe for concurrent use.
//
// Each received message is wrapped in an OpenTelemetry consumer span that
// records the subject and (if applicable) queue group, enabling distributed
// tracing across publishers and subscribers.
type Subscriber struct {
	conn *Connection

	mu            sync.RWMutex    // mu protects subscriptions and closed from concurrent access.
	subscriptions []*subscription // subscriptions is the list of active subscriptions for cleanup.
	closed        bool            // closed prevents new subscriptions after Close is called.

	// middleware is the composed middleware chain applied to all message handlers.
	middleware SubscribeMiddleware
}

// NewSubscriber creates a new Subscriber bound to the given Connection. The
// subscriber starts with no middleware; use UseMiddleware to add logging,
// error recovery, or message filtering.
func NewSubscriber(conn *Connection) *Subscriber {
	return &Subscriber{
		conn:          conn,
		subscriptions: make([]*subscription, 0),
	}
}

// UseMiddleware adds middleware to the subscriber's chain. Multiple calls
// compose middleware in order: the first middleware added wraps the outermost
// layer. Middleware is applied when a subscription is created, so existing
// subscriptions are not affected by subsequent UseMiddleware calls.
func (s *Subscriber) UseMiddleware(mw SubscribeMiddleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.middleware == nil {
		s.middleware = mw
	} else {
		existing := s.middleware
		s.middleware = chainSubscribeMiddleware(existing, mw)
	}
}

// chainSubscribeMiddleware chains two subscribe middleware functions into one.
// The resulting middleware calls first(second(next)), so first is outermost.
func chainSubscribeMiddleware(first, second SubscribeMiddleware) SubscribeMiddleware {
	return func(next MessageHandler) MessageHandler {
		return first(second(next))
	}
}

// Subscribe creates a subscription to the specified subject. The subject may
// contain NATS wildcards (* for single token, > for multi-token). Each received
// message is processed by the handler, wrapped in any configured middleware and
// an OpenTelemetry consumer span.
func (s *Subscriber) Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error) {
	return s.subscribe(ctx, subject, "", handler)
}

// QueueSubscribe creates a queue subscription that load-balances messages across
// all subscribers in the same queue group. Only one subscriber in the group
// receives each message, enabling horizontal scaling of message consumers.
func (s *Subscriber) QueueSubscribe(ctx context.Context, subject string, queue string, handler MessageHandler) (Subscription, error) {
	return s.subscribe(ctx, subject, queue, handler)
}

// subscribe is the shared implementation for Subscribe and QueueSubscribe. It
// creates an OpenTelemetry span for the subscribe operation, verifies the
// connection, wraps the handler with middleware, builds the NATS callback,
// and registers the resulting subscription for lifecycle tracking.
func (s *Subscriber) subscribe(ctx context.Context, subject, queue string, handler MessageHandler) (Subscription, error) {
	spanName := "nats.subscribe"
	spanAttrs := []tracer.StartSpanOption{
		tracer.WithAttributes(
			tracer.StringAttr(attrMessagingDestination, subject),
		),
	}
	if queue != "" {
		spanName = "nats.queue_subscribe"
		spanAttrs = []tracer.StartSpanOption{
			tracer.WithAttributes(
				tracer.StringAttr(attrMessagingDestination, subject),
				tracer.StringAttr("messaging.queue", queue),
			),
		}
	}

	ctx, span := tracer.Start(ctx, spanName, spanAttrs...)
	defer span.End()

	if err := s.checkReady(span); err != nil {
		return nil, err
	}

	nc := s.conn.Conn()
	if nc == nil {
		span.SetError(utils.ErrNotConnected)
		return nil, utils.ErrNotConnected
	}

	subCtx, cancel := context.WithCancel(ctx)
	wrappedHandler := s.wrapHandler(handler)

	msgCallback := s.buildMsgCallback(subCtx, queue, wrappedHandler)

	var natsSub *nats.Subscription
	var err error
	if queue != "" {
		natsSub, err = nc.QueueSubscribe(subject, queue, msgCallback)
	} else {
		natsSub, err = nc.Subscribe(subject, msgCallback)
	}

	if err != nil {
		cancel()
		span.SetError(err)
		logger.Error(ctx, "failed to create NATS subscription",
			logger.String("subject", subject),
			logger.String("queue", queue),
			logger.Err(err),
		)
		return nil, err
	}

	sub := &subscription{
		natsSub: natsSub,
		subject: subject,
		queue:   queue,
		cancel: cancel,
	}

	s.mu.Lock()
	s.subscriptions = append(s.subscriptions, sub)
	s.mu.Unlock()

	s.conn.Touch()

	span.SetOK()
	logger.Info(ctx, "NATS subscription created",
		logger.String("subject", subject),
		logger.String("queue", queue),
	)

	return sub, nil
}

// checkReady verifies the subscriber is open and the connection is active.
// Returns ErrConnectionClosed if the subscriber has been closed, or
// ErrNotConnected if the underlying NATS connection is down.
func (s *Subscriber) checkReady(span tracer.Span) error {
	s.mu.RLock()
	closed := s.closed
	s.mu.RUnlock()

	if closed {
		span.SetError(utils.ErrConnectionClosed)
		return utils.ErrConnectionClosed
	}

	if !s.conn.IsConnected() {
		span.SetError(utils.ErrNotConnected)
		return utils.ErrNotConnected
	}

	return nil
}

// buildMsgCallback creates the NATS message callback that bridges from NATS's
// push-based delivery to the MessageHandler interface. Each invocation checks
// the subscription context, creates a consumer span with subject/queue attributes,
// and delegates to the wrapped handler. The span is ended after the handler returns.
func (s *Subscriber) buildMsgCallback(subCtx context.Context, queue string, handler MessageHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		select {
		case <-subCtx.Done():
			return
		default:
		}

		attrs := []tracer.StartSpanOption{
			tracer.WithAttributes(
				tracer.StringAttr(attrMessagingDestination, msg.Subject),
			),
		}
		if queue != "" {
			attrs = []tracer.StartSpanOption{
				tracer.WithAttributes(
					tracer.StringAttr(attrMessagingDestination, msg.Subject),
					tracer.StringAttr("messaging.queue", queue),
				),
			}
		}

		msgCtx, msgSpan := tracer.StartConsumer(subCtx, "nats.receive", attrs...)

		if err := handler(msgCtx, msg); err != nil {
			msgSpan.SetError(err)
		} else {
			msgSpan.SetOK()
		}
		msgSpan.End()
	}
}

// Unsubscribe removes a specific subscription from the tracked list and drains
// it to allow in-flight messages to complete. Uses swap-and-truncate to remove
// the subscription from the slice in O(1) time.
func (s *Subscriber) Unsubscribe(sub Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.subscriptions {
		if existing.subject == sub.Subject() && existing.natsSub == sub.(*subscription).natsSub {
			// Remove from slice (swap with last, then truncate)
			s.subscriptions[i] = s.subscriptions[len(s.subscriptions)-1]
			s.subscriptions = s.subscriptions[:len(s.subscriptions)-1]
			return sub.Drain()
		}
	}
	return nil
}

// Close marks the subscriber as closed, cancels all subscription contexts, and
// drains every active subscription to allow in-flight messages to complete.
// After Close returns, all subsequent Subscribe calls will return ErrConnectionClosed.
func (s *Subscriber) Close(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.subscriber.close")
	defer span.End()

	s.mu.Lock()
	defer s.mu.Unlock()

	count := len(s.subscriptions)
	s.closed = true

	var drainErrs []error
	for _, sub := range s.subscriptions {
		if sub.cancel != nil {
			sub.cancel()
		}
		if err := sub.Drain(); err != nil {
			drainErrs = append(drainErrs, err)
		}
	}

	if len(drainErrs) > 0 {
		logger.Warn(ctx, "errors draining subscriptions during close",
			logger.Int("error_count", len(drainErrs)),
		)
	}

	s.subscriptions = nil

	span.SetOK()
	logger.Info(ctx, "NATS subscriber closed",
		logger.Int("subscriptions_closed", count),
	)

	return nil
}

// wrapHandler wraps a handler with the current middleware chain. The middleware
// snapshot is taken at subscription time, so subsequent UseMiddleware calls do
// not affect existing subscriptions.
func (s *Subscriber) wrapHandler(handler MessageHandler) MessageHandler {
	s.mu.RLock()
	mw := s.middleware
	s.mu.RUnlock()

	if mw != nil {
		return mw(handler)
	}
	return handler
}

// subscription implements the Subscription interface by wrapping a nats.Subscription
// with a cancellable context. When Unsubscribe or Drain is called, the context is
// cancelled first to stop the message callback from processing new messages, then
// the underlying NATS subscription is cleaned up.
type subscription struct {
	natsSub *nats.Subscription // natsSub is the underlying NATS subscription.
	subject string             // subject is the NATS subject pattern.
	queue   string             // queue is the queue group name (empty for non-queue subscriptions).
	cancel  context.CancelFunc // cancel stops the subscription's message processing.
}

// Subject returns the NATS subject pattern this subscription listens on.
func (s *subscription) Subject() string {
	return s.subject
}

// Unsubscribe cancels the subscription context and immediately removes the
// subscription from the NATS server. Pending messages may be dropped.
func (s *subscription) Unsubscribe() error {
	// Cancel the context to stop processing new messages
	if s.cancel != nil {
		s.cancel()
	}

	if s.natsSub == nil {
		return nil
	}

	return s.natsSub.Unsubscribe()
}

// Drain cancels the subscription context and gracefully drains the underlying
// NATS subscription, allowing all in-flight messages to be processed before
// the subscription is fully removed.
func (s *subscription) Drain() error {
	// Cancel the context to stop processing new messages
	if s.cancel != nil {
		s.cancel()
	}

	if s.natsSub == nil {
		return nil
	}

	return s.natsSub.Drain()
}

// IsValid returns true if the underlying NATS subscription is still active.
func (s *subscription) IsValid() bool {
	if s.natsSub == nil {
		return false
	}
	return s.natsSub.IsValid()
}
