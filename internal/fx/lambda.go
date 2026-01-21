package fx

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda/handler"
	"go.uber.org/fx"
)

func LambdaModule(lambdaName string) fx.Option {
	return fx.Module("lambda",
		fx.Provide(
			//handlers
			fx.Annotate(
				handler.NewExampleSQSHandler,
				fx.As(new(lambda.LambdaHandler)),
				fx.ResultTags(`group:"lambda-handlers"`),
			),

			//registry
			fx.Annotate(
				lambda.NewRegistry,
				fx.ParamTags(`group:"lambda-handlers"`),
			),

			//runner
			lambda.NewRunner,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, runner *lambda.Runner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						runner.Run(lambdaName)
						return nil
					},
				})
			},
		),
	)
}
