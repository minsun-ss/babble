package main

import (
	"babel/handlers"
	"fmt"
	"html/template"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func webserver() {
	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", handlers.ServeZipFile)
	http.HandleFunc("/nah", handlers.HandleLibraryPage)
	http.HandleFunc("/info/", handlers.LibraryHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", handlers.HandleMenuItem())
	})

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	// webserver()
	// let's figure out orm now
	dsn := "myuser:mypassword@tcp(host.docker.internal:3306)/babel?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	var result []map[string]interface{}
	db.Raw("SELECT library, description FROM babel.libraries").Scan(&result)
	fmt.Printf("%+v\n", result)

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
}
