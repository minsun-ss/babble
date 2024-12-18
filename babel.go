package main

import (
	"babel/handlers"
	"fmt"
	"html/template"
	"net/http"
)

func webserver() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/info", handlers.HandleLibraryPage)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", handlers.HandleMenuItem())
	})

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	webserver()
}
