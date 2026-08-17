package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   METRICS COLLECTOR TESTS
   ======================================================================================================== */

var (
	metricsPublishFailed = "publish failed"
	metricsProcessFailed = "process failed"
	metricsPublishError  = "publish error"
	metricsProcessError  = "process error"
)

func TestNewMetricsCollector(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()

	assertions.NotNil(mc, "MetricsCollector should not be nil")

	metrics := mc.Collect()
	assertions.Equal(int64(0), metrics.PublishCount, "Initial PublishCount should be 0")
	assertions.Equal(int64(0), metrics.ReceiveCount, "Initial ReceiveCount should be 0")
}

func TestMetricsCollectorInterceptPublish(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())

	err := wrapped(context.Background(), newTestMsg(testSubject, testPayloadData))

	assertions.NoError(err, "Handler should not return error")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.PublishCount, "PublishCount should be 1")
	assertions.Equal(int64(0), metrics.PublishErrors, "PublishErrors should be 0")
	assertions.Greater(metrics.PublishLatency, time.Duration(0), "PublishLatency should be set")
	assertions.Equal(int64(len(testPayloadData)), metrics.BytesSent, "BytesSent should match data size")
}

func TestMetricsCollectorInterceptPublishError(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	handler := func(ctx context.Context, msg *nats.Msg) error {
		return errors.New(metricsPublishFailed)
	}

	wrapped := mc.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.Error(err, "Handler should return error")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.PublishCount, "PublishCount should be 1")
	assertions.Equal(int64(1), metrics.PublishErrors, "PublishErrors should be 1")
}

func TestMetricsCollectorInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptSubscribe()(noopSubscribeHandler())

	err := wrapped(context.Background(), newTestMsg(testSubject, testPayloadData))

	assertions.NoError(err, "Handler should not return error")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.ReceiveCount, "ReceiveCount should be 1")
	assertions.Equal(int64(0), metrics.ReceiveErrors, "ReceiveErrors should be 0")
	assertions.Greater(metrics.ProcessTime, time.Duration(0), "ProcessTime should be set")
	assertions.Equal(int64(len(testPayloadData)), metrics.BytesReceived, "BytesReceived should match data size")
}

func TestMetricsCollectorInterceptSubscribeError(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	handler := func(ctx context.Context, msg *nats.Msg) error {
		return errors.New(metricsProcessFailed)
	}

	wrapped := mc.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.Error(err, "Handler should return error")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.ReceiveCount, "ReceiveCount should be 1")
	assertions.Equal(int64(1), metrics.ReceiveErrors, "ReceiveErrors should be 1")
}

func TestMetricsCollectorReset(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.PublishCount, "PublishCount should be 1 before reset")

	mc.Reset()

	metrics = mc.Collect()
	assertions.Equal(int64(0), metrics.PublishCount, "PublishCount should be 0 after reset")
	assertions.Equal(int64(0), metrics.BytesSent, "BytesSent should be 0 after reset")
}

/* ========================================================================================================
   CALLBACK TESTS
   ======================================================================================================== */

