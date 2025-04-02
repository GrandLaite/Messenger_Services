package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var msgCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "message_service_requests_total",
		Help: "Total number of requests received by message service",
	},
	[]string{"path", "method"},
)

func init() {
	prometheus.MustRegister(msgCounter)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msgCounter.WithLabelValues(r.URL.Path, r.Method).Inc()
		next.ServeHTTP(w, r)
	})
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
