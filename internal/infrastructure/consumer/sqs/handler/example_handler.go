package handler

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/integration"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/consumer/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ExampleHandler struct {
	log                   logger.Logger
	exampleAPIIntegration integration.ExampleAPIIntegration
}

func NewExampleHandler(log logger.Logger, exampleAPIIntegration integration.ExampleAPIIntegration) *ExampleHandler {
	return &ExampleHandler{
		log:                   log.WithFields(map[string]any{"handler": "ExampleHandler"}),
		exampleAPIIntegration: exampleAPIIntegration,
	}
}

func (h *ExampleHandler) Handle(
	ctx context.Context,
	msg types.Message,
) (sqs.HandleResult, error) {
	h.log.Info("ExampleHandler handled message", "messageId", msg.MessageId)

	_, err := h.exampleAPIIntegration.Create(ctx)
	if err != nil {
		h.log.Error("Failed to create example", "error", err)
		return sqs.HandleDLQ, err
	}

	return sqs.HandleSuccess, nil
}
