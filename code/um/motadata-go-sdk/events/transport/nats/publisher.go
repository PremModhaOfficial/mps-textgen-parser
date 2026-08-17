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
	"github.com/nats-io/nats.go/jetstream"
)

// Publisher implements core.Publisher for NATS
type Publisher struct {
	conn *Connection

	mu     sync.RWMutex
	closed bool

	// Middleware
	middleware middleware.PublishMiddleware
}

// NewPublisher creates a new NATS publisher
func NewPublisher(conn *Connection) *Publisher {
	return &Publisher{
		conn: conn,
	}
}

// UseMiddleware adds middleware to the publisher
func (p *Publisher) UseMiddleware(mw middleware.PublishMiddleware) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.middleware == nil {
		p.middleware = mw
	} else {
		existing := p.middleware
		p.middleware = middleware.Chain(existing, mw)
	}
}

// Publish sends a message to the specified subject
func (p *Publisher) Publish(ctx context.Context, subject string, msg *core.Message) error {
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		return core.ErrConnectionClosed
	}
	mw := p.middleware
	// Keep lock during handler creation to ensure consistent middleware state
	handler := p.doPublish
	if mw != nil {
		handler = mw(handler)
	}
	p.mu.RUnlock()

	return handler(ctx, subject, msg)
}

