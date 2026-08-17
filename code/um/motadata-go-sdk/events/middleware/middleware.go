package middleware

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// PublishMiddleware intercepts publish operations
type PublishMiddleware func(PublishHandler) PublishHandler

// PublishHandler handles a publish operation
type PublishHandler func(ctx context.Context, subject string, msg *core.Message) error

// SubscribeMiddleware intercepts subscribe operations
type SubscribeMiddleware func(SubscribeHandler) SubscribeHandler

// SubscribeHandler handles a received message
type SubscribeHandler func(ctx context.Context, msg *core.Message) error

// Chain chains multiple publish middlewares
func Chain(middlewares ...PublishMiddleware) PublishMiddleware {
	return func(nextHandler PublishHandler) PublishHandler {
		for middlewareIndex := len(middlewares) - 1; middlewareIndex >= 0; middlewareIndex-- {
			nextHandler = middlewares[middlewareIndex](nextHandler)
		}
		return nextHandler
	}
}

// ChainSubscribe chains multiple subscribe middlewares
func ChainSubscribe(middlewares ...SubscribeMiddleware) SubscribeMiddleware {
	return func(nextHandler SubscribeHandler) SubscribeHandler {
		for middlewareIndex := len(middlewares) - 1; middlewareIndex >= 0; middlewareIndex-- {
			nextHandler = middlewares[middlewareIndex](nextHandler)
		}
		return nextHandler
	}
}

// Stack holds a collection of middlewares
type Stack struct {
	publish   []PublishMiddleware
	subscribe []SubscribeMiddleware
}

// NewStack creates a new middleware stack
func NewStack() *Stack {
	return &Stack{
		publish:   make([]PublishMiddleware, 0),
		subscribe: make([]SubscribeMiddleware, 0),
	}
}

// UsePublish adds a publish middleware
func (stack *Stack) UsePublish(middleware PublishMiddleware) {
	stack.publish = append(stack.publish, middleware)
}

// UseSubscribe adds a subscribe middleware
func (stack *Stack) UseSubscribe(middleware SubscribeMiddleware) {
	stack.subscribe = append(stack.subscribe, middleware)
}

// PublishChain returns the chained publish middlewares
func (stack *Stack) PublishChain() PublishMiddleware {
	return Chain(stack.publish...)
}

// SubscribeChain returns the chained subscribe middlewares
func (stack *Stack) SubscribeChain() SubscribeMiddleware {
	return ChainSubscribe(stack.subscribe...)
}

// WrapPublish wraps a publish handler with the middleware stack
func (stack *Stack) WrapPublish(handler PublishHandler) PublishHandler {
	return stack.PublishChain()(handler)
}

// WrapSubscribe wraps a subscribe handler with the middleware stack
func (stack *Stack) WrapSubscribe(handler SubscribeHandler) SubscribeHandler {
	return stack.SubscribeChain()(handler)
}

// Interceptor provides both publish and subscribe interception
type Interceptor interface {
	// InterceptPublish returns a publish middleware
	InterceptPublish() PublishMiddleware

	// InterceptSubscribe returns a subscribe middleware
	InterceptSubscribe() SubscribeMiddleware
}

// UseInterceptor adds an interceptor to the stack
func (stack *Stack) UseInterceptor(interceptor Interceptor) {
	stack.UsePublish(interceptor.InterceptPublish())
	stack.UseSubscribe(interceptor.InterceptSubscribe())
}
