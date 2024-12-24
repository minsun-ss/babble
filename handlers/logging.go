package handlers

import (
	"log/slog"
	"net/http"
	"time"
)

// Middleware struct is middleware handler to manage logging
type Middleware struct {
	handler http.Handler
}

// ServeHTTP handles request by passing real handler and logging
// relevant details about the request: e.g., latency
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	mw.handler.ServeHTTP(w, r)
	slog.Info("Stats", "path", r.URL.Path,
		"duration", time.Since(start),
	)
}

// Constructs a new middleware handler
func NewMiddleware(handler http.Handler) *Middleware {
	return &Middleware{handler}
}
