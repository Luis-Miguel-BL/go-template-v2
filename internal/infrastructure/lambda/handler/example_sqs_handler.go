package handler

import (
	"context"
	"encoding/json"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/aws/aws-lambda-go/events"
)

type ExampleSQSHandler struct {
	log logger.Logger
}

func NewExampleSQSHandler(log logger.Logger) *ExampleSQSHandler {
	exampleHandler := ExampleSQSHandler{
		log: log.WithFields(map[string]any{"handler": "ExampleSQSHandler"}),
	}

	return &exampleHandler
}

func (h *ExampleSQSHandler) LambdaName() string {
	return "ExampleSQSHandler"
}

func (h *ExampleSQSHandler) Handle(ctx context.Context, event events.SQSEvent) (res events.SQSEventResponse, err error) {
	e, err := json.Marshal(event)
	if err != nil {
		h.log.Error("Failed to marshal SQS event", "error", err)
		return res, err
	}
	h.log.Info("ExampleSQSHandler handled message " + string(e))
	return res, nil
}

func (h *ExampleSQSHandler) SampleEvent() events.SQSEvent {
	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "example-message-id",
				ReceiptHandle: "example-receipt-handle",
				Body:          "example-message-body",
			},
		},
	}
}
