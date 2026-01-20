package repository

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/domain/lead/model"
	_aws "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging"
	dynamo_model "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/persistence/model/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBLeadRepository struct {
	tableName      string
	dynamoDBClient *_aws.DynamoDBClient
	dispatcher     *messaging.AggregateRootEventDispatcher
}

func NewDynamoDBLeadRepository(tableName string, dispatcher *messaging.AggregateRootEventDispatcher, dynamoDBClient *_aws.DynamoDBClient) *DynamoDBLeadRepository {
	return &DynamoDBLeadRepository{
		tableName:      tableName,
		dynamoDBClient: dynamoDBClient,
		dispatcher:     dispatcher,
	}
}

func (r *DynamoDBLeadRepository) Save(ctx context.Context, lead *model.Lead) (err error) {
	dynamoModel := dynamo_model.Lead{}
	item, err := dynamoModel.ToRepo(*lead)
	if err != nil {
		return err
	}

	_, err = r.dynamoDBClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	if err != nil {
		return err
	}

	r.dispatcher.PublishUncommitedEvents(ctx, lead)

	return nil
}

func (r *DynamoDBLeadRepository) GetByID(ctx context.Context, leadID string) (lead *model.Lead, err error) {
	dynamoLead := dynamo_model.Lead{}
	output, err := r.dynamoDBClient.Query(
		ctx,
		&dynamodb.QueryInput{
			TableName:                aws.String(r.tableName),
			KeyConditionExpression:   aws.String("#pk = :pk"),
			ExpressionAttributeNames: map[string]string{"#pk": "PK"},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":pk": &types.AttributeValueMemberS{
					Value: dynamo_model.MakeLeadPK(leadID),
				},
			},
		},
	)
	if err != nil {
		return lead, err
	}
	if output.Count == 0 {
		return lead, domain.EntityNotFoundError("lead", leadID)
	}

	lead, err = dynamoLead.ToDomain(output.Items)
	if err != nil {
		return lead, err
	}

	return lead, nil
}
