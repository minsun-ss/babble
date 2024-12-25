package handlers

import (
	"babel/models"
	"embed"
	"html/template"
	"log/slog"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// generates the information to generate the webpage menu
func generateLibraryInfo(db *gorm.DB, libraryName string) models.PageLibraryData {
	var rawLibraryList []models.DBLibraryItem

	query := `SELECT description,
	concat(version_major, ".", version_minor, ".", version_patch) as version
	from babel.docs d
	join babel.doc_history dh
	on d.name = dh.name
	where d.name="` + libraryName + `"
	ORDER BY version_major desc, version_minor desc, version_patch desc`

	db.Raw(query).Scan(&rawLibraryList)

	var versions []models.PageLibraryLink
	var description string

	for _, item := range rawLibraryList {
		description = item.Description
		link := models.PageLibraryLink{
			Version: item.Version,
			Link:    "/docs/" + description + "/" + item.Version + "/",
		}
		versions = append(versions, link)
		slog.Debug("items", "version", item.Version, "description", item.Description)
	}

	library := models.PageLibraryData{
		Library:     libraryName,
		Description: description,
		Links:       versions,
	}
	return library
}

// LibraryHandler handles the full list of versions available and briefly
// describes the library
func InfoHandler(db *gorm.DB, babelFS embed.FS) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/info/")

		data := generateLibraryInfo(db, path)

		page := template.Must(template.ParseFS(babelFS, "templates/library.html"))
		page.ExecuteTemplate(w, "library", data)
	}
}
