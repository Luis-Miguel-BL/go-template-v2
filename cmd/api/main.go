package main

import (
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/modules"
	"go.uber.org/fx"
)

func main() {
	var wg sync.WaitGroup
	app := fx.New(
		modules.RootModule(&wg),
		modules.HttpModule(&wg),
		modules.ApplicationModule,
	)
	app.Run()
	wg.Wait()
}