func (p *Publisher) doPublish(ctx context.Context, subject string, msg *core.Message) error {
	ctx, span := tracer.StartProducer(ctx, "nats.publish",
		tracer.WithAttributes(
			tracer.StringAttr("messaging.destination", subject),
			tracer.StringAttr("message.id", msg.ID),
		),
	)
	defer span.End()

	if !p.conn.IsConnected() {
		span.SetError(core.ErrNotConnected)
		return core.ErrNotConnected
	}

	nc := p.conn.Conn()
	if nc == nil {
		span.SetError(core.ErrNotConnected)
		return core.ErrNotConnected
	}

	// Build NATS message
	natsMsg := &nats.Msg{
		Subject: subject,
		Data:    msg.Data,
	}

	// Add headers
	if len(msg.Headers) > 0 {
		natsMsg.Header = make(nats.Header)
		for k, v := range msg.Headers {
			for _, val := range v {
				natsMsg.Header.Add(k, val)
			}
		}
	}

	// Add message ID header
	if msg.ID != "" {
		if natsMsg.Header == nil {
			natsMsg.Header = make(nats.Header)
		}
		natsMsg.Header.Set(core.HeaderMessageID, msg.ID)
	}

	// Add reply subject
	if msg.Reply != "" {
		natsMsg.Reply = msg.Reply
	}

	p.conn.Touch()

	if err := nc.PublishMsg(natsMsg); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to publish NATS message",
			logger.String("subject", subject),
			logger.String("message_id", msg.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Debug(ctx, "NATS message published",
		logger.String("subject", subject),
		logger.String("message_id", msg.ID),
		logger.Int("data_size", len(msg.Data)),
	)
	return nil
}

// PublishAsync sends a message asynchronously
func (p *Publisher) PublishAsync(ctx context.Context, subject string, msg *core.Message) core.PubAckFuture {
	future := &pubAckFuture{
		ok:  make(chan *core.PubAck, 1),
		err: make(chan error, 1),
	}

	// Capture middleware under lock before starting goroutine
	p.mu.RLock()
	if p.closed {
		p.mu.RUnlock()
		future.err <- core.ErrConnectionClosed
		return future
	}
	mw := p.middleware
	p.mu.RUnlock()

	go func() {
		js := p.conn.JetStream()
		if js == nil {
			// Fall back to regular publish (which applies middleware)
			err := p.Publish(ctx, subject, msg)
			if err != nil {
				future.err <- err
			} else {
				future.ok <- &core.PubAck{}
			}
			return
		}

		// Apply middleware for JetStream publish
		handler := p.doJetStreamPublish(js, future)
		if mw != nil {
			handler = mw(handler)
		}

		if err := handler(ctx, subject, msg); err != nil {
			// Error already sent to future by doJetStreamPublish
			return
		}
	}()

	return future
}

// doJetStreamPublish returns a handler for JetStream publishing
func (p *Publisher) doJetStreamPublish(js jetstream.JetStream, future *pubAckFuture) func(ctx context.Context, subject string, msg *core.Message) error {
	return func(ctx context.Context, subject string, msg *core.Message) error {
		if !p.conn.IsConnected() {
			future.err <- core.ErrNotConnected
			return core.ErrNotConnected
		}

		// Build headers for JetStream message
		headers := make(nats.Header)
		if len(msg.Headers) > 0 {
			for k, v := range msg.Headers {
				for _, val := range v {
					headers.Add(k, val)
				}
			}
		}

		if msg.ID != "" {
			headers.Set(core.HeaderMessageID, msg.ID)
		}

		p.conn.Touch()

		// Use jetstream Publish with options
		var opts []jetstream.PublishOpt
		if len(headers) > 0 {
			// JetStream handles headers differently - set msg ID for deduplication
			if msgID := headers.Get(core.HeaderMessageID); msgID != "" {
				opts = append(opts, jetstream.WithMsgID(msgID))
			}
		}

		ack, err := js.Publish(ctx, subject, msg.Data, opts...)
		if err != nil {
			future.err <- err
			return err
		}

		future.ok <- &core.PubAck{
			Stream:   ack.Stream,
			Sequence: ack.Sequence,
			Domain:   ack.Domain,
		}
		return nil
	}
}

// Close closes the publisher, flushing any pending messages
func (publisher *Publisher) Close(ctx context.Context) error {
	publisher.mu.Lock()
	publisher.closed = true
	publisher.mu.Unlock()

	// Flush any pending messages to ensure they're sent
	if publisher.conn != nil && publisher.conn.IsConnected() {
		nc := publisher.conn.Conn()
		if nc != nil {
			// Use context timeout for flush, or default 5 seconds
			flushTimeout := 5 * time.Second
			if deadline, ok := ctx.Deadline(); ok {
				if timeout := time.Until(deadline); timeout > 0 {
					flushTimeout = timeout
				}
			}
			if err := nc.FlushTimeout(flushTimeout); err != nil {
				return err
			}
		}
	}

	return nil
}

// Request sends a message and waits for a reply (request-reply pattern)
func (publisher *Publisher) Request(ctx context.Context, subject string, msg *core.Message) (*core.Message, error) {

	publisher.mu.RLock()

	if publisher.closed {

		publisher.mu.RUnlock()

		return nil, core.ErrConnectionClosed
	}

	publisher.mu.RUnlock()

	if !publisher.conn.IsConnected() {

		return nil, core.ErrNotConnected
	}

	nc := publisher.conn.Conn()

	if nc == nil {

		return nil, core.ErrNotConnected
	}

	// Build NATS message
	natsMsg := &nats.Msg{
		Subject: subject,
		Data:    msg.Data,
	}

	// Add headers
	if len(msg.Headers) > 0 {
		natsMsg.Header = make(nats.Header)
		for k, v := range msg.Headers {
			for _, val := range v {
				natsMsg.Header.Add(k, val)
			}
		}
	}

	// Add message ID header
	if msg.ID != "" {

		if natsMsg.Header == nil {

			natsMsg.Header = make(nats.Header)
		}

		natsMsg.Header.Set(core.HeaderMessageID, msg.ID)
	}

	publisher.conn.Touch()

	// Calculate timeout from context or use default
	timeout := 30 * time.Second

	if deadline, ok := ctx.Deadline(); ok {

		timeout = time.Until(deadline)
	}

	// Send request and wait for reply
	replyMsg, err := nc.RequestMsg(natsMsg, timeout)

	if err != nil {

		if err == nats.ErrTimeout {

			return nil, core.ErrRequestTimeout
		}

		if err == nats.ErrNoResponders {

			return nil, core.ErrNoReply
		}

		return nil, err
	}

	// Convert reply to core.Message
	reply := &core.Message{
		Subject:   replyMsg.Subject,
		Data:      replyMsg.Data,
		Timestamp: time.Now(),
	}

	// Convert reply headers
	if len(replyMsg.Header) > 0 {
		reply.Headers = make(core.Headers)
		for k, v := range replyMsg.Header {
			reply.Headers[k] = v
		}

		// Extract message ID
		if msgID := replyMsg.Header.Get(core.HeaderMessageID); msgID != "" {
			reply.ID = msgID
		}
	}

	return reply, nil
}

// RequestWithTimeout sends a request with an explicit timeout
func (publisher *Publisher) RequestWithTimeout(ctx context.Context, subject string, msg *core.Message, timeout time.Duration) (*core.Message, error) {

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return publisher.Request(ctx, subject, msg)
}

// pubAckFuture implements core.PubAckFuture
type pubAckFuture struct {
	ok  chan *core.PubAck
	err chan error
}

func (f *pubAckFuture) Ok() <-chan *core.PubAck {
	return f.ok
}

func (f *pubAckFuture) Err() <-chan error {
	return f.err
}

func (f *pubAckFuture) Wait(ctx context.Context) (*core.PubAck, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ack := <-f.ok:
		return ack, nil
	case err := <-f.err:
		return nil, err
	}
}

