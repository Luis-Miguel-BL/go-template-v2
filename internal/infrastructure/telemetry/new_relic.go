package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/newrelic/go-agent/v3/integrations/nrecho-v4"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type NewRelicTelemetry struct {
	cfg *config.Config

	newRelicApp *newrelic.Application
}

func NewNewRelicTelemetry(cfg *config.Config) (*NewRelicTelemetry, error) {
	newRelicTelemetry := &NewRelicTelemetry{
		cfg: cfg,
	}

	err := newRelicTelemetry.initNewRelicApp()
	if err != nil {
		return nil, err
	}

	telemetry.SetTelemetry(newRelicTelemetry)

	return newRelicTelemetry, nil
}

func (n *NewRelicTelemetry) AddAttributes(ctx context.Context, attrs map[string]any) context.Context {
	txn := newrelic.FromContext(ctx)

	for k, v := range attrs {
		txn.AddAttribute(k, v)
	}
	return ctx
}

func (n *NewRelicTelemetry) AddEvent(ctx context.Context, event telemetry.Event) {
	eventName := fmt.Sprintf("%s%s", n.cfg.Monitor.NewRelicConfig.CustomEventPrefix, event.Name())
	attrs := event.Attributes()

	n.newRelicApp.RecordCustomEvent(eventName, attrs)
}

func (n *NewRelicTelemetry) RecordError(ctx context.Context, err error) {
	txn := newrelic.FromContext(ctx)
	txn.NoticeError(err)
}

func (n *NewRelicTelemetry) TraceIDFromContext(ctx context.Context) string {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		return ""
	}

	return txn.GetTraceMetadata().TraceID
}

func (n *NewRelicTelemetry) RecordMetric(ctx context.Context, metric telemetry.Metric) {
	n.newRelicApp.RecordCustomMetric(metric.Name(), float64(metric.Value()))
}

func (n *NewRelicTelemetry) StartSpan(ctx context.Context, name string) (context.Context, telemetry.Span) {
	txn := newrelic.FromContext(ctx)
	if txn == nil {
		txn = n.newRelicApp.StartTransaction(name)
	}

	span := txn.StartSegment(name)

	newCtx := newrelic.NewContext(ctx, txn)

	return newCtx, &NewRelicSpan{segment: span}
}

type NewRelicSpan struct {
	segment *newrelic.Segment
}

func (s *NewRelicSpan) End() {
	s.segment.End()
}

func (s *NewRelicSpan) SetAttributes(attrs map[string]any) {
	for k, v := range attrs {
		s.segment.AddAttribute(k, v)
	}
}

func (s *NewRelicSpan) RecordError(err error) {
	s.segment.AddAttribute("error", err.Error())
}

func (n *NewRelicTelemetry) Shutdown(ctx context.Context) error {
	n.newRelicApp.Shutdown(time.Duration(n.cfg.Monitor.NewRelicConfig.ShutdownTimeoutSeconds) * time.Second)

	return nil
}

func (n *NewRelicTelemetry) GetServerMiddlewares() []any {
	return []any{
		nrecho.Middleware(n.newRelicApp),
	}
}

func (n *NewRelicTelemetry) NewHttpTransport() http.RoundTripper {
	return newrelic.NewRoundTripper(http.DefaultTransport)
}

func (n *NewRelicTelemetry) initNewRelicApp() error {
	nrl, err := newrelic.NewApplication(
		newrelic.ConfigAppName(n.cfg.App.Name),
		newrelic.ConfigLicense(n.cfg.Monitor.NewRelicConfig.AppKey),
		newrelic.ConfigEnabled(n.cfg.Monitor.Enabled),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return err
	}

	n.newRelicApp = nrl
	return nil
}
