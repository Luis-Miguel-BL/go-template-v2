package observability

import (
	"context"
	"sync"
)

type Observability interface {
	// Tracing
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	AddAttributes(ctx context.Context, attrs map[string]any) context.Context
	AddEvent(ctx context.Context, event Event)
	RecordError(ctx context.Context, err error)
	TraceIDFromContext(ctx context.Context) string

	// Metrics
	RecordMetric(ctx context.Context, metric Metric)

	GetServerMiddlewares() []any

	Shutdown(ctx context.Context) error
}

type Span interface {
	End()
	SetAttributes(attrs map[string]any)
	RecordError(err error)
}

var (
	globalObservability Observability
	onceObservability   sync.Once
)

func GetObservability() Observability {
	return globalObservability
}

func SetObservability(tracer Observability) {
	onceObservability.Do(func() {
		globalObservability = tracer
	})
}

func StartSpan(ctx context.Context, name string) (context.Context, Span) {
	return globalObservability.StartSpan(ctx, name)
}
func AddAttributes(ctx context.Context, attrs map[string]any) context.Context {
	return globalObservability.AddAttributes(ctx, attrs)
}
func AddEvent(ctx context.Context, event Event) {
	globalObservability.AddEvent(ctx, event)
}
func RecordError(ctx context.Context, err error) {
	globalObservability.RecordError(ctx, err)
}
func TraceIDFromContext(ctx context.Context) string {
	return globalObservability.TraceIDFromContext(ctx)
}
