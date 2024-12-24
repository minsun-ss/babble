package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/segmentio/ksuid"
)

// Middleware struct is middleware handler to manage logging
type Middleware struct {
	handler http.Handler
}

const RequestIDHeader = "X-Request-ID"

// ServeHTTP handles request by passing real handler and logging
// relevant details about the request: e.g., latency. Also adding
// a ksuid request header to track the state of the request
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	requestId := r.Header.Get(RequestIDHeader)
	if requestId == "" {
		requestId = ksuid.New().String()
		r.Header.Set(RequestIDHeader, requestId)
	}

	// populating writer and contexts
	w.Header().Set(RequestIDHeader, requestId)

	start := time.Now()
	slog.Info("Entering request", "correlationId", requestId, "path", r.URL.Path)

	mw.handler.ServeHTTP(w, r)

	slog.Info("Request stats", "correlationId", requestId, "path", r.URL.Path,
		"duration", time.Since(start),
	)
}

// Constructs a new middleware handler
func NewMiddleware(handler http.Handler) *Middleware {
	return &Middleware{handler}
}
