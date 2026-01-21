package observability

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otel_metric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	otel_trace "go.opentelemetry.io/otel/trace"
)

type OtelTelemetry struct {
	cfg *config.Config

	tracer otel_trace.Tracer
	meter  otel_metric.Meter

	resource       *resource.Resource
	tracerProvider *trace.TracerProvider
	meterProvider  *metric.MeterProvider

	countersMap   map[string]otel_metric.Int64Counter
	countersMu    sync.Mutex
	gaugesMap     map[string]otel_metric.Int64Gauge
	gaugesMu      sync.Mutex
	histogramsMap map[string]otel_metric.Int64Histogram
	histogramsMu  sync.Mutex
}

type OtelSpan struct {
	span otel_trace.Span
}

func NewOtelTelemetry(cfg *config.Config) (*OtelTelemetry, error) {
	ctx := context.Background()
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.App.Name),
	)

	otelTelemetry := &OtelTelemetry{
		cfg:           cfg,
		resource:      res,
		countersMap:   make(map[string]otel_metric.Int64Counter),
		countersMu:    sync.Mutex{},
		gaugesMap:     make(map[string]otel_metric.Int64Gauge),
		gaugesMu:      sync.Mutex{},
		histogramsMap: make(map[string]otel_metric.Int64Histogram),
		histogramsMu:  sync.Mutex{},
	}

	err := otelTelemetry.initTracer(ctx)
	if err != nil {
		return nil, err
	}

	err = otelTelemetry.initMeter(ctx)
	if err != nil {
		return nil, err
	}

	observability.SetObservability(otelTelemetry)

	return otelTelemetry, nil
}

func (o *OtelTelemetry) Shutdown(ctx context.Context) error {
	if err := o.tracerProvider.Shutdown(ctx); err != nil {
		return err
	}
	if err := o.meterProvider.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (o *OtelTelemetry) TraceIDFromContext(ctx context.Context) string {
	span := otel_trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	sc := span.SpanContext()
	if !sc.HasTraceID() {
		return ""
	}
	return sc.TraceID().String()
}

func (o *OtelTelemetry) RecordError(ctx context.Context, err error) {
	span := otel_trace.SpanFromContext(ctx)
	if span != nil {
		span.RecordError(err)
	}
}

func (o *OtelTelemetry) AddAttributes(ctx context.Context, attrs map[string]any) context.Context {
	span := otel_trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(convertToOtelAttributes(attrs)...)
	}

	ctx = newContextWithSpanAttributes(ctx, attrs)
	return ctx
}

func (o *OtelTelemetry) RecordMetric(ctx context.Context, metric observability.Metric) {
	switch metric.Type() {
	case observability.MetricTypeCounter:
		o.recordCounter(ctx, metric.Name(), metric.Value(), metric.Attributes())
	case observability.MetricTypeGauge:
		o.recordGauge(ctx, metric.Name(), metric.Value(), metric.Attributes())
	case observability.MetricTypeHistogram:
		o.recordHistogram(ctx, metric.Name(), metric.Value(), metric.Attributes())
	}
}

// Using New Relic to add custom events
func (o *OtelTelemetry) AddEvent(ctx context.Context, event observability.Event) {
	eventName := fmt.Sprintf("%s%s", o.cfg.Monitor.NewRelicConfig.CustomEventPrefix, event.Name())
	attrs := event.Attributes()

	span := otel_trace.SpanFromContext(ctx)
	if span != nil {
		span.AddEvent(eventName, otel_trace.WithAttributes(convertToOtelAttributes(attrs)...))
	}
}

func (o *OtelTelemetry) StartSpan(ctx context.Context, name string) (context.Context, observability.Span) {
	var opts []otel_trace.SpanStartOption
	attrs := getSpanAttributesFromContext(ctx)
	if len(attrs) > 0 {
		opts = append(opts, otel_trace.WithAttributes(convertToOtelAttributes(attrs)...))
	}
	ctx, span := o.tracer.Start(ctx, name, opts...)

	return ctx, &OtelSpan{span: span}
}

func (s *OtelSpan) End() {
	s.span.End()
}

func (s *OtelSpan) SetAttributes(attrs map[string]any) {
	s.span.SetAttributes(convertToOtelAttributes(attrs)...)
}

func (s *OtelSpan) RecordError(err error) {
	s.span.RecordError(err)
}

func (o *OtelTelemetry) GetServerMiddlewares() []any {
	return []any{
		otelecho.Middleware(o.cfg.App.Name),
	}
}

