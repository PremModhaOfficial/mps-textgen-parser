package events

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// Default health monitor constants control the timing and capacity of health
// checks and round-trip time (RTT) measurement.
const (
	// defaultHealthCheckInterval is the default interval between health checks.
	// Each check verifies connection state, collects NATS statistics, and
	// measures RTT to the server.
	defaultHealthCheckInterval = 10 * time.Second

	// defaultRTTSamples is the default number of RTT samples retained in the
	// circular buffer for average/min/max calculations.
	defaultRTTSamples = 10

	// rttFlushTimeout is the timeout for the NATS Flush call used as the RTT
	// measurement probe. A flush round-trips to the server and back, providing
	// a reasonable approximation of network latency.
	rttFlushTimeout = 5 * time.Second
)

// HealthMonitor periodically checks connection health, collects NATS statistics
// (message counts, byte throughput, reconnections), and measures round-trip time
// (RTT) to the server. It maintains a circular buffer of RTT samples for
// rolling average, min, and max calculations.
//
// Health state changes (healthy <-> unhealthy) trigger the onHealthChange callback,
// which can be used for alerting or metrics. The monitor runs a background goroutine
// started by Start and stopped by Stop.
type HealthMonitor struct {
	conn *Connection // conn is the monitored NATS connection.

	// Configuration
	checkInterval time.Duration // checkInterval is the time between health check ticks.
	rttSamples    int           // rttSamples is the capacity of the RTT circular buffer.

	// Statistics
	stats   ConnectionStats // stats holds the latest aggregated connection statistics.
	statsMu sync.RWMutex    // statsMu protects stats and rttHist from concurrent access.
	rttHist []time.Duration // rttHist is a circular buffer of recent RTT measurements.
	rttIdx  int             // rttIdx is the write cursor into the rttHist circular buffer.

	// Lifecycle
	done   <-chan struct{}     // done is closed when cancel is called, terminating the monitor loop.
	cancel context.CancelFunc // cancel terminates the background goroutine.
	wg     sync.WaitGroup     // wg tracks the monitor goroutine for graceful shutdown.

	// Callbacks
	onHealthChange func(healthy bool, status HealthStatus) // onHealthChange fires on health state transitions.
	lastHealthy    atomic.Bool                             // lastHealthy stores the previous health state for change detection.
}

// ConnectionStats holds aggregated connection statistics collected by the
// HealthMonitor from the underlying NATS connection. Statistics are updated
// on each health check tick and are accessible via HealthMonitor.Stats().
type ConnectionStats struct {
	// Counters track cumulative message and error counts since connection.
	MessagesPublished  uint64 // MessagesPublished is the total number of messages sent.
	MessagesReceived   uint64 // MessagesReceived is the total number of messages received.
	BytesSent          uint64 // BytesSent is the total outbound bytes (payload + headers).
	BytesReceived      uint64 // BytesReceived is the total inbound bytes.
	Reconnects         uint64 // Reconnects is the number of automatic reconnections.
	Errors             uint64 // Errors is the total error count across all categories.
	PublishErrors      uint64 // PublishErrors counts publish-specific failures.
	SubscriptionErrors uint64 // SubscriptionErrors counts subscription-specific failures.

	// Latency holds RTT statistics derived from the circular buffer.
	LastRTT    time.Duration // LastRTT is the most recent RTT measurement.
	AvgRTT     time.Duration // AvgRTT is the rolling average across buffered samples.
	MinRTT     time.Duration // MinRTT is the lowest RTT observed since startup.
	MaxRTT     time.Duration // MaxRTT is the highest RTT observed since startup.
	RTTSamples int           // RTTSamples is the total number of RTT samples collected.

	// Connection info captures server and uptime metadata.
	ConnectedAt    time.Time // ConnectedAt is the timestamp of the initial connection.
	LastDisconnect time.Time // LastDisconnect is the timestamp of the most recent disconnection.
	LastReconnect  time.Time // LastReconnect is the timestamp of the most recent reconnection.
	UptimeSeconds  float64   // UptimeSeconds is the elapsed time since ConnectedAt.
	CurrentServer  string    // CurrentServer is the URL of the currently connected NATS server.

	// Health summarises the connection health at the last check.
	Healthy         bool      // Healthy is true if the connection passed the last health check.
	LastHealthCheck time.Time // LastHealthCheck is the timestamp of the most recent check.
}

