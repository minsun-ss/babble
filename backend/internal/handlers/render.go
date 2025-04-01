package handlers

import (
	"babel/backend/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

// generateLibraryList generates the list of libraries and teams that the libraries belong to
func generateLibraryList(db *gorm.DB) []models.JsonIndexMenuItem {
	var dbMenuList []models.DBIndexMenuItem

	db.Raw(`
		SELECT project_team, name
		FROM babel.docs
		WHERE hidden=0
		ORDER BY project_team;`).Scan(&dbMenuList)

	// marshal it into a json
	var jsonMenuList []models.JsonIndexMenuItem
	for _, item := range dbMenuList {
		slog.Debug("loaded menu item", "project_name", item.ProjectTeam, "library", item.Library)
		item := models.JsonIndexMenuItem{ProjectTeam: item.ProjectTeam, Library: item.Library}
		jsonMenuList = append(jsonMenuList, item)
	}
	return jsonMenuList
}

// IndexMenuHandler is the GET endpoint to generate menus for the front end
func IndexMenuHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := generateLibraryList(db)
		// log some record
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to marshal data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}

func generateLibraryLinks(db *gorm.DB, libraryName string) []models.JsonLibraryMenuItem {
	slog.Debug("Attempting to generate library links", "library", libraryName)
	var dbLibraryList []models.DBLibraryMenuItem

	query := `SELECT d.name, d.project_team, description,
	concat(version_major, ".", version_minor, ".", version_patch) as version
	from babel.docs d
	join babel.doc_history dh
	on d.name = dh.name
	where d.name="` + libraryName + `"
	ORDER BY version_major desc, version_minor desc, version_patch desc`

	db.Raw(query).Scan(&dbLibraryList)

	// marshal it into a json
	var jsonLibraryMenuItem []models.JsonLibraryMenuItem
	for _, item := range dbLibraryList {
		slog.Debug("loading library items", "library", item.Library, "project_team", item.ProjectTeam, "library_description", item.LibraryDescription, "version", item.Version)

		item := models.JsonLibraryMenuItem{Library: item.Library, ProjectTeam: item.ProjectTeam, LibraryDescription: item.LibraryDescription, Version: item.Version}
		jsonLibraryMenuItem = append(jsonLibraryMenuItem, item)
	}
	return jsonLibraryMenuItem
}

// LibraryLinksHandler is the GET endpoint to generate library links for the front end
func LibraryLinksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/links/")

		// find out library name
		values := strings.Split(path, "/")
		if len(values) < 1 {
			slog.Error("not enough data to parse docs")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		libraryName := values[0]

		data := generateLibraryLinks(db, libraryName)

		// log some record
		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to marshal data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
