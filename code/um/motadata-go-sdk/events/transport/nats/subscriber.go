package nats

import (
	"context"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
	"github.com/nats-io/nats.go"
)

// Subscriber implements core.Subscriber for NATS
type Subscriber struct {
	conn *Connection

	mu            sync.RWMutex
	subscriptions []*subscription
	closed        bool

	// Middleware
	middleware middleware.SubscribeMiddleware
}

// NewSubscriber creates a new NATS subscriber
func NewSubscriber(conn *Connection) *Subscriber {
	return &Subscriber{
		conn:          conn,
		subscriptions: make([]*subscription, 0),
	}
}

// UseMiddleware adds middleware to the subscriber
func (s *Subscriber) UseMiddleware(mw middleware.SubscribeMiddleware) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.middleware == nil {
		s.middleware = mw
	} else {
		existing := s.middleware
		s.middleware = middleware.ChainSubscribe(existing, mw)
	}
}

// Subscribe creates a subscription to the specified subject
func (subscriber *Subscriber) Subscribe(ctx context.Context, subject string, handler core.MessageHandler) (core.Subscription, error) {
	ctx, span := tracer.Start(ctx, "nats.subscribe",
		tracer.WithAttributes(
			tracer.StringAttr("messaging.destination", subject),
		),
	)
	defer span.End()

	subscriber.mu.Lock()
	if subscriber.closed {
		subscriber.mu.Unlock()
		span.SetError(core.ErrConnectionClosed)
		return nil, core.ErrConnectionClosed
	}
	subscriber.mu.Unlock()

	if !subscriber.conn.IsConnected() {
		span.SetError(core.ErrNotConnected)
		return nil, core.ErrNotConnected
	}

	nc := subscriber.conn.Conn()
	if nc == nil {
		span.SetError(core.ErrNotConnected)
		return nil, core.ErrNotConnected
	}

	// Create a subscription context that can be cancelled
	subCtx, cancel := context.WithCancel(ctx)

	// Wrap handler with middleware
	wrappedHandler := subscriber.wrapHandler(handler)

	// Create NATS subscription
	natsSub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		// Check if subscription context is still valid
		select {
		case <-subCtx.Done():
			return
		default:
		}

		coreMsg := subscriber.convertMessage(msg)

		// Start a consumer span for this message
		msgCtx, msgSpan := tracer.StartConsumer(subCtx, "nats.receive",
			tracer.WithAttributes(
				tracer.StringAttr("messaging.destination", msg.Subject),
				tracer.StringAttr("message.id", coreMsg.ID),
			),
		)

		// Inject headers into context for tracing/tenant propagation
		msgCtx = core.InjectContext(msgCtx, coreMsg.Headers)

		if err := wrappedHandler(msgCtx, coreMsg); err != nil {
			msgSpan.SetError(err)
		} else {
			msgSpan.SetOK()
		}
		msgSpan.End()
	})

	if err != nil {
		cancel()
		span.SetError(err)
		logger.Error(ctx, "failed to create NATS subscription",
			logger.String("subject", subject),
			logger.Err(err),
		)
		return nil, err
	}

	sub := &subscription{
		natsSub: natsSub,
		subject: subject,
		ctx:     subCtx,
		cancel:  cancel,
	}

	subscriber.mu.Lock()
	subscriber.subscriptions = append(subscriber.subscriptions, sub)
	subscriber.mu.Unlock()

	subscriber.conn.Touch()

	span.SetOK()
	logger.Info(ctx, "NATS subscription created",
		logger.String("subject", subject),
	)

	return sub, nil
}

