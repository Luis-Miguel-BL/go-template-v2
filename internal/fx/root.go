package fx

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	httpclient_factory "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient/factory"
	infra_logger "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/logger"
	infra_telemetry "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/telemetry"
	"go.uber.org/fx"
)

var RootModule = fx.Module("root",
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

		// fx.Annotate(infra_telemetry.NewNewRelicTelemetry, fx.As(new(telemetry.Telemetry))),
		fx.Annotate(infra_telemetry.NewOtelTelemetry, fx.As(new(telemetry.Telemetry))),

		// logger
		fx.Annotate(infra_logger.NewZapLogger, fx.As(new(logger.Logger))),

		// http clients
		httpclient_factory.NewHTTPClientFactory,
	),

	fx.Invoke(
		func(lc fx.Lifecycle, obs telemetry.Telemetry) {
			lc.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					return obs.Shutdown(ctx)
				},
			})
		},
	),
)
