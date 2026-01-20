package modules

import (
	"context"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging/sqs"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging/sqs/handler"
	"go.uber.org/fx"
)

func WorkerModule(wg *sync.WaitGroup) fx.Option {
	return fx.Module("worker",
		fx.Provide(
			// handlers
			handler.NewLeadCreatedHandler,

			// consumers
			fx.Annotate(
				newLeadCreatedConsumer,
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

func newLeadCreatedConsumer(cfg *config.Config, client *aws.SQSClient, handler *handler.LeadCreatedHandler, logger logger.Logger) *sqs.Consumer {
	consumerConfig := sqs.ConsumerConfig{
		QueueURL: cfg.Worker.SQSQueueURL,
	}
	return sqs.NewConsumer(consumerConfig, client, handler, logger)
}
