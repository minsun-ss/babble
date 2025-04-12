package babel

import (
	"babel/backend/internal/handlers"
	"context"
	"fmt"
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

	// Create a Huma API with the HTTP adapter
	api := humago.New(mux, *config.ApiCfg)

	// Register Huma endpoints
	huma.Get(api, "/api/v1/{name}", func(ctx context.Context, input *struct {
		Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
	}) (*handlers.GreetingOutput, error) {
		resp := &handlers.GreetingOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	// api.Register(huma.Operation{
	// 	OperationID: "get-docs",
	// 	Summary:     "Get documentation",
	// 	Method:      http.MethodGet,
	// 	Path:        "/api/docs",
	// 	Handler: func(ctx huma.Context) {
	// 		// Your handler implementation
	// 		ctx.Response().JSON(http.StatusOK, map[string]interface{}{
	// 			"message": "Documentation endpoint",
	// 		})
	// 	},
	// })

	// Mount the Huma API on your mux
	// This registers all Huma routes on your ServeMux
	// humaHandler := api.Handler()
	// mux.Handle("/api/", http.StripPrefix("/api", humaHandler))

	// end user endpoints
	// mux.HandleFunc("/api/v1/list/", handlers.BabelAPIListHandler(config.DB))
	// mux.HandleFunc("/api/v1/docs/", handlers.BabelAPIListHandler(config.DB))

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
