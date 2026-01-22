// internal/infrastructure/aws/lambda/runner.go
package lambda

import (
	"context"
	"fmt"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Runner struct {
	registry  *Registry
	cfg       *config.Config
	log       logger.Logger
	telemetry telemetry.Telemetry
}

func NewRunner(cfg *config.Config, log logger.Logger, registry *Registry, telemetry telemetry.Telemetry) *Runner {
	return &Runner{
		cfg:       cfg,
		log:       log,
		registry:  registry,
		telemetry: telemetry,
	}
}

func (r *Runner) Run(lambdaName string) {
	if lambdaName == "" {
		r.log.Error("lambda name is required to run the lambda runner")
		return
	}

	handler, err := r.registry.Get(lambdaName)
	if err != nil {
		r.log.Error(err.Error())
		return
	}

	if r.cfg.IsLocal() {
		r.runLocal(handler)
		return
	}

	lambda.Start(r.instrument(handler, lambdaName))
}

func (r *Runner) instrument(handler any, lambdaName string) any {
	return func(ctx context.Context, event any) (any, error) {
		ctx, span := r.telemetry.StartSpan(ctx, "lambda.handler."+lambdaName)
		defer span.End()

		switch h := handler.(type) {
		case func(context.Context, any) (any, error):
			return h(ctx, event)
		default:
			r.log.Error("unsupported handler type for instrumentation")
			return nil, nil
		}
	}
}

type localRunnable[Event any, Response any] interface {
	SampleEvent() Event
	Handle(ctx context.Context, event Event) (Response, error)
}

func (r *Runner) runLocal(handler any) {
	r.log.Info("Running lambda in local mode")

	switch h := handler.(type) {
	case localRunnable[events.SQSEvent, events.SQSEventResponse]:
		event := h.SampleEvent()
		resp, err := h.Handle(context.Background(), event)
		if err != nil {
			r.log.Error("local execution error: " + err.Error())
			return
		}
		r.log.Info(fmt.Sprintf("local execution finished successfully: %v", resp))

	default:
		r.log.Error("handler does not support local execution (localRunnable not implemented)")
	}
}
