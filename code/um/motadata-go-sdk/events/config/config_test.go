package config

import (
	"os"
	"testing"
	"time"
)

// ============== Config Tests ==============

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg == nil {
		t.Fatal("DefaultConfig returned nil")
	}

	if len(cfg.Servers) != 1 || cfg.Servers[0] != "nats://localhost:4222" {
		t.Errorf("unexpected default servers: %v", cfg.Servers)
	}

	if cfg.ConnectTimeout != 10*time.Second {
		t.Errorf("unexpected connect timeout: %v", cfg.ConnectTimeout)
	}

	if cfg.DrainTimeout != 30*time.Second {
		t.Errorf("unexpected drain timeout: %v", cfg.DrainTimeout)
	}

	if !cfg.Stream.Enabled {
		t.Error("expected stream to be enabled by default")
	}

	if cfg.IdleTimeout != 30*time.Minute {
		t.Errorf("unexpected idle timeout: %v", cfg.IdleTimeout)
	}

	if cfg.CleanupInterval != 1*time.Minute {
		t.Errorf("unexpected cleanup interval: %v", cfg.CleanupInterval)
	}
}

func TestDefaultReconnectConfig(t *testing.T) {
	cfg := DefaultReconnectConfig()

	if cfg.MaxAttempts != -1 {
		t.Errorf("expected -1 for infinite retries, got %d", cfg.MaxAttempts)
	}

	if cfg.InitialInterval != 100*time.Millisecond {
		t.Errorf("unexpected initial interval: %v", cfg.InitialInterval)
	}

	if cfg.MaxInterval != 30*time.Second {
		t.Errorf("unexpected max interval: %v", cfg.MaxInterval)
	}

	if cfg.Multiplier != 2.0 {
		t.Errorf("unexpected multiplier: %v", cfg.Multiplier)
	}

	if cfg.Jitter != 0.1 {
		t.Errorf("unexpected jitter: %v", cfg.Jitter)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := DefaultRetryConfig()

	if cfg.MaxAttempts != 3 {
		t.Errorf("expected 3 max attempts, got %d", cfg.MaxAttempts)
	}

	if cfg.InitialInterval != 100*time.Millisecond {
		t.Errorf("unexpected initial interval: %v", cfg.InitialInterval)
	}

	if cfg.MaxInterval != 5*time.Second {
		t.Errorf("unexpected max interval: %v", cfg.MaxInterval)
	}
}

func TestDefaultPublishConfig(t *testing.T) {
	cfg := DefaultPublishConfig()

	if cfg.AckTimeout != 5*time.Second {
		t.Errorf("unexpected ack timeout: %v", cfg.AckTimeout)
	}

	if !cfg.EnableDeduplication {
		t.Error("expected deduplication to be enabled")
	}

	if cfg.DeduplicationWindow != 2*time.Minute {
		t.Errorf("unexpected deduplication window: %v", cfg.DeduplicationWindow)
	}

	if cfg.Retry.MaxAttempts != 3 {
		t.Errorf("unexpected retry max attempts: %d", cfg.Retry.MaxAttempts)
	}
}

func TestDefaultSubscribeConfig(t *testing.T) {
	cfg := DefaultSubscribeConfig()

	if cfg.MaxConcurrent != 10 {
		t.Errorf("unexpected max concurrent: %d", cfg.MaxConcurrent)
	}

	if cfg.AckWait != 30*time.Second {
		t.Errorf("unexpected ack wait: %v", cfg.AckWait)
	}

	if cfg.BatchSize != 100 {
		t.Errorf("unexpected batch size: %d", cfg.BatchSize)
	}

	if cfg.BatchWait != 100*time.Millisecond {
		t.Errorf("unexpected batch wait: %v", cfg.BatchWait)
	}
}

func TestConfigClone(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Name = "test-client"
	cfg.TLS = &TLSConfig{
		Enabled:  true,
		CertFile: "/path/to/cert",
	}

	cloned := cfg.Clone()

	// Verify clone is not same pointer
	if cloned == cfg {
		t.Error("clone should be a different pointer")
	}

	// Verify values are copied
	if cloned.Name != cfg.Name {
		t.Errorf("name not cloned: %s != %s", cloned.Name, cfg.Name)
	}

	// Verify servers slice is copied
	if &cloned.Servers[0] == &cfg.Servers[0] {
		t.Error("servers slice should be a copy")
	}

	// Verify TLS is copied
	if cloned.TLS == cfg.TLS {
		t.Error("TLS config should be a copy")
	}
	if cloned.TLS.CertFile != cfg.TLS.CertFile {
		t.Errorf("TLS cert file not cloned")
	}

	// Modify clone and verify original is unchanged
	cloned.Name = "modified"
	if cfg.Name == "modified" {
		t.Error("modifying clone affected original")
	}
}

func TestConfigCloneNil(t *testing.T) {
	var cfg *Config
	cloned := cfg.Clone()
	if cloned != nil {
		t.Error("cloning nil config should return nil")
	}
}

func TestConfigCloneWithoutTLS(t *testing.T) {
	cfg := DefaultConfig()
	cloned := cfg.Clone()

	if cloned.TLS != nil {
		t.Error("TLS should be nil when original has no TLS")
	}
}

// ============== Options Tests ==============

func TestWithServers(t *testing.T) {
	cfg := DefaultConfig()
	WithServers("nats://server1:4222", "nats://server2:4222")(cfg)

	if len(cfg.Servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(cfg.Servers))
	}
	if cfg.Servers[0] != "nats://server1:4222" {
		t.Errorf("unexpected first server: %s", cfg.Servers[0])
	}
}

