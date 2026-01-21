package aws

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient struct {
	*dynamodb.Client
	obs observability.Observability
}

func NewDynamoDBClient(awsClient *AWSClient, obs observability.Observability) *DynamoDBClient {
	dynamodbClient := dynamodb.NewFromConfig(*awsClient.awsConfig)

	return &DynamoDBClient{
		Client: dynamodbClient,
		obs:    obs,
	}
}

func (d *DynamoDBClient) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	ctx, span := d.obs.StartSpan(ctx, "dynamodb.query")
	defer span.End()

	result, err := d.Client.Query(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}

func (d *DynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	ctx, span := d.obs.StartSpan(ctx, "dynamodb.put_item")
	defer span.End()

	result, err := d.Client.PutItem(ctx, params, optFns...)
	if err != nil {
		span.RecordError(err)
	}
	return result, err
}
