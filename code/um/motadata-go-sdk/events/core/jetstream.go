package core

import (
	"context"
	"errors"
	"time"
)

// ============================================================================
// JetStream Interfaces - Production-Level Stream Processing
// ============================================================================

// StreamManager manages JetStream streams
type StreamManager interface {
	// CreateStream creates or updates a stream
	CreateStream(ctx context.Context, cfg StreamConfig) (Stream, error)

	// UpdateStream updates an existing stream
	UpdateStream(ctx context.Context, cfg StreamConfig) (Stream, error)

	// DeleteStream deletes a stream
	DeleteStream(ctx context.Context, name string) error

	// GetStream returns a stream by name
	GetStream(ctx context.Context, name string) (Stream, error)

	// ListStreams lists all streams
	ListStreams(ctx context.Context) ([]Stream, error)

	// StreamNames returns all stream names
	StreamNames(ctx context.Context) ([]string, error)
}

// Stream represents a JetStream stream
type Stream interface {
	// Info returns stream information
	Info(ctx context.Context) (*StreamInfo, error)

	// Purge removes all messages from the stream
	Purge(ctx context.Context, opts ...StreamPurgeOption) error

	// DeleteMessage deletes a specific message
	DeleteMessage(ctx context.Context, seq uint64) error

	// GetMessage retrieves a specific message
	GetMessage(ctx context.Context, seq uint64) (*RawStreamMessage, error)

	// GetLastMessageForSubject gets the last message for a subject
	GetLastMessageForSubject(ctx context.Context, subject string) (*RawStreamMessage, error)

	// CreateConsumer creates a consumer for this stream
	CreateConsumer(ctx context.Context, cfg ConsumerConfig) (Consumer, error)

	// UpdateConsumer updates an existing consumer
	UpdateConsumer(ctx context.Context, cfg ConsumerConfig) (Consumer, error)

	// DeleteConsumer deletes a consumer
	DeleteConsumer(ctx context.Context, name string) error

	// Consumer returns a consumer by name
	Consumer(ctx context.Context, name string) (Consumer, error)

	// ListConsumers lists all consumers
	ListConsumers(ctx context.Context) ([]ConsumerInfo, error)

	// ConsumerNames returns all consumer names
	ConsumerNames(ctx context.Context) ([]string, error)
}

// ============================================================================
// Stream Configuration
// ============================================================================

