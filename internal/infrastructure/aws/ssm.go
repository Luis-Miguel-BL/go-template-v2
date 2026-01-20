package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMClient struct {
	*ssm.Client
}

func NewSSMClient(awsClient *AWSClient) *SSMClient {
	ssmClient := ssm.NewFromConfig(*awsClient.awsConfig)

	return &SSMClient{
		Client: ssmClient,
	}

}

func (c *SSMClient) GetParametersByPath(
	ctx context.Context,
	path string,
) (map[string]string, error) {
	out := make(map[string]string)

	var nextToken *string

	for {
		resp, err := c.Client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
			Path:           aws.String(path),
			WithDecryption: aws.Bool(true),
			Recursive:      aws.Bool(true),
			NextToken:      nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, param := range resp.Parameters {
			if param.Name == nil || param.Value == nil {
				continue
			}

			out[*param.Name] = *param.Value
		}

		if resp.NextToken == nil {
			break
		}

		nextToken = resp.NextToken
	}

	return out, nil
}
