package handlers

import (
	"log/slog"
	"net/http"

	"gorm.io/gorm"
)

// LivenessHandler is the /healthz endpoint check. The liveness check checks
// to make sure the database connection is still alive.
func LivenessHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			slog.Error("Failed to retrieve the db connection", "err", err)
		}

		err = sqlDB.Ping()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Header().Set("Content-Type", "text/plain")
			http.Error(w, "Service Unavailable Error", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "Service Healthy", http.StatusOK)
	}
}
