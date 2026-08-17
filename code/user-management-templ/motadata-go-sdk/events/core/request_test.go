package core

import (
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

var (
	testEndpointName       = "test"
	testEndpointSubject    = "test.subject"
	testBenchName          = "bench"
	testBenchSubject       = "bench.subject"
	testBenchParallelName  = "bench-parallel"
	testBenchParallelSubj  = "bench.parallel"
	testCustomSubject      = "custom.subject"
	testCustomQueueGroup   = "custom-q"
	testGroupQueueGroup    = "group-q"
	testXCustomHeader      = "X-Custom"
	testHeaderValue        = "value"
	testContentTypeHeader  = "Content-Type"
	testApplicationJSON    = "application/json"
	testXRequestIDHeader   = "X-Request-Id"
	testRequestIDValue     = "123"
	testXFirstHeader       = "X-First"
	testXSecondHeader      = "X-Second"
	testEmptyName          = "empty"
	testEmptySubject       = "empty.subject"
	testCustomStatsName    = "custom-stats"
	testConcurrentName     = "concurrent"
	testConcurrentSubject  = "concurrent.subject"
	testMetricCustomValue  = "custom-value"
	testBadRequestCode     = "400"
	testBadRequestDesc     = "bad request"
	testServiceUnavailCode = "503"
	testServiceUnavailDesc = "service unavailable"
	testInternalErrorCode  = "500"
	testInternalErrorDesc  = "internal error"
	testErrorDesc          = "error"
	testTestError          = "test error"
	testTestJSON           = `{"name":"test","value":42}`
)

/* ========================================================================================================
   CONTROL SUBJECT TESTS
   ======================================================================================================== */

func TestControlSubject(t *testing.T) {

	testCases := []struct {
		name     string
		verb     Verb
		svcName  string
		svcID    string
		expected string
	}{
		{"ping all services", PingVerb, "", "", "$SRV.PING"},
		{"ping by name", PingVerb, "MyService", "", "$SRV.PING.MyService"},
		{"ping by name and id", PingVerb, "MyService", "abc123", "$SRV.PING.MyService.abc123"},
		{"stats all services", StatsVerb, "", "", "$SRV.STATS"},
		{"stats with name and id", StatsVerb, "MySvc", "id1", "$SRV.STATS.MySvc.id1"},
		{"info by name", InfoVerb, "TestService", "", "$SRV.INFO.TestService"},
		{"info all", InfoVerb, "", "", "$SRV.INFO"},
		{"unknown verb", Verb(999), "", "", "$SRV.UNKNOWN"},
		{"unknown verb with name", Verb(999), "Svc", "", "$SRV.UNKNOWN.Svc"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)

			result := ControlSubject(tc.verb, tc.svcName, tc.svcID)

			assertions.Equal(tc.expected, result, "ControlSubject should match expected")
		})
	}
}

/* ========================================================================================================
   VERB STRING TESTS
   ======================================================================================================== */

func TestVerbString(t *testing.T) {

	testCases := []struct {
		verb     Verb
		expected string
	}{
		{PingVerb, "PING"},
		{StatsVerb, "STATS"},
		{InfoVerb, "INFO"},
		{Verb(999), "UNKNOWN"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {

			assertions := assert.New(t)

			assertions.Equal(tc.expected, tc.verb.String(), "Verb string should match")
		})
	}
}

/* ========================================================================================================
   SERVICE ERROR TESTS
   ======================================================================================================== */

func TestServiceError(t *testing.T) {

	assertions := assert.New(t)

	err := &ServiceError{
		Code:        testBadRequestCode,
		Description: testBadRequestDesc,
		Data:        []byte("details"),
	}

	assertions.Equal(testBadRequestDesc, err.Error(), "Error should return description")
}

