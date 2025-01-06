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

	// attempting to serve on 80
	slog.Info("Starting webserver...", "port", 80)
	go func() {
		err := http.ListenAndServe(":80", middlewareMux)
		if err != nil {
			log.Fatal("HTTP server error:", err)
		}
	}()

	// attempting to serve on 443 now
	go func() {
		cert_path := config.Cfg.GetString("CERT_PATH")
		cert_key := config.Cfg.GetString("KEY_PATH")
		if cert_path == "" || cert_key == "" {
			slog.Error("Serving https failing due to missing certs...", "port", 443)
		} else {
			slog.Info("Starting https webserver...", "port", 443)
			err = http.ListenAndServeTLS(":443", cert_path, cert_key, middlewareMux)
			if err != nil {
				slog.Warn("HTTPS server error, only serving 80 for now:", "err", err)
			}
		}
	}()

	select {}
}
