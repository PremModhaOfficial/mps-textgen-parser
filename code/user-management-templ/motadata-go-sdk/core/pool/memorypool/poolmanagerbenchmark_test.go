package memorypool

import (
	"testing"

	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

func BenchmarkPoolManagerStringPoolAcquireRelease(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireStringPool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseStringPool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerStringPoolAcquireReleaseGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireStringPool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseStringPool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerInt64PoolAcquireRelease(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireInt64Pool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseInt64Pool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerInt64PoolAcquireReleaseGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireInt64Pool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseInt64Pool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerFloat64PoolAcquireRelease(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireFloat64Pool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseFloat64Pool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerFloat64PoolAcquireReleaseGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireFloat64Pool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseFloat64Pool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerBytePoolAcquireRelease(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireBytePool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseBytePool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerBytePoolAcquireReleaseGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 100, false, true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		poolIndex, _ := pm.AcquireBytePool(50)
		if poolIndex != NotAvailable {
			pm.ReleaseBytePool(poolIndex)
		}
	}
}

func BenchmarkPoolManagerAllPoolsMixedOperations(b *testing.B) {
	b.ResetTimer()

	pm := NewPoolManager(getPoolConfig(21, 100, true, false))

	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			poolIndex, _ := pm.AcquireStringPool(30)
			if poolIndex != NotAvailable {
				pm.ReleaseStringPool(poolIndex)
			}
		case 1:
			poolIndex, _ := pm.AcquireInt64Pool(40)
			if poolIndex != NotAvailable {
				pm.ReleaseInt64Pool(poolIndex)
			}
		case 2:
			poolIndex, _ := pm.AcquireFloat64Pool(50)
			if poolIndex != NotAvailable {
				pm.ReleaseFloat64Pool(poolIndex)
			}
		case 3:
			poolIndex, _ := pm.AcquireBytePool(60)
			if poolIndex != NotAvailable {
				pm.ReleaseBytePool(poolIndex)
			}
		}
	}
}

func BenchmarkPoolManagerAllPoolsMixedOperationsGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(20, 200, true, true))

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		switch i % 4 {
		case 0:
			poolIndex, _ := pm.AcquireStringPool(25)
			if poolIndex != NotAvailable {
				pm.ReleaseStringPool(poolIndex)
			}
		case 1:
			poolIndex, _ := pm.AcquireInt64Pool(30)
			if poolIndex != NotAvailable {
				pm.ReleaseInt64Pool(poolIndex)
			}
		case 2:
			poolIndex, _ := pm.AcquireFloat64Pool(35)
			if poolIndex != NotAvailable {
				pm.ReleaseFloat64Pool(poolIndex)
			}
		case 3:
			poolIndex, _ := pm.AcquireBytePool(40)
			if poolIndex != NotAvailable {
				pm.ReleaseBytePool(poolIndex)
			}
		}
	}
}

func BenchmarkPoolManagerGetPoolStats(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 50, false, false))

	stringIndex, _ := pm.AcquireStringPool(25)
	int64Index, _ := pm.AcquireInt64Pool(25)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.GetPoolStats()
	}

	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)
}

func BenchmarkPoolManagerGetPoolStatsGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 50, false, true))

	stringIndex, _ := pm.AcquireStringPool(25)
	int64Index, _ := pm.AcquireInt64Pool(25)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pm.GetPoolStats()
	}

	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)
}

func BenchmarkPoolManagerShrinkAllPools(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 50, true, false))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		poolIndex1, _ := pm.AcquireStringPool(100)
		poolIndex2, _ := pm.AcquireInt64Pool(80)
		pm.ReleaseStringPool(poolIndex1)
		pm.ReleaseInt64Pool(poolIndex2)
		b.StartTimer()

		pm.ShrinkAllPools()
	}
}

func BenchmarkPoolManagerShrinkAllPoolsGlobal(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(10, 50, true, true))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		poolIndex1, _ := pm.AcquireStringPool(100)
		poolIndex2, _ := pm.AcquireInt64Pool(80)
		pm.ReleaseStringPool(poolIndex1)
		pm.ReleaseInt64Pool(poolIndex2)
		b.StartTimer()

		pm.ShrinkAllPools()
	}
}

func BenchmarkPoolManagerConcurrentStringPool(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(20, 100, false, true))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			poolIndex, _ := pm.AcquireStringPool(50)
			if poolIndex != NotAvailable {
				pm.ReleaseStringPool(poolIndex)
			}
		}
	})
}

func BenchmarkPoolManagerConcurrentMixedPools(b *testing.B) {

	pm := NewPoolManager(getPoolConfig(20, 100, false, true))

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				poolIndex, _ := pm.AcquireStringPool(25)
				if poolIndex != NotAvailable {
					pm.ReleaseStringPool(poolIndex)
				}
			case 1:
				poolIndex, _ := pm.AcquireInt64Pool(30)
				if poolIndex != NotAvailable {
					pm.ReleaseInt64Pool(poolIndex)
				}
			case 2:
				poolIndex, _ := pm.AcquireFloat64Pool(35)
				if poolIndex != NotAvailable {
					pm.ReleaseFloat64Pool(poolIndex)
				}
			case 3:
				poolIndex, _ := pm.AcquireBytePool(40)
				if poolIndex != NotAvailable {
					pm.ReleaseBytePool(poolIndex)
				}
			}
			i++
		}
	})
}

func getPoolConfig(poolSize, poolLength int, expandable, global bool) *PoolConfig {

	return &PoolConfig{
		PoolSize:   poolSize,
		PoolLength: poolLength,
		Expandable: expandable,
		Global:     global,
	}
}
