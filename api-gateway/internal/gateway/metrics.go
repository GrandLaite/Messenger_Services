package gateway

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var requestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_gateway_requests_total",
		Help: "Total number of requests handled by the API Gateway",
	},
	[]string{"path", "method"},
)

func init() {
	prometheus.MustRegister(requestCount)
}

func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	requestCount.WithLabelValues(r.URL.Path, r.Method).Inc()
	promhttp.Handler().ServeHTTP(w, r)
}
