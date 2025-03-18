package server

import (
	"context"
	"time"

	"go.uber.org/fx"
)

func RunFiberServer(lc fx.Lifecycle, p ServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := p.App.Listen(":4000"); err != nil {
					panic(err) // Handle errors appropriately in production code
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			return p.App.Shutdown() // Gracefully shutdown Fiber
		},
	})
}
