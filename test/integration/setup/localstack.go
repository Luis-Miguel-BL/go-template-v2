package setup

import (
	"context"
	"os"
	"sync"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

type localStackContainer struct {
	*localstack.LocalStackContainer
	endpoint string
	region   string
}

var once sync.Once
var lsContainer *localStackContainer

func provideLocalStack() *localStackContainer {
	var err error
	once.Do(func() {
		ctx := context.Background()
		lsContainer, err = startLocalStackContainer(ctx)
		if err != nil {
			panic(err)
		}

		os.Setenv("APP_ENVIRONMENT", "test")
		os.Setenv("APP_CONFIG_PATH", "../../../config")
		os.Setenv("APP_AWS_REGION", lsContainer.region)
		os.Setenv("APP_AWS_ENDPOINT", lsContainer.endpoint)
		os.Setenv("APP_AWS_ACCESS_KEY_ID", "test")
		os.Setenv("APP_AWS_SECRET_ACCESS_KEY", "test")
		os.Setenv("APP_AWS_SESSION_TOKEN", "test")
	})

	return lsContainer
}

func startLocalStackContainer(ctx context.Context) (*localStackContainer, error) {
	var err error

	lsContainer, err := localstack.Run(
		ctx,
		"localstack/localstack:latest",
		testcontainers.WithEnv(map[string]string{
			"SERVICES": "sqs,dynamodb",
			"DEBUG":    "1",
		}),
	)

	if err != nil {
		return nil, err
	}

	endpoint, err := lsContainer.Endpoint(ctx, "http")
	if err != nil {
		return nil, err
	}

	return &localStackContainer{
		LocalStackContainer: lsContainer,
		endpoint:            endpoint,
		region:              "sa-east-1",
	}, nil
}
