package main

import (
	"babel/handlers"
	"net/http"
	"os"

	"babel/utils"
	"flag"

	"log/slog"
)

func webserver(config *utils.Config) {
	dba := utils.DBPool(config)

	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/info/", handlers.LibraryHandler(dba))
	http.HandleFunc("/", handlers.IndexHandler(dba))
	// http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/docs/", handlers.ServeZipFileHandler(dba))

	http.ListenAndServe(":23456", nil)
}

var verbose int

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
	config := utils.GetConfig()

	// cli()
	slog.Info("Starting webserver...", "port", 23456)
	webserver(config)
	// let's figure out orm now

	// dba := db.DBPool(config)
	// handlers.GenerateLibraryInfo(dba, "traderpythonlib")
	// db.Stuff()
}
