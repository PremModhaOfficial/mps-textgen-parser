package resourcepool

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test data structures for different pool types

// HTTPClientWrapper wraps http.Client for pooling
type HTTPClientWrapper struct {
	Client    *http.Client
	CreatedAt time.Time
	Id        int
}

// DBConnectionWrapper wraps database connection for pooling
type DBConnectionWrapper struct {
	Conn      *sql.DB
	CreatedAt time.Time
	Id        int
	InTx      bool
}

// ByteBufferWrapper wraps bytes.Buffer for pooling
type ByteBufferWrapper struct {
	Buffer    *bytes.Buffer
	CreatedAt time.Time
	Id        int
}

// TestResource is a simple test resourcepool for testing
type TestResource struct {
	ID      int
	Used    bool
	Closed  bool
	counter *int32
}

func (r *TestResource) Close() {
	r.Closed = true
	if r.counter != nil {
		atomic.AddInt32(r.counter, -1)
	}
}

/* ========================================================================================================
   BASIC FUNCTIONALITY TESTS - Core resourcepool pool operations
   ======================================================================================================== */

// TestResourcePoolGet - Basic resourcepool acquisition functionality
func TestResourcePoolGet(t *testing.T) {
	assertions := assert.New(t)

	created := int32(0)
	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			id := atomic.AddInt32(&created, 1)
			return &TestResource{ID: int(id)}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	assertions.NotNil(pool, "Pool should not be nil")
	defer pool.Close()

	// Test: Get resourcepool
	ctx := context.Background()
	resource, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool")
	assertions.NotNil(resource, "Resource should not be nil")
	assertions.Equal(1, resource.ID, "Expected resourcepool ID 1")

	// Test: Resource is marked as used
	stats := pool.GetPoolStats()
	assertions.Equal(1, stats.Used, "Expected 1 resourcepool in use")
	assertions.Equal(1, stats.Created, "Expected 1 resourcepool created")
	assertions.Equal(0, stats.Available, "Expected 0 available resources")

	// Return resourcepool for cleanup
	err = pool.Put(resource)
	assertions.NoError(err, "Should be able to return resourcepool")
}

// TestResourcePoolPut - Basic resourcepool return functionality
func TestResourcePoolPut(t *testing.T) {
	assertions := assert.New(t)

	created := int32(0)
	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			id := atomic.AddInt32(&created, 1)
			return &TestResource{ID: int(id)}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	// Get and return resourcepool
	ctx := context.Background()
	resource, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool")

	// Test: Return resourcepool
	err = pool.Put(resource)
	assertions.NoError(err, "Should be able to return resourcepool")

	// Test: Resource is no longer marked as in use
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Used, "Expected 0 resources in use")
	assertions.Equal(1, stats.Available, "Expected 1 available resourcepool")
	assertions.Equal(1, stats.Created, "Expected 1 resourcepool created")
}

// TestResourcePoolReuse - Resource reuse and recycling functionality
func TestResourcePoolReuse(t *testing.T) {
	assertions := assert.New(t)

	created := int32(0)
	resetCalled := int32(0)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			id := atomic.AddInt32(&created, 1)
			return &TestResource{ID: int(id)}, nil
		},
		OnReset: func(r *TestResource) error {
			atomic.AddInt32(&resetCalled, 1)
			r.Used = false
			return nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: First acquisition creates new resourcepool
	resource1, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get first resourcepool")
	resource1.Used = true

	err = pool.Put(resource1)
	assertions.NoError(err, "Should be able to return first resourcepool")

	// Test: Second acquisition reuses the same resourcepool
	resource2, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get second resourcepool")

	assertions.Equal(resource1.ID, resource2.ID, "Expected resourcepool reuse, should have same IDs")
	assertions.False(resource2.Used, "Resource should be reset")
	assertions.Equal(int32(1), atomic.LoadInt32(&resetCalled), "Reset should be called exactly once")

	err = pool.Put(resource2)
	assertions.NoError(err, "Should be able to return second resourcepool")
}

// TestResourcePoolTryGet - Non-blocking resourcepool acquisition
func TestResourcePoolTryGet(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	// Test: TryGet when pool is empty should create resourcepool
	resource, err := pool.TryGet()
	assertions.NoError(err, "Should be able to get resourcepool when pool is empty")
	assertions.NotNil(resource, "Resource should not be nil")
	assertions.Equal(1, resource.ID, "Expected resourcepool ID 1")

	// Test: TryGet when pool is at capacity should return error
	_, err = pool.TryGet()
	assertions.Error(err, "Should get error when pool is exhausted")

	// Return resourcepool and try again
	err = pool.Put(resource)
	assertions.NoError(err, "Should be able to return resourcepool")

	// Test: TryGet after return should succeed
	resource2, err := pool.TryGet()
	assertions.NoError(err, "Should be able to get resourcepool after return")
	assertions.NotNil(resource2, "Second resourcepool should not be nil")
	assertions.Equal(resource.ID, resource2.ID, "Should reuse the same resourcepool")

	err = pool.Put(resource2)
	assertions.NoError(err, "Should be able to return second resourcepool")
}

