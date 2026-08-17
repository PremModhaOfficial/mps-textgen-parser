package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ---------------------------------------- ResolveProtocol Tests -------------------------------------------------- */

func TestResolveProtocol(t *testing.T) {
	testCases := []struct {
		name             string
		inputProtocol    string
		expectedProtocol ExporterProtocol
		expectError      bool
	}{
		{
			name:             "empty defaults to grpc",
			inputProtocol:    "",
			expectedProtocol: ProtocolGRPC,
			expectError:      false,
		},
		{
			name:             "grpc lowercase",
			inputProtocol:    "grpc",
			expectedProtocol: ProtocolGRPC,
			expectError:      false,
		},
		{
			name:             "grpc uppercase",
			inputProtocol:    "GRPC",
			expectedProtocol: ProtocolGRPC,
			expectError:      false,
		},
		{
			name:             "grpc mixed case",
			inputProtocol:    "GrPc",
			expectedProtocol: ProtocolGRPC,
			expectError:      false,
		},
		{
			name:             "http lowercase",
			inputProtocol:    "http",
			expectedProtocol: ProtocolHTTP,
			expectError:      false,
		},
		{
			name:             "http uppercase",
			inputProtocol:    "HTTP",
			expectedProtocol: ProtocolHTTP,
			expectError:      false,
		},
		{
			name:             "http/protobuf",
			inputProtocol:    "http/protobuf",
			expectedProtocol: ProtocolHTTP,
			expectError:      false,
		},
		{
			name:             "http/protobuf uppercase",
			inputProtocol:    "HTTP/PROTOBUF",
			expectedProtocol: ProtocolHTTP,
			expectError:      false,
		},
		{
			name:             "unsupported protocol",
			inputProtocol:    "websocket",
			expectedProtocol: "",
			expectError:      true,
		},
		{
			name:             "invalid protocol",
			inputProtocol:    "invalid",
			expectedProtocol: "",
			expectError:      true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assertions := assert.New(t)

			resolvedProtocol, err := ResolveProtocol(testCase.inputProtocol)

			if testCase.expectError {
				assertions.Error(err, "ResolveProtocol should return error for %q", testCase.inputProtocol)
				assertions.Contains(err.Error(), "unsupported OTEL protocol")
			} else {
				assertions.NoError(err, "ResolveProtocol should not return error for %q", testCase.inputProtocol)
				assertions.Equal(testCase.expectedProtocol, resolvedProtocol)
			}
		})
	}
}

/* ---------------------------------------- GRPCDialOptions Tests -------------------------------------------------- */

func TestGRPCDialOptions(t *testing.T) {
	t.Run("insecure true returns options", func(t *testing.T) {
		assertions := assert.New(t)

		options := GRPCDialOptions(true)
		assertions.NotNil(options, "GRPCDialOptions(true) should return options")
		assertions.Len(options, 1, "Should return exactly one dial option")
	})

	t.Run("insecure false returns nil", func(t *testing.T) {
		assertions := assert.New(t)

		options := GRPCDialOptions(false)
		assertions.Nil(options, "GRPCDialOptions(false) should return nil")
	})
}

func TestGRPCDialOption(t *testing.T) {
	t.Run("insecure true returns option", func(t *testing.T) {
		assertions := assert.New(t)

		option := GRPCDialOption(true)
		assertions.NotNil(option, "GRPCDialOption(true) should return an option")
	})

	t.Run("insecure false returns nil", func(t *testing.T) {
		assertions := assert.New(t)

		option := GRPCDialOption(false)
		assertions.Nil(option, "GRPCDialOption(false) should return nil")
	})
}

/* ---------------------------------------- AggregateErrors Tests -------------------------------------------------- */

