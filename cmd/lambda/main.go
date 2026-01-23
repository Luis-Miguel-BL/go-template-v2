package main

import (
	"sync"

	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
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
		_fx.LambdaModule(FunctionName),
		_fx.ApplicationModule(&wg),
	)
	app.Run()
	wg.Wait()
}
