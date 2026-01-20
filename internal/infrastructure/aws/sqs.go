package aws

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	*sqs.Client
}

func NewSQSClient(awsClient *AWSClient) *SQSClient {
	sqsClient := sqs.NewFromConfig(*awsClient.awsConfig)

	return &SQSClient{
		Client: sqsClient,
	}

}

func (s *SQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	ctx, span := observability.GetObservability().StartSpan(ctx, "sqs.send_message")
	defer span.End()

	result, err := s.Client.SendMessage(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}

func (s *SQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	ctx, span := observability.GetObservability().StartSpan(ctx, "sqs.receive_message")
	defer span.End()

	result, err := s.Client.ReceiveMessage(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}

func (s *SQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	ctx, span := observability.GetObservability().StartSpan(ctx, "sqs.delete_message")
	defer span.End()

	result, err := s.Client.DeleteMessage(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}