// HealthMonitorConfig configures the health monitor's check frequency and
// RTT buffer capacity.
type HealthMonitorConfig struct {
	CheckInterval time.Duration // CheckInterval sets how often health checks run (default: 10s).
	RTTSamples    int           // RTTSamples sets the circular buffer size for RTT history (default: 10).
}

// DefaultHealthMonitorConfig returns a configuration with sensible defaults:
// 10-second check interval and 10 RTT samples.
func DefaultHealthMonitorConfig() HealthMonitorConfig {
	return HealthMonitorConfig{
		CheckInterval: defaultHealthCheckInterval,
		RTTSamples:    defaultRTTSamples,
	}
}

// NewHealthMonitor creates a new HealthMonitor for the given Connection.
// Zero-valued configuration fields are replaced with defaults. The monitor
// does not start automatically; call Start to begin the background check loop.
// MinRTT is initialised to time.Hour so the first real measurement always wins.
func NewHealthMonitor(conn *Connection, cfg HealthMonitorConfig) *HealthMonitor {
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = defaultHealthCheckInterval
	}
	if cfg.RTTSamples == 0 {
		cfg.RTTSamples = defaultRTTSamples
	}

	ctx, cancel := context.WithCancel(context.Background())
	hm := &HealthMonitor{
		conn:          conn,
		checkInterval: cfg.CheckInterval,
		rttSamples:    cfg.RTTSamples,
		rttHist:       make([]time.Duration, cfg.RTTSamples),
		done:          ctx.Done(),
		cancel:        cancel,
	}
	hm.stats.MinRTT = time.Hour // Initialize to high value

	return hm
}

// Start launches the background health check goroutine. It ticks at the
// configured CheckInterval and runs until Stop is called.
func (hm *HealthMonitor) Start() {
	hm.wg.Add(1)
	go hm.monitorLoop()
}

// Stop cancels the background goroutine and blocks until it exits.
func (hm *HealthMonitor) Stop() {
	hm.cancel()
	hm.wg.Wait()
}

// OnHealthChange sets a callback that fires on health state transitions
// (healthy -> unhealthy or vice versa). The callback receives the new health
// state and the full HealthStatus. Only transitions trigger the callback;
// consecutive checks with the same state are suppressed.
func (hm *HealthMonitor) OnHealthChange(cb func(healthy bool, status HealthStatus)) {
	hm.onHealthChange = cb
}

// Stats returns a snapshot of the current connection statistics. The returned
// struct is a copy, so callers may inspect it without holding the lock.
func (hm *HealthMonitor) Stats() ConnectionStats {
	hm.statsMu.RLock()
	defer hm.statsMu.RUnlock()
	return hm.stats
}

// IsHealthy returns the most recent health state using an atomic load,
// making it safe to call from any goroutine without locking.
func (hm *HealthMonitor) IsHealthy() bool {
	return hm.lastHealthy.Load()
}

// monitorLoop runs the periodic health check loop on a background goroutine.
// It ticks at checkInterval and terminates when the context is cancelled.
func (hm *HealthMonitor) monitorLoop() {
	defer hm.wg.Done()

	ticker := time.NewTicker(hm.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-hm.done:
			return
		case <-ticker.C:
			hm.performHealthCheck()
		}
	}
}

// performHealthCheck runs a single health check cycle: retrieves the connection's
// HealthStatus, updates NATS statistics, measures RTT if connected, detects
// health state changes, and fires the onHealthChange callback if the state changed.
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

// updateStats pulls current statistics (message counts, byte throughput,
// reconnections, server URL) from the underlying nats.Conn and stores them
// in the stats struct. Also recalculates uptime from ConnectedAt.
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

