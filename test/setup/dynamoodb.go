package setup

import (
	"context"

	_aws "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var leadTable = "lead_single_table"

func createTables(ctx context.Context, client *_aws.DynamoDBClient) (tableName string, err error) {
	pkAttribute := "PK"
	skAttribute := "SK"

	_, err = client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(leadTable),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String(pkAttribute), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(skAttribute), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String(pkAttribute), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String(skAttribute), KeyType: types.KeyTypeRange},
		},
		BillingMode: types.BillingModePayPerRequest,
	})

	return leadTable, err
}

func deleteTables(ctx context.Context, client *_aws.DynamoDBClient) error {
	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(leadTable),
	})

	return err
}