func TestWithName(t *testing.T) {
	cfg := DefaultConfig()
	WithName("my-client")(cfg)

	if cfg.Name != "my-client" {
		t.Errorf("expected name 'my-client', got '%s'", cfg.Name)
	}
}

func TestWithConnectTimeout(t *testing.T) {
	cfg := DefaultConfig()
	WithConnectTimeout(5 * time.Second)(cfg)

	if cfg.ConnectTimeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cfg.ConnectTimeout)
	}
}

func TestWithDrainTimeout(t *testing.T) {
	cfg := DefaultConfig()
	WithDrainTimeout(15 * time.Second)(cfg)

	if cfg.DrainTimeout != 15*time.Second {
		t.Errorf("expected 15s drain timeout, got %v", cfg.DrainTimeout)
	}
}

func TestWithTLS(t *testing.T) {
	cfg := DefaultConfig()
	tlsCfg := &TLSConfig{Enabled: true, CertFile: "/cert.pem"}
	WithTLS(tlsCfg)(cfg)

	if cfg.TLS != tlsCfg {
		t.Error("TLS config not set correctly")
	}
}

func TestWithTLSFiles(t *testing.T) {
	cfg := DefaultConfig()
	WithTLSFiles("/cert.pem", "/key.pem", "/ca.pem")(cfg)

	if cfg.TLS == nil {
		t.Fatal("TLS config should not be nil")
	}
	if !cfg.TLS.Enabled {
		t.Error("TLS should be enabled")
	}
	if cfg.TLS.CertFile != "/cert.pem" {
		t.Errorf("unexpected cert file: %s", cfg.TLS.CertFile)
	}
	if cfg.TLS.KeyFile != "/key.pem" {
		t.Errorf("unexpected key file: %s", cfg.TLS.KeyFile)
	}
	if cfg.TLS.CAFile != "/ca.pem" {
		t.Errorf("unexpected CA file: %s", cfg.TLS.CAFile)
	}
}

func TestWithTLSSkipVerify(t *testing.T) {
	// Test with nil TLS
	cfg := DefaultConfig()
	WithTLSSkipVerify()(cfg)

	if cfg.TLS == nil {
		t.Fatal("TLS config should be created")
	}
	if !cfg.TLS.SkipVerify {
		t.Error("SkipVerify should be true")
	}

	// Test with existing TLS
	cfg2 := DefaultConfig()
	cfg2.TLS = &TLSConfig{Enabled: true}
	WithTLSSkipVerify()(cfg2)
	if !cfg2.TLS.SkipVerify {
		t.Error("SkipVerify should be true on existing TLS")
	}
}

func TestWithReconnect(t *testing.T) {
	cfg := DefaultConfig()
	rc := ReconnectConfig{
		MaxAttempts:     5,
		InitialInterval: 200 * time.Millisecond,
	}
	WithReconnect(rc)(cfg)

	if cfg.Reconnect.MaxAttempts != 5 {
		t.Errorf("expected 5 max attempts, got %d", cfg.Reconnect.MaxAttempts)
	}
}

