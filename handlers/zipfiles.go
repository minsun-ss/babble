package handlers

import (
	"archive/zip"
	"babel/db"
	"babel/models"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
)

func ServeZipFile(res http.ResponseWriter, req *http.Request) {
	check := html.EscapeString(req.URL.Path)
	fmt.Printf(check)
	filename := "docs/output.zip"
	if filename == "" {
		fmt.Println("where's my file?")
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Shit my file")
	}

	defer zr.Close()
	prefix := "/docs/"

	// http.FileServer(http.FS(zr)).ServeHTTP(res, req)
	stripped := http.StripPrefix(prefix, http.FileServer(http.FS(zr)))
	stripped.ServeHTTP(res, req)
}

func ServeZipFileHandler(db *db.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Printf("Url path: %s\n", req.URL.Path)

		// fetch information about the file
		path := strings.TrimPrefix(req.URL.Path, "/docs/")
		values := strings.Split(path, "/")
		library, version := values[0], values[1]

		fmt.Println(path)
		fmt.Printf("%s %s\n", library, version)
		prefix := "/docs/" + library + "/" + version + "/"
		fmt.Println(prefix)

		// var zipped_file models.DBLibraryZip
		var zipresult models.DBLibraryZip
		db.Raw(`select html from babel.doc_history
			where name="traderpythonlib" and version_major="4"
			and version_minor="0" and version_patch="0"`).Scan(&zipresult)

		fmt.Println("printing 4.0.0")
		zipped_file2 := zipresult.DataZip

		tmpFile, err := os.CreateTemp("", "zip")
		if err != nil {
			fmt.Printf("%v", err)
		}
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(zipped_file2)
		if err != nil {
			fmt.Printf("%v", err)
		}
		tmpFile.Close()

		zr, err := zip.OpenReader(tmpFile.Name())

		if err != nil {
			fmt.Printf("Shit my file: %v", err)
		}

		defer zr.Close()

		done := make(chan bool)
		defer func() {
			os.Remove(tmpFile.Name())
			log.Printf("Temp file %s deleted", tmpFile.Name())
			done <- true
		}()

		go func() {
			select {
			case <-req.Context().Done():
				log.Println("Client disconnected")
			case <-done:
				log.Println("Handler completed")
			}
		}()

		http.StripPrefix(prefix, http.FileServer(http.FS(zr))).ServeHTTP(res, req)
	}
}
