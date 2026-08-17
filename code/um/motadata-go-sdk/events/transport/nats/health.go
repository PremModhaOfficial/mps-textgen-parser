package nats

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// HealthMonitor monitors connection health and collects statistics
type HealthMonitor struct {
	conn *Connection

	// Configuration
	checkInterval time.Duration
	rttSamples    int

	// Statistics
	stats    ConnectionStats
	statsMu  sync.RWMutex
	rttHist  []time.Duration
	rttIndex int

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Callbacks
	onHealthChange func(healthy bool, status core.HealthStatus)
	lastHealthy    atomic.Bool
}

// ConnectionStats holds connection statistics
type ConnectionStats struct {
	// Counters
	MessagesPublished   uint64
	MessagesReceived    uint64
	BytesSent           uint64
	BytesReceived       uint64
	Reconnects          uint64
	Errors              uint64
	PublishErrors       uint64
	SubscriptionErrors  uint64

	// Latency
	LastRTT    time.Duration
	AvgRTT     time.Duration
	MinRTT     time.Duration
	MaxRTT     time.Duration
	RTTSamples int

	// Connection info
	ConnectedAt      time.Time
	LastDisconnect   time.Time
	LastReconnect    time.Time
	UptimeSeconds    float64
	CurrentServer    string

	// Health
	Healthy         bool
	LastHealthCheck time.Time
}

// HealthMonitorConfig configures the health monitor
type HealthMonitorConfig struct {
	CheckInterval time.Duration // How often to check health (default: 10s)
	RTTSamples    int           // Number of RTT samples to keep (default: 10)
}

// DefaultHealthMonitorConfig returns default configuration
func DefaultHealthMonitorConfig() HealthMonitorConfig {
	return HealthMonitorConfig{
		CheckInterval: 10 * time.Second,
		RTTSamples:    10,
	}
}

// NewHealthMonitor creates a new health monitor for a connection
func NewHealthMonitor(conn *Connection, cfg HealthMonitorConfig) *HealthMonitor {
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = 10 * time.Second
	}
	if cfg.RTTSamples == 0 {
		cfg.RTTSamples = 10
	}

	ctx, cancel := context.WithCancel(context.Background())
	hm := &HealthMonitor{
		conn:          conn,
		checkInterval: cfg.CheckInterval,
		rttSamples:    cfg.RTTSamples,
		rttHist:       make([]time.Duration, cfg.RTTSamples),
		ctx:           ctx,
		cancel:        cancel,
	}
	hm.stats.MinRTT = time.Hour // Initialize to high value

	return hm
}

// Start begins health monitoring
func (hm *HealthMonitor) Start() {
	hm.wg.Add(1)
	go hm.monitorLoop()
}

// Stop stops health monitoring
func (hm *HealthMonitor) Stop() {
	hm.cancel()
	hm.wg.Wait()
}

// OnHealthChange sets a callback for health state changes
func (hm *HealthMonitor) OnHealthChange(cb func(healthy bool, status core.HealthStatus)) {
	hm.onHealthChange = cb
}

// Stats returns current connection statistics
func (hm *HealthMonitor) Stats() ConnectionStats {
	hm.statsMu.RLock()
	defer hm.statsMu.RUnlock()
	return hm.stats
}

// IsHealthy returns current health status
func (hm *HealthMonitor) IsHealthy() bool {
	return hm.lastHealthy.Load()
}

// monitorLoop runs the health check loop
func (hm *HealthMonitor) monitorLoop() {
	defer hm.wg.Done()

	ticker := time.NewTicker(hm.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hm.ctx.Done():
			return
		case <-ticker.C:
			hm.performHealthCheck()
		}
	}
}

// performHealthCheck runs a health check
func (hm *HealthMonitor) performHealthCheck() {
	status := hm.conn.Health()

	// Update stats from NATS connection
	hm.updateStats()

	// Measure RTT if connected
	if hm.conn.IsConnected() {
		hm.measureRTT()
	}

	// Check for health state change
	wasHealthy := hm.lastHealthy.Load()
	nowHealthy := status.Healthy

	if wasHealthy != nowHealthy {
		hm.lastHealthy.Store(nowHealthy)
		if hm.onHealthChange != nil {
			hm.onHealthChange(nowHealthy, status)
		}
	}

	hm.statsMu.Lock()
	hm.stats.Healthy = nowHealthy
	hm.stats.LastHealthCheck = time.Now()
	hm.statsMu.Unlock()
}

