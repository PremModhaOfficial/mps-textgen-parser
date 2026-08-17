package metrics

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"

	"github.com/stretchr/testify/assert"
)

const (
	initShouldNotFail      = "Init should not fail"
	testGaugeDesc          = "Test gauge"
	testHistogramDesc      = "Test histogram"
	failedToInitMetricsFmt = "Failed to init metrics: %v"
	concurrentTestDesc     = "Concurrent test"
	nilLabelsTestDesc      = "Nil labels test"
)

func TestInit(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false // No actual export

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	// Registry should be available
	metricsRegistry := R()
	assertions.NotNil(metricsRegistry, "R() should not return nil after Init")
}

func TestInitFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	os.Setenv("METRICS_ENABLED", "true")
	os.Setenv("METRICS_SERVICE_NAME", "test-service")
	defer func() {
		os.Unsetenv("METRICS_ENABLED")
		os.Unsetenv("METRICS_SERVICE_NAME")
	}()

	initError := InitFromEnv()
	assertions.NoError(initError, "InitFromEnv should not fail")
	defer Shutdown(context.Background())
}

func TestCounter(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create counter
	requestCounter := NewCounter("test_requests", "Test request counter")
	assertions.NotNil(requestCounter, "Counter should not be nil")

	// Increment
	requestCounter.Inc(testContext)
	requestCounter.Inc(testContext, Labels{"method": "GET"})

	// Add
	requestCounter.Add(testContext, 5)
	requestCounter.Add(testContext, 10, Labels{"method": "POST"})
}

func TestGauge(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create gauge
	activeConnectionsGauge := NewGauge("active_connections", "Active connections")
	assertions.NotNil(activeConnectionsGauge, "Gauge should not be nil")

	// Set
	activeConnectionsGauge.Set(testContext, 10)
	activeConnectionsGauge.Set(testContext, 20, Labels{"pool": "main"})

	// Inc/Dec
	activeConnectionsGauge.Inc(testContext)
	activeConnectionsGauge.Dec(testContext)

	// Add
	activeConnectionsGauge.Add(testContext, 5)
	activeConnectionsGauge.Add(testContext, -3, Labels{"pool": "secondary"})
}

func TestHistogram(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create histogram
	requestDurationHistogram := NewHistogram("request_duration", "Request duration")
	assertions.NotNil(requestDurationHistogram, "Histogram should not be nil")

	// Observe values
	requestDurationHistogram.Observe(testContext, 10.5)
	requestDurationHistogram.Observe(testContext, 25.0, Labels{"endpoint": "/api/users"})

	// Observe duration
	operationStartTime := time.Now()
	time.Sleep(10 * time.Millisecond)
	requestDurationHistogram.ObserveDuration(testContext, operationStartTime)
}

func TestTimer(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Start timer
	operationTimer := NewTimer(testContext, "operation_duration", "Operation duration")
	assertions.NotNil(operationTimer, "Timer should not be nil")

	// Do some work
	time.Sleep(5 * time.Millisecond)

	// Stop timer
	operationTimer.Stop()
}

func TestTimerWithLabels(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Start timer with initial labels
	databaseQueryTimer := NewTimer(testContext, "db_query_duration", "Database query duration",
		Labels{"database": "postgres"})
	assertions.NotNil(databaseQueryTimer, "Timer with labels should not be nil")

	// Do some work
	time.Sleep(5 * time.Millisecond)

	// Stop with additional labels
	databaseQueryTimer.StopWithLabels(Labels{"query_type": "select"})
}

func TestNamespace(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create namespaced registry
	httpNamespaceMetrics := Namespace("http")
	databaseNamespaceMetrics := Namespace("database")

	assertions.NotNil(httpNamespaceMetrics, "HTTP namespace should not be nil")
	assertions.NotNil(databaseNamespaceMetrics, "Database namespace should not be nil")

	// Metrics should have namespace prefix
	httpRequestsCounter := httpNamespaceMetrics.Counter("requests", "HTTP requests")
	databaseQueriesCounter := databaseNamespaceMetrics.Counter("queries", "Database queries")

	assertions.Equal("http_requests", httpRequestsCounter.Name(), "HTTP counter should have namespace prefix")
	assertions.Equal("database_queries", databaseQueriesCounter.Name(), "Database counter should have namespace prefix")

	httpRequestsCounter.Inc(testContext)
	databaseQueriesCounter.Inc(testContext)
}

