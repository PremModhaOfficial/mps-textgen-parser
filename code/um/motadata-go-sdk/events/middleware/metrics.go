package middleware

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
)

// Metrics holds collected metrics
type Metrics struct {
	// Publish metrics
	PublishCount   int64
	PublishErrors  int64
	PublishLatency time.Duration

	// Subscribe metrics
	ReceiveCount  int64
	ReceiveErrors int64
	ProcessTime   time.Duration

	// Bytes
	BytesSent     int64
	BytesReceived int64
}

// MetricsCollector collects metrics
type MetricsCollector struct {
	mu sync.RWMutex

	// Publish
	publishCount   atomic.Int64
	publishErrors  atomic.Int64
	publishLatency atomic.Int64 // nanoseconds

	// Subscribe
	receiveCount  atomic.Int64
	receiveErrors atomic.Int64
	processTime   atomic.Int64 // nanoseconds

	// Bytes
	bytesSent     atomic.Int64
	bytesReceived atomic.Int64

	// Callbacks
	onPublish  func(subject string, latency time.Duration, err error)
	onReceive  func(subject string, processTime time.Duration, err error)
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{}
}

// OnPublish sets a callback for publish events
func (m *MetricsCollector) OnPublish(cb func(subject string, latency time.Duration, err error)) {
	m.mu.Lock()
	m.onPublish = cb
	m.mu.Unlock()
}

// OnReceive sets a callback for receive events
func (m *MetricsCollector) OnReceive(cb func(subject string, processTime time.Duration, err error)) {
	m.mu.Lock()
	m.onReceive = cb
	m.mu.Unlock()
}

// Collect returns the current metrics
func (m *MetricsCollector) Collect() Metrics {
	return Metrics{
		PublishCount:   m.publishCount.Load(),
		PublishErrors:  m.publishErrors.Load(),
		PublishLatency: time.Duration(m.publishLatency.Load()),
		ReceiveCount:   m.receiveCount.Load(),
		ReceiveErrors:  m.receiveErrors.Load(),
		ProcessTime:    time.Duration(m.processTime.Load()),
		BytesSent:      m.bytesSent.Load(),
		BytesReceived:  m.bytesReceived.Load(),
	}
}

// Reset resets all metrics
func (m *MetricsCollector) Reset() {
	m.publishCount.Store(0)
	m.publishErrors.Store(0)
	m.publishLatency.Store(0)
	m.receiveCount.Store(0)
	m.receiveErrors.Store(0)
	m.processTime.Store(0)
	m.bytesSent.Store(0)
	m.bytesReceived.Store(0)
}

// InterceptPublish returns the publish middleware
func (m *MetricsCollector) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			start := time.Now()

			err := next(ctx, subject, msg)

			latency := time.Since(start)

			m.publishCount.Add(1)
			m.publishLatency.Add(int64(latency))

			if msg != nil && msg.Data != nil {
				m.bytesSent.Add(int64(len(msg.Data)))
			}

			if err != nil {
				m.publishErrors.Add(1)
			}

			// Callback
			m.mu.RLock()
			cb := m.onPublish
			m.mu.RUnlock()
			if cb != nil {
				cb(subject, latency, err)
			}

			return err
		}
	}
}

// InterceptSubscribe returns the subscribe middleware
func (m *MetricsCollector) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			start := time.Now()

			m.receiveCount.Add(1)

			if msg != nil && msg.Data != nil {
				m.bytesReceived.Add(int64(len(msg.Data)))
			}

			err := next(ctx, msg)

			processTime := time.Since(start)
			m.processTime.Add(int64(processTime))

			if err != nil {
				m.receiveErrors.Add(1)
			}

			// Callback
			m.mu.RLock()
			cb := m.onReceive
			m.mu.RUnlock()
			if cb != nil {
				subject := ""
				if msg != nil {
					subject = msg.Subject
				}
				cb(subject, processTime, err)
			}

			return err
		}
	}
}

// MetricsMiddleware is a simple metrics middleware
func MetricsMiddleware() *MetricsCollector {

	return NewMetricsCollector()
}

// OTELMetricsMiddleware provides OpenTelemetry-based metrics collection
type OTELMetricsMiddleware struct {
	publishCounter   metrics.Counter
	publishErrors    metrics.Counter
	publishDuration  metrics.Histogram
	receiveCounter   metrics.Counter
	receiveErrors    metrics.Counter
	receiveDuration  metrics.Histogram
	bytesSent        metrics.Counter
	bytesReceived    metrics.Counter
}

