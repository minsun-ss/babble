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
// details
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	mw.handler.ServeHTTP(w, r)
	slog.Info("Request Time", "duration", time.Since(start))
}

// Constructs a new middleware handler
func NewMiddleware(handler http.Handler) *Middleware {
	return &Middleware{handler}
}
