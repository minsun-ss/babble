package babel

import (
	"babel/internal/handlers"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
)

// sets up webserver and appropriate handlers
func Webserver(config *Config) {
	mux := http.NewServeMux()

	// static files require moving moving down a subfolder to be
	// appropriately referenced
	static, err := fs.Sub(config.BabelFS, "assets")
	if err != nil {
		// for this particular error, yes, full webserver failure preferred
		log.Fatal("static assets embedding failed:", err)
	}
	fs := http.FileServer(http.FS(static))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints
	mux.HandleFunc("/", handlers.IndexHandler(config.DB, config.BabelFS))
	mux.HandleFunc("/info/", handlers.InfoHandler(config.DB, config.BabelFS))
	mux.HandleFunc("/docs/", handlers.DocsHandler(config.DB))
	mux.Handle("/metrics", handlers.HandleMetrics())

	// liveness check
	mux.HandleFunc("/healthz", handlers.LivenessHandler)

	middlewareMux := handlers.NewMiddleware(mux)

	slog.Info("Starting webserver...", "port", 23456)
	http.ListenAndServe(":23456", middlewareMux)
}
