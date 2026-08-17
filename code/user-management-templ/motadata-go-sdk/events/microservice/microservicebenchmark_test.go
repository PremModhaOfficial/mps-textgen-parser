package microservice

import (
	"context"
	"testing"
)

func BenchmarkRegistryRegisterService(b *testing.B) {
	r := NewRegistry(testRegistryCfg)
	defer r.Shutdown(context.Background())

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})
		}
	})
}

func BenchmarkRegistryGetService(b *testing.B) {
	r := NewRegistry(testRegistryCfg)
	defer r.Shutdown(context.Background())
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.GetService(testServiceID1)
		}
	})
}

func BenchmarkRegistryQueryServices(b *testing.B) {
	r := NewRegistry(testRegistryCfg)
	defer r.Shutdown(context.Background())

	for i := 0; i < 100; i++ {
		r.RegisterService(&Info{ID: "svc" + itoa(i), Name: testServiceTest, Status: StatusRunning, Tags: []string{"prod"}})
	}
	q := &Query{OnlyHealthy: true, Tags: []string{"prod"}}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.QueryServices(q)
		}
	})
}

func BenchmarkQueryMatches(b *testing.B) {
	service := &Info{
		ID: testServiceID1, Name: testServiceName, Type: ServiceTypeAPI, Status: StatusRunning,
		Capabilities: []string{"read", "write"}, Tags: []string{"production"},
		Metadata: map[string]string{"env": "prod"}, Instances: []*Instance{{ID: testInstanceID1}},
	}
	q := &Query{OnlyHealthy: true, Capabilities: []string{"read"}, Tags: []string{"production"}}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Matches(service)
		}
	})
}

func BenchmarkInfoClone(b *testing.B) {
	info := &Info{
		ID: testServiceID1, Name: testServiceName,
		Dependencies: []string{"dep1", "dep2", "dep3"}, Capabilities: []string{"cap1", "cap2"},
		Tags: []string{"tag1", "tag2", "tag3"}, Metadata: map[string]string{"key1": "value1", "key2": "value2"},
		Instances: []*Instance{{ID: testInstanceID1, Host: testHost}, {ID: testInstanceID2, Host: testHost}},
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			info.Clone()
		}
	})
}

func BenchmarkStatusString(b *testing.B) {
	var s string

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s = StatusRunning.String()
		}
	})
	_ = s
}