func TestAggregateErrors(t *testing.T) {
	t.Run("empty slice returns nil", func(t *testing.T) {
		assertions := assert.New(t)

		result := AggregateErrors([]error{})
		assertions.Nil(result, "AggregateErrors with empty slice should return nil")
	})

	t.Run("nil slice returns nil", func(t *testing.T) {
		assertions := assert.New(t)

		result := AggregateErrors(nil)
		assertions.Nil(result, "AggregateErrors with nil slice should return nil")
	})

	t.Run("all nil errors returns nil", func(t *testing.T) {
		assertions := assert.New(t)

		result := AggregateErrors([]error{nil, nil, nil})
		assertions.Nil(result, "AggregateErrors with all nil errors should return nil")
	})

	t.Run("single error returns that error", func(t *testing.T) {
		assertions := assert.New(t)

		singleError := errors.New("single error")
		result := AggregateErrors([]error{singleError})
		assertions.Equal(singleError, result, "AggregateErrors with single error should return that error")
	})

	t.Run("single non-nil among nils returns that error", func(t *testing.T) {
		assertions := assert.New(t)

		singleError := errors.New("single error")
		result := AggregateErrors([]error{nil, singleError, nil})
		assertions.Equal(singleError, result, "Should return the single non-nil error")
	})

	t.Run("multiple errors aggregates them", func(t *testing.T) {
		assertions := assert.New(t)

		error1 := errors.New("error one")
		error2 := errors.New("error two")
		error3 := errors.New("error three")

		result := AggregateErrors([]error{error1, error2, error3})
		assertions.NotNil(result, "AggregateErrors should return a combined error")
		assertions.Contains(result.Error(), "OTEL shutdown encountered errors")
		assertions.Contains(result.Error(), "error one")
		assertions.Contains(result.Error(), "error two")
		assertions.Contains(result.Error(), "error three")
	})

	t.Run("multiple errors with nils aggregates non-nils", func(t *testing.T) {
		assertions := assert.New(t)

		error1 := errors.New("error one")
		error2 := errors.New("error two")

		result := AggregateErrors([]error{nil, error1, nil, error2, nil})
		assertions.NotNil(result, "AggregateErrors should return a combined error")
		assertions.Contains(result.Error(), "error one")
		assertions.Contains(result.Error(), "error two")
	})
}

/* ---------------------------------------- ShutdownCollector Tests -------------------------------------------------- */

func TestNewShutdownCollector(t *testing.T) {
	assertions := assert.New(t)

	collector := NewShutdownCollector()
	assertions.NotNil(collector, "NewShutdownCollector should return a collector")
	assertions.False(collector.HasErrors(), "New collector should have no errors")
	assertions.Equal(0, collector.Count(), "New collector should have zero count")
	assertions.Nil(collector.Error(), "New collector Error() should return nil")
}

func TestShutdownCollectorCollect(t *testing.T) {
	t.Run("collect nil error does nothing", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		collector.Collect(nil)

		assertions.False(collector.HasErrors())
		assertions.Equal(0, collector.Count())
	})

	t.Run("collect single error", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		testError := errors.New("test error")
		collector.Collect(testError)

		assertions.True(collector.HasErrors())
		assertions.Equal(1, collector.Count())
		assertions.Equal(testError, collector.Error())
	})

	t.Run("collect multiple errors", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		error1 := errors.New("error 1")
		error2 := errors.New("error 2")

		collector.Collect(error1)
		collector.Collect(error2)

		assertions.True(collector.HasErrors())
		assertions.Equal(2, collector.Count())
		assertions.Contains(collector.Error().Error(), "error 1")
		assertions.Contains(collector.Error().Error(), "error 2")
	})
}

func TestShutdownCollectorCollectFunc(t *testing.T) {
	t.Run("collect from nil func", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		collector.CollectFunc(nil)

		assertions.False(collector.HasErrors())
	})

	t.Run("collect from func returning nil", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		collector.CollectFunc(func() error {
			return nil
		})

		assertions.False(collector.HasErrors())
	})

	t.Run("collect from func returning error", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		testError := errors.New("shutdown error")
		collector.CollectFunc(func() error {
			return testError
		})

		assertions.True(collector.HasErrors())
		assertions.Equal(1, collector.Count())
		assertions.Equal(testError, collector.Error())
	})

	t.Run("collect from multiple funcs", func(t *testing.T) {
		assertions := assert.New(t)

		collector := NewShutdownCollector()
		collector.CollectFunc(func() error {
			return errors.New("error 1")
		})
		collector.CollectFunc(func() error {
			return nil
		})
		collector.CollectFunc(func() error {
			return errors.New("error 2")
		})

		assertions.True(collector.HasErrors())
		assertions.Equal(2, collector.Count())
	})
}