func (o *OtelTelemetry) recordCounter(ctx context.Context, name string, value int64, attrs map[string]any) {
	o.countersMu.Lock()
	counter, ok := o.countersMap[name]
	if !ok {
		counter, _ = o.meter.Int64Counter(name)
		o.countersMap[name] = counter
	}
	o.countersMu.Unlock()
	otelAttrs := convertToOtelAttributes(attrs)
	counter.Add(ctx, value, otel_metric.WithAttributes(otelAttrs...))
}

func (o *OtelTelemetry) recordGauge(ctx context.Context, name string, value int64, attrs map[string]any) {
	o.gaugesMu.Lock()
	gauge, ok := o.gaugesMap[name]
	if !ok {
		gauge, _ = o.meter.Int64Gauge(name)
		o.gaugesMap[name] = gauge
	}
	o.gaugesMu.Unlock()

	otelAttrs := convertToOtelAttributes(attrs)
	gauge.Record(ctx, value, otel_metric.WithAttributes(otelAttrs...))
}

func (o *OtelTelemetry) recordHistogram(ctx context.Context, name string, value int64, attrs map[string]any) {
	o.histogramsMu.Lock()
	histogram, ok := o.histogramsMap[name]
	if !ok {
		histogram, _ = o.meter.Int64Histogram(name)
		o.histogramsMap[name] = histogram
	}
	o.histogramsMu.Unlock()

	otelAttrs := convertToOtelAttributes(attrs)
	histogram.Record(ctx, value, otel_metric.WithAttributes(otelAttrs...))
}

func (o *OtelTelemetry) initTracer(ctx context.Context) error {
	traceOpts := []trace.TracerProviderOption{
		trace.WithResource(o.resource),
	}

	if o.cfg.Monitor.Enabled {
		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(o.cfg.Monitor.NewRelicConfig.Endpoint),
			otlptracehttp.WithHeaders(map[string]string{
				"api-key": o.cfg.Monitor.NewRelicConfig.AppKey,
			}),
		)
		if err != nil {
			return err
		}

		traceOpts = append(traceOpts, trace.WithBatcher(exporter))
	}

	tp := trace.NewTracerProvider(traceOpts...)
	otel.SetTracerProvider(tp)

	o.tracerProvider = tp
	o.tracer = tp.Tracer(o.cfg.App.Name)

	return nil
}

func (o *OtelTelemetry) initMeter(ctx context.Context) error {
	metricOpts := []metric.Option{
		metric.WithResource(o.resource),
	}
	if o.cfg.Monitor.Enabled {
		exporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(o.cfg.Monitor.NewRelicConfig.Endpoint),
			otlpmetrichttp.WithHeaders(map[string]string{
				"api-key": o.cfg.Monitor.NewRelicConfig.AppKey,
			}),
		)
		if err != nil {
			return err
		}
		metricOpts = append(metricOpts, metric.WithReader(metric.NewPeriodicReader(exporter)))
	}

	mp := metric.NewMeterProvider(
		metricOpts...,
	)
	otel.SetMeterProvider(mp)

	o.meterProvider = mp
	o.meter = mp.Meter(o.cfg.App.Name)

	return nil
}

func convertToOtelAttributes(attrs map[string]any) []attribute.KeyValue {
	var otelAttrs []attribute.KeyValue
	for k, v := range attrs {
		otelAttrs = append(otelAttrs, attribute.String(k, fmt.Sprintf("%v", v)))
	}
	return otelAttrs
}

type contextKeySpanAttrs string

const contextKeySpanAttrsKey contextKeySpanAttrs = "span-attributes"

func newContextWithSpanAttributes(ctx context.Context, attrs map[string]any) context.Context {
	parentAttrs, ok := ctx.Value(contextKeySpanAttrsKey).(map[string]any)
	if ok && parentAttrs != nil {
		merged := make(map[string]any, len(parentAttrs)+len(attrs))
		maps.Copy(merged, parentAttrs)
		maps.Copy(merged, attrs)
		return context.WithValue(ctx, contextKeySpanAttrsKey, merged)
	}

	return context.WithValue(ctx, contextKeySpanAttrsKey, attrs)

}

func getSpanAttributesFromContext(ctx context.Context) map[string]any {
	attrs, ok := ctx.Value(contextKeySpanAttrsKey).(map[string]any)
	if !ok {
		return nil
	}
	return attrs
}
