package main

import (
	"archive/zip"
	"fmt"
	"html"
	"html/template"
	"net/http"
)

type PageData struct {
	Options  []Option
	Selected string
}

type Option struct {
	Value string
	Label string
}

type Library struct {
}

type MenuItem struct {
	Title    string
	Link     string
	Children []MenuItem
	More     string
}

func serve(res http.ResponseWriter, req *http.Request) {
	check := html.EscapeString(req.URL.Path)
	fmt.Printf(check)
	filename := "docs/test.zip"
	if filename == "" {
		fmt.Println("where's my file?")
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Shit my file")
	}

	defer zr.Close()
	http.StripPrefix("/docs/", http.FileServer(http.FS(zr))).ServeHTTP(res, req)
}

func webserver() {
	tmpl := template.Must(template.ParseGlob("templates/*.html"))
	http.HandleFunc("/docs/", serve)

	// Menu data structure
	menu := []MenuItem{
		{
			Title: "Menu",
			Link:  "#",
		},
		{
			Title: "traderpythonlib",
			Link:  "/docs",
			Children: []MenuItem{
				{Title: "Latest", Link: "/docs"},
				{Title: "1.29.0", Link: "/products/new"},
				{Title: "1.28.0", Link: "/products/categories"},
			},
			More: "/docs",
		},
		{
			Title: "deskbot",
			Link:  "#",
			Children: []MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "3.0.0", Link: "/users/new"},
				{Title: "2.9.0", Link: "/users/groups"},
			},
		},
		{
			Title: "fndmoodeng",
			Link:  "#",
			Children: []MenuItem{
				{Title: "Latest", Link: "/users"},
				{Title: "1.0.0", Link: "/users/new"},
				{Title: "0.9.0", Link: "/users/groups"},
			},
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", menu)
	})

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func main() {
	webserver()
}
