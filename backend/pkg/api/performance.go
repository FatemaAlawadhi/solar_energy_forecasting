package api

import (
	"backend/pkg/db"
	"backend/pkg/db/queries"
	structs "backend/pkg/struct"
	"encoding/json"
	"fmt"
	"net/http"
)

func Performance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := structs.PerformanceResponse{}

	// Get Monthly Generation
	rows, err := db.Database.Query(queries.GetMonthlyGeneration)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying monthly generation: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var gen structs.Generation
		err := rows.Scan(
			&gen.Year, &gen.Month, &gen.LocationID, &gen.LocationName,
			&gen.ActualKWH, &gen.TheoreticalKWH,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning monthly generation: %v", err), http.StatusInternalServerError)
			return
		}
		response.MonthlyGeneration = append(response.MonthlyGeneration, gen)
	}

	// Get Monthly Performance
	rows, err = db.Database.Query(queries.GetMonthlyPerformance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying monthly performance: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var perf structs.Performance
		err := rows.Scan(
			&perf.Year, &perf.Month, &perf.LocationID, &perf.LocationName,
			&perf.PerformanceRatio, &perf.CapacityFactor, &perf.OutputPerPV,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning monthly performance: %v", err), http.StatusInternalServerError)
			return
		}
		response.MonthlyPerformance = append(response.MonthlyPerformance, perf)
	}

	// Get Yearly Performance
	rows, err = db.Database.Query(queries.GetYearlyPerformance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying yearly performance: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var perf structs.Performance
		err := rows.Scan(
			&perf.Year, &perf.LocationID, &perf.LocationName,
			&perf.PerformanceRatio, &perf.CapacityFactor, &perf.OutputPerPV,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning yearly performance: %v", err), http.StatusInternalServerError)
			return
		}
		response.YearlyPerformance = append(response.YearlyPerformance, perf)
	}

	// Get Overall Performance
	rows, err = db.Database.Query(queries.GetOverallPerformance)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error querying overall performance: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var perf structs.Performance
		err := rows.Scan(
			&perf.StartYear, &perf.EndYear, &perf.LocationID, &perf.LocationName,
			&perf.PerformanceRatio, &perf.CapacityFactor, &perf.OutputPerPV,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error scanning overall performance: %v", err), http.StatusInternalServerError)
			return
		}
		response.OverallPerformance = append(response.OverallPerformance, perf)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
