package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Handler interface {
	Handle(ctx context.Context, msg types.Message) (HandleResult, error)
}

type HandleResult string

const (
	HandleSuccess HandleResult = "success"
	HandleRetry   HandleResult = "retry"
	HandleDLQ     HandleResult = "dlq"
)
