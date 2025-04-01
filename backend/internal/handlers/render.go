package handlers

import (
	"babel/backend/internal/models"
	"encoding/json"
	"fmt"
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
	indexList := make(map[string][]string)
	for _, item := range dbMenuList {
		slog.Debug("loaded menu item", "project_name", item.ProjectTeam, "library", item.Library)
		indexList[item.ProjectTeam] = append(indexList[item.ProjectTeam], item.Library)
	}
	for project, libraries := range indexList {
		item := models.JsonIndexMenuItem{ProjectTeam: project, Libraries: libraries}
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

func generateLibraryLinks(db *gorm.DB, libraryName string) ([]models.JsonLibraryMenuItem, error) {
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

	// return error if nothing was retrieved from the database
	if len(dbLibraryList) == 0 {
		return nil, fmt.Errorf("no links were retrieved from database for library: %s", libraryName)
	}

	// marshal it into a json
	var jsonLibraryMenuItem []models.JsonLibraryMenuItem
	var library, projectTeam, libraryDescription string
	var versions []string
	for _, item := range dbLibraryList {
		slog.Debug("loading library items", "library", item.Library, "project_team", item.ProjectTeam, "library_description", item.LibraryDescription, "version", item.Version)
		library = item.Library
		projectTeam = item.ProjectTeam
		libraryDescription = item.LibraryDescription
		versions = append(versions, item.Version)
	}
	item := models.JsonLibraryMenuItem{Library: library, ProjectTeam: projectTeam, LibraryDescription: libraryDescription, Versions: versions}
	jsonLibraryMenuItem = append(jsonLibraryMenuItem, item)

	return jsonLibraryMenuItem, nil
}

// LibraryLinksHandler is the GET endpoint to generate library links for the front end
func LibraryLinksHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/links/")
		slog.Debug("links ", "path", path)

		// find out library name
		values := strings.Split(path, "/")
		if len(values) < 1 || values[0] == "" {
			slog.Error("not enough data to parse links")
			http.Error(w, "Internal Server Error: not enough data to parse links path", http.StatusInternalServerError)
			return
		}
		libraryName := values[0]

		data, err := generateLibraryLinks(db, libraryName)
		if err != nil {
			slog.Error("failure to fetch data", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Internal Server Error": "failure to fetch data"}`))
			return
		}
		slog.Debug("retrieved from db", "count", len(data))

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(data)
		if err != nil {
			slog.Error("failure to marshal data", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"Internal Server Error": "failed to marshal data"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	}
}
