package api

import (
	"encoding/json"
	"net/http"
	"backend/pkg/db/queries"
)

func SystemConfiguration(w http.ResponseWriter, r *http.Request) {
	locations, err := queries.GetLocationData()
	if err != nil {
		http.Error(w, "Error fetching location data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(locations); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
