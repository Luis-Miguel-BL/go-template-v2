package setup

import (
	"context"

	_aws "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var queueName = "example-queue"
var queueURL = ""

func createQueues(ctx context.Context, client *_aws.SQSClient) (queueURL string, err error) {
	out, err := client.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	queueURL = *out.QueueUrl

	return queueURL, nil
}

func deleteQueues(ctx context.Context, client *_aws.SQSClient) error {
	_, err := client.DeleteQueue(ctx, &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	})

	return err
}
