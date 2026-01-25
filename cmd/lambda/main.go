package main

import (
	"context"
	"sync"

	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda"
	"go.uber.org/fx"
)

var (
	// FunctionName defined in build time for lambda function
	FunctionName string
)

func main() {
	var wg sync.WaitGroup
	app := fx.New(
		_fx.RootModule,
		_fx.LambdaModule,
		_fx.ApplicationModule(&wg),
		fx.Invoke(
			func(lc fx.Lifecycle, runner *lambda.Runner) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						runner.Run(FunctionName)
						return nil
					},
				})
			},
		),
	)
	app.Run()
	wg.Wait()
}