func TestWithMaxReconnects(t *testing.T) {
	cfg := DefaultConfig()
	WithMaxReconnects(10)(cfg)

	if cfg.Reconnect.MaxAttempts != 10 {
		t.Errorf("expected 10 max reconnects, got %d", cfg.Reconnect.MaxAttempts)
	}
}

func TestWithStream(t *testing.T) {
	cfg := DefaultConfig()
	sc := StreamConfig{Enabled: true, Domain: "hub"}
	WithStream(sc)(cfg)

	if cfg.Stream.Domain != "hub" {
		t.Errorf("expected domain 'hub', got '%s'", cfg.Stream.Domain)
	}
}

func TestWithStreamDisabled(t *testing.T) {
	cfg := DefaultConfig()
	WithStreamDisabled()(cfg)

	if cfg.Stream.Enabled {
		t.Error("stream should be disabled")
	}
}

func TestWithStreamDomain(t *testing.T) {
	cfg := DefaultConfig()
	WithStreamDomain("leaf")(cfg)

	if cfg.Stream.Domain != "leaf" {
		t.Errorf("expected domain 'leaf', got '%s'", cfg.Stream.Domain)
	}
}

func TestWithIdleTimeout(t *testing.T) {
	cfg := DefaultConfig()
	WithIdleTimeout(1 * time.Hour)(cfg)

	if cfg.IdleTimeout != 1*time.Hour {
		t.Errorf("expected 1h idle timeout, got %v", cfg.IdleTimeout)
	}
}

func TestWithCleanupInterval(t *testing.T) {
	cfg := DefaultConfig()
	WithCleanupInterval(5 * time.Minute)(cfg)

	if cfg.CleanupInterval != 5*time.Minute {
		t.Errorf("expected 5m cleanup interval, got %v", cfg.CleanupInterval)
	}
}

func TestApply(t *testing.T) {
	cfg := DefaultConfig()
	result := Apply(cfg,
		WithName("test"),
		WithConnectTimeout(5*time.Second),
		WithStreamDisabled(),
	)

	if result != cfg {
		t.Error("Apply should return same config pointer")
	}
	if cfg.Name != "test" {
		t.Errorf("name not applied: %s", cfg.Name)
	}
	if cfg.ConnectTimeout != 5*time.Second {
		t.Errorf("timeout not applied: %v", cfg.ConnectTimeout)
	}
	if cfg.Stream.Enabled {
		t.Error("stream should be disabled")
	}
}

// ============== Publish Options Tests ==============

func TestWithRetry(t *testing.T) {
	cfg := DefaultPublishConfig()
	rc := RetryConfig{MaxAttempts: 5}
	WithRetry(rc)(&cfg)

	if cfg.Retry.MaxAttempts != 5 {
		t.Errorf("expected 5 max attempts, got %d", cfg.Retry.MaxAttempts)
	}
}

func TestWithMaxRetries(t *testing.T) {
	cfg := DefaultPublishConfig()
	WithMaxRetries(10)(&cfg)

	if cfg.Retry.MaxAttempts != 10 {
		t.Errorf("expected 10 max retries, got %d", cfg.Retry.MaxAttempts)
	}
}

func TestWithAckTimeout(t *testing.T) {
	cfg := DefaultPublishConfig()
	WithAckTimeout(10 * time.Second)(&cfg)

	if cfg.AckTimeout != 10*time.Second {
		t.Errorf("expected 10s ack timeout, got %v", cfg.AckTimeout)
	}
}

func TestWithDeduplication(t *testing.T) {
	cfg := DefaultPublishConfig()
	WithDeduplication(false)(&cfg)

	if cfg.EnableDeduplication {
		t.Error("deduplication should be disabled")
	}
}

func TestWithDeduplicationWindow(t *testing.T) {
	cfg := DefaultPublishConfig()
	WithDeduplicationWindow(5 * time.Minute)(&cfg)

	if cfg.DeduplicationWindow != 5*time.Minute {
		t.Errorf("expected 5m window, got %v", cfg.DeduplicationWindow)
	}
}

// ============== Subscribe Options Tests ==============

