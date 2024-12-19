package main

import (
	"babel/db"
	"babel/handlers"
	"fmt"
	"net/http"

	"babel/utils"
)

func webserver(config *utils.Config) {
	dba := db.DBPool(config)

	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	// set up redirects and templates
	// tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/info/", handlers.LibraryHandler(dba))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.ExecuteTemplate(w, "index.html", handlers.GenerateMenuFields(dba))
	// })
	http.HandleFunc("/", handlers.IndexHandler(dba))

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	config := utils.GetConfig()

	webserver(config)
	// let's figure out orm now

	// dba := db.DBPool(config)
	// handlers.GenerateLibraryInfo(dba, "traderpythonlib")
	// db.Stuff()
}
