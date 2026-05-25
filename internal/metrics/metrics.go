package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gatekeeper_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	RateLimitedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "gatekeeper_rate_limit_blocked_total",
			Help: "Total rate limited requests",
		},
	)
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "gatekeeper_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: []float64{
			0.005,
			0.01,
			0.025,
			0.05,
			0.1,
			0.25,
			0.5,
			1,
			2,
			5,
		},
	},
	[]string{"method", "route"},
)

func Init() {
	prometheus.MustRegister(RequestsTotal)
	prometheus.MustRegister(RateLimitedTotal)
	prometheus.MustRegister(RequestDuration)
}