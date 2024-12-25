/*
	Package main implements a web server that serves versioned user-facing

documentation for TA-managed libraries.

The webserver is organized into several components:
  - handlers: HTTP request handlers, health checks and middleware
  - models: Data structures and database models
  - config: Application configuration management
  - templates: HTML templates
  - static: Static files such as css and htmx js files
  - utils: Miscellaneous functions
*/
package main

import (
	"babel/config"
	"babel/handlers"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"log/slog"
)

//go:embed static/* templates/*.html
var babelFS embed.FS

var (
	logLevel slog.Level
)

// sets up webserver and appropriate handlers
func webserver(config *config.Config) {
	mux := http.NewServeMux()

	// static files require moving moving down a subfolder to be
	// appropriately referenced
	static, err := fs.Sub(config.BabelFS, "static")
	if err != nil {
		// for this particular error, yes, full webserver failure preferred
		log.Fatal("static assets embedding failed:", err)
	}
	fs := http.FileServer(http.FS(static))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints
	mux.HandleFunc("/", handlers.IndexHandler(config.DB, babelFS))
	mux.HandleFunc("/info/", handlers.InfoHandler(config.DB, babelFS))
	mux.HandleFunc("/docs/", handlers.DocsHandler(config.DB))

	// liveness check
	mux.HandleFunc("/healthz", handlers.LivenessHandler)

	middlewareMux := handlers.NewMiddleware(mux)

	slog.Info("Starting webserver...", "port", 23456)
	http.ListenAndServe(":23456", middlewareMux)
}

func main() {
	config := config.NewConfig(babelFS)
	webserver(config)
}
