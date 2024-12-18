package handlers

import (
	"babel/db"
	"babel/models"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

func HandleLibraryPage(res http.ResponseWriter, req *http.Request) {
	data := models.LibraryData{
		Library:     "traderpythonlib",
		Description: "A trading library",
		Links: []models.LibraryLink{
			{Version: "1.2.3", Link: "/docs"},
			{Version: "1.2.1", Link: "/docs"},
		},
	}
	page := template.Must(template.ParseFiles("templates/library.html"))
	page.ExecuteTemplate(res, "library", data)
}

func LibraryHandler(res http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/info/")
	fmt.Println(path)

	data := db.FetchAllLibraryInfo(path)

	page := template.Must(template.ParseFiles("templates/library.html"))
	page.ExecuteTemplate(res, "library", data)
}
