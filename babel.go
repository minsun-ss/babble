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
	"os"

	"flag"

	"log/slog"
)

//go:embed static/*
var staticFS embed.FS

//go:embed templates/*.html
var templatesFS embed.FS

// sets up webserver and appropriate handlers
func webserver(config *config.Config) {
	mux := http.NewServeMux()

	// static files require moving moving down a subfolder to be
	// appropriately referenced
	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		// for this particular error, yes, full webserver failure preferred
		log.Fatal("static assets embedding failed:", err)
	}
	fs := http.FileServer(http.FS(static))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints
	mux.HandleFunc("/", handlers.IndexHandler(config.DB, templatesFS))
	mux.HandleFunc("/info/", handlers.InfoHandler(config.DB, templatesFS))
	mux.HandleFunc("/docs/", handlers.DocsHandler(config.DB))

	// liveness check
	mux.HandleFunc("/healthz", handlers.LivenessHandler)

	middlewareMux := handlers.NewMiddleware(mux)

	slog.Info("Starting webserver...", "port", 23456)
	http.ListenAndServe(":23456", middlewareMux)
}

// Sets up the logging function and customizes verbosity as needed
func init() {
	vFlag := flag.Bool("v", false, "verbosity level")
	vvFlag := flag.Bool("vv", false, "verbosity level 2")
	vvvFlag := flag.Bool("vvv", false, "verbosity level 3")
	flag.Parse()

	var logLevel slog.Level

	if *vFlag {
		logLevel = slog.LevelWarn
	} else if *vvFlag {
		logLevel = slog.LevelInfo
	} else if *vvvFlag {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelError
	}

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Adds source file and line number
	})

	logger := slog.New(textHandler)
	slog.SetDefault(logger)
}

func main() {
	config := config.NewConfig()
	webserver(config)
}
