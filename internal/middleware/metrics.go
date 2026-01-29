package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/v-kuu/mini-marketplace/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler, path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{
			ResponseWriter: w,
			status: http.StatusOK,
		}

		start := time.Now()
		metrics.HttpInFlight.Inc()
		defer metrics.HttpInFlight.Dec()

		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rec.status)

		metrics.HttpRequestsTotal.WithLabelValues(
			r.Method,
			path,
			status,
		).Inc()

		metrics.HttpRequestDuration.WithLabelValues(
			r.Method,
			path,
		).Observe(duration)
	})
}
