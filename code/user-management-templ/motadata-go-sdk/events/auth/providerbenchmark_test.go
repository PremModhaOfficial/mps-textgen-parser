package auth

import (
	"context"
	"testing"
)

func BenchmarkStaticProviderAuthenticate(b *testing.B) {

	ctx := context.Background()
	creds := Token(testTokenValue)
	p := NewStaticProvider(creds)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.Authenticate(ctx)
	}
}

func BenchmarkManagedProviderAuthenticateParallel(b *testing.B) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	provider := NewManagedProvider(cm, testTenantID)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			provider.Authenticate(ctx)
		}
	})
}

func BenchmarkParseRegistrationPayload(b *testing.B) {

	data := []byte(`{"tenant_id":"` + testTenantID + `","jwt":"` + testJWTTokenLiteral + `","seed":"` + testSeedValue + `"}`)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRegistrationPayload(data)
	}
}
