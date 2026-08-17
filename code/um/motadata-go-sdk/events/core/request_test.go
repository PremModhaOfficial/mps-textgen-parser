package core

import (
	"testing"
	"time"
)

func TestControlSubject(t *testing.T) {
	tests := []struct {
		name     string
		verb     Verb
		svcName  string
		svcID    string
		expected string
	}{
		{
			name:     "ping all services",
			verb:     PingVerb,
			svcName:  "",
			svcID:    "",
			expected: "$SRV.PING",
		},
		{
			name:     "ping by name",
			verb:     PingVerb,
			svcName:  "MyService",
			svcID:    "",
			expected: "$SRV.PING.MyService",
		},
		{
			name:     "ping by name and id",
			verb:     PingVerb,
			svcName:  "MyService",
			svcID:    "abc123",
			expected: "$SRV.PING.MyService.abc123",
		},
		{
			name:     "stats all services",
			verb:     StatsVerb,
			svcName:  "",
			svcID:    "",
			expected: "$SRV.STATS",
		},
		{
			name:     "info by name",
			verb:     InfoVerb,
			svcName:  "TestService",
			svcID:    "",
			expected: "$SRV.INFO.TestService",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ControlSubject(tt.verb, tt.svcName, tt.svcID)
			if result != tt.expected {
				t.Errorf("ControlSubject(%v, %q, %q) = %q, want %q",
					tt.verb, tt.svcName, tt.svcID, result, tt.expected)
			}
		})
	}
}

