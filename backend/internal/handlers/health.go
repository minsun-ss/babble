package handlers

import (
	"net/http"
)

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// TODO: do i need a readiness probe?
