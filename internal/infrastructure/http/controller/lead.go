package controller

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/dto"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/usecase"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/payload"
	"github.com/labstack/echo/v4"
	// "go.opentelemetry.io/otel"
	// "go.opentelemetry.io/otel/attribute"
	// "go.opentelemetry.io/otel/metric"
	// "go.opentelemetry.io/otel/trace"
)

type LeadController struct {
	createLeadUseCase *usecase.CreateLead
}

func NewLeadController(createLeadUseCase *usecase.CreateLead) *LeadController {
	return &LeadController{createLeadUseCase}
}

func (c *LeadController) Create(ctx echo.Context) error {
	context := ctx.Request().Context()

	log := logger.FromContext(context)

	log.Debug("Creating lead")
	t := observability.GetObservability()
	_, span := t.StartSpan(context, "LeadController.Create")
	defer span.End()

	// span := trace.SpanFromContext(context)
	// if span != nil {
	// 	span.SetAttributes(attribute.String("user.id", "123"))
	// 	span.SetAttributes(attribute.String("operation", "create_lead"))
	// 	span.AddEvent("TemplateLeadCreated1")
	// 	span.AddEvent("TemplateLeadCreated", trace.WithAttributes(attribute.String("lead_id", "abc-123")))
	// }

	// meter := otel.Meter("api")
	// gauge, _ := meter.Int64Gauge("teste_gauge", metric.WithDescription("Number of active users"), metric.WithUnit("1"))
	// counter, _ := meter.Int64Counter("requests_total", metric.WithDescription("Total number of requests"), metric.WithUnit("1"))
	// histogram, _ := meter.Int64Histogram("request_duration_ms", metric.WithDescription("Duration of requests"), metric.WithUnit("ms"))

	// gauge.Record(context, 42, metric.WithAttributes(attribute.String("endpoint", "/leads")))
	// counter.Add(context, 1)
	// histogram.Record(context, 150, metric.WithAttributes(attribute.String("endpoint", "/leads")))

	reqBody := payload.CreateLeadRequest{}
	err := ctx.Bind(&reqBody)
	if err != nil {
		return Error(ctx, err)
	}

	output, err := c.createLeadUseCase.Execute(context, dto.CreateLeadInput(reqBody))
	if err != nil {
		return Error(ctx, err)
	}
	return Ok(ctx, payload.CreateLeadResponse{
		LeadID:      output.LeadID,
		AccessToken: output.AccessToken.Token,
		ExpiresIn:   output.AccessToken.ExpiresIn,
	})
}
