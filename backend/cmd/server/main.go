package main

import (
	"backend/pkg/db"
	"fmt"
	"backend/pkg/api"
	"net/http"
	"backend/pkg/data"
)

func enableCORS(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers for all responses including errors
        w.Header().Set("Access-Control-Allow-Origin", "*")  // Allow all origins in development
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Max-Age", "3600")
        
        // Handle preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Call the actual handler
        handler(w, r)
    }
}

func main() {
	fmt.Println("APP Started")
    // Then initialize the new database
    db.InitializeDb()
    defer db.Database.Close()
	data.FillDb()


	http.HandleFunc("/api/environment-impact", enableCORS(api.EnvironmentalImpact))
	http.HandleFunc("/api/weather-impact", enableCORS(api.WeatherImpact))
	http.HandleFunc("/api/total-power-generation", enableCORS(api.TotalPowerGeneration))
	http.HandleFunc("/api/awali-power-generation", enableCORS(api.AwaliPowerGeneration))
	http.HandleFunc("/api/uob-power-generation", enableCORS(api.UOBPowerGeneration))
	http.HandleFunc("/api/refinery-power-generation", enableCORS(api.RefineryPowerGeneration))
	http.HandleFunc("/api/performance", enableCORS(api.Performance))
	http.HandleFunc("/api/system-configuration", enableCORS(api.SystemConfiguration))

	//Start the server on port 8080
	fmt.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	fmt.Println("App ended")
}
