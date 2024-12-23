package main

import (
	"babel/db"
	"babel/handlers"
	"fmt"
	"net/http"
	"strings"

	"babel/utils"
	"flag"

	"github.com/sirupsen/logrus"
)

var verbose int
var logger = logrus.New()

func webserver(config *utils.Config) {
	dba := db.DBPool(config)

	// set up static endpoint to serve styles
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/info/", handlers.LibraryHandler(dba))
	http.HandleFunc("/", handlers.IndexHandler(dba))
	// http.HandleFunc("/docs/", handlers.ServeZipFile)
	fmt.Println("Now setting up the ServeZipFileHandler")
	http.HandleFunc("/docs/", handlers.ServeZipFileHandler(dba))

	fmt.Println("Starting server at :23456...")
	http.ListenAndServe(":23456", nil)
}

func cli() {
	flag.Func("v", "verbosity level (use -v, -vv, -vvv)", func(s string) error {
		verbose = strings.Count(s, "v")
		return nil
	})
	flag.Parse()

	// fmt.Println(verbose)

	// switch verbose {
	// case 1:
	// 	log.SetLevel(log.InfoLevel)
	// case 2:
	// 	log.SetLevel(log.DebugLevel)
	// case 3:
	// 	log.SetLevel(log.TraceLevel)
	// default:
	// 	log.SetLevel(log.WarnLevel)
	// }
}
func main() {
	config := utils.GetConfig()

	// cli()
	webserver(config)
	// let's figure out orm now

	// dba := db.DBPool(config)
	// handlers.GenerateLibraryInfo(dba, "traderpythonlib")
	// db.Stuff()
}
