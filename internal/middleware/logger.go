package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(logg *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         200,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start).Milliseconds()
			requestID := r.Header.Get("X-Request-ID")

			logg.Info(
				"http request",
				"method", r.Method,
				"path", chi.RouteContext(r.Context()).RoutePattern(),
				"status", rw.status,
				"duration_ms", duration,
				"request_id", requestID,
			)
		})
	}
}