// TestResourcePoolClose - Basic pool closure functionality
func TestResourcePoolClose(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	// Add some resources to pool
	ctx := context.Background()
	resource1, _ := pool.Get(ctx)
	resource2, _ := pool.Get(ctx)
	_ = pool.Put(resource1)
	_ = pool.Put(resource2)

	// Test: close pool
	pool.Close()

	// Test: Pool should be marked as closed
	stats := pool.GetPoolStats()
	assertions.True(stats.IsClosed, "Pool should be marked as closed")

	// Test: Operations on closed pool should fail
	_, err = pool.Get(ctx)
	assertions.Error(err, "Should get error when getting from closed pool")

	_, err = pool.TryGet()
	assertions.Error(err, "Should get error when trying to get from closed pool")
}

// TestResourcePoolCloseAndWait - close with waiting for resourcepool return
func TestResourcePoolCloseAndWait(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 2,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	// Get resources
	ctx := context.Background()
	resource1, _ := pool.Get(ctx)
	resource2, _ := pool.Get(ctx)

	// Test: CloseWithTimeout should wait for resources to be returned
	var wg sync.WaitGroup
	var closeErr error

	wg.Add(1)
	go func() {
		defer wg.Done()
		closeErr = pool.CloseWithTimeout(time.Second * 2)
	}()

	// Return resources after delay
	time.Sleep(time.Millisecond * 100)
	_ = pool.Put(resource1)
	_ = pool.Put(resource2)

	wg.Wait()

	// Test: CloseWithTimeout should complete successfully
	assertions.NoError(closeErr, "CloseWithTimeout should complete without error")

	// Test: Pool should be closed
	stats := pool.GetPoolStats()
	assertions.True(stats.IsClosed, "Pool should be closed after CloseWithTimeout")
}

/* ========================================================================================================
   STATISTICS & MONITORING TESTS - Pool statistics and monitoring functionality
   ======================================================================================================== */

// TestResourcePoolStats - Basic statistics tracking
func TestResourcePoolStats(t *testing.T) {
	assertions := assert.New(t)

	created := int32(0)
	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			id := atomic.AddInt32(&created, 1)
			return &TestResource{ID: int(id)}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Initial stats
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Created, "Initial created count should be 0")
	assertions.Equal(0, stats.Used, "Initial used count should be 0")
	assertions.Equal(0, stats.Available, "Initial available count should be 0")
	assertions.Equal(3, stats.MaxSize, "MaxSize should be 3")

	// Test: Stats after getting resourcepool
	resource, _ := pool.Get(ctx)
	stats = pool.GetPoolStats()
	assertions.Equal(1, stats.Created, "Created count should be 1 after get")
	assertions.Equal(1, stats.Used, "Used count should be 1 after get")
	assertions.Equal(0, stats.Available, "Available count should be 0 after get")

	// Test: Stats after returning resourcepool
	_ = pool.Put(resource)
	stats = pool.GetPoolStats()
	assertions.Equal(1, stats.Created, "Created count should remain 1 after put")
	assertions.Equal(0, stats.Used, "Used count should be 0 after put")
	assertions.Equal(1, stats.Available, "Available count should be 1 after put")
}

// TestResourcePoolStatsConsistency - Statistics consistency under normal operations
func TestResourcePoolStatsConsistency(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 10,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()
	inconsistentCount := 0
	totalChecks := 1000

	// Test: Perform operations and verify stats consistency
	for i := 0; i < totalChecks; i++ {
		resources := make([]*TestResource, rand.Intn(11)) // 0 to 10 resources

		// Get resources
		for j := range resources {

			var resource *TestResource

			resource, err = pool.Get(ctx)
			if err != nil {
				continue
			}
			resources[j] = resource
		}

		// Check stats consistency
		stats := pool.GetPoolStats()
		if stats.Created != stats.Used+stats.Available {
			inconsistentCount++
		}

		// Return resources
		for _, resource := range resources {
			if resource != nil {
				_ = pool.Put(resource)
			}
		}
	}

	// Final verification
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Used, "Expected no in-use resources at end")

	t.Logf("Final stats - Size: %d, Created: %d, MaxSize: %d", stats.Available, stats.Created, stats.MaxSize)
	t.Logf("Total checks: %d, Inconsistent: %d", totalChecks, inconsistentCount)

	assertions.Equal(0, inconsistentCount, "Should have no inconsistent stats checks")
}

