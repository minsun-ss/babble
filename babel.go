package main

import (
	"babel/handlers"
	"fmt"
	"html/template"
	"net/http"
)

func webserver() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", handlers.HandleMenuItem())
	})

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	webserver()
}