func TestServiceMetrics(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create service metrics
	userServiceMetrics := NewServiceMetrics("user_service")
	assertions.NotNil(userServiceMetrics, "ServiceMetrics should not be nil")

	// Record a request
	operationStartTime := time.Now()
	time.Sleep(5 * time.Millisecond)
	operationDuration := time.Since(operationStartTime)

	userServiceMetrics.RecordRequest(testContext, "get_user", operationDuration, nil)

	// Record an error
	userServiceMetrics.RecordRequest(testContext, "create_user", operationDuration, context.DeadlineExceeded)

	// Use individual metrics
	userServiceMetrics.Active("list_users").Inc(testContext)
	userServiceMetrics.Active("list_users").Dec(testContext)
}

func TestServiceMetricsTimer(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	orderServiceMetrics := NewServiceMetrics("order_service")
	assertions.NotNil(orderServiceMetrics, "ServiceMetrics should not be nil")

	// Use timer
	processOrderTimer := orderServiceMetrics.Timer(testContext, "process_order")
	assertions.NotNil(processOrderTimer, "ServiceMetrics Timer should not be nil")
	time.Sleep(5 * time.Millisecond)
	processOrderTimer.Stop()
}

func TestLabels(t *testing.T) {
	assertions := assert.New(t)

	firstLabelsSet := Labels{"a": "1", "b": "2"}
	secondLabelsSet := Labels{"b": "3", "c": "4"}

	mergedLabels := firstLabelsSet.Merge(secondLabelsSet)

	assertions.Equal("1", mergedLabels["a"], "merged[a] should be '1'")
	assertions.Equal("3", mergedLabels["b"], "merged[b] should be '3' (second set overrides first)")
	assertions.Equal("4", mergedLabels["c"], "merged[c] should be '4'")
}

func TestDefaultConfig(t *testing.T) {
	assertions := assert.New(t)

	metricsConfig := DefaultConfig()

	assertions.True(metricsConfig.Enabled, "Enabled should be true by default")
	assertions.False(metricsConfig.OTELEnabled, "OTELEnabled should be false by default")
	assertions.False(metricsConfig.PrometheusEnabled, "PrometheusEnabled should be false by default")
	assertions.Equal(15*time.Second, metricsConfig.ExportInterval, "ExportInterval should be 15s by default")
}

func TestConfigFromEnv(t *testing.T) {
	assertions := assert.New(t)

	os.Setenv("METRICS_ENABLED", "true")
	os.Setenv("SERVICE_NAME", "env-service")
	os.Setenv("METRICS_OTEL_ENABLED", "true")
	os.Setenv("METRICS_OTEL_ENDPOINT", "collector:4317")
	os.Setenv("METRICS_PROMETHEUS_ENABLED", "true")
	os.Setenv("METRICS_PROMETHEUS_PORT", "8080")
	os.Setenv("METRICS_EXPORT_INTERVAL", "30s")

	defer func() {
		os.Unsetenv("METRICS_ENABLED")
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("METRICS_OTEL_ENABLED")
		os.Unsetenv("METRICS_OTEL_ENDPOINT")
		os.Unsetenv("METRICS_PROMETHEUS_ENABLED")
		os.Unsetenv("METRICS_PROMETHEUS_PORT")
		os.Unsetenv("METRICS_EXPORT_INTERVAL")
	}()

	metricsConfig := ConfigFromEnv()

	assertions.True(metricsConfig.Enabled, "Enabled should be true")
	assertions.Equal("env-service", metricsConfig.ServiceName, "ServiceName should be 'env-service'")
	assertions.True(metricsConfig.OTELEnabled, "OTELEnabled should be true")
	assertions.Equal("collector:4317", metricsConfig.OTELEndpoint, "OTELEndpoint should be 'collector:4317'")
	assertions.True(metricsConfig.PrometheusEnabled, "PrometheusEnabled should be true")
	assertions.Equal(8080, metricsConfig.PrometheusPort, "PrometheusPort should be 8080")
	assertions.Equal(30*time.Second, metricsConfig.ExportInterval, "ExportInterval should be 30s")
}

func TestRegistryStats(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	// Create some metrics
	NewCounter("test_counter1", "Test counter 1")
	NewCounter("test_counter2", "Test counter 2")
	NewGauge("test_gauge", testGaugeDesc)
	NewHistogram("test_histogram", testHistogramDesc)

	registryStats := R().Stats()

	assertions.Equal(2, registryStats.Counters, "Should have 2 counters")
	assertions.Equal(1, registryStats.Gauges, "Should have 1 gauge")
	assertions.Equal(1, registryStats.Histograms, "Should have 1 histogram")
}

