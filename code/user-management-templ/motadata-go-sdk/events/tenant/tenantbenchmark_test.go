package tenant

import (
	"context"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"testing"
)

func BenchmarkRegistryRegister(b *testing.B) {
	r := NewRegistry()
	info := &Info{ID: testTenantID1, Name: testTenantTest}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.mu.Lock()
			r.tenants = make(map[string]*Info)
			r.mu.Unlock()
			r.Register(info)
		}
	})
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()
	r.Register(&Info{ID: testTenantID1, Name: testTenantTest})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Get(testTenantID1)
		}
	})
}

func BenchmarkRegistryList(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 100; i++ {
		r.Register(&Info{ID: "tenant-" + string(rune('a'+i%26))})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.List()
		}
	})
}

func BenchmarkManagerConnect(b *testing.B) {
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass(testUser, testPass)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Connect(ctx, testTenantID1, creds)
		}
	})
}

func BenchmarkManagerGet(b *testing.B) {
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass(testUser, testPass)
	m.Connect(ctx, testTenantID1, creds)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Get(testTenantID1)
		}
	})
}
