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

func GenerateLibraryInfo(db *db.DB, library_name string) models.LibraryData {
	var raw_librarylist []models.DBLibraryItem

	query := `SELECT description,
	concat(version_major, ".", version_minor, ".", version_patch) as version
	from babel.docs d
	join babel.doc_history dh
	on d.name = dh.name
	where d.name="` + library_name + `"
	ORDER BY version_major desc, version_minor desc, version_patch desc`

	db.Raw(query).Scan(&raw_librarylist)

	var versions []models.LibraryLink
	var description string

	for _, item := range raw_librarylist {
		description = item.Description
		link := models.LibraryLink{
			Version: item.Version,
			Link:    "/docs/" + description + "/" + item.Version,
		}
		versions = append(versions, link)
		fmt.Printf("%s %s\n", item.Version, item.Description)
	}

	library := models.LibraryData{
		Library:     library_name,
		Description: description,
		Links:       versions,
	}
	return library
}

func LibraryHandler(db *db.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/info/")
		fmt.Println(path)

		data := GenerateLibraryInfo(db, path)

		page := template.Must(template.ParseFiles("templates/library.html"))
		page.ExecuteTemplate(res, "library", data)
	}
}
