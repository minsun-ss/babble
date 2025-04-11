package handlers

import (
	"archive/zip"
	"babel/backend/internal/models"
	"babel/backend/internal/utils"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// the function generateDocsData retrieves the specific library documentation from the database
// and returns the result
func generateDocsData(db *gorm.DB, library string, version string) (*models.DBLibraryZip, error) {
	var dbZipResult models.DBLibraryZip

	values := strings.Split(version, ".")
	major, err := strconv.Atoi(values[0])
	if err != nil {
		return &dbZipResult, fmt.Errorf("failed to parse library major version: %w", err)
	}
	minor, err := strconv.Atoi(values[1])
	if err != nil {
		return &dbZipResult, fmt.Errorf("failed to parse library minor version: %w", err)
	}
	patch, err := strconv.Atoi(values[2])
	if err != nil {
		return &dbZipResult, fmt.Errorf("failed to parse library patch version: %w", err)
	}

	db.Raw(`SELECT html from babel.doc_history
		WHERE name=@library and version_major=@major
		and version_minor=@minor and version_patch=@patch`,
		sql.Named("library", library),
		sql.Named("major", major),
		sql.Named("minor", minor),
		sql.Named("patch", patch),
	).Scan(&dbZipResult)

	return &dbZipResult, nil
}

func DocsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/docs/")

		// find out library name and version from path and error out if
		// unavailable
		values := strings.Split(path, "/")
		if len(values) < 2 {
			slog.Error("not enough data to parse docs")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		library, version := values[0], values[1]

		// fix the prefix to strip out and retrieve correlationId
		prefix := "/docs/" + library + "/" + version + "/"
		requestId := r.Header.Get(RequestIDHeader)

		slog.Debug("configuration", "correlationId", requestId,
			"path", path, "library", library, "version", version,
			"prefix", prefix)

		result, err := generateDocsData(db, library, version)
		if err != nil {
			slog.Error("failed to parse docs zip", "error", err, "path", path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		zipResult := result.DataZip

		tmpFile, err := os.CreateTemp("", "zip")
		if err != nil {
			slog.Error("creating temporary zip file failed",
				"correlationId", requestId,
				"error", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(zipResult)
		if err != nil {
			slog.Error("Writing buffer into zip file failed",
				"correlationId", requestId,
				"error", err)
			http.Error(w, "Zip File Error", http.StatusInternalServerError)
			return
		}
		tmpFile.Close()

		zr, err := zip.OpenReader(tmpFile.Name())
		if err != nil {
			slog.Error("openreader failed on temp zip file",
				"correlationId", requestId,
				"error", err)
			http.Error(w, "Invalid File Error", http.StatusInternalServerError)
			return
		}
		defer zr.Close()

		ch := utils.NewCleanupFileHandler(tmpFile, requestId)
		ch.MonitorContext(r.Context())
		ch.CleanupFile()

		http.StripPrefix(prefix, http.FileServer(http.FS(zr))).ServeHTTP(w, r)
	}
}
