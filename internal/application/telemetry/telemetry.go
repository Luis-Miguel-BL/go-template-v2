package telemetry

import (
	"context"
	"net/http"
	"sync"
)

type Telemetry interface {
	// Tracing
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	AddAttributes(ctx context.Context, attrs map[string]any) context.Context
	AddEvent(ctx context.Context, event Event)
	RecordError(ctx context.Context, err error)
	TraceIDFromContext(ctx context.Context) string

	// Metrics
	RecordMetric(ctx context.Context, metric Metric)

	GetServerMiddlewares() []any
	NewHttpTransport() http.RoundTripper

	Shutdown(ctx context.Context) error
}

type Span interface {
	End()
	SetAttributes(attrs map[string]any)
	RecordError(err error)
}

var (
	globalTelemetry Telemetry
	onceTelemetry   sync.Once
)

func GetTelemetry() Telemetry {
	return globalTelemetry
}

func SetTelemetry(tracer Telemetry) {
	onceTelemetry.Do(func() {
		globalTelemetry = tracer
	})
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return globalTelemetry.StartSpan(ctx, name)
}
func AddAttributes(ctx context.Context, attrs map[string]any) context.Context {
	return globalTelemetry.AddAttributes(ctx, attrs)
}
func AddEvent(ctx context.Context, event Event) {
	globalTelemetry.AddEvent(ctx, event)
}
func RecordError(ctx context.Context, err error) {
	globalTelemetry.RecordError(ctx, err)
}
func TraceIDFromContext(ctx context.Context) string {
	return globalTelemetry.TraceIDFromContext(ctx)
}
