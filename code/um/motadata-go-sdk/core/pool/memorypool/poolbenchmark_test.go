package memorypool

import (
	"testing"

	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

func BenchmarkStringMemoryPoolAcquireRelease(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 100, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkStringMemoryPoolAcquireReleaseSmall(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(5)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkStringMemoryPoolAcquireReleaseLarge(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 1000, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(500)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkStringMemoryPoolAcquireReleaseExpandable(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(100)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkInt64MemoryPoolAcquireRelease(b *testing.B) {
	memoryPool := newMemoryPool[int64](10, 100, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkInt64MemoryPoolAcquireReleaseSmall(b *testing.B) {
	memoryPool := newMemoryPool[int64](10, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(5)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkInt64MemoryPoolAcquireReleaseLarge(b *testing.B) {
	memoryPool := newMemoryPool[int64](10, 1000, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(500)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkInt64MemoryPoolAcquireReleaseExpandable(b *testing.B) {
	memoryPool := newMemoryPool[int64](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(100)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkFloat64MemoryPoolAcquireRelease(b *testing.B) {
	memoryPool := newMemoryPool[float64](10, 100, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkFloat64MemoryPoolAcquireReleaseSmall(b *testing.B) {
	memoryPool := newMemoryPool[float64](10, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(5)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkFloat64MemoryPoolAcquireReleaseLarge(b *testing.B) {
	memoryPool := newMemoryPool[float64](10, 1000, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(500)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkFloat64MemoryPoolAcquireReleaseExpandable(b *testing.B) {
	memoryPool := newMemoryPool[float64](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(100)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkByteMemoryPoolAcquireRelease(b *testing.B) {
	memoryPool := newMemoryPool[byte](10, 100, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkByteMemoryPoolAcquireReleaseSmall(b *testing.B) {
	memoryPool := newMemoryPool[byte](10, 20, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(5)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkByteMemoryPoolAcquireReleaseLarge(b *testing.B) {
	memoryPool := newMemoryPool[byte](10, 1000, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(500)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkByteMemoryPoolAcquireReleaseExpandable(b *testing.B) {
	memoryPool := newMemoryPool[byte](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(100)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}
}

func BenchmarkStringMemoryPoolGetPool(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 100, false)
	indices := make([]int, 10)

	for i := 0; i < 10; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		indices[i] = poolIndex
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memoryPool.getPool(indices[i%10])
	}

	for _, poolIndex := range indices {
		memoryPool.releasePool(poolIndex)
	}
}

func BenchmarkStringMemoryPoolGetUsedPools(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 100, false)
	indices := make([]int, 5)

	for i := 0; i < 5; i++ {
		poolIndex, _ := memoryPool.acquirePool(50)
		indices[i] = poolIndex
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memoryPool.getUsedPools()
	}

	for _, poolIndex := range indices {
		memoryPool.releasePool(poolIndex)
	}
}

func BenchmarkStringMemoryPoolExpandPool(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		poolIndex, _ := memoryPool.acquirePool(25)
		b.StartTimer()

		memoryPool.expandPool(poolIndex, 100)

		b.StopTimer()
		memoryPool.releasePool(poolIndex)
		b.StartTimer()
	}
}

func BenchmarkStringMemoryPoolShrinkPool(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 50, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		poolIndex, _ := memoryPool.acquirePool(100)
		memoryPool.releasePool(poolIndex)
		b.StartTimer()

		memoryPool.shrinkPool()
	}
}

func BenchmarkStringMemoryPoolTestPoolLeak(b *testing.B) {
	memoryPool := newMemoryPool[string](10, 50, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		memoryPool.acquirePool(25)
		memoryPool.acquirePool(30)
		b.StartTimer()

		memoryPool.testPoolLeak()
	}
}

func BenchmarkStringMemoryPoolConcurrentOperations(b *testing.B) {
	memoryPool := newMemoryPool[string](20, 100, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := memoryPool.acquirePool(50)
			if poolIndex != NotAvailable {
				memoryPool.releasePool(poolIndex)
			}
		}
	})
}

func BenchmarkInt64MemoryPoolConcurrentOperations(b *testing.B) {
	memoryPool := newMemoryPool[int64](20, 100, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := memoryPool.acquirePool(50)
			if poolIndex != NotAvailable {
				memoryPool.releasePool(poolIndex)
			}
		}
	})
}

func BenchmarkFloat64MemoryPoolConcurrentOperations(b *testing.B) {
	memoryPool := newMemoryPool[float64](20, 100, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := memoryPool.acquirePool(50)
			if poolIndex != NotAvailable {
				memoryPool.releasePool(poolIndex)
			}
		}
	})
}

func BenchmarkByteMemoryPoolConcurrentOperations(b *testing.B) {
	memoryPool := newMemoryPool[byte](20, 100, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := memoryPool.acquirePool(50)
			if poolIndex != NotAvailable {
				memoryPool.releasePool(poolIndex)
			}
		}
	})
}

func BenchmarkMixedMemoryPoolOperations(b *testing.B) {
	stringPool := newMemoryPool[string](5, 50, true)
	int64Pool := newMemoryPool[int64](5, 50, true)
	float64Pool := newMemoryPool[float64](5, 50, true)
	bytePool := newMemoryPool[byte](5, 50, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			poolIndex, _ := stringPool.acquirePool(25)
			if poolIndex != NotAvailable {
				stringPool.releasePool(poolIndex)
			}
		case 1:
			poolIndex, _ := int64Pool.acquirePool(30)
			if poolIndex != NotAvailable {
				int64Pool.releasePool(poolIndex)
			}
		case 2:
			poolIndex, _ := float64Pool.acquirePool(35)
			if poolIndex != NotAvailable {
				float64Pool.releasePool(poolIndex)
			}
		case 3:
			poolIndex, _ := bytePool.acquirePool(40)
			if poolIndex != NotAvailable {
				bytePool.releasePool(poolIndex)
			}
		}
	}
}

func BenchmarkStringMemoryPoolHighContentionSinglePool(b *testing.B) {
	memoryPool := newMemoryPool[string](2, 100, false)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := memoryPool.acquirePool(50)
			if poolIndex != NotAvailable {
				memoryPool.releasePool(poolIndex)
			}
		}
	})
}

func BenchmarkStringMemoryPoolMemoryIntensiveOperations(b *testing.B) {
	memoryPool := newMemoryPool[string](5, 10000, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := memoryPool.acquirePool(5000)
		if poolIndex != NotAvailable {
			buffer := memoryPool.getPool(poolIndex)
			for j := 0; j < 100; j++ {
				if j < len(buffer) {
					buffer[j] = "test"
				}
			}
			memoryPool.releasePool(poolIndex)
		}
	}
}
