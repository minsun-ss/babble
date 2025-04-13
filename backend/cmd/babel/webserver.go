package babel

import (
	"babel/backend/internal/handlers"
	"io/fs"
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// Webserver sets up the webserver and sets up appropriate handlers
func Webserver(config *Config) {
	mux := http.NewServeMux()

	// static files require moving moving down a subfolder to be
	// appropriately referenced - this endpoint is to serve the internal
	// frontend
	static, err := fs.Sub(config.BabelFS, "assets")
	if err != nil {
		// for this particular error, yes, full webserver failure preferred
		log.Fatal("static assets embedding failed:", err)
	}
	fs := http.FileServer(http.FS(static))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints - some are semi deprecated but are here for reasons
	mux.HandleFunc("/", handlers.IndexHandler(config.DB, config.BabelFS))
	mux.HandleFunc("/info/", handlers.InfoHandler(config.DB, config.BabelFS))

	// endpoints passed through to front end without handling
	mux.HandleFunc("/libraries/", handlers.DocsHandler(config.DB))

	// frontend endpoints - these are specifically for the frontend to use
	mux.HandleFunc("/internal/menu/", handlers.IndexMenuHandler(config.DB))
	mux.HandleFunc("/internal/links/", handlers.LibraryLinksHandler(config.DB))

	// Create a Huma API with the HTTP adapter & register endpoints
	api := humago.New(mux, *config.ApiCfg)
	api_grp := huma.NewGroup(api, "/api/v1")
	// huma.Register(api_grp, handlers.GreetingOperation(), handlers.Greeting)
	huma.Register(api_grp, handlers.ListOperation(), handlers.APIList)
	huma.Register(api_grp, handlers.LibraryRetriveOperation(), handlers.APIList)
	huma.Register(api_grp, handlers.LibraryRetriveOperation(), handlers.APIList)

	// liveness check & prometheus
	mux.HandleFunc("/healthz", handlers.LivenessHandler(config.DB))
	mux.Handle("/metrics", handlers.MetricsHandler())

	middlewareMux := handlers.NewMiddlewareHandler(mux)

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
