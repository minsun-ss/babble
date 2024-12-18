package main

import (
	"babel/db"
	"babel/handlers"
	"fmt"
	"html/template"
	"net/http"
)

func webserver() {
	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/nah", handlers.HandleLibraryPage)
	http.HandleFunc("/info/", handlers.LibraryHandler)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.ExecuteTemplate(w, "index.html", handlers.HandleMenuItem())
	// })

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", db.Stuff())
	})

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	webserver()
	// let's figure out orm now

	// db.Stuff()
}
