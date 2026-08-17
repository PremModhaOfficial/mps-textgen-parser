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

func TestInit(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false // No actual export

	initError := Init(metricsConfig)
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
	defer Shutdown(context.Background())

	// Create some metrics
	NewCounter("test_counter1", "Test counter 1")
	NewCounter("test_counter2", "Test counter 2")
	NewGauge("test_gauge", "Test gauge")
	NewHistogram("test_histogram", "Test histogram")

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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
	assertions.NoError(initError, "Init should not fail")
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
		t.Fatalf("Failed to init metrics: %v", err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("concurrent_counter", "Concurrent test")

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
		t.Fatalf("Failed to init metrics: %v", err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	gauge := NewGauge("concurrent_gauge", "Concurrent test")

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
		t.Fatalf("Failed to init metrics: %v", err)
	}
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("concurrent_histogram", "Concurrent test")

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
		t.Fatalf("Failed to init metrics: %v", err)
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
	counter := NewCounter("nil_labels_counter", "Nil labels test")
	gauge := NewGauge("nil_labels_gauge", "Nil labels test")
	histogram := NewHistogram("nil_labels_histogram", "Nil labels test")

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
	gauge, err := provider.Gauge(MetricOpts{Name: "test_gauge", Description: "Test gauge"})
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
		MetricOpts: MetricOpts{Name: "test_histogram", Description: "Test histogram", Unit: "ms"},
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
