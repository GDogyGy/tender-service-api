package prometheusMiddleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "acme",
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		},
		[]string{"handler", "method", "status"},
	)

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
	}, []string{"handler", "method", "status"})
)

type PrometheusMiddleware struct {
}

func NewPrometheusMiddleware() *PrometheusMiddleware {
	return &PrometheusMiddleware{}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func init() {
	prometheus.MustRegister(httpRequestsTotal, httpDuration)
}

func (h *PrometheusMiddleware) PrometheusMiddleware(handlerName string, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)

		httpDuration.WithLabelValues(
			handlerName,
			r.Method,
			status,
		).Observe(duration)

		httpRequestsTotal.WithLabelValues(
			handlerName,
			r.Method,
			status,
		).Inc()
	})
}