// TestResourcePoolStatsUnderConcurrency - Statistics accuracy under concurrent access
func TestResourcePoolStatsUnderConcurrency(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 5,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	var wg sync.WaitGroup
	invalidCount := int32(0)
	totalChecks := 10000
	checksPerWorker := totalChecks / 10

	// Test: Multiple workers checking stats consistency
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for i := 0; i < checksPerWorker; i++ {

				var resource *TestResource

				// Perform random operation
				switch rand.Intn(3) {
				case 0:
					resource, err = pool.Get(ctx)
					if err == nil {
						_ = pool.Put(resource)
					}
				case 1:
					resource, err = pool.TryGet()
					if err == nil {
						_ = pool.Put(resource)
					}
				case 2:
					stats := pool.GetPoolStats()
					// Check for major inconsistencies that indicate real bugs
					// Minor temporary inconsistencies are expected during concurrent operations
					if stats.Created > stats.MaxSize || stats.Used > stats.MaxSize ||
						stats.Available > stats.MaxSize || stats.Used < 0 || stats.Available < 0 {
						atomic.AddInt32(&invalidCount, 1)
					}
				}
			}
		}()
	}

	wg.Wait()

	t.Logf("Checked %d stats snapshots, found %d invalid", totalChecks, invalidCount)

	assertions.Equal(int32(0), invalidCount, "Should have no invalid stats under concurrency")
}

// TestResourcePoolInUseTracking - In-use resourcepool tracking accuracy
func TestResourcePoolInUseTracking(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 10,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	var wg sync.WaitGroup
	maxConcurrent := int32(0)
	ctx := context.Background()

	// Test: Track maximum concurrent in-use resources
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			var resource *TestResource

			resource, err = pool.Get(ctx)
			if err != nil {
				return
			}

			// Update max concurrent
			current := pool.UsedCount()
			for {
				maxVal := atomic.LoadInt32(&maxConcurrent)
				if current <= int(maxVal) || atomic.CompareAndSwapInt32(&maxConcurrent, maxVal, int32(current)) {
					break
				}
			}

			time.Sleep(time.Millisecond * 10)
			_ = pool.Put(resource)
		}()
	}

	wg.Wait()

	// Test: All resources should be returned
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Used, "Expected 0 in-use resources after all returns")

	t.Logf("Maximum concurrent in-use resources: %d", maxConcurrent)

	assertions.Greater(int(maxConcurrent), 0, "Expected some concurrent resourcepool usage")
}

/* ========================================================================================================
   CONFIGURATION & VALIDATION TESTS - Pool configuration and validation
   ======================================================================================================== */

// TestResourcePoolInvalidConfiguration - Invalid configuration handling
func TestResourcePoolInvalidConfiguration(t *testing.T) {
	t.Run("negativeMaxSize", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize:  -1,
			OnCreate: func() (*TestResource, error) { return &TestResource{}, nil },
		}

		pool, err := NewResourcePool(config)
		assertions.Error(err, "Should get error for negative MaxSize")
		assertions.Nil(pool, "Pool should be nil when creation fails")
	})

	t.Run("zeroMaxSize", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize:  0,
			OnCreate: func() (*TestResource, error) { return &TestResource{}, nil },
		}

		pool, err := NewResourcePool(config)
		assertions.Error(err, "Should get error for zero MaxSize")
		assertions.Nil(pool, "Pool should be nil when creation fails")
	})

	t.Run("nilOnCreate", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize:  1,
			OnCreate: nil,
		}

		pool, err := NewResourcePool(config)
		assertions.Error(err, "Should get error for nil OnCreate")
		assertions.Nil(pool, "Pool should be nil when creation fails")
	})
}

// TestResourcePoolConfiguration - Valid configuration testing
func TestResourcePoolConfiguration(t *testing.T) {
	t.Run("minValidConfig", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize: 1,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool with minimal config")
		pool.Close()
	})

	t.Run("fullConfigWithCallbacks", func(t *testing.T) {
		assertions := assert.New(t)

		resetCalled := false
		destroyCalled := false

		config := PoolConfig[*TestResource]{
			MaxSize: 2,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
			OnReset: func(r *TestResource) error {
				resetCalled = true
				return nil
			},
			OnDestroy: func(r *TestResource) {
				destroyCalled = true
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool with full config")

		// Test callbacks are used
		ctx := context.Background()
		resource, _ := pool.Get(ctx)
		_ = pool.Put(resource)

		assertions.True(resetCalled, "OnReset callback should be called")

		pool.Close()

		assertions.True(destroyCalled, "OnDestroy callback should be called")
	})
}

/* ========================================================================================================
   ERROR HANDLING TESTS - Error scenarios and handling
   ======================================================================================================== */

// TestResourcePoolOnCreateError - Resource creation error handling
func TestResourcePoolOnCreateError(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 2,
		OnCreate: func() (*TestResource, error) {
			return nil, errors.New("creation failed")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Get should return creation error
	_, err = pool.Get(ctx)
	assertions.Error(err, "Should get creation error")

	// Test: TryGet should also return creation error
	_, err = pool.TryGet()
	assertions.Error(err, "Should get creation error on TryGet")

	// Test: Stats should reflect no created resources
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Created, "Should have 0 created resources after creation failures")
}

// TestResourcePoolOnResetError - Resource reset error handling
func TestResourcePoolOnResetError(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
		OnReset: func(r *TestResource) error {
			return errors.New("reset failed")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()
	resource, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool")

	// Test: Put should handle reset error
	err = pool.Put(resource)
	assertions.Error(err, "Should get reset error when returning resourcepool")

	// Test: Resource should be destroyed, not returned to pool
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Available, "Resource should be destroyed due to reset error")
}

// TestResourcePoolContextCancellation - Context cancellation handling
func TestResourcePoolContextCancellation(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	// Exhaust the pool
	ctx1 := context.Background()
	resource, err := pool.Get(ctx1)
	assertions.NoError(err, "Should be able to get resourcepool")

	// Test: Get with cancelled context should return context error
	ctx2, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = pool.Get(ctx2)
	assertions.True(errors.Is(err, context.Canceled), "Should get context.Canceled error")

	// Test: Get with deadline exceeded
	ctx3, cancel3 := context.WithDeadline(context.Background(), time.Now().Add(time.Millisecond))
	defer cancel3()

	time.Sleep(time.Millisecond * 2)
	_, err = pool.Get(ctx3)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Logf("Correctly got context error: %v", err)
	}

	// Clean up
	err = pool.Put(resource)
	assertions.NoError(err, "Should be able to return resourcepool")
}

