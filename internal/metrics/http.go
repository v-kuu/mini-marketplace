package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "marketplace",
			Subsystem: "http",
			Name: "requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "marketplace",
			Subsystem: "http",
			Name: "request_duration_seconds",
			Help: "HTTP request latency",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	HttpInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "marketplace",
			Subsystem: "http",
			Name: "in_flight_requests",
			Help: "Current number of in-flight HTTP requests",
		},
	)
)