func TestWithQueueGroup(t *testing.T) {
	cfg := DefaultSubscribeConfig()
	WithQueueGroup("workers")(&cfg)

	if cfg.QueueGroup != "workers" {
		t.Errorf("expected queue group 'workers', got '%s'", cfg.QueueGroup)
	}
}

func TestWithMaxConcurrent(t *testing.T) {
	cfg := DefaultSubscribeConfig()
	WithMaxConcurrent(20)(&cfg)

	if cfg.MaxConcurrent != 20 {
		t.Errorf("expected 20 max concurrent, got %d", cfg.MaxConcurrent)
	}
}

func TestWithAckWait(t *testing.T) {
	cfg := DefaultSubscribeConfig()
	WithAckWait(1 * time.Minute)(&cfg)

	if cfg.AckWait != 1*time.Minute {
		t.Errorf("expected 1m ack wait, got %v", cfg.AckWait)
	}
}

func TestWithBatch(t *testing.T) {
	cfg := DefaultSubscribeConfig()
	WithBatch(50, 200*time.Millisecond)(&cfg)

	if cfg.BatchSize != 50 {
		t.Errorf("expected batch size 50, got %d", cfg.BatchSize)
	}
	if cfg.BatchWait != 200*time.Millisecond {
		t.Errorf("expected batch wait 200ms, got %v", cfg.BatchWait)
	}
}

// ============== Loader Tests ==============

func TestNewLoader(t *testing.T) {
	loader := NewLoader("TEST_")
	if loader == nil {
		t.Fatal("NewLoader returned nil")
	}
	if loader.prefix != "TEST_" {
		t.Errorf("expected prefix 'TEST_', got '%s'", loader.prefix)
	}
}

func TestDefaultLoader(t *testing.T) {
	loader := DefaultLoader()
	if loader.prefix != "EVENTS_" {
		t.Errorf("expected prefix 'EVENTS_', got '%s'", loader.prefix)
	}
}

func TestLoaderGetEnv(t *testing.T) {
	os.Setenv("TEST_KEY", "value")
	defer os.Unsetenv("TEST_KEY")

	loader := NewLoader("TEST_")
	if loader.getEnv("KEY") != "value" {
		t.Error("getEnv failed to retrieve value")
	}
}

func TestLoaderGetBool(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{"true", true},
		{"TRUE", true},
		{"True", true},
		{"1", true},
		{"yes", true},
		{"YES", true},
		{"false", false},
		{"0", false},
		{"no", false},
		{"", false},
	}

	loader := NewLoader("TEST_BOOL_")
	for _, tt := range tests {
		os.Setenv("TEST_BOOL_KEY", tt.value)
		if loader.getBool("KEY") != tt.expected {
			t.Errorf("getBool(%q) = %v, expected %v", tt.value, loader.getBool("KEY"), tt.expected)
		}
		os.Unsetenv("TEST_BOOL_KEY")
	}
}

func TestLoaderGetInt(t *testing.T) {
	loader := NewLoader("TEST_INT_")

	os.Setenv("TEST_INT_KEY", "42")
	defer os.Unsetenv("TEST_INT_KEY")

	if loader.getInt("KEY") != 42 {
		t.Errorf("expected 42, got %d", loader.getInt("KEY"))
	}

	// Test invalid value
	os.Setenv("TEST_INT_KEY", "invalid")
	if loader.getInt("KEY") != 0 {
		t.Error("expected 0 for invalid int")
	}

	// Test empty value
	os.Unsetenv("TEST_INT_KEY")
	if loader.getInt("KEY") != 0 {
		t.Error("expected 0 for missing value")
	}
}

func TestLoaderGetFloat(t *testing.T) {
	loader := NewLoader("TEST_FLOAT_")

	os.Setenv("TEST_FLOAT_KEY", "3.14")
	defer os.Unsetenv("TEST_FLOAT_KEY")

	if loader.getFloat("KEY") != 3.14 {
		t.Errorf("expected 3.14, got %f", loader.getFloat("KEY"))
	}

	// Test invalid value
	os.Setenv("TEST_FLOAT_KEY", "invalid")
	if loader.getFloat("KEY") != 0 {
		t.Error("expected 0 for invalid float")
	}
}

