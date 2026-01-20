package handler

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type LeadCreatedHandler struct {
	log logger.Logger
}

func NewLeadCreatedHandler(log logger.Logger) *LeadCreatedHandler {
	return &LeadCreatedHandler{
		log: log.WithFields(map[string]any{"handler": "LeadCreatedHandler"}),
	}
}

func (h *LeadCreatedHandler) Handle(
	ctx context.Context,
	msg types.Message,
) (sqs.HandleResult, error) {
	h.log.Info("LeadCreatedHandler handled message", "messageId", msg.MessageId)

	return sqs.HandleSuccess, nil
}
