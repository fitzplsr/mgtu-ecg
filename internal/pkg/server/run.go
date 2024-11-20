package server

import (
	"context"
	"go.uber.org/fx"
	"time"
)

func RunFiberServer(lc fx.Lifecycle, p ServerParams) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := p.App.Listen(":3000"); err != nil {
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