func TestHelperMetrics(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Request counter
	apiRequestsCounter := RequestCounter("api")
	assertions.NotNil(apiRequestsCounter, "RequestCounter should not be nil")
	apiRequestsCounter.Inc(testContext, Labels{"method": "GET"})

	// Error counter
	apiErrorsCounter := ErrorCounter("api")
	assertions.NotNil(apiErrorsCounter, "ErrorCounter should not be nil")
	apiErrorsCounter.Inc(testContext, Labels{"code": "500"})

	// Duration histogram
	apiDurationHistogram := DurationHistogram("api")
	assertions.NotNil(apiDurationHistogram, "DurationHistogram should not be nil")
	apiDurationHistogram.Observe(testContext, 25.5)

	// Size histogram
	requestSizeHistogram := SizeHistogram("request")
	assertions.NotNil(requestSizeHistogram, "SizeHistogram should not be nil")
	requestSizeHistogram.Observe(testContext, 1024)

	// Active gauge
	activeConnectionsGauge := ActiveGauge("connections")
	assertions.NotNil(activeConnectionsGauge, "ActiveGauge should not be nil")
	activeConnectionsGauge.Inc(testContext)
	activeConnectionsGauge.Dec(testContext)
}

func TestTimeFunc(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Time a function
	functionExecuted := false
	TimeFunc(testContext, "slow_operation", "Slow operation duration", func() {
		time.Sleep(5 * time.Millisecond)
		functionExecuted = true
	})

	assertions.True(functionExecuted, "Function should have been executed")
}

func TestTimeFuncWithLabels(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Time a function with labels
	functionExecuted := false
	TimeFuncWithLabels(testContext, "labeled_operation", "Labeled operation duration",
		Labels{"type": "critical"}, func() {
			time.Sleep(5 * time.Millisecond)
			functionExecuted = true
		})

	assertions.True(functionExecuted, "Function should have been executed")
}

func TestMetricReuse(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	initError := Init(metricsConfig)
	assertions.NoError(initError, initShouldNotFail)
	defer Shutdown(context.Background())

	// Create same metric twice - should return same instance
	firstReuseCounter := NewCounter("reuse_counter", "Reuse test")
	secondReuseCounter := NewCounter("reuse_counter", "Reuse test")

	assertions.Equal(firstReuseCounter.Name(), secondReuseCounter.Name(), "Counter names should match for reused metric")
}

func TestDisabledMetrics(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = false // Disabled

	initError := Init(metricsConfig)
	assertions.NoError(initError, "Init should not fail with disabled metrics")
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Should still work (no-op)
	disabledCounter := NewCounter("disabled_counter", "Test")
	assertions.NotPanics(func() {
		disabledCounter.Inc(testContext)
	}, "Disabled counter should not panic on Inc")
}

/* ---------------------------------------- Cardinality Tests -------------------------------------------------- */

func TestHighCardinalityLabels(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("high_cardinality_counter", "High cardinality test")

	// Simulate high cardinality - 1000 unique label combinations
	for i := 0; i < 1000; i++ {
		labels := Labels{
			"user_id":    fmt.Sprintf("user-%d", i),
			"session_id": fmt.Sprintf("session-%d", i),
		}
		counter.Inc(testContext, labels)
	}

	// Should complete without error or memory issues
}

func TestLabelValueVariations(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("label_variations", "Label variation test")

	// Empty label value
	counter.Inc(testContext, Labels{"key": ""})

	// Long label value
	longValue := make([]byte, 1000)
	for i := range longValue {
		longValue[i] = 'a'
	}
	counter.Inc(testContext, Labels{"key": string(longValue)})

	// Special characters
	counter.Inc(testContext, Labels{"key": "value with spaces"})
	counter.Inc(testContext, Labels{"key": "value\nwith\nnewlines"})
	counter.Inc(testContext, Labels{"key": "value\twith\ttabs"})
}

/* ---------------------------------------- Concurrent Safety Tests -------------------------------------------------- */

