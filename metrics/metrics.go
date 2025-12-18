package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal counts total HTTP requests
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lights_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	// HTTPRequestDuration measures HTTP request duration
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "lights_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	// LightOperationsTotal counts light control operations
	LightOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "lights_operations_total",
			Help: "Total number of light control operations",
		},
		[]string{"operation", "result"},
	)

	// ActiveConnections tracks current active connections
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "lights_http_active_connections",
			Help: "Number of active HTTP connections",
		},
	)
)