// updateStats updates statistics from the NATS connection
func (hm *HealthMonitor) updateStats() {
	nc := hm.conn.Conn()
	if nc == nil {
		return
	}

	natsStats := nc.Stats()

	hm.statsMu.Lock()
	defer hm.statsMu.Unlock()

	hm.stats.MessagesPublished = natsStats.OutMsgs
	hm.stats.MessagesReceived = natsStats.InMsgs
	hm.stats.BytesSent = natsStats.OutBytes
	hm.stats.BytesReceived = natsStats.InBytes
	hm.stats.Reconnects = natsStats.Reconnects
	hm.stats.CurrentServer = nc.ConnectedUrl()

	// Calculate uptime
	if !hm.stats.ConnectedAt.IsZero() {
		hm.stats.UptimeSeconds = time.Since(hm.stats.ConnectedAt).Seconds()
	}
}

// measureRTT measures round-trip time to the server
func (hm *HealthMonitor) measureRTT() {
	nc := hm.conn.Conn()
	if nc == nil {
		return
	}

	// Use NATS flush with timeout as RTT measurement
	start := time.Now()
	err := nc.FlushTimeout(5 * time.Second)
	if err != nil {
		return // Don't record failed RTT
	}
	rtt := time.Since(start)

	hm.statsMu.Lock()
	defer hm.statsMu.Unlock()

	// Update RTT history
	hm.rttHist[hm.rttIndex] = rtt
	hm.rttIndex = (hm.rttIndex + 1) % len(hm.rttHist)

	// Update stats
	hm.stats.LastRTT = rtt
	hm.stats.RTTSamples++

	// Update min/max
	if rtt < hm.stats.MinRTT {
		hm.stats.MinRTT = rtt
	}
	if rtt > hm.stats.MaxRTT {
		hm.stats.MaxRTT = rtt
	}

	// Calculate average
	var total time.Duration
	count := 0
	for _, r := range hm.rttHist {
		if r > 0 {
			total += r
			count++
		}
	}
	if count > 0 {
		hm.stats.AvgRTT = total / time.Duration(count)
	}
}

// RecordPublish records a publish operation
func (hm *HealthMonitor) RecordPublish(bytes int) {
	hm.statsMu.Lock()
	hm.stats.MessagesPublished++
	hm.stats.BytesSent += uint64(bytes)
	hm.statsMu.Unlock()
}

// RecordReceive records a receive operation
func (hm *HealthMonitor) RecordReceive(bytes int) {
	hm.statsMu.Lock()
	hm.stats.MessagesReceived++
	hm.stats.BytesReceived += uint64(bytes)
	hm.statsMu.Unlock()
}

// RecordError records an error
func (hm *HealthMonitor) RecordError(errType string) {
	hm.statsMu.Lock()
	hm.stats.Errors++
	switch errType {
	case "publish":
		hm.stats.PublishErrors++
	case "subscription":
		hm.stats.SubscriptionErrors++
	}
	hm.statsMu.Unlock()
}

// RecordConnect records a connection event
func (hm *HealthMonitor) RecordConnect() {
	hm.statsMu.Lock()
	hm.stats.ConnectedAt = time.Now()
	hm.statsMu.Unlock()
}

// RecordDisconnect records a disconnection event
func (hm *HealthMonitor) RecordDisconnect() {
	hm.statsMu.Lock()
	hm.stats.LastDisconnect = time.Now()
	hm.statsMu.Unlock()
}

// RecordReconnect records a reconnection event
func (hm *HealthMonitor) RecordReconnect() {
	hm.statsMu.Lock()
	hm.stats.LastReconnect = time.Now()
	hm.stats.Reconnects++
	hm.statsMu.Unlock()
}

// HealthChecker provides a simple health check interface
type HealthChecker interface {
	Check(ctx context.Context) error
}

// SimpleHealthChecker implements HealthChecker
type SimpleHealthChecker struct {
	conn *Connection
}

// NewHealthChecker creates a simple health checker
func NewHealthChecker(conn *Connection) *SimpleHealthChecker {
	return &SimpleHealthChecker{conn: conn}
}

// Check performs a health check
func (hc *SimpleHealthChecker) Check(ctx context.Context) error {
	if !hc.conn.IsConnected() {
		return core.ErrNotConnected
	}

	nc := hc.conn.Conn()
	if nc == nil {
		return core.ErrNotConnected
	}

	// Try to flush to verify connection is alive
	return nc.FlushTimeout(5 * time.Second)
}
