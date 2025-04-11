package handlers

import (
	"babel/backend/internal/models"
	"fmt"
	"net/http"
)

// listDocs generates a list from docs given a hash
func listDocs(hash_key string) {
	fmt.Println("Surprise!")
}

// updateDocs should update an existing doc or add to it if it's not there
func updateDocs(hash_key string, update_item models.JsonUpdateItem) {

	fmt.Println("Not surprise!")
}

// ReceiveUpdate is the handler endpoint
func ReceiveUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
