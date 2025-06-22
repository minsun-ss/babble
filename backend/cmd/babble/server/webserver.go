package main

import (
	"babble/backend/internal/api"
	"babble/backend/internal/config"
	"babble/backend/internal/handlers"

	"io/fs"
	"log"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// Webserver sets up the webserver and sets up appropriate handlers
func Webserver(config *config.Config) {
	mux := http.NewServeMux()

	// static files require moving moving down a subfolder to be
	// appropriately referenced - this endpoint is to serve the internal
	// frontend
	static, err := fs.Sub(config.BabbleFS, "assets")
	if err != nil {
		// for this particular error, yes, full webserver failure preferred
		log.Fatal("static assets embedding failed:", err)
	}
	fs := http.FileServer(http.FS(static))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints - some are semi deprecated but are here for reasons
	mux.HandleFunc("/", handlers.IndexHandler(config.DB, config.BabbleFS))
	mux.HandleFunc("/info/", handlers.InfoHandler(config.DB, config.BabbleFS))

	// endpoints passed through to front end without handling
	mux.HandleFunc("/libraries/", handlers.DocsHandler(config.DB))

	// frontend endpoints - these are specifically for the frontend to use
	mux.HandleFunc("/internal/menu/", handlers.IndexMenuHandler(config.DB))
	mux.HandleFunc("/internal/links/", handlers.LibraryLinksHandler(config.DB))

	// Create a Huma API with the HTTP adapter & register endpoints
	babbleAPI := humago.New(mux, *config.ApiCfg)
	babbleAPIGroup := huma.NewGroup(babbleAPI, "/api/v1")
	huma.Register(babbleAPIGroup, api.ListLibrariesOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.PostLibrariesOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.GetLibraryOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.GetLibraryVersionOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.PatchLibraryOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.PatchLibraryVersionOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.DeleteLibraryOperation(), handlers.APIListHandler)
	huma.Register(babbleAPIGroup, api.DeleteLibraryVersionOperation(), handlers.APIListHandler)

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
