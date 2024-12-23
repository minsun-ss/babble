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

// func ServeZipFile(res http.ResponseWriter, req *http.Request) {
// 	check := html.EscapeString(req.URL.Path)
// 	fmt.Printf(check)
// 	filename := "docs/output.zip"
// 	if filename == "" {
// 		fmt.Println("where's my file?")
// 	}
// 	zr, err := zip.OpenReader(filename)
// 	if err != nil {
// 		fmt.Println("Shit my file")
// 	}

// 	defer zr.Close()
// 	prefix := "/docs/"

// 	// http.FileServer(http.FS(zr)).ServeHTTP(res, req)
// 	stripped := http.StripPrefix(prefix, http.FileServer(http.FS(zr)))
// 	stripped.ServeHTTP(res, req)
// }

func ServeZipFileHandler(db *utils.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		slog.Info("Url path", "url", req.URL.Path)

		// fetch information about the file
		path := strings.TrimPrefix(req.URL.Path, "/docs/")
		values := strings.Split(path, "/")
		library, version := values[0], values[1]

		slog.Debug(path)
		slog.Debug("Configuration", "library", library, "version", version)
		prefix := "/docs/" + library + "/" + version + "/"
		slog.Debug("Prefix used", "prefix", prefix)

		// var zipped_file models.DBLibraryZip
		var zipresult models.DBLibraryZip
		db.Raw(`select html from babel.doc_history
			where name="traderpythonlib" and version_major="4"
			and version_minor="0" and version_patch="0"`).Scan(&zipresult)

		zipped_file2 := zipresult.DataZip

		tmpFile, err := os.CreateTemp("", "zip")
		if err != nil {
			slog.Error("creating temporary zip file failed", "error", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(zipped_file2)
		if err != nil {
			slog.Error("Writing zipped file failed", "error", err)
		}
		tmpFile.Close()

		zr, err := zip.OpenReader(tmpFile.Name())

		if err != nil {
			slog.Error("Shit my zipped file failed to open in OpenReader", "error", err)
		}

		defer zr.Close()

		done := make(chan bool)
		defer func() {
			os.Remove(tmpFile.Name())
			slog.Info("Temp file deleted", "filename", tmpFile.Name())
			done <- true
		}()

		go func() {
			select {
			case <-req.Context().Done():
				slog.Info("Client disconnected")
			case <-done:
				slog.Info("Handler completed")
			}
		}()

		http.StripPrefix(prefix, http.FileServer(http.FS(zr))).ServeHTTP(res, req)
	}
}
