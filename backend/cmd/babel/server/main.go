/*
Package main implements a web server that serves versioned user-facing documentation for TA-managed libraries.

The webserver is organized into several components:
  - handlers: HTTP request handlers, health checks and middleware
  - models: Data structures and database models
  - config: Application configuration management
  - templates: HTML templates
  - static: Static files such as css and htmx js files
  - utils: Miscellaneous functions

The location of main.go is largely to support the use of the embed directive,
which does not support relative paths or in fact any path outside the current
directory and its subdirectories.
*/
package main

import (
	"babel/backend/internal/config"
	"embed"
	"flag"
	"os"

	"log/slog"
)

//go:embed assets
var babelFS embed.FS

var (
	logLevel slog.Level
)

// Init sets up the logging function and customizes verbosity as needed
func init() {
	// turns out golang doesn't really have count flags per se
	vFlag := flag.Bool("v", false, "verbosity level 1")
	vvFlag := flag.Bool("vv", false, "verbosity level 2")
	vvvFlag := flag.Bool("vvv", false, "verbosity level 3")
	flag.Parse()

	if *vFlag {
		logLevel = slog.LevelWarn
	} else if *vvFlag {
		logLevel = slog.LevelInfo
	} else if *vvvFlag {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true, // Adds source file and line number
	})

	logger := slog.New(textHandler)
	slog.SetDefault(logger)
}

func main() {
	config := config.NewConfig(babelFS)
	Webserver(config)
}