func TestServiceErrorWithData(t *testing.T) {

	assertions := assert.New(t)

	err := &ServiceError{
		Code:        testServiceUnavailCode,
		Description: testServiceUnavailDesc,
		Data:        []byte(`{"retry_after":30}`),
	}

	assertions.Equal(testServiceUnavailDesc, err.Error(), "Error should return description")
	assertions.Equal(testServiceUnavailCode, err.Code, "Code should be preserved")
	assertions.NotNil(err.Data, "Data should be preserved")
}

/* ========================================================================================================
   REQUEST HANDLER FUNC TESTS
   ======================================================================================================== */

func TestRequestHandlerFuncHandle(t *testing.T) {

	assertions := assert.New(t)

	var called bool
	handler := RequestHandlerFunc(func(r Request) {
		called = true
	})

	handler.Handle(nil)

	assertions.True(called, "Handle should invoke the underlying function")
}

/* ========================================================================================================
   RESPONSE TESTS
   ======================================================================================================== */

func TestResponseIsError(t *testing.T) {

	t.Run("no error", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Data:    []byte("ok"),
			Headers: make(nats.Header),
		}

		assertions.False(resp.IsError(), "IsError should be false for successful response")
	})

	t.Run("with error", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Error: &ServiceError{
				Code:        testInternalErrorCode,
				Description: testInternalErrorDesc,
			},
		}

		assertions.True(resp.IsError(), "IsError should be true for error response")
	})
}

func TestResponseJSON(t *testing.T) {

	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("success", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Data: []byte(testTestJSON),
		}

		var data testData
		err := resp.JSON(&data)

		assertions.NoError(err, "JSON should parse successfully")
		assertions.Equal(testEndpointName, data.Name, "Name should match")
		assertions.Equal(42, data.Value, "Value should match")
	})

	t.Run("with error", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Error: &ServiceError{
				Code:        testBadRequestCode,
				Description: testBadRequestDesc,
			},
		}

		var data testData
		err := resp.JSON(&data)

		assertions.Error(err, "JSON should return error when response has error")
	})

	t.Run("invalid data", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Data: []byte("not-json"),
		}

		var data testData
		err := resp.JSON(&data)

		assertions.Error(err, "JSON should fail for invalid data")
	})

	t.Run("nil data", func(t *testing.T) {

		assertions := assert.New(t)

		resp := &Response{
			Data: nil,
		}

		var data testData
		err := resp.JSON(&data)

		assertions.Error(err, "JSON should fail for nil data")
	})
}

/* ========================================================================================================
   ENDPOINT STATS TESTS
   ======================================================================================================== */

func TestEndpointStats(t *testing.T) {

	assertions := assert.New(t)

	endpoint := &Endpoint{
		Name:       testEndpointName,
		Subject:    testEndpointSubject,
		QueueGroup: "q",
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)
	endpoint.RecordRequest(20*time.Millisecond, nil)
	endpoint.RecordRequest(30*time.Millisecond, &ServiceError{Description: testTestError})

	stats := endpoint.Stats()

	assertions.Equal(int64(3), stats.NumRequests, "NumRequests should be 3")
	assertions.Equal(int64(1), stats.NumErrors, "NumErrors should be 1")
	assertions.Equal(60*time.Millisecond, stats.TotalProcessingTime, "TotalProcessingTime should be 60ms")
	assertions.Equal(20*time.Millisecond, stats.AverageProcessingTime, "AverageProcessingTime should be 20ms")
	assertions.Equal(testTestError, stats.LastError, "LastError should match")
}

func TestEndpointStatsZeroRequests(t *testing.T) {

	assertions := assert.New(t)

	endpoint := &Endpoint{
		Name:    testEmptyName,
		Subject: testEmptySubject,
	}

	stats := endpoint.Stats()

	assertions.Equal(int64(0), stats.NumRequests, "NumRequests should be 0")
	assertions.Equal(time.Duration(0), stats.AverageProcessingTime, "Average should be 0 for no requests")
}

