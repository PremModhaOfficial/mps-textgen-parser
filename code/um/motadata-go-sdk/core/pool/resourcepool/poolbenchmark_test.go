package resourcepool

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"
)

/* ========================================================================================================
   BENCHMARK TESTS
   ======================================================================================================== */

// BenchmarkResourcePoolGetPutOperation benchmarks Get/Put operations
func BenchmarkResourcePoolGetPutOperation(b *testing.B) {
	config := PoolConfig[*TestResource]{
		MaxSize: 100,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{}, nil
		},
	}

	pool, _ := NewResourcePool(config)
	defer pool.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resource, err := pool.Get(ctx)
			if err != nil {
				b.Fatal(err)
			}
			_ = pool.Put(resource)
		}
	})
}

// BenchmarkResourcePoolStats benchmarks Stats operations
func BenchmarkResourcePoolStats(b *testing.B) {
	config := PoolConfig[*TestResource]{
		MaxSize: 100,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{}, nil
		},
	}

	pool, _ := NewResourcePool(config)
	defer pool.Close()

	ctx := context.Background()

	// Pre-populate pool
	resources := make([]*TestResource, 500)
	for i := 0; i < len(resources); i++ {
		resources[i], _ = pool.Get(ctx)
	}
	for i := 0; i < len(resources)/2; i++ {
		_ = pool.Put(resources[i])
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.GetPoolStats()
		}
	})
}

// BenchmarkResourcePoolHighContention benchmarks under high contention
func BenchmarkResourcePoolHighContention(b *testing.B) {
	config := PoolConfig[*TestResource]{
		MaxSize: 10,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{}, nil
		},
		OnReset: func(r *TestResource) error {
			r.Used = false
			return nil
		},
	}

	pool, _ := NewResourcePool(config)
	defer pool.Close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resource, err := pool.Get(ctx)
			if err != nil {
				b.Fatal(err)
			}
			resource.Used = true
			// Simulate some work
			time.Sleep(time.Nanosecond)
			_ = pool.Put(resource)
		}
	})
}

func BenchmarkResourcePoolHTTPClientGetPutOperation(b *testing.B) {

	config := PoolConfig[*HTTPClientWrapper]{
		MaxSize: 10,
		OnCreate: func() (*HTTPClientWrapper, error) {
			return &HTTPClientWrapper{
				Client: &http.Client{},
			}, nil
		},
		OnReset: func(client *HTTPClientWrapper) error {
			return nil
		},
	}

	pool, err := NewResourcePool(config)

	if err != nil {
		b.Fatalf("Failed to create pool: %v", err)
	}

	defer pool.Close()

	ctx := context.Background()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

			var client *HTTPClientWrapper

			client, err = pool.Get(ctx)

			if err != nil {
				b.Fatalf("Get failed: %v", err)
			}

			err = pool.Put(client)

			if err != nil {
				b.Fatalf("Put failed: %v", err)
			}
		}
	})
}

func BenchmarkResourcePoolByteBufferGetPutOperation(b *testing.B) {

	config := PoolConfig[*ByteBufferWrapper]{
		MaxSize: 20,
		OnCreate: func() (*ByteBufferWrapper, error) {
			return &ByteBufferWrapper{
				Buffer: &bytes.Buffer{},
			}, nil
		},
		OnReset: func(wrapper *ByteBufferWrapper) error {
			wrapper.Buffer.Reset()
			return nil
		},
	}

	pool, err := NewResourcePool(config)

	if err != nil {
		b.Fatalf("Failed to create pool: %v", err)
	}

	defer pool.Close()

	ctx := context.Background()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {

		for pb.Next() {

			var wrapper *ByteBufferWrapper

			wrapper, err = pool.Get(ctx)

			if err != nil {
				b.Fatalf("Get failed: %v", err)
			}

			// Simulate buffer usage
			wrapper.Buffer.WriteString("benchmark test data")

			err = pool.Put(wrapper)

			if err != nil {
				b.Fatalf("Put failed: %v", err)
			}
		}
	})
}
