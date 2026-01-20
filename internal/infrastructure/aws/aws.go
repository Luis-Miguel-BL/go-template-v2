package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AWSClient struct {
	cfg       AWSConfig
	awsConfig *aws.Config
}

type AWSConfig struct {
	Region   string
	Endpoint string
}

func NewAWSClient(cfg AWSConfig) *AWSClient {
	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
	}

	if cfg.Endpoint != "" {
		opts = append(opts, config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(
				func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:           cfg.Endpoint,
						SigningRegion: region,
					}, nil
				},
			),
		))
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		panic(err)
	}

	return &AWSClient{
		cfg:       cfg,
		awsConfig: &awsCfg,
	}
}