func TestMetricsCollectorOnPublishCallback(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()

	var callbackSubject string
	var callbackLatency time.Duration
	var callbackErr error
	callCount := 0

	mc.OnPublish(func(subject string, latency time.Duration, err error) {
		callbackSubject = subject
		callbackLatency = latency
		callbackErr = err
		callCount++
	})

	wrapped := mc.InterceptPublish()(noopPublishHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.Equal(1, callCount, "Callback should be called once")
	assertions.Equal(testSubject, callbackSubject, "Callback subject should match")
	assertions.Greater(callbackLatency, time.Duration(0), "Callback latency should be positive")
	assertions.NoError(callbackErr, "Callback error should be nil")
}

func TestMetricsCollectorOnReceiveCallback(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()

	var callbackSubject string
	var callbackProcessTime time.Duration
	callCount := 0

	mc.OnReceive(func(subject string, processTime time.Duration, err error) {
		callbackSubject = subject
		callbackProcessTime = processTime
		callCount++
	})

	wrapped := mc.InterceptSubscribe()(noopSubscribeHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.Equal(1, callCount, "Callback should be called once")
	assertions.Equal(testSubject, callbackSubject, "Callback subject should match")
	assertions.Greater(callbackProcessTime, time.Duration(0), "Callback process time should be positive")
}

/* ========================================================================================================
   PER-SUBJECT METRICS TESTS
   ======================================================================================================== */

func TestNewPerSubjectMetrics(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()

	assertions.NotNil(psm, "PerSubjectMetrics should not be nil")
}

func TestPerSubjectMetricsGet(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()

	mc1 := psm.Get(testSubject1)
	mc2 := psm.Get(testSubject2)
	mc1Again := psm.Get(testSubject1)

	assertions.NotNil(mc1, "MetricsCollector should not be nil")
	assertions.NotNil(mc2, "MetricsCollector should not be nil")
	assertions.Same(mc1, mc1Again, "Same subject should return same collector")
	assertions.NotSame(mc1, mc2, "Different subjects should return different collectors")
}

func TestPerSubjectMetricsInterceptPublish(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()
	wrapped := psm.InterceptPublish()(noopPublishHandler())
	ctx := context.Background()

	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject2, testPayload))

	allMetrics := psm.All()

	assertions.Equal(2, len(allMetrics), "Should have metrics for 2 subjects")
	assertions.Equal(int64(2), allMetrics[testSubject1].PublishCount, "subject1 should have 2 publishes")
	assertions.Equal(int64(1), allMetrics[testSubject2].PublishCount, "subject2 should have 1 publish")
}

func TestPerSubjectMetricsInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()
	wrapped := psm.InterceptSubscribe()(noopSubscribeHandler())
	ctx := context.Background()

	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject2, testPayload))

	allMetrics := psm.All()

	assertions.Equal(2, len(allMetrics), "Should have metrics for 2 subjects")
	assertions.Equal(int64(2), allMetrics[testSubject1].ReceiveCount, "subject1 should have 2 receives")
	assertions.Equal(int64(1), allMetrics[testSubject2].ReceiveCount, "subject2 should have 1 receive")
}

/* ========================================================================================================
   HELPER FUNCTION TESTS
   ======================================================================================================== */

func TestMetricsMiddlewareFactory(t *testing.T) {
	assertions := assert.New(t)

	mc := MetricsMiddleware()

	assertions.NotNil(mc, "MetricsMiddleware() should return MetricsCollector")
}

/* ========================================================================================================
   NULL MESSAGE HANDLING TESTS
   ======================================================================================================== */

func TestMetricsCollectorNilMessage(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())

	err := wrapped(context.Background(), nil)

	assertions.NoError(err, "Should handle nil message without panic")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.PublishCount, "PublishCount should be 1")
	assertions.Equal(int64(0), metrics.BytesSent, "BytesSent should be 0 for nil message")
}

func TestMetricsCollectorSubscribeNilMessage(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptSubscribe()(noopSubscribeHandler())

	err := wrapped(context.Background(), nil)

	assertions.NoError(err, "Should handle nil message without panic")

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.ReceiveCount, "ReceiveCount should be 1")
	assertions.Equal(int64(0), metrics.BytesReceived, "BytesReceived should be 0 for nil message")
}

/* ========================================================================================================
   CONCURRENCY TESTS
   ======================================================================================================== */

func TestMetricsCollectorConcurrency(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())
	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)

	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			wrapped(ctx, msg)
			done <- true
		}()
	}

	for i := 0; i < 100; i++ {
		<-done
	}

	metrics := mc.Collect()
	assertions.Equal(int64(100), metrics.PublishCount, "PublishCount should be 100 after concurrent publishes")
}
