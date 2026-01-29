package metrics

import "github.com/prometheus/client_golang/prometheus"

func Register() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		HttpRequestDuration,
		HttpInFlight,
	)
}
