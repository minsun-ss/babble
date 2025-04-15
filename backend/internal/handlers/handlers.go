/*
Package handlers contains all the handlers used for all outgoing endpoints
in the application. A list of all implemented handlers is in cmd/babel/webserver.go.
*/
package handlers

import (
	"github.com/prometheus/client_golang/prometheus"
)

const RequestIDHeader = "X-Request-ID"

var (
	registry prometheus.Registry

	requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "babel",
		Name:      "requests_total",
		Help:      "Total number of requests received",
	})

	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace: "babel",
		Name:      "request_duration_milliseconds",
		Help:      "Duration of babel request in milliseconds",
		Buckets:   []float64{100, 500, 1000, 2000, 5000},
	})

	requestLatency = prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "babel",
		Name:      "request_latency_milliseconds",
		Help:      "Request latency in milliseconds",
		Objectives: map[float64]float64{
			0.5: 0.05, 0.9: 0.01, 0.99: .001,
		},
	})
)

// init sets up the prometheus registry and registers custom metrics needed.
func init() {
	registry := prometheus.NewRegistry()
	prometheus.DefaultGatherer = registry
	registry.MustRegister(requestsTotal)
	registry.MustRegister(requestDuration)
	registry.MustRegister(requestLatency)
}
