package data

import (
	"backend/pkg/calculation"
	"fmt"
	"log"
	"os/exec"
	"backend/pkg/db"
)

func isTableEmpty(tableName string) bool {
	var count int
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	err := db.Database.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error checking if table %s is empty: %v", tableName, err)
		return false
	}
	return count == 0
}

func executePythonScript(scriptPath string) {
	cmd := exec.Command("python3", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing script %s: %v\nOutput: %s", scriptPath, err, output)
	} else {
		log.Printf("Output of script %s:\n%s", scriptPath, output)
	}
}

func FillDb() {
	// To fill daily_weather table
	if isTableEmpty("weather_daily") {
		log.Println("Filling table: daily_weather")
		FetchWeatherData()
	}

	// To fill monthly_weather table
	if isTableEmpty("weather_monthly") {
		log.Println("Filling table: monthly_weather")
		InsertMonthlyWeatherData()
	}

	// To fill locations table
	if isTableEmpty("locations") {
		log.Println("Filling table: locations")
		InitializeLocations()
	}

	// To fill monthly_generation table
	if isTableEmpty("monthly_generation") {
		log.Println("Filling table: monthly_generation")
		ImportEnergyData()
		executePythonScript("../../pkg/model/monthly/random_forest_model.py")
		calculation.CalculateTheorticalOutput()
	}

	// To fill feature importance
	if isTableEmpty("feature_importance") {
		log.Println("Filling table: feature_importance")
		executePythonScript("../../pkg/model/monthly/weather_only_model.py")
	}

	// To fill monthly_performance
	if isTableEmpty("monthly_performance") {
		log.Println("Filling table: monthly_performance")
		calculation.CalculateMonthlyPerformance()
	}

	// To fill yearly_performance
	if isTableEmpty("yearly_performance") {
		log.Println("Filling table: yearly_performance")
		calculation.CalculateYearlyPerformance()
	}

	// To fill overall_performance
	if isTableEmpty("overall_performance") {
		log.Println("Filling table: overall_performance")
		calculation.CalculateOverallPerformance()
	}
}
