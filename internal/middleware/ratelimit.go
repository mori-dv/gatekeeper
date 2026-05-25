package middleware

import (
	"net"
	"net/http"
	"strconv"

	"gatekeeper/internal/limiter"
	"gatekeeper/internal/metrics"
)

func RateLimit(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "invalid remote address", http.StatusInternalServerError)
				return
			}

			allowed, remaining, err := l.Allow(r.Context(), ip)
			if err != nil {
				http.Error(w, "rate limiter error", http.StatusInternalServerError)
				return
			}

			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(l.GetLimit()))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

			if !allowed {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				metrics.RateLimitedTotal.Inc()
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}