func TestEndpointStatsWithCustomHandler(t *testing.T) {

	assertions := assert.New(t)

	customData := map[string]string{"metric": testMetricCustomValue}

	endpoint := &Endpoint{
		Name:    testCustomStatsName,
		Subject: testCustomSubject,
		statsHandler: func(e *Endpoint) any {
			return customData
		},
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)

	stats := endpoint.Stats()

	assertions.Equal(customData, stats.Data, "Custom stats data should be returned")
	assertions.Equal(int64(1), stats.NumRequests, "NumRequests should be 1")
}

func TestEndpointReset(t *testing.T) {

	assertions := assert.New(t)

	endpoint := &Endpoint{
		Name:    testEndpointName,
		Subject: testEndpointSubject,
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)
	endpoint.RecordRequest(20*time.Millisecond, &ServiceError{Description: testErrorDesc})

	endpoint.Reset()
	stats := endpoint.Stats()

	assertions.Equal(int64(0), stats.NumRequests, "NumRequests after reset should be 0")
	assertions.Equal(int64(0), stats.NumErrors, "NumErrors after reset should be 0")
	assertions.Equal(time.Duration(0), stats.TotalProcessingTime, "TotalProcessingTime after reset should be 0")
	assertions.Empty(stats.LastError, "LastError after reset should be empty")
}

/* ========================================================================================================
   ENDPOINT CONCURRENT ACCESS TESTS
   ======================================================================================================== */

func TestEndpointConcurrentRecordAndStats(t *testing.T) {

	assertions := assert.New(t)

	endpoint := &Endpoint{
		Name:    testConcurrentName,
		Subject: testConcurrentSubject,
	}

	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			var err error
			if idx%5 == 0 {
				err = &ServiceError{Description: testTestError}
			}
			endpoint.RecordRequest(time.Duration(idx)*time.Microsecond, err)
		}(i)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			endpoint.Stats()
		}()
	}

	wg.Wait()

	stats := endpoint.Stats()

	assertions.Equal(int64(100), stats.NumRequests, "All requests should be recorded")
	assertions.Equal(int64(20), stats.NumErrors, "Errors should be counted correctly")
}

/* ========================================================================================================
   ENDPOINT OPTIONS TESTS
   ======================================================================================================== */

func TestEndpointOptions(t *testing.T) {

	t.Run("WithEndpointSubject", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &endpointConfig{}
		WithEndpointSubject(testCustomSubject)(cfg)

		assertions.Equal(testCustomSubject, cfg.subject, "Subject should match")
	})

	t.Run("WithEndpointMetadata", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &endpointConfig{}
		meta := map[string]string{"key": testHeaderValue}
		WithEndpointMetadata(meta)(cfg)

		assertions.Equal(testHeaderValue, cfg.metadata["key"], "Metadata should match")
	})

	t.Run("WithEndpointQueueGroup", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &endpointConfig{}
		WithEndpointQueueGroup(testCustomQueueGroup)(cfg)

		assertions.Equal(testCustomQueueGroup, cfg.queueGroup, "QueueGroup should match")
	})

	t.Run("WithEndpointQueueGroupDisabled", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &endpointConfig{}
		WithEndpointQueueGroupDisabled()(cfg)

		assertions.True(cfg.queueGroupDisabled, "queueGroupDisabled should be true")
	})
}

/* ========================================================================================================
   GROUP OPTIONS TESTS
   ======================================================================================================== */

func TestGroupOptions(t *testing.T) {

	t.Run("WithGroupQueueGroup", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &groupConfig{}
		WithGroupQueueGroup(testGroupQueueGroup)(cfg)

		assertions.Equal(testGroupQueueGroup, cfg.queueGroup, "QueueGroup should match")
	})

	t.Run("WithGroupQueueGroupDisabled", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &groupConfig{}
		WithGroupQueueGroupDisabled()(cfg)

		assertions.True(cfg.queueGroupDisabled, "queueGroupDisabled should be true")
	})
}

