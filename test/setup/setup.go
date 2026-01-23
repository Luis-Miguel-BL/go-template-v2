package setup

import (
	"context"
	"sync"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
	"github.com/Luis-Miguel-BL/go-lm-template/test/util"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type BaseTestSuite struct {
	suite.Suite
	*util.TestUtil

	ctx     context.Context
	wg      *sync.WaitGroup
	app     *fx.App
	readyCh chan struct{}
}

func (s *BaseTestSuite) SetupApp(wg *sync.WaitGroup, modules ...fx.Option) {
	s.ctx = context.Background()
	s.wg = wg
	s.readyCh = make(chan struct{})

	opts := []fx.Option{
		s.testModule(),
	}
	opts = append(opts, modules...)

	s.app = fx.New(opts...)

	go func() {
		s.app.Run()
	}()

	s.awaitReady()
}

func (s *BaseTestSuite) TearDownSuite() {
	if s.wg != nil {
		s.wg.Wait()
	}

	if s.app != nil {
		_ = s.app.Stop(s.ctx)
	}

}

func (s *BaseTestSuite) testModule() fx.Option {
	return fx.Module(
		"test",
		fx.Provide(
			provideLocalStack,
		),
		fx.Invoke(
			s.testLifecycle,
		),
	)
}

func (s *BaseTestSuite) testLifecycle(lc fx.Lifecycle,
	_ *localStackContainer,
	cfg *config.Config,
	httpClientFactory httpclient.HTTPClientFactory,
	dynamoDBClient *aws.DynamoDBClient,
	sqsClient *aws.SQSClient,
) {
	s.TestUtil = util.NewTestUtil(cfg, httpClientFactory)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			deleteTables(ctx, dynamoDBClient)
			deleteQueues(ctx, sqsClient)
			return nil
		},
		OnStart: func(ctx context.Context) error {
			tableName, err := createTables(ctx, dynamoDBClient)
			if err != nil {
				return err
			}
			cfg.AWS.DynamoDB.LeadTableName = tableName

			queueURL, err := createQueues(ctx, sqsClient)
			if err != nil {
				return err
			}
			cfg.Consumer.SQSQueueURL = queueURL

			close(s.readyCh)
			return nil
		},
	})
}

func (s *BaseTestSuite) awaitReady() {
	timeout := 10 * time.Second
	select {
	case <-s.readyCh:
	case <-time.After(timeout):
		s.T().Fatal("app did not become ready within timeout")
	}
}
