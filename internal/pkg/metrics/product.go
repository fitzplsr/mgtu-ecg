package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	EcgUploadsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "mgtu_ecg",
			Subsystem: "business",
			Name:      "ecg_uploads_total",
			Help:      "Всего загружено ЭКГ-файлов",
		})
	EcgProcessedTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "mgtu_ecg",
			Subsystem: "business",
			Name:      "ecg_processed_total",
			Help:      "Результат обработки ЭКГ",
		},
		[]string{"status"}, // success | failed
	)
	EcgProcessingSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "mgtu_ecg",
			Subsystem: "business",
			Name:      "ecg_processing_seconds",
			Help:      "Время обработки одной записи",
			Buckets:   prometheus.DefBuckets,
		})
)