// NewOTELMetricsMiddleware creates a new OpenTelemetry metrics middleware
func NewOTELMetricsMiddleware() *OTELMetricsMiddleware {

	return &OTELMetricsMiddleware{
		publishCounter:   metrics.NewCounter("events_publish_total", "Total number of messages published"),
		publishErrors:    metrics.NewCounter("events_publish_errors_total", "Total number of publish errors"),
		publishDuration:  metrics.DurationHistogram("events_publish"),
		receiveCounter:   metrics.NewCounter("events_receive_total", "Total number of messages received"),
		receiveErrors:    metrics.NewCounter("events_receive_errors_total", "Total number of receive errors"),
		receiveDuration:  metrics.DurationHistogram("events_receive"),
		bytesSent:        metrics.NewCounter("events_bytes_sent_total", "Total bytes sent"),
		bytesReceived:    metrics.NewCounter("events_bytes_received_total", "Total bytes received"),
	}
}

// InterceptPublish returns the publish middleware with OpenTelemetry metrics
func (otelMetrics *OTELMetricsMiddleware) InterceptPublish() PublishMiddleware {

	return func(next PublishHandler) PublishHandler {

		return func(ctx context.Context, subject string, msg *core.Message) error {

			start := time.Now()

			err := next(ctx, subject, msg)

			duration := time.Since(start)

			labels := metrics.Labels{
				"subject": subject,
			}

			otelMetrics.publishCounter.Inc(ctx, labels)
			otelMetrics.publishDuration.Observe(ctx, float64(duration.Milliseconds()), labels)

			if msg != nil && msg.Data != nil {

				otelMetrics.bytesSent.Add(ctx, float64(len(msg.Data)), labels)
			}

			if err != nil {

				otelMetrics.publishErrors.Inc(ctx, labels)
			}

			return err
		}
	}
}

// InterceptSubscribe returns the subscribe middleware with OpenTelemetry metrics
func (otelMetrics *OTELMetricsMiddleware) InterceptSubscribe() SubscribeMiddleware {

	return func(next SubscribeHandler) SubscribeHandler {

		return func(ctx context.Context, msg *core.Message) error {

			start := time.Now()

			subject := ""

			if msg != nil {

				subject = msg.Subject

				otelMetrics.bytesReceived.Add(ctx, float64(len(msg.Data)), metrics.Labels{"subject": subject})
			}

			otelMetrics.receiveCounter.Inc(ctx, metrics.Labels{"subject": subject})

			err := next(ctx, msg)

			duration := time.Since(start)

			labels := metrics.Labels{
				"subject": subject,
			}

			otelMetrics.receiveDuration.Observe(ctx, float64(duration.Milliseconds()), labels)

			if err != nil {

				otelMetrics.receiveErrors.Inc(ctx, labels)
			}

			return err
		}
	}
}

// PerSubjectMetrics collects metrics per subject
type PerSubjectMetrics struct {
	mu      sync.RWMutex
	metrics map[string]*MetricsCollector
}

// NewPerSubjectMetrics creates per-subject metrics
func NewPerSubjectMetrics() *PerSubjectMetrics {
	return &PerSubjectMetrics{
		metrics: make(map[string]*MetricsCollector),
	}
}

// Get returns metrics for a subject
func (p *PerSubjectMetrics) Get(subject string) *MetricsCollector {
	p.mu.RLock()
	m, ok := p.metrics[subject]
	p.mu.RUnlock()

	if ok {
		return m
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Double-check after acquiring write lock
	if m, ok = p.metrics[subject]; ok {
		return m
	}

	m = NewMetricsCollector()
	p.metrics[subject] = m
	return m
}

// All returns all collected metrics
func (p *PerSubjectMetrics) All() map[string]Metrics {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make(map[string]Metrics, len(p.metrics))
	for subject, collector := range p.metrics {
		result[subject] = collector.Collect()
	}
	return result
}

// InterceptPublish returns the publish middleware
func (p *PerSubjectMetrics) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			m := p.Get(subject)
			handler := m.InterceptPublish()(next)
			return handler(ctx, subject, msg)
		}
	}
}

// InterceptSubscribe returns the subscribe middleware
func (p *PerSubjectMetrics) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			subject := ""
			if msg != nil {
				subject = msg.Subject
			}
			m := p.Get(subject)
			handler := m.InterceptSubscribe()(next)
			return handler(ctx, msg)
		}
	}
}
