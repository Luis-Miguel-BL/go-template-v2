package main

import (
	"sync"

	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"go.uber.org/fx"
)

func main() {
	var wg sync.WaitGroup
	app := fx.New(
		_fx.RootModule,
		_fx.ConsumerModule(&wg),
		_fx.ApplicationModule(&wg),
	)
	app.Run()
	wg.Wait()
}
