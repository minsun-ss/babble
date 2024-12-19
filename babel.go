package main

import (
	"babel/db"
	"babel/handlers"
	"fmt"
	"net/http"
	"strings"

	"babel/utils"
	"flag"
	"log"
)

var verbose int

func webserver(config *utils.Config) {
	dba := db.DBPool(config)

	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	// http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/info/", handlers.LibraryHandler(dba))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.ExecuteTemplate(w, "index.html", handlers.GenerateMenuFields(dba))
	// })
	http.HandleFunc("/", handlers.IndexHandler(dba))
	http.HandleFunc("/docs/", handlers.ServeZipFileHandler(dba))
	// http.HandleFunc("/check/", handlers.ServeZipFile)

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	flag.Func("v", "verbosity level (use -v, -vv, -vvv)", func(s string) error {
		verbose = strings.Count(s, "v")
		return nil
	})
	flag.Parse()

	config := utils.GetConfig()
	fmt.Println(verbose)
	switch verbose {
	case 1:
		log.SetLevel(log.InfoLevel)
	case 2:
		log.SetLevel(log.DebugLevel)
	case 3:
		log.SetLevel(log.TraceLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	webserver(config)
	// let's figure out orm now

	// dba := db.DBPool(config)
	// handlers.GenerateLibraryInfo(dba, "traderpythonlib")
	// db.Stuff()
}