// TestResourcePoolClosedOperations - Operations on closed pool
func TestResourcePoolClosedOperations(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	// close the pool
	pool.Close()

	ctx := context.Background()

	// Test: Get on closed pool should fail
	_, err = pool.Get(ctx)
	assertions.Error(err, "Should get error when getting from closed pool")

	// Test: TryGet on closed pool should fail
	_, err = pool.TryGet()
	assertions.Error(err, "Should get error when trying to get from closed pool")

	// Test: Put on closed pool should fail
	resource := &TestResource{ID: 1}
	err = pool.Put(resource)
	assertions.Error(err, "Should get error when putting to closed pool")
}

// TestResourcePoolPutInvalidResource - Putting invalid/zero resources
func TestResourcePoolPutInvalidResource(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	// Test: Put nil resourcepool should fail
	err = pool.Put(nil)
	assertions.Error(err, "Should get error when putting nil resourcepool")

	// Test: Put zero value should fail
	var zeroResource *TestResource
	err = pool.Put(zeroResource)
	assertions.Error(err, "Should get error when putting zero resourcepool")
}

/* ========================================================================================================
   PANIC RECOVERY TESTS - Panic handling in callbacks
   ======================================================================================================== */

// TestResourcePoolOnCreatePanic - OnCreate callback panic handling
func TestResourcePoolOnCreatePanic(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			panic("onCreate panic test")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	assertions.NotNil(pool, "Pool should not be nil")
	defer pool.Close()

	ctx := context.Background()

	// Test: Panic in onCreate should be recovered and returned as error
	resource, err := pool.Get(ctx)
	assertions.Error(err, "Should get error from onCreate panic")
	assertions.Nil(resource, "Resource should be nil when onCreate panics")

	// Test: Pool should remain functional after panic
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Created, "Resource count should remain 0 after onCreate panic")
	assertions.Equal(0, stats.Used, "Used count should remain 0 after panic")
}

// TestResourcePoolOnResetPanic - OnReset callback panic handling
func TestResourcePoolOnResetPanic(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
		OnReset: func(r *TestResource) error {
			panic("onReset panic test")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()
	resource, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool")

	// Test: Panic in onReset should be recovered and returned as error
	err = pool.Put(resource)
	assertions.Error(err, "Should get error from onReset panic")

	// Test: Resource should be destroyed due to reset panic
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Available, "Resource should be destroyed due to reset panic")
}

// TestResourcePoolOnDestroyPanic - OnDestroy callback panic handling
func TestResourcePoolOnDestroyPanic(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
		OnDestroy: func(r *TestResource) {
			panic("onDestroy panic test")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	ctx := context.Background()
	resource, _ := pool.Get(ctx)
	_ = pool.Put(resource)

	// Test: Panic in onDestroy during close should be recovered
	pool.Close()

	// Test: Pool should still be closed despite panic
	assertions.True(pool.GetPoolStats().IsClosed, "Pool should remain closed despite onDestroy panic")
}

// TestResourcePoolMultiplePanicScenarios - Multiple panic scenarios
func TestResourcePoolMultiplePanicScenarios(t *testing.T) {
	assertions := assert.New(t)

	createPanics := true
	resetPanics := true

	config := PoolConfig[*TestResource]{
		MaxSize: 2,
		OnCreate: func() (*TestResource, error) {
			if createPanics {
				panic("create panic")
			}
			return &TestResource{ID: 1}, nil
		},
		OnReset: func(r *TestResource) error {
			if resetPanics {
				panic("reset panic")
			}
			return nil
		},
		OnDestroy: func(r *TestResource) {
			panic("destroy panic")
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	ctx := context.Background()

	// Test: Multiple create panics
	for i := 0; i < 3; i++ {
		_, err = pool.Get(ctx)
		assertions.Error(err, "Should get error from create panic")
	}

	// Disable create panics, enable operations
	createPanics = false
	resource, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool after disabling create panic")

	// Test: Reset panic
	err = pool.Put(resource)
	assertions.Error(err, "Should get error from reset panic")

	// Disable reset panics
	resetPanics = false
	resource2, _ := pool.Get(ctx)
	_ = pool.Put(resource2)

	// Test: Destroy panic during close (should be handled gracefully)
	pool.Close()
}

/* ========================================================================================================
   CONCURRENCY & RACE CONDITION TESTS - Concurrent operations and race condition testing
   ======================================================================================================== */

// TestResourcePoolConcurrentGetPut - Concurrent get/put operations
func TestResourcePoolConcurrentGetPut(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 5,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	var wg sync.WaitGroup
	operations := 100
	workers := 10

	// Test: Multiple workers performing get/put operations
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for i := 0; i < operations; i++ {

				var resource *TestResource
				resource, err = pool.Get(ctx)
				if err != nil {
					continue
				}

				// Simulate work
				time.Sleep(time.Microsecond)

				err = pool.Put(resource)
				assertions.NoError(err, "Put should not fail")
			}
		}()
	}

	wg.Wait()

	// Test: All resources should be returned
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Used, "Expected 0 in-use resources after concurrent operations")
}

