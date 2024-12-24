package handlers

import (
	"archive/zip"
	"babel/models"
	"babel/utils"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

func ServeZipFileHandler(db *utils.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch information about the file
		path := strings.TrimPrefix(r.URL.Path, "/docs/")
		values := strings.Split(path, "/")
		library, version := values[0], values[1]
		prefix := "/docs/" + library + "/" + version + "/"
		requestId := r.Header.Get(RequestIDHeader)

		slog.Debug("configuration", "correlationId", requestId,
			"path", path,
			"library", library,
			"version", version, "prefix", prefix)

		var dbZipResult models.DBLibraryZip
		db.Raw(`select html from babel.doc_history
			where name="traderpythonlib" and version_major="4"
			and version_minor="0" and version_patch="0"`).Scan(&dbZipResult)

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
