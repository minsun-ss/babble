package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// Middleware struct is middleware container for the handler that
// manages logging our prometheus metrics
type Middleware struct {
	handler http.Handler
}

// ServeHTTP handles request by passing real handler and logging
// relevant details about the request: e.g., latency. Also adding
// a ksuid request header to track the state of the request
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get(RequestIDHeader)
	requestsTotal.Inc()
	if requestId == "" {
		requestId = ksuid.New().String()
		r.Header.Set(RequestIDHeader, requestId)
	}

	// populating writer and contexts
	w.Header().Set(RequestIDHeader, requestId)

	start := time.Now()
	slog.Info("Entering request", "correlationId", requestId, "path", r.URL.Path)

	mw.handler.ServeHTTP(w, r)

	duration := time.Since(start)
	requestDuration.Observe(float64(duration.Milliseconds()))
	requestLatency.Observe(float64(duration.Milliseconds()))
	slog.Info("Request stats", "correlationId", requestId, "path", r.URL.Path,
		"duration", time.Since(start),
	)
}

// NewMiddlewareHandler is wrapper around the new middleware handler
func NewMiddlewareHandler(handler http.Handler) *Middleware {
	return &Middleware{handler}
}

// LivenessHandler is the /healthz endpoint check. The liveness check checks
// to make sure the database connection is still alive.
func LivenessHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("Failed to retrieve the db connection", "err", err)
		}

		err = sqlDB.Ping()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "Service Unavailable Error", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, "Service Healthy")
	}
}

// MetricsHandler handles the custom prometheus metrics for the babel backend service
func MetricsHandler() http.Handler {
	return promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{})
}