func TestLoaderGetDuration(t *testing.T) {
	loader := NewLoader("TEST_DUR_")

	os.Setenv("TEST_DUR_KEY", "5s")
	defer os.Unsetenv("TEST_DUR_KEY")

	if loader.getDuration("KEY") != 5*time.Second {
		t.Errorf("expected 5s, got %v", loader.getDuration("KEY"))
	}

	// Test invalid value
	os.Setenv("TEST_DUR_KEY", "invalid")
	if loader.getDuration("KEY") != 0 {
		t.Error("expected 0 for invalid duration")
	}
}

func TestLoaderParseServers(t *testing.T) {
	loader := NewLoader("")

	tests := []struct {
		input    string
		expected []string
	}{
		{"nats://a:4222", []string{"nats://a:4222"}},
		{"nats://a:4222,nats://b:4222", []string{"nats://a:4222", "nats://b:4222"}},
		{"nats://a:4222, nats://b:4222", []string{"nats://a:4222", "nats://b:4222"}},
		{" nats://a:4222 , nats://b:4222 ", []string{"nats://a:4222", "nats://b:4222"}},
		{"", []string{}},
		{",,,", []string{}},
	}

	for _, tt := range tests {
		result := loader.parseServers(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseServers(%q) = %v, expected %v", tt.input, result, tt.expected)
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("parseServers(%q)[%d] = %q, expected %q", tt.input, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set up environment
	os.Setenv("EVENTS_SERVERS", "nats://test:4222")
	os.Setenv("EVENTS_NAME", "test-client")
	os.Setenv("EVENTS_CONNECT_TIMEOUT", "5s")
	os.Setenv("EVENTS_DRAIN_TIMEOUT", "10s")
	defer func() {
		os.Unsetenv("EVENTS_SERVERS")
		os.Unsetenv("EVENTS_NAME")
		os.Unsetenv("EVENTS_CONNECT_TIMEOUT")
		os.Unsetenv("EVENTS_DRAIN_TIMEOUT")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if len(cfg.Servers) != 1 || cfg.Servers[0] != "nats://test:4222" {
		t.Errorf("unexpected servers: %v", cfg.Servers)
	}
	if cfg.Name != "test-client" {
		t.Errorf("unexpected name: %s", cfg.Name)
	}
	if cfg.ConnectTimeout != 5*time.Second {
		t.Errorf("unexpected connect timeout: %v", cfg.ConnectTimeout)
	}
	if cfg.DrainTimeout != 10*time.Second {
		t.Errorf("unexpected drain timeout: %v", cfg.DrainTimeout)
	}
}

func TestLoadFromEnvWithPrefix(t *testing.T) {
	os.Setenv("MY_SERVERS", "nats://custom:4222")
	defer os.Unsetenv("MY_SERVERS")

	cfg, err := LoadFromEnvWithPrefix("MY_")
	if err != nil {
		t.Fatalf("LoadFromEnvWithPrefix failed: %v", err)
	}

	if len(cfg.Servers) != 1 || cfg.Servers[0] != "nats://custom:4222" {
		t.Errorf("unexpected servers: %v", cfg.Servers)
	}
}

func TestLoadFromEnvWithTLS(t *testing.T) {
	os.Setenv("EVENTS_SERVERS", "nats://test:4222")
	os.Setenv("EVENTS_TLS_ENABLED", "true")
	os.Setenv("EVENTS_TLS_SKIP_VERIFY", "true")
	defer func() {
		os.Unsetenv("EVENTS_SERVERS")
		os.Unsetenv("EVENTS_TLS_ENABLED")
		os.Unsetenv("EVENTS_TLS_SKIP_VERIFY")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if cfg.TLS == nil {
		t.Fatal("TLS config should not be nil")
	}
	if !cfg.TLS.Enabled {
		t.Error("TLS should be enabled")
	}
	if !cfg.TLS.SkipVerify {
		t.Error("TLS SkipVerify should be true")
	}
}

func TestLoadFromEnvWithReconnect(t *testing.T) {
	os.Setenv("EVENTS_SERVERS", "nats://test:4222")
	os.Setenv("EVENTS_RECONNECT_MAX_ATTEMPTS", "5")
	os.Setenv("EVENTS_RECONNECT_INITIAL", "200ms")
	os.Setenv("EVENTS_RECONNECT_MAX", "10s")
	os.Setenv("EVENTS_RECONNECT_MULTIPLIER", "1.5")
	os.Setenv("EVENTS_RECONNECT_JITTER", "0.2")
	defer func() {
		os.Unsetenv("EVENTS_SERVERS")
		os.Unsetenv("EVENTS_RECONNECT_MAX_ATTEMPTS")
		os.Unsetenv("EVENTS_RECONNECT_INITIAL")
		os.Unsetenv("EVENTS_RECONNECT_MAX")
		os.Unsetenv("EVENTS_RECONNECT_MULTIPLIER")
		os.Unsetenv("EVENTS_RECONNECT_JITTER")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if cfg.Reconnect.MaxAttempts != 5 {
		t.Errorf("unexpected max attempts: %d", cfg.Reconnect.MaxAttempts)
	}
	if cfg.Reconnect.InitialInterval != 200*time.Millisecond {
		t.Errorf("unexpected initial interval: %v", cfg.Reconnect.InitialInterval)
	}
	if cfg.Reconnect.MaxInterval != 10*time.Second {
		t.Errorf("unexpected max interval: %v", cfg.Reconnect.MaxInterval)
	}
	if cfg.Reconnect.Multiplier != 1.5 {
		t.Errorf("unexpected multiplier: %v", cfg.Reconnect.Multiplier)
	}
	if cfg.Reconnect.Jitter != 0.2 {
		t.Errorf("unexpected jitter: %v", cfg.Reconnect.Jitter)
	}
}

func TestLoadFromEnvWithStream(t *testing.T) {
	os.Setenv("EVENTS_SERVERS", "nats://test:4222")
	os.Setenv("EVENTS_STREAM_DISABLED", "true")
	os.Setenv("EVENTS_STREAM_DOMAIN", "hub")
	os.Setenv("EVENTS_STREAM_PREFIX", "prefix")
	defer func() {
		os.Unsetenv("EVENTS_SERVERS")
		os.Unsetenv("EVENTS_STREAM_DISABLED")
		os.Unsetenv("EVENTS_STREAM_DOMAIN")
		os.Unsetenv("EVENTS_STREAM_PREFIX")
	}()

	cfg, err := LoadFromEnv()
	if err != nil {
		t.Fatalf("LoadFromEnv failed: %v", err)
	}

	if cfg.Stream.Enabled {
		t.Error("stream should be disabled")
	}
	if cfg.Stream.Domain != "hub" {
		t.Errorf("unexpected domain: %s", cfg.Stream.Domain)
	}
	if cfg.Stream.Prefix != "prefix" {
		t.Errorf("unexpected prefix: %s", cfg.Stream.Prefix)
	}
}

func TestLoadFromEnvWithOptions(t *testing.T) {
	os.Setenv("EVENTS_SERVERS", "nats://test:4222")
	defer os.Unsetenv("EVENTS_SERVERS")

	loader := DefaultLoader()
	cfg, err := loader.LoadFromEnvWithOptions(
		WithName("custom-name"),
		WithConnectTimeout(15*time.Second),
	)

	if err != nil {
		t.Fatalf("LoadFromEnvWithOptions failed: %v", err)
	}

	if cfg.Name != "custom-name" {
		t.Errorf("expected name 'custom-name', got '%s'", cfg.Name)
	}
	if cfg.ConnectTimeout != 15*time.Second {
		t.Errorf("expected 15s timeout, got %v", cfg.ConnectTimeout)
	}
}

// ============== Validation Tests ==============

func TestValidate(t *testing.T) {
	// Valid config
	cfg := DefaultConfig()
	if err := Validate(cfg); err != nil {
		t.Errorf("valid config should not return error: %v", err)
	}
}

func TestValidateNoServers(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Servers = nil

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for no servers")
	}
}

func TestValidateEmptyServer(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Servers = []string{"nats://valid:4222", ""}

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for empty server")
	}
}

func TestValidateInvalidServerScheme(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Servers = []string{"http://invalid:4222"}

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for invalid server scheme")
	}
}

func TestValidateValidServerSchemes(t *testing.T) {
	validSchemes := []string{
		"nats://localhost:4222",
		"tls://localhost:4222",
		"ws://localhost:4222",
		"wss://localhost:4222",
	}

	for _, server := range validSchemes {
		cfg := DefaultConfig()
		cfg.Servers = []string{server}
		if err := Validate(cfg); err != nil {
			t.Errorf("scheme %s should be valid: %v", server, err)
		}
	}
}

func TestValidateConnectTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ConnectTimeout = 0

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for zero connect timeout")
	}
}

func TestValidateDrainTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.DrainTimeout = 0

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for zero drain timeout")
	}
}

