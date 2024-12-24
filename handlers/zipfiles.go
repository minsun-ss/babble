package handlers

import (
	"archive/zip"
	"babel/config"
	"babel/models"
	"babel/utils"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func DocsHandler(db *config.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch information about the file
		path := strings.TrimPrefix(r.URL.Path, "/docs/")

		values := strings.Split(path, "/")
		library, version := values[0], values[1]
		values = strings.Split(version, ".")
		major, err := strconv.Atoi(values[0])
		if err != nil {
			slog.Error("failed to parse library major version", "error", err, "path", path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		minor, err := strconv.Atoi(values[1])
		if err != nil {
			slog.Error("failed to parse library minor version", "error", err, "path", path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		patch, err := strconv.Atoi(values[2])
		if err != nil {
			slog.Error("failed to parse library patch version", "error", err, "path", path)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		prefix := "/docs/" + library + "/" + version + "/"
		requestId := r.Header.Get(RequestIDHeader)

		slog.Debug("configuration", "correlationId", requestId,
			"path", path, "library", library, "version", version,
			"prefix", prefix, "major", major, "minor", minor, "patch", patch)

		var dbZipResult models.DBLibraryZip
		db.Raw(`SELECT html from babel.doc_history
			WHERE name=@library and version_major=@major
			and version_minor=@minor and version_patch=@patch`,
			sql.Named("library", library),
			sql.Named("major", major),
			sql.Named("minor", minor),
			sql.Named("patch", patch),
		).Scan(&dbZipResult)

		zipResult := dbZipResult.DataZip

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
		}
		tmpFile.Close()

		zr, err := zip.OpenReader(tmpFile.Name())
		if err != nil {
			slog.Error("OpenReader failed on temp zip file",
				"correlationId", requestId,
				"error", err)
		}
		defer zr.Close()

		ch := utils.NewCleanupFileHandler(tmpFile, requestId)
		ch.MonitorContext(r.Context())
		ch.CleanupFile()

		http.StripPrefix(prefix, http.FileServer(http.FS(zr))).ServeHTTP(w, r)
	}
}
