package handlers

import (
	"html/template"
	"net/http"
)

type LibraryLink struct {
	Version string
	Link    string
}

type LibraryData struct {
	Library     string
	Description string
	Links       []LibraryLink
}

func HandleLibraryPage(res http.ResponseWriter, req *http.Request) {
	data := LibraryData{
		Library:     "traderpythonlib",
		Description: "A trading library",
		Links: []LibraryLink{
			{Version: "1.2.3", Link: "/docs"},
			{Version: "1.2.1", Link: "/docs"},
		},
	}
	page := template.Must(template.ParseFiles("templates/library.html"))
	page.ExecuteTemplate(res, "library", data)
}
