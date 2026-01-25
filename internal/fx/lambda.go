package fx

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda/handler"
	"go.uber.org/fx"
)

var LambdaModule = fx.Module("lambda",
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
)
