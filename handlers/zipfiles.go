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
	filename := "docs/test.zip"
	if filename == "" {
		fmt.Println("where's my file?")
	}
	zr, err := zip.OpenReader(filename)
	if err != nil {
		fmt.Println("Shit my file")
	}

	defer zr.Close()
	http.StripPrefix("/check/whatever/3.2.1", http.FileServer(http.FS(zr))).ServeHTTP(res, req)
}

func ServeZipFileHandler(db *db.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// fmt.Println(req.URL.)
		path := strings.TrimPrefix(req.URL.Path, "/docs/")
		values := strings.Split(path, "/")

		library, version := values[0], values[1]

		fmt.Println(path)
		fmt.Printf("%s %s", library, version)

		// var zipped_file models.DBLibraryZip
		var zipresult models.DBLibraryZip
		db.Raw(`select html from babel.doc_history
			where name="traderpythonlib" and version_major="3"
			and version_minor="2" and version_patch="1"`).Scan(&zipresult)

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

		http.StripPrefix(req.URL.Path, http.FileServer(http.FS(zr))).ServeHTTP(res, req)
	}
}
