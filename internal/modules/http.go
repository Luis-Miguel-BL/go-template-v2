package modules

import (
	"context"
	"sync"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/service"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/auth/jwt"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/controller"
	"go.uber.org/fx"
)

func HttpModule(wg *sync.WaitGroup) fx.Option {
	return fx.Module("http",
		fx.Provide(
			http.NewServer,
			controller.NewLeadController,
			controller.NewAuthController,
			jwt.NewJWTHelper[service.TokenClaims],
		),
		fx.Invoke(
			func(lc fx.Lifecycle, server *http.Server) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						server.Run(ctx)
						return nil
					},
					OnStop: func(ctx context.Context) error {
						return server.Stop(ctx)
					},
				})
			},
		),
	)
}