// QueueSubscribe creates a queue subscription
func (subscriber *Subscriber) QueueSubscribe(ctx context.Context, subject string, queue string, handler core.MessageHandler) (core.Subscription, error) {
	ctx, span := tracer.Start(ctx, "nats.queue_subscribe",
		tracer.WithAttributes(
			tracer.StringAttr("messaging.destination", subject),
			tracer.StringAttr("messaging.queue", queue),
		),
	)
	defer span.End()

	subscriber.mu.Lock()
	if subscriber.closed {
		subscriber.mu.Unlock()
		span.SetError(core.ErrConnectionClosed)
		return nil, core.ErrConnectionClosed
	}
	subscriber.mu.Unlock()

	if !subscriber.conn.IsConnected() {
		span.SetError(core.ErrNotConnected)
		return nil, core.ErrNotConnected
	}

	nc := subscriber.conn.Conn()
	if nc == nil {
		span.SetError(core.ErrNotConnected)
		return nil, core.ErrNotConnected
	}

	// Create a subscription context that can be cancelled
	subCtx, cancel := context.WithCancel(ctx)

	// Wrap handler with middleware
	wrappedHandler := subscriber.wrapHandler(handler)

	// Create NATS queue subscription
	natsSub, err := nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		// Check if subscription context is still valid
		select {
		case <-subCtx.Done():
			return
		default:
		}

		coreMsg := subscriber.convertMessage(msg)

		// Start a consumer span for this message
		msgCtx, msgSpan := tracer.StartConsumer(subCtx, "nats.receive",
			tracer.WithAttributes(
				tracer.StringAttr("messaging.destination", msg.Subject),
				tracer.StringAttr("messaging.queue", queue),
				tracer.StringAttr("message.id", coreMsg.ID),
			),
		)

		// Inject headers into context for tracing/tenant propagation
		msgCtx = core.InjectContext(msgCtx, coreMsg.Headers)

		if err := wrappedHandler(msgCtx, coreMsg); err != nil {
			msgSpan.SetError(err)
		} else {
			msgSpan.SetOK()
		}
		msgSpan.End()
	})

	if err != nil {
		cancel()
		span.SetError(err)
		logger.Error(ctx, "failed to create NATS queue subscription",
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
		ctx:     subCtx,
		cancel:  cancel,
	}

	subscriber.mu.Lock()
	subscriber.subscriptions = append(subscriber.subscriptions, sub)
	subscriber.mu.Unlock()

	subscriber.conn.Touch()

	span.SetOK()
	logger.Info(ctx, "NATS queue subscription created",
		logger.String("subject", subject),
		logger.String("queue", queue),
	)

	return sub, nil
}

// Close closes all subscriptions
func (subscriber *Subscriber) Close(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.subscriber.close")
	defer span.End()

	subscriber.mu.Lock()
	defer subscriber.mu.Unlock()

	count := len(subscriber.subscriptions)
	subscriber.closed = true

	for _, sub := range subscriber.subscriptions {
		// Cancel the subscription context first
		if sub.cancel != nil {
			sub.cancel()
		}
		_ = sub.Drain()
	}

	subscriber.subscriptions = nil

	span.SetOK()
	logger.Info(ctx, "NATS subscriber closed",
		logger.Int("subscriptions_closed", count),
	)

	return nil
}

func (s *Subscriber) wrapHandler(handler core.MessageHandler) middleware.SubscribeHandler {
	h := func(ctx context.Context, msg *core.Message) error {
		return handler(ctx, msg)
	}

	s.mu.RLock()
	mw := s.middleware
	s.mu.RUnlock()

	if mw != nil {
		return mw(h)
	}
	return h
}

func (subscriber *Subscriber) convertMessage(msg *nats.Msg) *core.Message {

	coreMsg := &core.Message{
		Subject:   msg.Subject,
		Data:      msg.Data,
		Reply:     msg.Reply,
		Timestamp: time.Now(), // Set receive timestamp
	}

	// Convert headers
	if len(msg.Header) > 0 {
		coreMsg.Headers = make(core.Headers)
		for k, v := range msg.Header {
			coreMsg.Headers[k] = v
		}

		// Extract message ID
		if msgID := msg.Header.Get(core.HeaderMessageID); msgID != "" {
			coreMsg.ID = msgID
		}
	}

	return coreMsg
}

// subscription implements core.Subscription
type subscription struct {
	natsSub *nats.Subscription
	subject string
	queue   string
	ctx     context.Context
	cancel  context.CancelFunc
}

func (s *subscription) Subject() string {
	return s.subject
}

func (sub *subscription) Unsubscribe() error {

	// Cancel the context to stop processing new messages
	if sub.cancel != nil {

		sub.cancel()
	}

	if sub.natsSub == nil {

		return nil
	}

	return sub.natsSub.Unsubscribe()
}

func (sub *subscription) Drain() error {

	// Cancel the context to stop processing new messages
	if sub.cancel != nil {

		sub.cancel()
	}

	if sub.natsSub == nil {

		return nil
	}

	return sub.natsSub.Drain()
}

func (s *subscription) IsValid() bool {
	if s.natsSub == nil {
		return false
	}
	return s.natsSub.IsValid()
}
