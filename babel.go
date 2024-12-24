package main

import (
	"babel/config"
	"babel/handlers"
	"net/http"
	"os"

	"flag"

	"log/slog"
)

// sets up webserver and appropriate handlers
func webserver(config *config.Config) {
	mux := http.NewServeMux()

	// static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// http endpoints
	mux.HandleFunc("/", handlers.IndexHandler(config.DBpool))
	mux.HandleFunc("/info/", handlers.LibraryHandler(config.DBpool))
	mux.HandleFunc("/docs/", handlers.DocsHandler(config.DBpool))

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
