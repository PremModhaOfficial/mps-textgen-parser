package core

import (
	"bytes"
	"sync"
)

// ============================================================================
// Object Pools for Performance Optimization
// ============================================================================

// MessagePool is a pool of Message objects for reuse
type MessagePool struct {
	pool sync.Pool
}

// NewMessagePool creates a new MessagePool
func NewMessagePool() *MessagePool {
	return &MessagePool{
		pool: sync.Pool{
			New: func() any {
				return &Message{
					Headers: make(Headers),
				}
			},
		},
	}
}

// Get retrieves a Message from the pool
func (p *MessagePool) Get() *Message {
	return p.pool.Get().(*Message)
}

// Put returns a Message to the pool after resetting it
func (p *MessagePool) Put(msg *Message) {
	if msg == nil {
		return
	}
	// Reset the message
	msg.Subject = ""
	msg.Data = msg.Data[:0]
	msg.Reply = ""
	msg.ID = ""
	for k := range msg.Headers {
		delete(msg.Headers, k)
	}
	p.pool.Put(msg)
}

// HeadersPool is a pool of Headers objects for reuse
type HeadersPool struct {
	pool sync.Pool
}

// NewHeadersPool creates a new HeadersPool
func NewHeadersPool() *HeadersPool {
	return &HeadersPool{
		pool: sync.Pool{
			New: func() any {
				return make(Headers)
			},
		},
	}
}

// Get retrieves Headers from the pool
func (p *HeadersPool) Get() Headers {
	return p.pool.Get().(Headers)
}

// Put returns Headers to the pool after clearing
func (p *HeadersPool) Put(h Headers) {
	if h == nil {
		return
	}
	for k := range h {
		delete(h, k)
	}
	p.pool.Put(h)
}

// ByteBufferPool is a pool of byte buffers for reuse
type ByteBufferPool struct {
	pool sync.Pool
}

// NewByteBufferPool creates a new ByteBufferPool with the given initial size
func NewByteBufferPool(initialSize int) *ByteBufferPool {
	return &ByteBufferPool{
		pool: sync.Pool{
			New: func() any {
				return bytes.NewBuffer(make([]byte, 0, initialSize))
			},
		},
	}
}

// Get retrieves a buffer from the pool
func (p *ByteBufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put returns a buffer to the pool after resetting
func (p *ByteBufferPool) Put(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	buf.Reset()
	p.pool.Put(buf)
}

// ============================================================================
// Global Default Pools
// ============================================================================

var (
	defaultMessagePool = NewMessagePool()
	defaultHeadersPool = NewHeadersPool()
)

// AcquireMessage gets a Message from the default pool
func AcquireMessage() *Message {
	return defaultMessagePool.Get()
}

// AcquireMessageWithData gets a Message from the pool and sets its data
func AcquireMessageWithData(data []byte) *Message {
	msg := defaultMessagePool.Get()
	msg.Data = append(msg.Data, data...)
	return msg
}

// ReleaseMessage returns a Message to the default pool
func ReleaseMessage(msg *Message) {
	defaultMessagePool.Put(msg)
}

// AcquireHeaders gets Headers from the default pool
func AcquireHeaders() Headers {
	return defaultHeadersPool.Get()
}

// ReleaseHeaders returns Headers to the default pool
func ReleaseHeaders(h Headers) {
	defaultHeadersPool.Put(h)
}
