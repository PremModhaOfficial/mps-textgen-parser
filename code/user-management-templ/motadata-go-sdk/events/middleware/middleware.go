package middleware

import (
	"context"

	"github.com/nats-io/nats.go"
)

// PublishMiddleware is a higher-order function that wraps a PublishHandler to intercept
// outgoing publish operations. It enables cross-cutting concerns (tracing, metrics,
// retry, rate limiting) to be applied transparently without modifying business logic.
type PublishMiddleware func(PublishHandler) PublishHandler

// PublishHandler is the function signature for handling an outgoing NATS publish operation.
// Middleware chains terminate at a concrete handler that performs the actual publish.
type PublishHandler func(ctx context.Context, msg *nats.Msg) error

// SubscribeMiddleware is a higher-order function that wraps a SubscribeHandler to intercept
// incoming message processing. It allows cross-cutting concerns to be applied to the
// subscriber side, mirroring the publish middleware pattern.
type SubscribeMiddleware func(SubscribeHandler) SubscribeHandler

// SubscribeHandler is the function signature for processing an incoming NATS message.
// Middleware chains terminate at a concrete handler that performs the actual message processing.
type SubscribeHandler func(ctx context.Context, msg *nats.Msg) error

// Chain composes multiple publish middlewares into a single PublishMiddleware.
// Middlewares are applied in the order provided: the first middleware in the list
// is the outermost wrapper (executed first on the way in, last on the way out).
// Internally, the list is iterated in reverse so that wrapping produces the correct
// execution order, forming a layered onion architecture around the final handler.
func Chain(middlewares ...PublishMiddleware) PublishMiddleware {
	return func(nextHandler PublishHandler) PublishHandler {
		for middlewareIndex := len(middlewares) - 1; middlewareIndex >= 0; middlewareIndex-- {
			nextHandler = middlewares[middlewareIndex](nextHandler)
		}
		return nextHandler
	}
}

// ChainSubscribe composes multiple subscribe middlewares into a single SubscribeMiddleware.
// The execution order mirrors Chain: the first middleware in the list is the outermost
// wrapper and executes first. This enables a consistent layered architecture for both
// the publish and subscribe paths.
func ChainSubscribe(middlewares ...SubscribeMiddleware) SubscribeMiddleware {
	return func(nextHandler SubscribeHandler) SubscribeHandler {
		for middlewareIndex := len(middlewares) - 1; middlewareIndex >= 0; middlewareIndex-- {
			nextHandler = middlewares[middlewareIndex](nextHandler)
		}
		return nextHandler
	}
}

// Stack is a builder that accumulates publish and subscribe middlewares independently
// and composes them into chains on demand. It provides a convenient API for registering
// middlewares without needing to manage ordering manually -- middlewares execute in
// registration order (first registered = outermost wrapper).
type Stack struct {
	// publish holds the ordered list of publish-side middlewares.
	publish []PublishMiddleware
	// subscribe holds the ordered list of subscribe-side middlewares.
	subscribe []SubscribeMiddleware
}

// NewStack creates a new, empty middleware stack with pre-allocated slices
// for both publish and subscribe middlewares.
func NewStack() *Stack {
	return &Stack{
		publish:   make([]PublishMiddleware, 0),
		subscribe: make([]SubscribeMiddleware, 0),
	}
}

// UsePublish appends a publish middleware to the stack. Middlewares added first
// will be the outermost wrappers when the chain is built.
func (stack *Stack) UsePublish(middleware PublishMiddleware) {
	stack.publish = append(stack.publish, middleware)
}

// UseSubscribe appends a subscribe middleware to the stack. Middlewares added first
// will be the outermost wrappers when the chain is built.
func (stack *Stack) UseSubscribe(middleware SubscribeMiddleware) {
	stack.subscribe = append(stack.subscribe, middleware)
}

// PublishChain composes all registered publish middlewares into a single
// PublishMiddleware using Chain. The returned middleware can be applied
// to any PublishHandler.
func (stack *Stack) PublishChain() PublishMiddleware {
	return Chain(stack.publish...)
}

// SubscribeChain composes all registered subscribe middlewares into a single
// SubscribeMiddleware using ChainSubscribe. The returned middleware can be
// applied to any SubscribeHandler.
func (stack *Stack) SubscribeChain() SubscribeMiddleware {
	return ChainSubscribe(stack.subscribe...)
}

// WrapPublish is a convenience method that builds the publish middleware chain
// and immediately applies it to the given handler, returning a fully wrapped handler
// ready for use.
func (stack *Stack) WrapPublish(handler PublishHandler) PublishHandler {
	return stack.PublishChain()(handler)
}

// WrapSubscribe is a convenience method that builds the subscribe middleware chain
// and immediately applies it to the given handler, returning a fully wrapped handler
// ready for use.
func (stack *Stack) WrapSubscribe(handler SubscribeHandler) SubscribeHandler {
	return stack.SubscribeChain()(handler)
}

// Interceptor is implemented by types that provide both publish and subscribe
// middleware from a single logical concern (e.g., TracingMiddleware, MetricsCollector).
// This allows a single component to be registered with the Stack for both paths
// via UseInterceptor, ensuring symmetric instrumentation.
type Interceptor interface {
	// InterceptPublish returns a PublishMiddleware that wraps outgoing publish operations.
	InterceptPublish() PublishMiddleware

	// InterceptSubscribe returns a SubscribeMiddleware that wraps incoming message processing.
	InterceptSubscribe() SubscribeMiddleware
}

// UseInterceptor registers both the publish and subscribe middlewares provided by
// the given Interceptor. This is the preferred way to add dual-sided middleware
// (e.g., tracing, metrics) so that both publish and subscribe paths are instrumented
// consistently.
func (stack *Stack) UseInterceptor(interceptor Interceptor) {
	stack.UsePublish(interceptor.InterceptPublish())
	stack.UseSubscribe(interceptor.InterceptSubscribe())
}
