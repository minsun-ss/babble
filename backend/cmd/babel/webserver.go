package babel

import (
	"babel/backend/internal/handlers"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
)

// Webserver sets up the webserver and sets up appropriate handlers
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

	// prometheus endpoint
	mux.Handle("/metrics", handlers.HandleMetrics())

	// frontend endpoints
	mux.HandleFunc("/api/menu/", handlers.IndexMenuHandler(config.DB))
	mux.HandleFunc("/api/links/", handlers.LibraryLinksHandler(config.DB))

	// liveness check
	mux.HandleFunc("/healthz", handlers.LivenessHandler)

	middlewareMux := handlers.NewMiddleware(mux)

	// attempting to serve on 80
	slog.Info("Starting webserver...", "port", 80)
	go func() {
		err := http.ListenAndServe(":80", middlewareMux)
		if err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	select {}
}