func TestValidateReconnectIntervals(t *testing.T) {
	cfg := DefaultConfig()

	// Zero initial interval
	cfg.Reconnect.InitialInterval = 0
	if err := Validate(cfg); err == nil {
		t.Error("expected error for zero initial interval")
	}

	// Max < Initial
	cfg.Reconnect.InitialInterval = 1 * time.Second
	cfg.Reconnect.MaxInterval = 500 * time.Millisecond
	if err := Validate(cfg); err == nil {
		t.Error("expected error for max < initial")
	}
}

func TestValidateReconnectMultiplier(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Reconnect.Multiplier = 0.5

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for multiplier < 1")
	}
}

func TestValidateReconnectJitter(t *testing.T) {
	cfg := DefaultConfig()

	// Negative jitter
	cfg.Reconnect.Jitter = -0.1
	if err := Validate(cfg); err == nil {
		t.Error("expected error for negative jitter")
	}

	// Jitter > 1
	cfg.Reconnect.Jitter = 1.1
	if err := Validate(cfg); err == nil {
		t.Error("expected error for jitter > 1")
	}
}

func TestValidateIdleTimeout(t *testing.T) {
	cfg := DefaultConfig()
	cfg.IdleTimeout = 0

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for zero idle timeout")
	}
}

func TestValidateCleanupInterval(t *testing.T) {
	cfg := DefaultConfig()
	cfg.CleanupInterval = 0

	err := Validate(cfg)
	if err == nil {
		t.Error("expected error for zero cleanup interval")
	}
}