// measureRTT measures round-trip time to the NATS server using a flush probe.
// The result is stored in a circular buffer (rttHist) and used to compute
// rolling average, min, and max RTT. Failed flush attempts are silently
// discarded to avoid polluting the statistics.
func (hm *HealthMonitor) measureRTT() {
	nc := hm.conn.Conn()
	if nc == nil {
		return
	}

	// Use NATS flush with timeout as RTT measurement
	start := time.Now()
	err := nc.FlushTimeout(rttFlushTimeout)
	if err != nil {
		return // Do not record failed RTT
	}
	rtt := time.Since(start)

	hm.statsMu.Lock()
	defer hm.statsMu.Unlock()

	// Update RTT history
	hm.rttHist[hm.rttIdx] = rtt
	hm.rttIdx = (hm.rttIdx + 1) % len(hm.rttHist)

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

// RecordPublish increments the published message count and adds the given
// byte count to BytesSent. Thread-safe for concurrent publishers.
func (hm *HealthMonitor) RecordPublish(bytes int) {
	hm.statsMu.Lock()
	hm.stats.MessagesPublished++
	hm.stats.BytesSent += uint64(bytes)
	hm.statsMu.Unlock()
}

// RecordReceive increments the received message count and adds the given
// byte count to BytesReceived. Thread-safe for concurrent subscribers.
func (hm *HealthMonitor) RecordReceive(bytes int) {
	hm.statsMu.Lock()
	hm.stats.MessagesReceived++
	hm.stats.BytesReceived += uint64(bytes)
	hm.statsMu.Unlock()
}

// RecordError increments the total error counter and the category-specific
// counter based on errType ("publish" or "subscription"). Unrecognised
// error types only increment the total.
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

// RecordConnect records the timestamp of a successful connection event.
// Called by the Connection's connected handler to track uptime.
func (hm *HealthMonitor) RecordConnect() {
	hm.statsMu.Lock()
	hm.stats.ConnectedAt = time.Now()
	hm.statsMu.Unlock()
}

// RecordDisconnect records the timestamp of a disconnection event.
// Called by the Connection's disconnected handler.
func (hm *HealthMonitor) RecordDisconnect() {
	hm.statsMu.Lock()
	hm.stats.LastDisconnect = time.Now()
	hm.statsMu.Unlock()
}

// RecordReconnect records the timestamp of a reconnection event and increments
// the reconnect counter. Called by the Connection's reconnected handler.
func (hm *HealthMonitor) RecordReconnect() {
	hm.statsMu.Lock()
	hm.stats.LastReconnect = time.Now()
	hm.stats.Reconnects++
	hm.statsMu.Unlock()
}

// HealthChecker provides a minimal health check interface suitable for
// integration with health check frameworks (e.g., Kubernetes liveness probes).
type HealthChecker interface {
	// Check verifies the connection is alive and returns an error if unhealthy.
	Check(ctx context.Context) error
}

// SimpleHealthChecker implements HealthChecker by verifying the NATS connection
// is active and performing a flush round-trip to confirm server reachability.
type SimpleHealthChecker struct {
	conn *Connection // conn is the NATS connection to check.
}

// NewHealthChecker creates a SimpleHealthChecker that checks the given
// Connection's liveness via flush probe.
func NewHealthChecker(conn *Connection) *SimpleHealthChecker {
	return &SimpleHealthChecker{conn: conn}
}

// Check verifies the connection is alive by checking the connection state and
// performing a NATS flush with a 5-second timeout. The flush round-trips to
// the server, confirming both network and server liveness. Returns
// ErrNotConnected if the connection is down, or the flush error if the
// round-trip fails.
func (hc *SimpleHealthChecker) Check(_ context.Context) error {
	if !hc.conn.IsConnected() {
		return utils.ErrNotConnected
	}

	nc := hc.conn.Conn()
	if nc == nil {
		return utils.ErrNotConnected
	}

	// Try to flush to verify connection is alive
	return nc.FlushTimeout(5 * time.Second)
}