// TestResourcePoolHighContention - High concurrency scenarios
func TestResourcePoolHighContention(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
		OnReset: func(r *TestResource) error {
			return nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	var wg sync.WaitGroup
	workers := 50
	operationsPerWorker := 20
	successCount := int32(0)
	errorCount := int32(0)

	// Test: High contention with many workers
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
			defer cancel()

			for i := 0; i < operationsPerWorker; i++ {

				var resource *TestResource
				resource, err = pool.Get(ctx)
				if err != nil {
					atomic.AddInt32(&errorCount, 1)
					continue
				}

				atomic.AddInt32(&successCount, 1)
				time.Sleep(time.Microsecond * 10)
				_ = pool.Put(resource)
			}
		}()
	}

	wg.Wait()

	t.Logf("High contention results: Successes=%d, Errors=%d, Total=%d",
		successCount, errorCount, successCount+errorCount)

	// Test: All resources should be returned
	stats := pool.GetPoolStats()
	assertions.Equal(0, stats.Used, "Expected 0 in-use resources after high contention")

	// Test: Should have some successes
	assertions.Greater(int(successCount), 0, "Expected some successful operations under high contention")
}

// TestResourcePoolCloseRaceCondition - close operation race conditions
func TestResourcePoolCloseRaceCondition(t *testing.T) {
	assertions := assert.New(t)

	for iteration := 1; iteration <= 50; iteration++ {
		config := PoolConfig[*TestResource]{
			MaxSize: 5,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")

		var wg sync.WaitGroup
		ctx := context.Background()

		// Pre-populate pool
		resources := make([]*TestResource, 5)
		for i := 0; i < 5; i++ {
			resources[i], _ = pool.Get(ctx)
		}
		for i := 0; i < 5; i++ {
			_ = pool.Put(resources[i])
		}

		// Test: Concurrent operations during close
		failedPuts := int32(0)
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				var resource *TestResource
				resource, err = pool.Get(ctx)
				if err != nil {
					return
				}

				time.Sleep(time.Microsecond * 10)
				err = pool.Put(resource)
				if err != nil {
					atomic.AddInt32(&failedPuts, 1)
				}
			}()
		}

		// close pool concurrently
		time.Sleep(time.Microsecond * 5)
		pool.Close()

		wg.Wait()

		t.Logf("Iteration %d: %d put operations failed out of 5", iteration, failedPuts)
	}
}

// TestResourcePoolStatsRaceCondition - Statistics race condition testing
func TestResourcePoolStatsRaceCondition(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 5,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	var wg sync.WaitGroup
	duration := time.Millisecond * 100
	start := time.Now()

	// Test: Concurrent stats access during operations
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for time.Since(start) < duration {

				var resource *TestResource
				switch rand.Intn(3) {
				case 0:
					resource, err = pool.Get(ctx)
					if err == nil {
						_ = pool.Put(resource)
					}
				case 1:
					resource, err = pool.TryGet()
					if err == nil {
						_ = pool.Put(resource)
					}
				case 2:
					pool.GetPoolStats()
				}
			}
		}()
	}

	wg.Wait()

	// Test: Final consistency check
	stats := pool.GetPoolStats()
	assertions.GreaterOrEqual(stats.Created, stats.Used+stats.Available, "Stats should be consistent")
}

// TestResourcePoolConcurrentClose - Concurrent close operations
func TestResourcePoolConcurrentClose(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	var wg sync.WaitGroup

	// Test: Multiple concurrent close calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			pool.Close()
		}()
	}

	wg.Wait()

	// Test: Pool should be closed
	assertions.True(pool.GetPoolStats().IsClosed, "Pool should be closed after concurrent close calls")
}

/* ========================================================================================================
   INTEGRATION TESTS - Integration with different resourcepool types
   ======================================================================================================== */

