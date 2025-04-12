package handlers

import (
	"fmt"
	"io"
	"net/http"

	"gorm.io/gorm"
)

// listDocs generates a list from docs given a hash
func listDocs(hash_key string) {
	fmt.Println("Surprise!")
}

// updateDocs should update an existing doc or add to it if it's not there
func updateDocs(hash_key string, update_item io.Reader) {
	fmt.Println("Not surprise!")
}

// BabelAPIListHandler handles the end point to list all docs
func BabelAPIListHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET methods are allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Got past the get man;")
	}
}

func BabelAPIUpdateDocsHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET methods are allowed", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprintf(w, "Got past the get man;")
	}
}
