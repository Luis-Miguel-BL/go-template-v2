package aws

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient struct {
	*dynamodb.Client
	telemetry telemetry.Telemetry
}

func NewDynamoDBClient(awsClient *AWSClient, telemetry telemetry.Telemetry) *DynamoDBClient {
	dynamodbClient := dynamodb.NewFromConfig(*awsClient.awsConfig)

	return &DynamoDBClient{
		Client:    dynamodbClient,
		telemetry: telemetry,
	}
}

func (d *DynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	ctx, span := d.telemetry.StartSpan(ctx, "dynamodb.query")
	defer span.End()

	result, err := d.Client.Query(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}

func (d *DynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	ctx, span := d.telemetry.StartSpan(ctx, "dynamodb.put_item")
	defer span.End()

	result, err := d.Client.PutItem(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}