func TestConcurrentCounterSafety(t *testing.T) {
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	if err != nil {
		t.Fatalf(failedToInitMetricsFmt, err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("concurrent_counter", concurrentTestDesc)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			labels := Labels{"goroutine": fmt.Sprintf("%d", id)}
			for j := 0; j < 1000; j++ {
				counter.Inc(testContext, labels)
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentGaugeSafety(t *testing.T) {
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	if err != nil {
		t.Fatalf(failedToInitMetricsFmt, err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	gauge := NewGauge("concurrent_gauge", concurrentTestDesc)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				gauge.Set(testContext, float64(j))
				gauge.Inc(testContext)
				gauge.Dec(testContext)
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentHistogramSafety(t *testing.T) {
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	if err != nil {
		t.Fatalf(failedToInitMetricsFmt, err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("concurrent_histogram", concurrentTestDesc)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				histogram.Observe(testContext, float64(j%100))
			}
		}(i)
	}
	wg.Wait()
}

func TestConcurrentMetricCreation(t *testing.T) {
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	if err != nil {
		t.Fatalf(failedToInitMetricsFmt, err)
	}
	defer Shutdown(context.Background())

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// All goroutines try to create the same metric
			_ = NewCounter("shared_counter", "Shared counter")
			_ = NewGauge("shared_gauge", "Shared gauge")
			_ = NewHistogram("shared_histogram", "Shared histogram")
		}(i)
	}
	wg.Wait()
}

/* ---------------------------------------- Edge Case Tests -------------------------------------------------- */

func TestZeroValue(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("zero_counter", "Zero test")

	// Adding zero should work
	counter.Add(testContext, 0)
}

func TestNegativeCounterAdd(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("negative_counter", "Negative test")

	// Should not panic - negative values are ignored for counters
	assertions.NotPanics(func() {
		counter.Add(testContext, -10)
	})
}

func TestTimerDoubleStop(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	timer := NewTimer(testContext, "double_stop", "Double stop test")

	// First stop
	timer.Stop()

	// Second stop should not panic or record again
	assertions.NotPanics(func() {
		timer.Stop()
	})
}

func TestNilLabels(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("nil_labels_counter", nilLabelsTestDesc)
	gauge := NewGauge("nil_labels_gauge", nilLabelsTestDesc)
	histogram := NewHistogram("nil_labels_histogram", nilLabelsTestDesc)

	// Should not panic with nil labels
	assertions.NotPanics(func() {
		counter.Inc(testContext, nil)
		gauge.Set(testContext, 10, nil)
		histogram.Observe(testContext, 50, nil)
	})
}

func TestEmptyLabelsMap(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("empty_labels_counter", "Empty labels test")

	emptyLabels := Labels{}
	counter.Inc(testContext, emptyLabels)
}

/* ---------------------------------------- MustInit Tests -------------------------------------------------- */

func TestMustInit(t *testing.T) {
	t.Run("successful MustInit", func(t *testing.T) {
		assertions := assert.New(t)
		Shutdown(context.Background())
		defer Shutdown(context.Background())

		metricsConfig := DefaultConfig()
		metricsConfig.Enabled = true
		metricsConfig.OTELEnabled = false

		assertions.NotPanics(func() {
			MustInit(metricsConfig)
		})
	})
}

/* ---------------------------------------- InitFromUnifiedConfig Tests -------------------------------------------------- */

func TestInitFromUnifiedConfig(t *testing.T) {
	assertions := assert.New(t)
	Shutdown(context.Background())
	defer Shutdown(context.Background())

	cfg := config.DefaultConfig()
	cfg.Metrics.Enabled = true
	cfg.Metrics.OTELEnabled = false

	err := InitFromUnifiedConfig(cfg)
	assertions.NoError(err)
}

/* ---------------------------------------- NoopProvider Tests -------------------------------------------------- */

func TestNoopProviderGauge(t *testing.T) {
	assertions := assert.New(t)

	// Create a noop provider
	provider := &noopProvider{}

	// Create a gauge
	gauge, err := provider.Gauge(MetricOpts{Name: "test_gauge", Description: testGaugeDesc})
	assertions.NoError(err)
	assertions.NotNil(gauge)

	// Operations should not panic
	testContext := context.Background()
	assertions.NotPanics(func() {
		gauge.Set(testContext, 10)
		gauge.Inc(testContext)
		gauge.Dec(testContext)
		gauge.Add(testContext, 5)
	})
}

func TestNoopProviderHistogram(t *testing.T) {
	assertions := assert.New(t)

	// Create a noop provider
	provider := &noopProvider{}

	// Create a histogram
	histogram, err := provider.Histogram(HistogramOpts{
		MetricOpts: MetricOpts{Name: "test_histogram", Description: testHistogramDesc, Unit: "ms"},
		Buckets:    DefaultHistogramBuckets,
	})
	assertions.NoError(err)
	assertions.NotNil(histogram)

	// Operations should not panic
	testContext := context.Background()
	assertions.NotPanics(func() {
		histogram.Observe(testContext, 100)
		histogram.ObserveDuration(testContext, time.Now())
	})
}

/* ---------------------------------------- Namespace Without Init Tests -------------------------------------------------- */

func TestNamespaceWithoutInit(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing state
	Shutdown(context.Background())

	// Namespace should work even without initialization (uses noop)
	registry := Namespace("test-namespace")
	assertions.NotNil(registry)

	testContext := context.Background()
	counter := registry.Counter("test_counter", "Test counter")
	assertions.NotNil(counter)

	// Should not panic
	assertions.NotPanics(func() {
		counter.Inc(testContext)
	})
}

/* ---------------------------------------- Noop Metric Implementation Tests -------------------------------------------------- */

func TestNoopCounterMethods(t *testing.T) {
	assertions := assert.New(t)

	counter := &noopCounter{name: "test_noop_counter"}

	assertions.Equal("test_noop_counter", counter.Name(), "noopCounter Name() should return the name")

	testContext := context.Background()

	// All methods should be no-ops without panicking
	assertions.NotPanics(func() {
		counter.Inc(testContext)
		counter.Inc(testContext, Labels{"key": "value"})
		counter.Add(testContext, 5.0)
		counter.Add(testContext, 10.0, Labels{"method": "GET"})
	})
}

func TestNoopGaugeMethods(t *testing.T) {
	assertions := assert.New(t)

	gauge := &noopGauge{name: "test_noop_gauge"}

	assertions.Equal("test_noop_gauge", gauge.Name(), "noopGauge Name() should return the name")

	testContext := context.Background()

	// All methods should be no-ops without panicking
	assertions.NotPanics(func() {
		gauge.Set(testContext, 42.0)
		gauge.Set(testContext, 42.0, Labels{"pool": "main"})
		gauge.Inc(testContext)
		gauge.Inc(testContext, Labels{"key": "value"})
		gauge.Dec(testContext)
		gauge.Dec(testContext, Labels{"key": "value"})
		gauge.Add(testContext, 5.0)
		gauge.Add(testContext, -3.0, Labels{"key": "value"})
	})
}

func TestNoopHistogramMethods(t *testing.T) {
	assertions := assert.New(t)

	histogram := &noopHistogram{name: "test_noop_histogram"}

	assertions.Equal("test_noop_histogram", histogram.Name(), "noopHistogram Name() should return the name")

	testContext := context.Background()

	// All methods should be no-ops without panicking
	assertions.NotPanics(func() {
		histogram.Observe(testContext, 100.0)
		histogram.Observe(testContext, 200.0, Labels{"endpoint": "/api"})
		histogram.ObserveDuration(testContext, time.Now())
		histogram.ObserveDuration(testContext, time.Now(), Labels{"method": "GET"})
	})
}

/* ---------------------------------------- CounterWithUnit / GaugeWithUnit Tests -------------------------------------------------- */

func TestCounterWithUnit(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create counter with unit
	counter := R().CounterWithUnit("bytes_sent", "Total bytes sent", "bytes")
	assertions.NotNil(counter, "CounterWithUnit should not return nil")

	// Use the counter
	counter.Inc(testContext)
	counter.Add(testContext, 1024, Labels{"protocol": "tcp"})

	// Create same counter again - should return cached instance
	counterAgain := R().CounterWithUnit("bytes_sent", "Total bytes sent", "bytes")
	assertions.Equal(counter.Name(), counterAgain.Name(), "CounterWithUnit should return cached instance")
}

func TestCounterWithUnitAndLabels(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create counter with unit and default labels
	counter := R().CounterWithUnit("requests_with_unit", "Request count", "1", Labels{"service": "api"})
	assertions.NotNil(counter)

	counter.Inc(testContext)
}

func TestGaugeWithUnit(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create gauge with unit
	gauge := R().GaugeWithUnit("memory_usage", "Memory usage", "bytes")
	assertions.NotNil(gauge, "GaugeWithUnit should not return nil")

	// Use the gauge
	gauge.Set(testContext, 1048576)
	gauge.Inc(testContext)
	gauge.Dec(testContext)
	gauge.Add(testContext, 512, Labels{"region": "us-east"})

	// Create same gauge again - should return cached instance
	gaugeAgain := R().GaugeWithUnit("memory_usage", "Memory usage", "bytes")
	assertions.Equal(gauge.Name(), gaugeAgain.Name(), "GaugeWithUnit should return cached instance")
}

func TestGaugeWithUnitAndLabels(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create gauge with unit and default labels
	gauge := R().GaugeWithUnit("cpu_usage_with_labels", "CPU usage", "percent", Labels{"host": "server1"})
	assertions.NotNil(gauge)

	gauge.Set(testContext, 75.5)
}

/* ---------------------------------------- RegistryStats.String Tests -------------------------------------------------- */

func TestRegistryStatsString(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	// Create some metrics
	NewCounter("stats_string_counter", "Test counter")
	NewGauge("stats_string_gauge", testGaugeDesc)
	NewHistogram("stats_string_histogram", testHistogramDesc)

	stats := R().Stats()
	statsString := stats.String()

	assertions.Contains(statsString, "Registry{", "Stats string should contain Registry prefix")
	assertions.Contains(statsString, "counters=", "Stats string should contain counters count")
	assertions.Contains(statsString, "gauges=", "Stats string should contain gauges count")
	assertions.Contains(statsString, "histograms=", "Stats string should contain histograms count")
}

func TestRegistryStatsStringEmpty(t *testing.T) {
	assertions := assert.New(t)

	stats := RegistryStats{Counters: 0, Gauges: 0, Histograms: 0}
	statsString := stats.String()

	assertions.Equal("Registry{counters=0, gauges=0, histograms=0}", statsString)
}

/* ---------------------------------------- MetricType.String Tests -------------------------------------------------- */

func TestMetricTypeString(t *testing.T) {
	assertions := assert.New(t)

	assertions.Equal("counter", CounterType.String(), "CounterType.String() should return 'counter'")
	assertions.Equal("gauge", GaugeType.String(), "GaugeType.String() should return 'gauge'")
	assertions.Equal("histogram", HistogramType.String(), "HistogramType.String() should return 'histogram'")

	// Unknown metric type
	unknownType := MetricType(999)
	assertions.Equal("unknown", unknownType.String(), "Unknown MetricType.String() should return 'unknown'")
}

/* ---------------------------------------- MustInit Panic Tests -------------------------------------------------- */

func TestMustInitPanic(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	// Create a config that will cause Init to fail
	// OTELEnabled=true with an invalid endpoint should cause a failure
	badConfig := DefaultConfig()
	badConfig.Enabled = true
	badConfig.OTELEnabled = true
	badConfig.OTELEndpoint = "" // Empty endpoint
	badConfig.OTELProtocol = "invalid-protocol"

	assertions.Panics(func() {
		MustInit(badConfig)
	}, "MustInit should panic when Init fails")

	// Clean up in case it somehow succeeded
	Shutdown(context.Background())
}

/* ---------------------------------------- Registry.Shutdown Tests -------------------------------------------------- */

func TestRegistryShutdown(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)

	registry := R()
	assertions.NotNil(registry)

	// Registry shutdown should succeed
	shutdownErr := registry.Shutdown(context.Background())
	assertions.NoError(shutdownErr, "Registry.Shutdown should not return error")

	// Clean up global state
	Shutdown(context.Background())
}

func TestRegistryShutdownNoop(t *testing.T) {
	assertions := assert.New(t)

	// Create a registry with noop provider
	provider := &noopProvider{}
	registry := NewRegistry(provider, "test")

	shutdownErr := registry.Shutdown(context.Background())
	assertions.NoError(shutdownErr, "Noop registry shutdown should not return error")
}

/* ---------------------------------------- otel Metric Name Tests -------------------------------------------------- */

func TestOtelCounterName(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	counter := R().Counter("otel_name_test_counter", "Test counter for name")
	assertions.Equal("app_otel_name_test_counter", counter.Name(), "Counter Name() should return the full namespaced name")
}

func TestOtelGaugeName(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	gauge := R().Gauge("otel_name_test_gauge", "Test gauge for name")
	assertions.Equal("app_otel_name_test_gauge", gauge.Name(), "Gauge Name() should return the full namespaced name")
}

func TestOtelHistogramName(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	histogram := R().Histogram("otel_name_test_histogram", "Test histogram for name")
	assertions.Equal("app_otel_name_test_histogram", histogram.Name(), "Histogram Name() should return the full namespaced name")
}

/* ---------------------------------------- StopWithLabels nil timerLabels Tests -------------------------------------------------- */

func TestStopWithLabelsNilTimerLabels(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	// Create timer WITHOUT initial labels - timerLabels will be nil
	timer := NewTimer(testContext, "no_initial_labels_timer", "Timer without initial labels")
	assertions.NotNil(timer)

	// StopWithLabels should follow the nil timerLabels path
	assertions.NotPanics(func() {
		timer.StopWithLabels(Labels{"status": "success"})
	}, "StopWithLabels with nil timerLabels should not panic")
}

func TestStopWithLabelsDoubleStop(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)
	defer Shutdown(context.Background())

	testContext := context.Background()

	timer := NewTimer(testContext, "double_stop_labels_timer", "Double stop with labels")

	// First stop with labels
	timer.StopWithLabels(Labels{"status": "ok"})

	// Second stop should be a no-op
	assertions.NotPanics(func() {
		timer.StopWithLabels(Labels{"status": "duplicate"})
	}, "Double StopWithLabels should not panic")
}

/* ---------------------------------------- Provider Double Shutdown Tests -------------------------------------------------- */

func TestProviderDoubleShutdown(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	provider, err := NewOTELProvider(metricsConfig)
	assertions.NoError(err)
	assertions.NotNil(provider)

	// First shutdown
	shutdownErr := provider.Shutdown(context.Background())
	assertions.NoError(shutdownErr, "First shutdown should succeed")

	// Second shutdown should return nil (idempotent)
	shutdownErr2 := provider.Shutdown(context.Background())
	assertions.NoError(shutdownErr2, "Second shutdown should return nil")
}

/* ---------------------------------------- Global Shutdown Tests -------------------------------------------------- */

func TestGlobalShutdownWithoutInit(t *testing.T) {
	assertions := assert.New(t)

	// Ensure clean state
	Shutdown(context.Background())

	// Shutdown without init should not error
	err := Shutdown(context.Background())
	assertions.NoError(err, "Shutdown without Init should not return error")
}

func TestGlobalDoubleShutdown(t *testing.T) {
	assertions := assert.New(t)

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	err := Init(metricsConfig)
	assertions.NoError(err)

	// First shutdown
	err = Shutdown(context.Background())
	assertions.NoError(err, "First global shutdown should succeed")

	// Second shutdown should be safe
	err = Shutdown(context.Background())
	assertions.NoError(err, "Second global shutdown should succeed")
}

/* ---------------------------------------- NoopProvider Shutdown Tests -------------------------------------------------- */

func TestNoopProviderShutdown(t *testing.T) {
	assertions := assert.New(t)

	provider := &noopProvider{}

	err := provider.Shutdown(context.Background())
	assertions.NoError(err, "noopProvider Shutdown should return nil")
}

/* ---------------------------------------- Pool Wrapper Nil Safety Tests -------------------------------------------------- */

func TestPutLabelsWrapperNil(t *testing.T) {
	assert.NotPanics(t, func() {
		putLabelsWrapper(nil)
	}, "putLabelsWrapper(nil) should not panic")
}

func TestPutAttributesWrapperNil(t *testing.T) {
	assert.NotPanics(t, func() {
		putAttributesWrapper(nil)
	}, "putAttributesWrapper(nil) should not panic")
}

func TestPutPooledAttrsNil(t *testing.T) {
	assert.NotPanics(t, func() {
		putPooledAttrs(nil)
	}, "putPooledAttrs(nil) should not panic")
}

/* ---------------------------------------- mergeAllLabels Tests -------------------------------------------------- */

func TestMergeAllLabelsEmpty(t *testing.T) {
	assertions := assert.New(t)

	result := mergeAllLabels()
	assertions.Nil(result, "mergeAllLabels with no args should return nil")
}

func TestMergeAllLabelsMultiple(t *testing.T) {
	assertions := assert.New(t)

	result := mergeAllLabels(Labels{"a": "1"}, Labels{"b": "2"}, Labels{"a": "3"})
	assertions.Equal("3", result["a"], "Later labels should override earlier ones")
	assertions.Equal("2", result["b"], "Non-conflicting labels should be preserved")
}
