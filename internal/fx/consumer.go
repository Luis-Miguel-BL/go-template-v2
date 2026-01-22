package fx

import (
	"context"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/consumer/sqs"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/consumer/sqs/handler"
	"go.uber.org/fx"
)

func ConsumerModule(wg *sync.WaitGroup) fx.Option {
	return fx.Module("consumer",
		fx.Provide(
			// handlers
			handler.NewExampleHandler,

			// consumers
			fx.Annotate(
				newExampleConsumer,
				fx.ResultTags(`group:"sqs-consumers"`),
			),
			fx.Annotate(
				func(consumers []*sqs.Consumer) []*sqs.Consumer {
					return consumers
				},
				fx.ParamTags(`group:"sqs-consumers"`),
			),
		),
		fx.Invoke(
			func(lc fx.Lifecycle, consumers []*sqs.Consumer) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						for _, c := range consumers {
							go c.Start(wg)
						}
						return nil
					},
					OnStop: func(ctx context.Context) error {
						for _, c := range consumers {
							c.Stop()
						}
						return nil
					},
				})
			},
		),
	)
}

func newExampleConsumer(cfg *config.Config, client *aws.SQSClient, telemetry telemetry.Telemetry, handler *handler.ExampleHandler, logger logger.Logger) *sqs.Consumer {
	consumerConfig := sqs.ConsumerConfig{
		QueueURL: cfg.Worker.SQSQueueURL,
	}
	return sqs.NewConsumer(consumerConfig, client, telemetry, handler, logger)
}