/* ---------------------------------------- NewOTELResource Tests -------------------------------------------------- */

func TestNewOTELResource(t *testing.T) {
	t.Run("creates resource with valid service info", func(t *testing.T) {
		assertions := assert.New(t)

		serviceInfo := ServiceInfo{
			Name:        "test-service",
			Version:     "1.0.0",
			Environment: "production",
		}

		resource, err := NewOTELResource(serviceInfo)
		assertions.NoError(err, "NewOTELResource should not return error")
		assertions.NotNil(resource, "Resource should not be nil")
	})

	t.Run("creates resource with empty values", func(t *testing.T) {
		assertions := assert.New(t)

		serviceInfo := ServiceInfo{
			Name:        "",
			Version:     "",
			Environment: "",
		}

		resource, err := NewOTELResource(serviceInfo)
		assertions.NoError(err, "NewOTELResource should handle empty values")
		assertions.NotNil(resource, "Resource should not be nil")
	})

	t.Run("creates resource with special characters", func(t *testing.T) {
		assertions := assert.New(t)

		serviceInfo := ServiceInfo{
			Name:        "test-service_v2",
			Version:     "1.0.0-beta+build.123",
			Environment: "staging-us-east-1",
		}

		resource, err := NewOTELResource(serviceInfo)
		assertions.NoError(err, "NewOTELResource should handle special characters")
		assertions.NotNil(resource, "Resource should not be nil")
	})
}

func TestNewOTELResourceFromConfig(t *testing.T) {
	t.Run("creates resource from config values", func(t *testing.T) {
		assertions := assert.New(t)

		resource, err := NewOTELResourceFromConfig("my-service", "2.0.0", "development")
		assertions.NoError(err, "NewOTELResourceFromConfig should not return error")
		assertions.NotNil(resource, "Resource should not be nil")
	})

	t.Run("creates resource from empty config values", func(t *testing.T) {
		assertions := assert.New(t)

		resource, err := NewOTELResourceFromConfig("", "", "")
		assertions.NoError(err, "NewOTELResourceFromConfig should handle empty values")
		assertions.NotNil(resource, "Resource should not be nil")
	})
}

/* ---------------------------------------- Benchmark Tests -------------------------------------------------- */

func BenchmarkResolveProtocolGRPC(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ResolveProtocol("grpc")
	}
}

func BenchmarkResolveProtocolHTTP(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ResolveProtocol("http")
	}
}

func BenchmarkResolveProtocolEmpty(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = ResolveProtocol("")
	}
}

func BenchmarkGRPCDialOptionsInsecure(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GRPCDialOptions(true)
	}
}

func BenchmarkGRPCDialOptionsSecure(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = GRPCDialOptions(false)
	}
}

func BenchmarkAggregateErrorsEmpty(b *testing.B) {
	emptyErrors := []error{}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AggregateErrors(emptyErrors)
	}
}

func BenchmarkAggregateErrorsSingle(b *testing.B) {
	singleError := []error{errors.New("test error")}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AggregateErrors(singleError)
	}
}

func BenchmarkAggregateErrorsMultiple(b *testing.B) {
	multipleErrors := []error{
		errors.New("error 1"),
		errors.New("error 2"),
		errors.New("error 3"),
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = AggregateErrors(multipleErrors)
	}
}

func BenchmarkShutdownCollectorCollect(b *testing.B) {
	testError := errors.New("test error")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		collector := NewShutdownCollector()
		collector.Collect(testError)
		_ = collector.Error()
	}
}

func BenchmarkNewOTELResource(b *testing.B) {
	serviceInfo := ServiceInfo{
		Name:        "test-service",
		Version:     "1.0.0",
		Environment: "production",
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewOTELResource(serviceInfo)
	}
}
