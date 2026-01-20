package modules

import (
	"context"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/eventbus"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/subscriber"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	infra_logger "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/messaging"
	infra_observability "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/observability"
	"go.uber.org/fx"
)

func RootModule(wg *sync.WaitGroup) fx.Option {
	return fx.Module("root",
		fx.Provide(
			//config
			config.LoadBootstrapConfig,
			config.Load,

			//aws
			func(cfg *config.BootstrapConfig) *aws.AWSClient {
				return aws.NewAWSClient(aws.AWSConfig{
					Region:   cfg.AWS.Region,
					Endpoint: cfg.AWS.Endpoint,
				})
			},
			aws.NewDynamoDBClient,
			aws.NewSQSClient,
			aws.NewSSMClient,

			// fx.Annotate(infra_observability.NewNewRelicTelemetry, fx.As(new(observability.Observability))),
			fx.Annotate(infra_observability.NewOtelTelemetry, fx.As(new(observability.Observability))),

			// logger
			fx.Annotate(infra_logger.NewZapLogger, fx.As(new(logger.Logger))),

			// messaging & eventbus
			messaging.NewAggregateRootEventDispatcher,
			fx.Annotate(subscriber.NewMonitorSubscriber, fx.ResultTags(`group:"subscribers"`)),
			fx.Annotate(subscriber.NewMetricSubscriber, fx.ResultTags(`group:"subscribers"`)),
			fx.Annotate(
				func(eventSubscribers ...eventbus.EventSubscriber) eventbus.EventBus {
					return messaging.NewDomainEventBus(wg, eventSubscribers...)
				},
				fx.ParamTags(`group:"subscribers"`),
			),
		),

		fx.Invoke(
			func(lc fx.Lifecycle, telemetry observability.Observability) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return telemetry.Shutdown(ctx)
					},
				})
			},
		),
	)
}
