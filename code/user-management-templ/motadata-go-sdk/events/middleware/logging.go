package middleware

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

// LogLevel represents logging level
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Logger interface for logging
type Logger interface {
	Debug(msg string, fields map[string]any)
	Info(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

// LoggingMiddleware adds logging to publish/subscribe operations
type LoggingMiddleware struct {
	logger         Logger
	level          LogLevel
	logPayload     bool
	maxPayloadSize int
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger:         logger,
		level:          LogLevelInfo,
		logPayload:     false,
		maxPayloadSize: 1024,
	}
}

// WithLevel sets the log level
func (l *LoggingMiddleware) WithLevel(level LogLevel) *LoggingMiddleware {
	l.level = level
	return l
}

// WithPayload enables payload logging
func (l *LoggingMiddleware) WithPayload(enabled bool, maxSize int) *LoggingMiddleware {
	l.logPayload = enabled
	l.maxPayloadSize = maxSize
	return l
}

// InterceptPublish returns the publish middleware
func (l *LoggingMiddleware) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			start := time.Now()

			fields := map[string]any{"subject": msg.Subject}
			l.addContextFields(ctx, fields)
			l.addMessageFields(msg, fields)

			l.log(LogLevelDebug, "publishing message", fields)

			err := next(ctx, msg)
			fields["latency_ms"] = time.Since(start).Milliseconds()

			if err != nil {
				fields["error"] = err.Error()
				l.log(LogLevelError, "publish failed", fields)
			} else {
				l.log(LogLevelInfo, "message published", fields)
			}

			return err
		}
	}
}

// InterceptSubscribe returns the subscribe middleware
func (l *LoggingMiddleware) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			start := time.Now()

			fields := map[string]any{}
			l.addContextFields(ctx, fields)
			if msg != nil {
				fields["subject"] = msg.Subject
			}
			l.addMessageFields(msg, fields)

			l.log(LogLevelDebug, "processing message", fields)

			err := next(ctx, msg)
			fields["process_time_ms"] = time.Since(start).Milliseconds()

			if err != nil {
				fields["error"] = err.Error()
				l.log(LogLevelError, "message processing failed", fields)
			} else {
				l.log(LogLevelDebug, "message processed", fields)
			}

			return err
		}
	}
}

// addContextFields adds trace context and tenant ID fields if available.
func (l *LoggingMiddleware) addContextFields(ctx context.Context, fields map[string]any) {
	if tc, ok := core.TraceContextFromContext(ctx); ok {
		fields["trace_id"] = tc.TraceID
		fields["span_id"] = tc.SpanID
	}
	if tenantID, ok := core.TenantIDFromContext(ctx); ok {
		fields["tenant_id"] = tenantID
	}
}

// addMessageFields adds message-level fields including optional payload.
func (l *LoggingMiddleware) addMessageFields(msg *nats.Msg, fields map[string]any) {
	if msg == nil {
		return
	}

	fields["msg_id"] = msg.Header.Get(core.HeaderMessageID)
	fields["size"] = len(msg.Data)

	if l.logPayload && len(msg.Data) > 0 {
		payload := string(msg.Data)
		if len(payload) > l.maxPayloadSize {
			payload = payload[:l.maxPayloadSize] + "..."
		}
		fields["payload"] = payload
	}
}

func (l *LoggingMiddleware) log(level LogLevel, msg string, fields map[string]any) {
	if l.logger == nil || level < l.level {
		return
	}

	switch level {
	case LogLevelDebug:
		l.logger.Debug(msg, fields)
	case LogLevelInfo:
		l.logger.Info(msg, fields)
	case LogLevelWarn:
		l.logger.Warn(msg, fields)
	case LogLevelError:
		l.logger.Error(msg, fields)
	}
}

// Logging returns a logging middleware
func Logging(logger Logger) *LoggingMiddleware {
	return NewLoggingMiddleware(logger)
}

// NoOpLogger is a logger that does nothing (placeholder)
type NoOpLogger struct{}

func (NoOpLogger) Debug(_ string, _ map[string]any) { /* no-op by design */ }
func (NoOpLogger) Info(_ string, _ map[string]any)  { /* no-op by design */ }
func (NoOpLogger) Warn(_ string, _ map[string]any)  { /* no-op by design */ }
func (NoOpLogger) Error(_ string, _ map[string]any) { /* no-op by design */ }

// OTELLogger implements Logger interface using OpenTelemetry logger
type OTELLogger struct {
	logger *logger.Logger
}

// NewOTELLogger creates a new OpenTelemetry-based logger
func NewOTELLogger() *OTELLogger {

	return &OTELLogger{
		logger: logger.L().Named("events"),
	}
}

// Debug logs at debug level
func (otelLogger *OTELLogger) Debug(msg string, fields map[string]any) {

	otelLogger.logger.Debug(context.Background(), msg, mapToFields(fields)...)
}

// Info logs at info level
func (otelLogger *OTELLogger) Info(msg string, fields map[string]any) {

	otelLogger.logger.Info(context.Background(), msg, mapToFields(fields)...)
}

// Warn logs at warn level
func (otelLogger *OTELLogger) Warn(msg string, fields map[string]any) {

	otelLogger.logger.Warn(context.Background(), msg, mapToFields(fields)...)
}

// Error logs at error level
func (otelLogger *OTELLogger) Error(msg string, fields map[string]any) {

	otelLogger.logger.Error(context.Background(), msg, mapToFields(fields)...)
}

// mapToFields converts a map to logger.Field slice
func mapToFields(m map[string]any) []logger.Field {

	if m == nil {

		return nil
	}

	fields := make([]logger.Field, 0, len(m))

	for key, value := range m {

		fields = append(fields, logger.Any(key, value))
	}

	return fields
}

// OTELLoggingMiddleware creates a logging middleware using OpenTelemetry logger
func OTELLoggingMiddleware() *LoggingMiddleware {

	return NewLoggingMiddleware(NewOTELLogger())
}
