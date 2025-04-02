package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var authCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "auth_service_requests_total",
	Help: "Total number of requests received by auth service",
}, []string{"path", "method"})

func init() {
	prometheus.MustRegister(authCounter)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authCounter.WithLabelValues(r.URL.Path, r.Method).Inc()
		next.ServeHTTP(w, r)
	})
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
