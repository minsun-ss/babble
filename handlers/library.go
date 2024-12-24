package handlers

import (
	"babel/models"
	"babel/utils"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

// generates the information to generate the menu
func generateLibraryInfo(db *utils.DB, library_name string) models.LibraryData {
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
			Link:    "/docs/" + description + "/" + item.Version + "/",
		}
		versions = append(versions, link)
		slog.Debug("items", "version", item.Version, "description", item.Description)
	}

	library := models.LibraryData{
		Library:     library_name,
		Description: description,
		Links:       versions,
	}
	return library
}

func LibraryHandler(db *utils.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		path := strings.TrimPrefix(req.URL.Path, "/info/")
		slog.Info("Library handler", "path", path)

		data := generateLibraryInfo(db, path)

		page := template.Must(template.ParseFiles("templates/library.html"))
		page.ExecuteTemplate(res, "library", data)
	}
}
