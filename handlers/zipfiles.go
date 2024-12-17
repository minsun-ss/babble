package handlers

import (
	"archive/zip"
	"fmt"
	"html"
	"net/http"
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
	http.StripPrefix("/docs/", http.FileServer(http.FS(zr))).ServeHTTP(res, req)
}
