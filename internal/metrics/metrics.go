package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestCount    *prometheus.CounterVec
	errorCount      *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func NewMetrics(serviceName string) (*Metrics, *prometheus.Registry) {
	constLabels := prometheus.Labels{"service": serviceName}

	reg := prometheus.NewRegistry()

	m := &Metrics{
		requestCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "http_request_total",
			Help:        "Total number of http requests",
			ConstLabels: constLabels,
		}, []string{"method", "path", "status"}),

		errorCount: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name:        "http_errors_total",
			Help:        "Total number of http errors",
			ConstLabels: constLabels,
		}, []string{"method", "path", "status"}),

		requestDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:        "http_request_duration_seconds",
			Help:        "Http request duration in seconds",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: constLabels,
		}, []string{"method", "path"}),
	}

	reg.MustRegister(m.requestCount)
	reg.MustRegister(m.errorCount)
	reg.MustRegister(m.requestDuration)

	return m, reg
}

func (m *Metrics) RecordRequest(method, path string, statusCode int) {
	m.requestCount.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}

func (m *Metrics) RecordError(method, path string, statusCode int) {
	m.errorCount.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}

func (m *Metrics) RecordRequestDuration(method, path string, duration time.Duration) {
	m.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}
