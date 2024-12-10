package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"backend/pkg/db"
	structure "backend/pkg/struct"
)


type FeatureImportance struct {
	FeatureName     string  `json:"featureName"`
	ImportanceValue float64 `json:"importanceValue"`
}

type WeatherImpactResponse struct {
	WeatherData       []structure.WeatherImpactData  `json:"weatherData"`
	FeatureImportance []FeatureImportance  `json:"featureImportance"`
}

func WeatherImpact(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT 
			w.year,
			w.month,
			w.avg_sunshine_duration_seconds,
			w.avg_daylight_duration_seconds,
			w.min_temperature_C,
			w.avg_temperature_C,
			w.max_temperature_C,
			w.avg_solar_irradiance_wm2,
			w.avg_relative_humidity_percent,
			w.avg_cloud_cover_percent,
			w.avg_wind_speed_kmh,
			w.total_rainfall_mm,
			SUM(m.actual_kwh) as total_kwh
		FROM weather_monthly w
		LEFT JOIN monthly_generation m 
		ON w.year = m.year AND w.month = m.month
		GROUP BY w.year, w.month
		ORDER BY w.year, w.month;
	`

	rows, err := db.Database.Query(query)
	if err != nil {
		fmt.Println("Error querying weather impact data:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var weatherImpactData []structure.WeatherImpactData

	for rows.Next() {
		var data structure.WeatherImpactData
		err := rows.Scan(
			&data.Year,
			&data.Month,
			&data.AvgSunshineDuration,
			&data.AvgDaylightDuration,
			&data.MinTemperature,
			&data.AvgTemperature,
			&data.MaxTemperature,
			&data.AvgSolarIrradiance,
			&data.AvgRelativeHumidity,
			&data.AvgCloudCover,
			&data.AvgWindSpeed,
			&data.CumulativeRainfall,
			&data.TotalPowerGeneration,
		)
		if err != nil {
			fmt.Println("Error scanning weather impact data:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		weatherImpactData = append(weatherImpactData, data)
	}

	// feature importance query
	featureQuery := `
		SELECT 
			feature_name,
			importance_value as feature_importance
		FROM feature_importance
		ORDER BY importance_value DESC;
	`

	featureRows, err := db.Database.Query(featureQuery)
	if err != nil {
		fmt.Println("Error querying feature importance:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer featureRows.Close()

	var featureImportance []FeatureImportance
	for featureRows.Next() {
		var feature FeatureImportance
		err := featureRows.Scan(&feature.FeatureName, &feature.ImportanceValue)
		if err != nil {
			fmt.Println("Error scanning feature importance:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		featureImportance = append(featureImportance, feature)
	}

	response := WeatherImpactResponse{
		WeatherData:       weatherImpactData,
		FeatureImportance: featureImportance,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
