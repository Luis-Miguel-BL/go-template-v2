package handler

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/consumer/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type ExampleHandler struct {
	log logger.Logger
}

func NewExampleHandler(log logger.Logger) *ExampleHandler {
	return &ExampleHandler{
		log: log.WithFields(map[string]any{"handler": "ExampleHandler"}),
	}
}

func (h *ExampleHandler) Handle(
	ctx context.Context,
	msg types.Message,
) (sqs.HandleResult, error) {
	h.log.Info("ExampleHandler handled message", "messageId", msg.MessageId)

	return sqs.HandleSuccess, nil
}