// BatchPublisher supports batch publishing with auto-flush capabilities
type BatchPublisher struct {
	publisher *Publisher
	messages  []*batchMessage
	mu        sync.Mutex

	// Configuration
	maxBatchSize  int
	flushInterval time.Duration
	concurrent    bool

	// Auto-flush
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type batchMessage struct {
	subject string
	msg     *core.Message
}

// BatchPublisherOption configures a BatchPublisher
type BatchPublisherOption func(*BatchPublisher)

// WithMaxBatchSize sets the maximum batch size before auto-flush
func WithMaxBatchSize(size int) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.maxBatchSize = size
	}
}

// WithFlushInterval sets the auto-flush interval
func WithFlushInterval(d time.Duration) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.flushInterval = d
	}
}

// WithConcurrentFlush enables concurrent message publishing during flush
func WithConcurrentFlush(concurrent bool) BatchPublisherOption {
	return func(bp *BatchPublisher) {
		bp.concurrent = concurrent
	}
}

// NewBatchPublisher creates a new batch publisher
func NewBatchPublisher(publisher *Publisher, opts ...BatchPublisherOption) *BatchPublisher {
	ctx, cancel := context.WithCancel(context.Background())
	bp := &BatchPublisher{
		publisher:     publisher,
		messages:      make([]*batchMessage, 0, 100),
		maxBatchSize:  0, // No limit by default
		flushInterval: 0, // No auto-flush by default
		ctx:           ctx,
		cancel:        cancel,
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

// startAutoFlush starts the auto-flush timer
func (bp *BatchPublisher) startAutoFlush() {
	bp.wg.Add(1)
	go func() {
		defer bp.wg.Done()
		ticker := time.NewTicker(bp.flushInterval)
		defer ticker.Stop()

		for {
			select {
			case <-bp.ctx.Done():
				return
			case <-ticker.C:
				if bp.Count() > 0 {
					bp.Flush(bp.ctx)
				}
			}
		}
	}()
}

// Add adds a message to the batch, triggering auto-flush if max size is reached
func (bp *BatchPublisher) Add(subject string, msg *core.Message) error {
	bp.mu.Lock()
	bp.messages = append(bp.messages, &batchMessage{subject: subject, msg: msg})
	shouldFlush := bp.maxBatchSize > 0 && len(bp.messages) >= bp.maxBatchSize
	bp.mu.Unlock()

	if shouldFlush {
		return bp.Flush(bp.ctx)
	}
	return nil
}

// AddMultiple adds multiple messages to the batch
func (bp *BatchPublisher) AddMultiple(messages map[string]*core.Message) error {
	bp.mu.Lock()
	for subject, msg := range messages {
		bp.messages = append(bp.messages, &batchMessage{subject: subject, msg: msg})
	}
	shouldFlush := bp.maxBatchSize > 0 && len(bp.messages) >= bp.maxBatchSize
	bp.mu.Unlock()

	if shouldFlush {
		return bp.Flush(bp.ctx)
	}
	return nil
}

// Flush publishes all messages in the batch
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

// flushSequential publishes messages one at a time
func (bp *BatchPublisher) flushSequential(ctx context.Context, messages []*batchMessage) error {
	var errs []error
	for _, msg := range messages {
		if err := bp.publisher.Publish(ctx, msg.subject, msg.msg); err != nil {
			errs = append(errs, err)
		}
	}
	return bp.collectErrors(errs)
}

// flushConcurrent publishes messages concurrently
func (bp *BatchPublisher) flushConcurrent(ctx context.Context, messages []*batchMessage) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(messages))

	for _, msg := range messages {
		wg.Add(1)
		go func(m *batchMessage) {
			defer wg.Done()
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
	return bp.collectErrors(errs)
}

// collectErrors combines errors into a single error or nil
func (bp *BatchPublisher) collectErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	if len(errs) == 1 {
		return errs[0]
	}
	return core.NewMultiError(errs)
}

// Count returns the number of pending messages
func (bp *BatchPublisher) Count() int {
	bp.mu.Lock()
	defer bp.mu.Unlock()
	return len(bp.messages)
}

// Close stops auto-flush and flushes remaining messages
func (bp *BatchPublisher) Close(ctx context.Context) error {
	bp.cancel()
	bp.wg.Wait()

	// Final flush
	return bp.Flush(ctx)
}
