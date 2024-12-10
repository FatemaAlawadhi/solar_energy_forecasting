package api

import (
	"encoding/json"
	"net/http"
	"backend/pkg/db/queries"
	structure "backend/pkg/struct"
)

func TotalPowerGeneration(w http.ResponseWriter, r *http.Request) {
	response := structure.PowerGenerationResponse{
		LastMonth: queries.GetLastMonthPowerGeneration("Total System"),
		LastYear:  queries.GetLastYearPowerGeneration("Total System"),
		Forecast:  queries.GetPowerGenerationForecast("Total System"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func AwaliPowerGeneration(w http.ResponseWriter, r *http.Request) {
	response := structure.PowerGenerationResponse{
		LastMonth: queries.GetLastMonthPowerGeneration("Awali"),
		LastYear:  queries.GetLastYearPowerGeneration("Awali"),
		Forecast:  queries.GetPowerGenerationForecast("Awali"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func RefineryPowerGeneration(w http.ResponseWriter, r *http.Request) {
	response := structure.PowerGenerationResponse{
		LastMonth: queries.GetLastMonthPowerGeneration("Refinery"),
		LastYear:  queries.GetLastYearPowerGeneration("Refinery"),
		Forecast:  queries.GetPowerGenerationForecast("Refinery"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func UOBPowerGeneration(w http.ResponseWriter, r *http.Request) {
	response := structure.PowerGenerationResponse{
		LastMonth: queries.GetLastMonthPowerGeneration("UOB"),
		LastYear:  queries.GetLastYearPowerGeneration("UOB"),
		Forecast:  queries.GetPowerGenerationForecast("UOB"),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}