func TestValidatePublish(t *testing.T) {
	cfg := DefaultPublishConfig()
	if err := ValidatePublish(&cfg); err != nil {
		t.Errorf("valid publish config should not return error: %v", err)
	}

	// Negative max attempts
	cfg.Retry.MaxAttempts = -1
	if err := ValidatePublish(&cfg); err == nil {
		t.Error("expected error for negative max attempts")
	}

	// Zero ack timeout
	cfg.Retry.MaxAttempts = 3
	cfg.AckTimeout = 0
	if err := ValidatePublish(&cfg); err == nil {
		t.Error("expected error for zero ack timeout")
	}
}

func TestValidateSubscribe(t *testing.T) {
	cfg := DefaultSubscribeConfig()
	if err := ValidateSubscribe(&cfg); err != nil {
		t.Errorf("valid subscribe config should not return error: %v", err)
	}

	// Zero max concurrent
	cfg.MaxConcurrent = 0
	if err := ValidateSubscribe(&cfg); err == nil {
		t.Error("expected error for zero max concurrent")
	}

	// Zero ack wait
	cfg.MaxConcurrent = 10
	cfg.AckWait = 0
	if err := ValidateSubscribe(&cfg); err == nil {
		t.Error("expected error for zero ack wait")
	}
}

func TestValidateTLSNil(t *testing.T) {
	if err := ValidateTLS(nil); err != nil {
		t.Errorf("nil TLS should be valid: %v", err)
	}
}

func TestValidateTLSCertKeyMismatch(t *testing.T) {
	// Cert without key
	tls := &TLSConfig{
		Enabled:  true,
		CertFile: "/cert.pem",
	}
	if err := ValidateTLS(tls); err == nil {
		t.Error("expected error for cert without key")
	}

	// Key without cert
	tls = &TLSConfig{
		Enabled: true,
		KeyFile: "/key.pem",
	}
	if err := ValidateTLS(tls); err == nil {
		t.Error("expected error for key without cert")
	}
}

// ============== Benchmarks ==============

func BenchmarkDefaultConfig(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = DefaultConfig()
	}
}

func BenchmarkConfigClone(b *testing.B) {
	cfg := DefaultConfig()
	cfg.TLS = &TLSConfig{Enabled: true}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = cfg.Clone()
	}
}

func BenchmarkApplyOptions(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cfg := DefaultConfig()
		Apply(cfg,
			WithName("test"),
			WithConnectTimeout(5*time.Second),
			WithMaxReconnects(10),
		)
	}
}

func BenchmarkValidate(b *testing.B) {
	cfg := DefaultConfig()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Validate(cfg)
	}
}