func TestVerbString(t *testing.T) {
	tests := []struct {
		verb     Verb
		expected string
	}{
		{PingVerb, "PING"},
		{StatsVerb, "STATS"},
		{InfoVerb, "INFO"},
		{Verb(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.verb.String(); got != tt.expected {
				t.Errorf("Verb.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestServiceError(t *testing.T) {
	err := &ServiceError{
		Code:        "400",
		Description: "bad request",
		Data:        []byte("details"),
	}

	if err.Error() != "bad request" {
		t.Errorf("ServiceError.Error() = %q, want %q", err.Error(), "bad request")
	}
}

func TestResponseIsError(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		resp := &Response{
			Data:    []byte("ok"),
			Headers: make(Headers),
		}
		if resp.IsError() {
			t.Error("IsError() should be false for successful response")
		}
	})

	t.Run("with error", func(t *testing.T) {
		resp := &Response{
			Error: &ServiceError{
				Code:        "500",
				Description: "internal error",
			},
		}
		if !resp.IsError() {
			t.Error("IsError() should be true for error response")
		}
	})
}

func TestResponseJSON(t *testing.T) {
	type testData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	t.Run("success", func(t *testing.T) {
		resp := &Response{
			Data: []byte(`{"name":"test","value":42}`),
		}

		var data testData
		if err := resp.JSON(&data); err != nil {
			t.Errorf("JSON() error = %v", err)
		}
		if data.Name != "test" || data.Value != 42 {
			t.Errorf("JSON() data = %+v, want {Name:test Value:42}", data)
		}
	})

	t.Run("with error", func(t *testing.T) {
		resp := &Response{
			Error: &ServiceError{
				Code:        "400",
				Description: "bad request",
			},
		}

		var data testData
		if err := resp.JSON(&data); err == nil {
			t.Error("JSON() should return error when response has error")
		}
	})
}

func TestEndpointStats(t *testing.T) {
	endpoint := &Endpoint{
		Name:       "test",
		Subject:    "test.subject",
		QueueGroup: "q",
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)
	endpoint.RecordRequest(20*time.Millisecond, nil)
	endpoint.RecordRequest(30*time.Millisecond, &ServiceError{Description: "test error"})

	stats := endpoint.Stats()

	if stats.NumRequests != 3 {
		t.Errorf("NumRequests = %d, want 3", stats.NumRequests)
	}
	if stats.NumErrors != 1 {
		t.Errorf("NumErrors = %d, want 1", stats.NumErrors)
	}
	if stats.TotalProcessingTime != 60*time.Millisecond {
		t.Errorf("TotalProcessingTime = %v, want 60ms", stats.TotalProcessingTime)
	}
	if stats.AverageProcessingTime != 20*time.Millisecond {
		t.Errorf("AverageProcessingTime = %v, want 20ms", stats.AverageProcessingTime)
	}
	if stats.LastError != "test error" {
		t.Errorf("LastError = %q, want %q", stats.LastError, "test error")
	}
}

func TestEndpointReset(t *testing.T) {
	endpoint := &Endpoint{
		Name:    "test",
		Subject: "test.subject",
	}

	endpoint.RecordRequest(10*time.Millisecond, nil)
	endpoint.RecordRequest(20*time.Millisecond, &ServiceError{Description: "error"})

	endpoint.Reset()
	stats := endpoint.Stats()

	if stats.NumRequests != 0 {
		t.Errorf("NumRequests after reset = %d, want 0", stats.NumRequests)
	}
	if stats.NumErrors != 0 {
		t.Errorf("NumErrors after reset = %d, want 0", stats.NumErrors)
	}
	if stats.TotalProcessingTime != 0 {
		t.Errorf("TotalProcessingTime after reset = %v, want 0", stats.TotalProcessingTime)
	}
	if stats.LastError != "" {
		t.Errorf("LastError after reset = %q, want empty", stats.LastError)
	}
}

func TestEndpointOptions(t *testing.T) {
	t.Run("WithEndpointSubject", func(t *testing.T) {
		cfg := &endpointConfig{}
		WithEndpointSubject("custom.subject")(cfg)
		if cfg.subject != "custom.subject" {
			t.Errorf("subject = %q, want %q", cfg.subject, "custom.subject")
		}
	})

	t.Run("WithEndpointMetadata", func(t *testing.T) {
		cfg := &endpointConfig{}
		meta := map[string]string{"key": "value"}
		WithEndpointMetadata(meta)(cfg)
		if cfg.metadata["key"] != "value" {
			t.Errorf("metadata[key] = %q, want %q", cfg.metadata["key"], "value")
		}
	})

	t.Run("WithEndpointQueueGroup", func(t *testing.T) {
		cfg := &endpointConfig{}
		WithEndpointQueueGroup("custom-q")(cfg)
		if cfg.queueGroup != "custom-q" {
			t.Errorf("queueGroup = %q, want %q", cfg.queueGroup, "custom-q")
		}
	})

	t.Run("WithEndpointQueueGroupDisabled", func(t *testing.T) {
		cfg := &endpointConfig{}
		WithEndpointQueueGroupDisabled()(cfg)
		if !cfg.queueGroupDisabled {
			t.Error("queueGroupDisabled should be true")
		}
	})
}

func TestGroupOptions(t *testing.T) {
	t.Run("WithGroupQueueGroup", func(t *testing.T) {
		cfg := &groupConfig{}
		WithGroupQueueGroup("group-q")(cfg)
		if cfg.queueGroup != "group-q" {
			t.Errorf("queueGroup = %q, want %q", cfg.queueGroup, "group-q")
		}
	})

	t.Run("WithGroupQueueGroupDisabled", func(t *testing.T) {
		cfg := &groupConfig{}
		WithGroupQueueGroupDisabled()(cfg)
		if !cfg.queueGroupDisabled {
			t.Error("queueGroupDisabled should be true")
		}
	})
}

func TestRequestOptions(t *testing.T) {
	t.Run("WithRequestTimeout", func(t *testing.T) {
		cfg := &requestConfig{}
		WithRequestTimeout(5 * time.Second)(cfg)
		if cfg.timeout != 5*time.Second {
			t.Errorf("timeout = %v, want 5s", cfg.timeout)
		}
	})

	t.Run("WithMaxReplies", func(t *testing.T) {
		cfg := &requestConfig{}
		WithMaxReplies(10)(cfg)
		if cfg.maxReplies != 10 {
			t.Errorf("maxReplies = %d, want 10", cfg.maxReplies)
		}
	})

	t.Run("WithRequestHeaders", func(t *testing.T) {
		cfg := &requestConfig{}
		headers := make(Headers)
		headers.Set("X-Custom", "value")
		WithRequestHeaders(headers)(cfg)
		if cfg.headers.Get("X-Custom") != "value" {
			t.Errorf("headers[X-Custom] = %q, want %q", cfg.headers.Get("X-Custom"), "value")
		}
	})
}

func TestRespondOptions(t *testing.T) {
	t.Run("WithResponseHeaders", func(t *testing.T) {
		cfg := &respondConfig{}
		headers := make(Headers)
		headers.Set("Content-Type", "application/json")
		WithResponseHeaders(headers)(cfg)
		if cfg.headers.Get("Content-Type") != "application/json" {
			t.Errorf("headers[Content-Type] = %q, want %q",
				cfg.headers.Get("Content-Type"), "application/json")
		}
	})

	t.Run("WithResponseHeader", func(t *testing.T) {
		cfg := &respondConfig{}
		WithResponseHeader("X-Request-Id", "123")(cfg)
		if cfg.headers.Get("X-Request-Id") != "123" {
			t.Errorf("headers[X-Request-Id] = %q, want %q",
				cfg.headers.Get("X-Request-Id"), "123")
		}
	})
}
