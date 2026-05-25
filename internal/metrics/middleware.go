package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = 200
	}

	return r.ResponseWriter.Write(b)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		rec := &statusRecorder{
			ResponseWriter: w,
			status:         0,
		}

		next.ServeHTTP(rec, r)

		route := chi.RouteContext(r.Context()).RoutePattern()

		duration := time.Since(start).Seconds()

		RequestsTotal.WithLabelValues(
			r.Method,
			route,
			strconv.Itoa(rec.status),
		).Inc()

		RequestDuration.WithLabelValues(
			r.Method,
			route,
		).Observe(duration)
	})
}