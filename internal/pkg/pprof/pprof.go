package pprof

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
)

func StartPprof(lc fx.Lifecycle, logger *zap.Logger) {
	srv := &http.Server{Addr: "localhost:6060"}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Fatal("pprof server failed", zap.Error(err))
				}
			}()
			logger.Info("pprof server started", zap.String("address", srv.Addr))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutting down pprof server...")
			return srv.Shutdown(ctx)
		},
	})
}
