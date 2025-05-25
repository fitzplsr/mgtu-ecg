package metrics

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Count all http requests by status code, method and path.",
		},
		[]string{"method", "path", "status_code", "service"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of all HTTP requests by status code, method and path.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status_code", "service"},
	)

	httpRequestsInProgress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_in_progress_total",
			Help: "All the requests in progress",
		},
		[]string{"method", "service"},
	)
)

const service = "mgtu-ecg"

type MetricsMW struct {
	MW fiber.Handler
}

func NewMetricsMW() *MetricsMW {
	return &MetricsMW{MW: func(c *fiber.Ctx) error {
		method := string(append([]byte(nil), c.Method()...))
		path := string(append([]byte(nil), c.Path()...))

		httpRequestsInProgress.
			WithLabelValues(method, service).Inc()

		start := time.Now()
		err := c.Next() // вызываем хендлеры

		status := strconv.Itoa(c.Response().StatusCode())
		latency := time.Since(start).Seconds()

		httpRequestsTotal.
			WithLabelValues(method, path, status, service).Inc()

		httpRequestDuration.
			WithLabelValues(method, path, status, service).Observe(latency)

		httpRequestsInProgress.
			WithLabelValues(method, service).Dec()

		return err
	},
	}
}