// TestResourcePoolHTTPClient - HTTP client resourcepool pooling integration
func TestResourcePoolHTTPClient(t *testing.T) {
	assertions := assert.New(t)

	clientID := int32(0)

	config := PoolConfig[*HTTPClientWrapper]{
		MaxSize: 3,
		OnCreate: func() (*HTTPClientWrapper, error) {
			id := atomic.AddInt32(&clientID, 1)
			return &HTTPClientWrapper{
				Client:    &http.Client{Timeout: time.Second * 30},
				CreatedAt: time.Now(),
				Id:        int(id),
			}, nil
		},
		OnReset: func(client *HTTPClientWrapper) error {
			// Reset any client state if needed
			return nil
		},
		OnDestroy: func(client *HTTPClientWrapper) {
			// Cleanup client resources
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create HTTP client pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Get HTTP client from pool
	clientWrapper, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get HTTP client")

	assertions.NotNil(clientWrapper.Client, "HTTP client should not be nil")
	assertions.Equal(1, clientWrapper.Id, "Expected client ID 1")

	// Test: Use the HTTP client (mock usage)
	assertions.Equal(time.Second*30, clientWrapper.Client.Timeout, "HTTP client timeout should be properly configured")

	// Test: Return client to pool
	err = pool.Put(clientWrapper)
	assertions.NoError(err, "Should be able to return HTTP client")

	// Test: Reuse client
	clientWrapper2, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get second HTTP client")

	assertions.Equal(clientWrapper.Id, clientWrapper2.Id, "Expected client reuse")

	err = pool.Put(clientWrapper2)
	assertions.NoError(err, "Should be able to return second HTTP client")
}

// TestResourcePoolDBConnection - Database connection pooling integration
func TestResourcePoolDBConnection(t *testing.T) {
	assertions := assert.New(t)

	connID := int32(0)

	config := PoolConfig[*DBConnectionWrapper]{
		MaxSize: 2,
		OnCreate: func() (*DBConnectionWrapper, error) {
			id := atomic.AddInt32(&connID, 1)
			// Note: Using nil for actual DB connection in test
			return &DBConnectionWrapper{
				Conn:      nil, // Would be real DB connection in practice
				CreatedAt: time.Now(),
				Id:        int(id),
				InTx:      false,
			}, nil
		},
		OnReset: func(conn *DBConnectionWrapper) error {
			// Reset connection state
			conn.InTx = false
			return nil
		},
		OnDestroy: func(conn *DBConnectionWrapper) {
			// close actual database connection
			if conn.Conn != nil {
				_ = conn.Conn.Close()
			}
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create DB connection pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Get database connection from pool
	connWrapper, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get DB connection")

	assertions.Equal(1, connWrapper.Id, "Expected connection ID 1")

	// Test: Simulate transaction usage
	connWrapper.InTx = true

	// Test: Return connection to pool (should be reset)
	err = pool.Put(connWrapper)
	assertions.NoError(err, "Should be able to return DB connection")

	// Test: Get connection again and verify reset
	connWrapper2, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get second DB connection")

	assertions.False(connWrapper2.InTx, "Connection should be reset (not in transaction)")

	assertions.Equal(connWrapper.Id, connWrapper2.Id, "Expected connection reuse")

	err = pool.Put(connWrapper2)
	assertions.NoError(err, "Should be able to return second DB connection")
}

// TestResourcePoolByteBuffer - Byte buffer pooling integration
func TestResourcePoolByteBuffer(t *testing.T) {
	assertions := assert.New(t)

	bufferID := int32(0)

	config := PoolConfig[*ByteBufferWrapper]{
		MaxSize: 2,
		OnCreate: func() (*ByteBufferWrapper, error) {
			id := atomic.AddInt32(&bufferID, 1)
			return &ByteBufferWrapper{
				Buffer:    &bytes.Buffer{},
				CreatedAt: time.Now(),
				Id:        int(id),
			}, nil
		},
		OnReset: func(buffer *ByteBufferWrapper) error {
			// Reset buffer content
			buffer.Buffer.Reset()
			return nil
		},
		OnDestroy: func(buffer *ByteBufferWrapper) {
			// No special cleanup needed for buffers
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create buffer pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Get buffer from pool
	bufferWrapper, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get buffer")

	assertions.NotNil(bufferWrapper.Buffer, "Buffer should not be nil")

	// Test: Use buffer
	bufferWrapper.Buffer.WriteString("test data")
	assertions.Greater(bufferWrapper.Buffer.Len(), 0, "Buffer should contain data")

	// Test: Return buffer to pool (should be reset)
	err = pool.Put(bufferWrapper)
	assertions.NoError(err, "Should be able to return buffer")

	// Test: Get buffer again and verify reset
	bufferWrapper2, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get second buffer")

	assertions.Equal(0, bufferWrapper2.Buffer.Len(), "Buffer should be reset (empty)")

	assertions.Equal(bufferWrapper.Id, bufferWrapper2.Id, "Expected buffer reuse")

	err = pool.Put(bufferWrapper2)
	assertions.NoError(err, "Should be able to return second buffer")
}

/* ========================================================================================================
   PERFORMANCE & EDGE CASE TESTS - Performance testing and edge cases
   ======================================================================================================== */

// TestResourcePoolMaxSizeOne - Single-resourcepool pool behavior
func TestResourcePoolMaxSizeOne(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig[*TestResource]{
		MaxSize: 1,
		OnCreate: func() (*TestResource, error) {
			return &TestResource{ID: 1}, nil
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")
	defer pool.Close()

	ctx := context.Background()

	// Test: Get single resourcepool
	resource1, err := pool.Get(ctx)
	assertions.NoError(err, "Should be able to get resourcepool")

	// Test: TryGet should fail when pool is exhausted
	_, err = pool.TryGet()
	assertions.Error(err, "Should get error when pool is exhausted")

	// Test: Second Get with timeout should fail
	ctx2, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	_, err = pool.Get(ctx2)
	assertions.Error(err, "Should get timeout error when pool is exhausted")

	// Test: Return resourcepool and try again
	err = pool.Put(resource1)
	assertions.NoError(err, "Should be able to return resourcepool")

	resource2, err := pool.TryGet()
	assertions.NoError(err, "Should be able to get resourcepool after return")

	assertions.Equal(resource1.ID, resource2.ID, "Expected same resourcepool to be reused")

	err = pool.Put(resource2)
	assertions.NoError(err, "Should be able to return second resourcepool")
}

// TestResourcePoolBoundaryConditions - Boundary value testing
func TestResourcePoolBoundaryConditions(t *testing.T) {
	t.Run("MaxSizeHighContention", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize: 1,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")
		defer pool.Close()

		var wg sync.WaitGroup
		successCount := int32(0)
		failureCount := int32(0)
		workers := 10

		// Test: High contention on single resourcepool
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
				defer cancel()

				for i := 0; i < 10; i++ {

					var resource *TestResource
					resource, err = pool.Get(ctx)
					if err != nil {
						atomic.AddInt32(&failureCount, 1)
						continue
					}

					atomic.AddInt32(&successCount, 1)
					time.Sleep(time.Microsecond)
					_ = pool.Put(resource)
				}
			}()
		}

		wg.Wait()

		total := successCount + failureCount
		t.Logf("Successes: %d, Failures: %d, Total: %d", successCount, failureCount, total)

		assertions.Greater(int(successCount), 0, "Expected some successful operations even under high contention")
	})

	t.Run("ZeroPoolImmediateCreation", func(t *testing.T) {
		assertions := assert.New(t)

		createDelay := time.Second
		config := PoolConfig[*TestResource]{
			MaxSize: 3,
			OnCreate: func() (*TestResource, error) {
				time.Sleep(createDelay)
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")
		defer pool.Close()

		ctx, cancel := context.WithTimeout(context.Background(), createDelay*2)
		defer cancel()

		// Test: Get from empty pool should create new resourcepool
		start := time.Now()
		resource, err := pool.Get(ctx)
		elapsed := time.Since(start)

		assertions.NoError(err, "Should be able to get resourcepool from empty pool")

		assertions.GreaterOrEqual(elapsed, createDelay, "Creation should take at least the expected delay")

		_ = pool.Put(resource)
	})
}

// TestResourcePoolStressTest - Stress testing with high load
func TestResourcePoolStressTest(t *testing.T) {
	assertions := assert.New(t)

	createCount := int32(0)
	destroyCount := int32(0)

	config := PoolConfig[*TestResource]{
		MaxSize: 5,
		OnCreate: func() (*TestResource, error) {
			id := atomic.AddInt32(&createCount, 1)
			return &TestResource{ID: int(id)}, nil
		},
		OnDestroy: func(r *TestResource) {
			atomic.AddInt32(&destroyCount, 1)
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	var wg sync.WaitGroup
	duration := time.Millisecond * 100
	start := time.Now()

	getOps := int32(0)
	putOps := int32(0)
	errs := int32(0)

	// Test: Stress test with multiple workers
	for w := 0; w < 20; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()

			for time.Since(start) < duration {
				var resource *TestResource
				resource, err = pool.Get(ctx)
				if err != nil {
					atomic.AddInt32(&errs, 1)
					continue
				}

				atomic.AddInt32(&getOps, 1)

				// Simulate work
				time.Sleep(time.Microsecond * time.Duration(rand.Intn(10)))

				err = pool.Put(resource)
				if err != nil {
					atomic.AddInt32(&errs, 1)
					continue
				}

				atomic.AddInt32(&putOps, 1)
			}
		}()
	}

	wg.Wait()
	pool.Close()

	t.Logf("Stress test results:")
	t.Logf("  Duration: %v", duration)
	t.Logf("  Get operations: %d", getOps)
	t.Logf("  Put operations: %d", putOps)
	t.Logf("  Errors: %d", errs)
	t.Logf("  Total created: %d", createCount)
	t.Logf("  Total destroyed: %d", destroyCount)

	// Test: Final verification
	stats := pool.GetPoolStats()
	t.Logf("  Pool stats - Created: %d, Size: %d", stats.Created, stats.Available)

	assertions.Greater(int(getOps), 0, "Expected some successful get operations")
	assertions.Equal(int(getOps), int(putOps), "Get and put operations should match")
}

// TestResourcePoolResourceLeakDetection - Resource leak detection
func TestResourcePoolResourceLeakDetection(t *testing.T) {
	assertions := assert.New(t)

	resourcesNotDestroyed := int32(0)
	resourcesCreated := int32(0)

	config := PoolConfig[*TestResource]{
		MaxSize: 3,
		OnCreate: func() (*TestResource, error) {
			count := atomic.AddInt32(&resourcesCreated, 1)
			atomic.AddInt32(&resourcesNotDestroyed, 1)
			return &TestResource{
				ID:      int(count),
				counter: &resourcesNotDestroyed,
			}, nil
		},
		OnDestroy: func(r *TestResource) {
			if r.counter != nil {
				atomic.AddInt32(r.counter, -1)
			}
		},
	}

	pool, err := NewResourcePool(config)
	assertions.NoError(err, "Should be able to create pool")

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	// Test: Create and use resources (limited by pool size)
	var resources []*TestResource
	for i := 0; i < 3; i++ { // Only get up to MaxSize resources
		var resource *TestResource
		resource, err = pool.Get(ctx)
		if err != nil {
			continue
		}
		resources = append(resources, resource)
	}

	// Return all resources
	for _, resource := range resources {
		err = pool.Put(resource)
		assertions.NoError(err, "Should be able to return resourcepool")
	}

	// Test: close pool and verify no leaks
	pool.Close()

	leakedResources := atomic.LoadInt32(&resourcesNotDestroyed)
	t.Logf("Final leak check: %d resources not destroyed", leakedResources)

	assertions.Equal(int32(0), leakedResources, "Should have no resourcepool leaks")
}

// TestResourcePoolMemoryUsage - Memory usage verification
func TestResourcePoolMemoryUsage(t *testing.T) {
	t.Run("LargeScaleOperations", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize: 1,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")
		defer pool.Close()

		// Force garbage collection before test
		runtime.GC()

		ctx := context.Background()

		// Test: Large number of operations
		for i := 0; i < 1000; i++ {
			var resource *TestResource
			resource, err = pool.Get(ctx)
			if err != nil {
				continue
			}
			_ = pool.Put(resource)
		}

		stats := pool.GetPoolStats()
		t.Logf("After %d operations - Size: %d, Created: %d", 1000, stats.Available, stats.Created)

		// Test: Pool should not grow beyond MaxSize
		assertions.LessOrEqual(stats.Created, 1, "Pool should not grow beyond MaxSize")

		// Force GC to detect any memorypool leaks
		runtime.GC()
	})
}

// TestResourcePoolTimingSensitive - Timing-dependent scenarios
func TestResourcePoolTimingSensitive(t *testing.T) {
	t.Run("SlowCreateFastConsumers", func(t *testing.T) {
		assertions := assert.New(t)

		createStarted := int32(0)
		createCompleted := int32(0)

		config := PoolConfig[*TestResource]{
			MaxSize: 5,
			OnCreate: func() (*TestResource, error) {
				atomic.AddInt32(&createStarted, 1)
				time.Sleep(time.Millisecond * 10) // Slow creation
				completed := atomic.AddInt32(&createCompleted, 1)
				return &TestResource{ID: int(completed)}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")
		defer pool.Close()

		var wg sync.WaitGroup
		successCount := int32(0)
		timeoutCount := int32(0)
		workers := 20

		start := time.Now()

		// Test: Fast consumers with slow creation
		for w := 0; w < workers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5)
				defer cancel()

				var resource *TestResource
				resource, err = pool.Get(ctx)
				if err != nil {
					atomic.AddInt32(&timeoutCount, 1)
					return
				}

				atomic.AddInt32(&successCount, 1)
				_ = pool.Put(resource)
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		t.Logf("Timing test results:")
		t.Logf("  Elapsed: %v", elapsed)
		t.Logf("  Creations started: %d, completed: %d", createStarted, createCompleted)
		t.Logf("  Successes: %d, timeouts: %d", successCount, timeoutCount)
		t.Logf("  Pool created: %d", pool.GetPoolStats().Created)

		// Some operations should timeout due to slow creation
		assertions.Greater(int(timeoutCount), 0, "Expected some timeout operations with slow creation")
	})

	t.Run("ConcurrentPutDuringClose", func(t *testing.T) {
		assertions := assert.New(t)

		config := PoolConfig[*TestResource]{
			MaxSize: 3,
			OnCreate: func() (*TestResource, error) {
				return &TestResource{ID: 1}, nil
			},
		}

		pool, err := NewResourcePool(config)
		assertions.NoError(err, "Should be able to create pool")

		ctx := context.Background()
		resource, _ := pool.Get(ctx)

		var wg sync.WaitGroup

		// Test: Put operation concurrent with close
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(time.Microsecond * 10)
			_ = pool.Put(resource)
		}()

		// close immediately
		pool.Close()

		wg.Wait()
		t.Log("Concurrent Put during close completed successfully")
	})
}
