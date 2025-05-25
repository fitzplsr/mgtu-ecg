package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
)

func InvokeMetricsServer(lc fx.Lifecycle) {
	prometheus.MustRegister(
		EcgUploadsTotal,
		EcgProcessedTotal,
		EcgProcessingSeconds,

		httpRequestsTotal,
		httpRequestDuration,
		httpRequestsInProgress,
	)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				http.Handle("/metrics", promhttp.Handler())
				http.ListenAndServe(":2112", nil)
			}()
			return nil
		},
	})
}
