package handlers

import (
	"babel/models"
	"babel/utils"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

// generates the information to generate the webpage menu
func generateLibraryInfo(db *utils.DB, libraryName string) models.LibraryData {
	var rawLibraryList []models.DBLibraryItem

	query := `SELECT description,
	concat(version_major, ".", version_minor, ".", version_patch) as version
	from babel.docs d
	join babel.doc_history dh
	on d.name = dh.name
	where d.name="` + libraryName + `"
	ORDER BY version_major desc, version_minor desc, version_patch desc`

	db.Raw(query).Scan(&rawLibraryList)

	var versions []models.LibraryLink
	var description string

	for _, item := range rawLibraryList {
		description = item.Description
		link := models.LibraryLink{
			Version: item.Version,
			Link:    "/docs/" + description + "/" + item.Version + "/",
		}
		versions = append(versions, link)
		slog.Debug("items", "version", item.Version, "description", item.Description)
	}

	library := models.LibraryData{
		Library:     libraryName,
		Description: description,
		Links:       versions,
	}
	return library
}

func LibraryHandler(db *utils.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/info/")

		data := generateLibraryInfo(db, path)

		page := template.Must(template.ParseFiles("templates/library.html"))
		page.ExecuteTemplate(w, "library", data)
	}
}