// StreamConfig defines stream configuration
type StreamConfig struct {
	// Name is the stream name (required)
	Name string `json:"name"`

	// Description describes the stream
	Description string `json:"description,omitempty"`

	// Subjects are the subjects this stream listens to
	Subjects []string `json:"subjects"`

	// Retention policy
	Retention RetentionPolicy `json:"retention"`

	// MaxConsumers is the maximum number of consumers
	MaxConsumers int `json:"max_consumers,omitempty"`

	// MaxMsgs is the maximum number of messages
	MaxMsgs int64 `json:"max_msgs,omitempty"`

	// MaxBytes is the maximum total size
	MaxBytes int64 `json:"max_bytes,omitempty"`

	// MaxAge is the maximum age of messages
	MaxAge time.Duration `json:"max_age,omitempty"`

	// MaxMsgSize is the maximum message size
	MaxMsgSize int32 `json:"max_msg_size,omitempty"`

	// Discard policy when limits are reached
	Discard DiscardPolicy `json:"discard"`

	// Storage type
	Storage StorageType `json:"storage"`

	// Replicas is the number of replicas
	Replicas int `json:"num_replicas,omitempty"`

	// NoAck disables acknowledgments
	NoAck bool `json:"no_ack,omitempty"`

	// DuplicateWindow is the window for duplicate detection
	DuplicateWindow time.Duration `json:"duplicate_window,omitempty"`

	// Placement specifies cluster placement
	Placement *Placement `json:"placement,omitempty"`

	// Mirror configures stream mirroring
	Mirror *StreamSource `json:"mirror,omitempty"`

	// Sources configures stream sourcing
	Sources []*StreamSource `json:"sources,omitempty"`

	// AllowRollup allows rollup headers
	AllowRollup bool `json:"allow_rollup_hdrs,omitempty"`

	// DenyDelete denies message deletion
	DenyDelete bool `json:"deny_delete,omitempty"`

	// DenyPurge denies stream purge
	DenyPurge bool `json:"deny_purge,omitempty"`

	// AllowDirect allows direct get
	AllowDirect bool `json:"allow_direct,omitempty"`

	// Metadata contains stream metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// RetentionPolicy defines message retention
type RetentionPolicy int

const (
	// LimitsPolicy retains messages until limits are reached
	LimitsPolicy RetentionPolicy = iota
	// InterestPolicy retains messages while there are consumers
	InterestPolicy
	// WorkQueuePolicy for work queue semantics
	WorkQueuePolicy
)

// DiscardPolicy defines what to discard when limits are reached
type DiscardPolicy int

const (
	// DiscardOld discards old messages
	DiscardOld DiscardPolicy = iota
	// DiscardNew discards new messages
	DiscardNew
)

// StorageType defines storage backend
type StorageType int

const (
	// FileStorage stores to disk
	FileStorage StorageType = iota
	// MemoryStorage stores in memory
	MemoryStorage
)

// Placement specifies stream placement
type Placement struct {
	Cluster string   `json:"cluster,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

// StreamSource defines a source stream for mirroring/sourcing
type StreamSource struct {
	Name          string         `json:"name"`
	OptStartSeq   uint64         `json:"opt_start_seq,omitempty"`
	OptStartTime  *time.Time     `json:"opt_start_time,omitempty"`
	FilterSubject string         `json:"filter_subject,omitempty"`
	External      *ExternalStream `json:"external,omitempty"`
}

// ExternalStream defines an external stream reference
type ExternalStream struct {
	APIPrefix     string `json:"api,omitempty"`
	DeliverPrefix string `json:"deliver,omitempty"`
}

// ============================================================================
// Stream Information
// ============================================================================

// StreamInfo contains stream state information
type StreamInfo struct {
	Config  StreamConfig `json:"config"`
	Created time.Time    `json:"created"`
	State   StreamState  `json:"state"`
	Cluster *ClusterInfo `json:"cluster,omitempty"`
	Mirror  *StreamSourceInfo `json:"mirror,omitempty"`
	Sources []*StreamSourceInfo `json:"sources,omitempty"`
}

// StreamState contains stream state
type StreamState struct {
	Messages   uint64     `json:"messages"`
	Bytes      uint64     `json:"bytes"`
	FirstSeq   uint64     `json:"first_seq"`
	FirstTime  time.Time  `json:"first_ts"`
	LastSeq    uint64     `json:"last_seq"`
	LastTime   time.Time  `json:"last_ts"`
	Consumers  int        `json:"consumer_count"`
	Deleted    []uint64   `json:"deleted,omitempty"`
	NumDeleted int        `json:"num_deleted,omitempty"`
	NumSubjects uint64    `json:"num_subjects,omitempty"`
	Subjects   map[string]uint64 `json:"subjects,omitempty"`
}

// ClusterInfo contains cluster information
type ClusterInfo struct {
	Name     string        `json:"name,omitempty"`
	Leader   string        `json:"leader,omitempty"`
	Replicas []*PeerInfo   `json:"replicas,omitempty"`
}

// PeerInfo contains peer information
type PeerInfo struct {
	Name    string        `json:"name"`
	Current bool          `json:"current"`
	Active  time.Duration `json:"active"`
	Offline bool          `json:"offline,omitempty"`
	Lag     uint64        `json:"lag,omitempty"`
}

// StreamSourceInfo contains source stream information
type StreamSourceInfo struct {
	Name   string        `json:"name"`
	Lag    uint64        `json:"lag"`
	Active time.Duration `json:"active"`
	Error  *APIError     `json:"error,omitempty"`
}

// APIError represents a JetStream API error
type APIError struct {
	Code        int    `json:"code"`
	ErrorCode   uint16 `json:"err_code,omitempty"`
	Description string `json:"description"`
}

// RawStreamMessage represents a raw stream message
type RawStreamMessage struct {
	Subject  string
	Sequence uint64
	Header   Headers
	Data     []byte
	Time     time.Time
}

// StreamPurgeOption configures stream purge
type StreamPurgeOption func(*streamPurgeConfig)

type streamPurgeConfig struct {
	subject string
	seq     uint64
	keep    uint64
}

// WithPurgeSubject limits purge to a specific subject
func WithPurgeSubject(subject string) StreamPurgeOption {
	return func(c *streamPurgeConfig) {
		c.subject = subject
	}
}

// WithPurgeSequence purges messages up to a sequence
func WithPurgeSequence(seq uint64) StreamPurgeOption {
	return func(c *streamPurgeConfig) {
		c.seq = seq
	}
}

// WithPurgeKeep keeps a number of messages
func WithPurgeKeep(keep uint64) StreamPurgeOption {
	return func(c *streamPurgeConfig) {
		c.keep = keep
	}
}

// ============================================================================
// Consumer Configuration
// ============================================================================

// ConsumerConfig defines consumer configuration
type ConsumerConfig struct {
	// Name is the consumer name (required for durable consumers)
	Name string `json:"name,omitempty"`

	// Durable is the durable name (deprecated, use Name)
	Durable string `json:"durable_name,omitempty"`

	// Description describes the consumer
	Description string `json:"description,omitempty"`

	// DeliverPolicy defines where to start delivering
	DeliverPolicy DeliverPolicy `json:"deliver_policy"`

	// OptStartSeq is the starting sequence (for DeliverByStartSequence)
	OptStartSeq uint64 `json:"opt_start_seq,omitempty"`

	// OptStartTime is the starting time (for DeliverByStartTime)
	OptStartTime *time.Time `json:"opt_start_time,omitempty"`

	// AckPolicy defines acknowledgment policy
	AckPolicy AckPolicy `json:"ack_policy"`

	// AckWait is the acknowledgment wait time
	AckWait time.Duration `json:"ack_wait,omitempty"`

	// MaxDeliver is the maximum redelivery attempts
	MaxDeliver int `json:"max_deliver,omitempty"`

	// BackOff defines redelivery backoff
	BackOff []time.Duration `json:"backoff,omitempty"`

	// FilterSubject filters messages by subject
	FilterSubject string `json:"filter_subject,omitempty"`

	// FilterSubjects filters by multiple subjects
	FilterSubjects []string `json:"filter_subjects,omitempty"`

	// ReplayPolicy defines replay speed
	ReplayPolicy ReplayPolicy `json:"replay_policy"`

	// RateLimit limits message rate (bits per second)
	RateLimit uint64 `json:"rate_limit_bps,omitempty"`

	// SampleFrequency for sampling (e.g., "50%")
	SampleFrequency string `json:"sample_freq,omitempty"`

	// MaxWaiting is max pending pull requests
	MaxWaiting int `json:"max_waiting,omitempty"`

	// MaxAckPending is max unacknowledged messages
	MaxAckPending int `json:"max_ack_pending,omitempty"`

	// HeadersOnly delivers only headers
	HeadersOnly bool `json:"headers_only,omitempty"`

	// MaxRequestBatch is max batch size per pull request
	MaxRequestBatch int `json:"max_batch,omitempty"`

	// MaxRequestExpires is max duration for pull requests
	MaxRequestExpires time.Duration `json:"max_expires,omitempty"`

	// MaxRequestMaxBytes is max bytes per pull request
	MaxRequestMaxBytes int `json:"max_bytes,omitempty"`

	// InactiveThreshold removes inactive consumers
	InactiveThreshold time.Duration `json:"inactive_threshold,omitempty"`

	// Replicas is the number of replicas
	Replicas int `json:"num_replicas,omitempty"`

	// MemoryStorage uses memory storage
	MemoryStorage bool `json:"mem_storage,omitempty"`

	// Metadata contains consumer metadata
	Metadata map[string]string `json:"metadata,omitempty"`
}

// DeliverPolicy defines where to start delivering
type DeliverPolicy int

const (
	// DeliverAll delivers all messages
	DeliverAll DeliverPolicy = iota
	// DeliverLast delivers the last message
	DeliverLast
	// DeliverNew delivers only new messages
	DeliverNew
	// DeliverByStartSequence delivers from a sequence
	DeliverByStartSequence
	// DeliverByStartTime delivers from a time
	DeliverByStartTime
	// DeliverLastPerSubject delivers last per subject
	DeliverLastPerSubject
)

// AckPolicy defines acknowledgment behavior
type AckPolicy int

const (
	// AckNone requires no acknowledgment
	AckNone AckPolicy = iota
	// AckAll acknowledges all previous messages
	AckAll
	// AckExplicit requires explicit acknowledgment
	AckExplicit
)

// ReplayPolicy defines replay behavior
type ReplayPolicy int

const (
	// ReplayInstant replays as fast as possible
	ReplayInstant ReplayPolicy = iota
	// ReplayOriginal replays at original speed
	ReplayOriginal
)

// ============================================================================
// Consumer Interface - Pull Consumer Pattern
// ============================================================================

// Consumer represents a JetStream consumer
type Consumer interface {
	// Info returns consumer information
	Info(ctx context.Context) (*ConsumerInfo, error)

	// Fetch fetches messages (blocking)
	Fetch(batch int, opts ...FetchOption) (MessageBatch, error)

	// FetchBytes fetches messages up to max bytes
	FetchBytes(maxBytes int, opts ...FetchOption) (MessageBatch, error)

	// FetchNoWait fetches available messages without waiting
	FetchNoWait(batch int) (MessageBatch, error)

	// Next fetches a single message
	Next(opts ...FetchOption) (Msg, error)

	// Messages returns a continuous message iterator
	Messages(opts ...MessagesOption) (MessageIterator, error)

	// Consume starts consuming messages with a handler
	Consume(handler MsgHandler, opts ...ConsumeOption) (ConsumeContext, error)
}

// ConsumerInfo contains consumer state
type ConsumerInfo struct {
	Stream         string         `json:"stream_name"`
	Name           string         `json:"name"`
	Created        time.Time      `json:"created"`
	Config         ConsumerConfig `json:"config"`
	Delivered      SequenceInfo   `json:"delivered"`
	AckFloor       SequenceInfo   `json:"ack_floor"`
	NumAckPending  int            `json:"num_ack_pending"`
	NumRedelivered int            `json:"num_redelivered"`
	NumWaiting     int            `json:"num_waiting"`
	NumPending     uint64         `json:"num_pending"`
	Cluster        *ClusterInfo   `json:"cluster,omitempty"`
	PushBound      bool           `json:"push_bound,omitempty"`
}

// SequenceInfo contains sequence information
type SequenceInfo struct {
	Consumer uint64     `json:"consumer_seq"`
	Stream   uint64     `json:"stream_seq"`
	Last     *time.Time `json:"last_active,omitempty"`
}

// ============================================================================
// Message Interfaces
// ============================================================================

// Msg represents a JetStream message
type Msg interface {
	// Data returns the message data
	Data() []byte

	// Headers returns the message headers
	Headers() Headers

	// Subject returns the message subject
	Subject() string

	// Reply returns the reply subject
	Reply() string

	// Metadata returns JetStream metadata
	Metadata() (*MsgMetadata, error)

	// Ack acknowledges the message
	Ack() error

	// DoubleAck acknowledges and waits for confirmation
	DoubleAck(ctx context.Context) error

	// Nak negative-acknowledges the message
	Nak() error

	// NakWithDelay negative-acknowledges with delay
	NakWithDelay(delay time.Duration) error

	// Term terminates the message (no redelivery)
	Term() error

	// TermWithReason terminates with a reason
	TermWithReason(reason string) error

	// InProgress signals the message is being processed
	InProgress() error
}

// MsgMetadata contains JetStream message metadata
type MsgMetadata struct {
	Sequence     SequencePair
	NumDelivered uint64
	NumPending   uint64
	Timestamp    time.Time
	Stream       string
	Consumer     string
	Domain       string
}

// SequencePair contains consumer and stream sequence
type SequencePair struct {
	Consumer uint64
	Stream   uint64
}

// MsgHandler handles consumed messages
type MsgHandler func(msg Msg)

// ============================================================================
// Message Batch
// ============================================================================

// MessageBatch represents a batch of messages
type MessageBatch interface {
	// Messages returns a channel of messages
	Messages() <-chan Msg

	// Error returns any error that occurred
	Error() error
}

// ============================================================================
// Message Iterator
// ============================================================================

// MessageIterator provides continuous message iteration
type MessageIterator interface {
	// Next returns the next message
	Next() (Msg, error)

	// Stop stops the iterator
	Stop()
}

// MessagesOption configures message iteration
type MessagesOption func(*messagesConfig)

type messagesConfig struct {
	stopAfter     int
	stopAfterErr  error
	heartbeat     time.Duration
	withoutDelete bool
}

// WithStopAfter stops after n messages
func WithStopAfter(n int) MessagesOption {
	return func(c *messagesConfig) {
		c.stopAfter = n
	}
}

// WithIteratorHeartbeat sets the heartbeat interval
func WithIteratorHeartbeat(d time.Duration) MessagesOption {
	return func(c *messagesConfig) {
		c.heartbeat = d
	}
}

// ============================================================================
// Consume Context
// ============================================================================

// ConsumeContext provides control over message consumption
type ConsumeContext interface {
	// Stop stops consuming
	Stop()

	// Drain drains pending messages and stops
	Drain()
}

// ConsumeOption configures message consumption
type ConsumeOption func(*consumeConfig)

type consumeConfig struct {
	errHandler     ConsumeErrorHandler
	maxMessages    int
	maxBytes       int
	expires        time.Duration
	heartbeat      time.Duration
	notifyHandler  func()
	stopAfter      int
	stopAfterErr   error
	thresholdPct   int
	thresholdMsgs  int
	thresholdBytes int
}

// ConsumeErrorHandler handles consume errors
type ConsumeErrorHandler func(consumeCtx ConsumeContext, err error)

// WithConsumeErrorHandler sets the error handler
func WithConsumeErrorHandler(handler ConsumeErrorHandler) ConsumeOption {
	return func(c *consumeConfig) {
		c.errHandler = handler
	}
}

// WithConsumeMaxMessages sets max messages per batch
func WithConsumeMaxMessages(max int) ConsumeOption {
	return func(c *consumeConfig) {
		c.maxMessages = max
	}
}

// WithConsumeExpiry sets pull request expiry
func WithConsumeExpiry(d time.Duration) ConsumeOption {
	return func(c *consumeConfig) {
		c.expires = d
	}
}

// WithConsumeHeartbeat sets heartbeat interval
func WithConsumeHeartbeat(d time.Duration) ConsumeOption {
	return func(c *consumeConfig) {
		c.heartbeat = d
	}
}

// WithConsumeStopAfter stops after n messages
func WithConsumeStopAfter(n int) ConsumeOption {
	return func(c *consumeConfig) {
		c.stopAfter = n
	}
}

// ============================================================================
// Fetch Options
// ============================================================================

// FetchOption configures a fetch operation
type FetchOption func(*fetchConfig)

type fetchConfig struct {
	maxWait   time.Duration
	heartbeat time.Duration
}

// WithFetchMaxWait sets max wait time
func WithFetchMaxWait(d time.Duration) FetchOption {
	return func(c *fetchConfig) {
		c.maxWait = d
	}
}

// WithFetchHeartbeat sets heartbeat interval
func WithFetchHeartbeat(d time.Duration) FetchOption {
	return func(c *fetchConfig) {
		c.heartbeat = d
	}
}

// ============================================================================
// JetStream Errors
// ============================================================================

var (
	// ErrStreamNotFound indicates the stream does not exist
	ErrJetStreamNotEnabled = errors.New("jetstream not enabled")

	// ErrConsumerNotFound indicates the consumer does not exist
	ErrConsumerNotFound = errors.New("consumer not found")

	// ErrMsgNotFound indicates the message does not exist
	ErrMsgNotFound = errors.New("message not found")

	// ErrBadRequest indicates an invalid request
	ErrBadRequest = errors.New("bad request")

	// ErrConsumerNameRequired indicates consumer name is required
	ErrConsumerNameRequired = errors.New("consumer name required")

	// ErrStreamNameRequired indicates stream name is required
	ErrStreamNameRequired = errors.New("stream name required")

	// ErrMsgAlreadyAckd indicates the message was already acknowledged
	ErrMsgAlreadyAckd = errors.New("message already acknowledged")

	// ErrNoMessages indicates no messages are available
	ErrNoMessages = errors.New("no messages")

	// ErrMaxBytesExceeded indicates max bytes was exceeded
	ErrMaxBytesExceeded = errors.New("max bytes exceeded")

	// ErrOrderedConsumerReset indicates an ordered consumer reset
	ErrOrderedConsumerReset = errors.New("ordered consumer reset")
)