/* ========================================================================================================
   REQUEST OPTIONS TESTS
   ======================================================================================================== */

func TestRequestOptions(t *testing.T) {

	t.Run("WithRequestTimeout", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &requestConfig{}
		WithRequestTimeout(5 * time.Second)(cfg)

		assertions.Equal(5*time.Second, cfg.timeout, "Timeout should match")
	})

	t.Run("WithMaxReplies", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &requestConfig{}
		WithMaxReplies(10)(cfg)

		assertions.Equal(10, cfg.maxReplies, "MaxReplies should match")
	})

	t.Run("WithRequestHeaders", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &requestConfig{}
		headers := make(nats.Header)
		headers.Set(testXCustomHeader, testHeaderValue)
		WithRequestHeaders(headers)(cfg)

		assertions.Equal(testHeaderValue, cfg.headers.Get(testXCustomHeader), "Header should match")
	})

	t.Run("WithExpectedRTT", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &requestConfig{}
		WithExpectedRTT(100 * time.Millisecond)(cfg)

		assertions.Equal(100*time.Millisecond, cfg.expectedRTT, "ExpectedRTT should match")
	})
}

/* ========================================================================================================
   RESPOND OPTIONS TESTS
   ======================================================================================================== */

func TestRespondOptions(t *testing.T) {

	t.Run("WithResponseHeaders", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &respondConfig{}
		headers := make(nats.Header)
		headers.Set(testContentTypeHeader, testApplicationJSON)
		WithResponseHeaders(headers)(cfg)

		assertions.Equal(testApplicationJSON, cfg.headers.Get(testContentTypeHeader), "Header must match")
	})

	t.Run("WithResponseHeader", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &respondConfig{}
		WithResponseHeader(testXRequestIDHeader, testRequestIDValue)(cfg)

		assertions.Equal(testRequestIDValue, cfg.headers.Get(testXRequestIDHeader), "Header should be matched")
	})

	t.Run("WithResponseHeaderMultiple", func(t *testing.T) {

		assertions := assert.New(t)

		cfg := &respondConfig{}
		WithResponseHeader(testXFirstHeader, "1")(cfg)
		WithResponseHeader(testXSecondHeader, "2")(cfg)

		assertions.Equal("1", cfg.headers.Get(testXFirstHeader), "First header should be set")
		assertions.Equal("2", cfg.headers.Get(testXSecondHeader), "Second header should be set")
	})
}

/* ========================================================================================================
   BENCHMARK TESTS
   ======================================================================================================== */

func BenchmarkControlSubject(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ControlSubject(PingVerb, "MyService", "abc123")
	}
}

func BenchmarkEndpointRecordRequest(b *testing.B) {

	endpoint := &Endpoint{
		Name:    testBenchName,
		Subject: testBenchSubject,
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		endpoint.RecordRequest(10*time.Millisecond, nil)
	}
}

func BenchmarkEndpointRecordRequestParallel(b *testing.B) {

	endpoint := &Endpoint{
		Name:    testBenchParallelName,
		Subject: testBenchParallelSubj,
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			endpoint.RecordRequest(10*time.Millisecond, nil)
		}
	})
}

func BenchmarkEndpointStats(b *testing.B) {

	endpoint := &Endpoint{
		Name:    testBenchName,
		Subject: testBenchSubject,
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)
	endpoint.RecordRequest(20*time.Millisecond, &ServiceError{Description: "err"})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		endpoint.Stats()
	}
}

func BenchmarkResponseJSON(b *testing.B) {

	resp := &Response{
		Data: []byte(testTestJSON),
	}

	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var data testData
		resp.JSON(&data)
	}
}

func BenchmarkVerbString(b *testing.B) {

	var s string

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s = PingVerb.String()
		s = StatsVerb.String()
		s = InfoVerb.String()
	}

	_ = s
}
