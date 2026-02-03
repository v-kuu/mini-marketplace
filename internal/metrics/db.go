package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	DbSemaphoreWaitDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "marketplace",
			Subsystem: "db",
			Name: "semaphore_wait_duration_seconds",
			Help: "Time spent waiting for database semaphore",
			Buckets: prometheus.DefBuckets,
		},
	)

	DbSemaphoreInUse = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "marketplace",
			Subsystem: "db",
			Name: "semaphore_in_use",
			Help: "Current number of database operations in progress",
		},
	)
)
