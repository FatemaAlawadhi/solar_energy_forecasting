// data/random_data.go
package data

import (
	"backend/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"math"
	structure "backend/pkg/struct"
)

func FetchWeatherData() {
	// Clear the table before inserting new data
	_, err := db.Database.Exec("DELETE FROM weather_daily") 
	if err != nil {
		log.Fatalf("Error clearing weather table: %v", err)
	}

	startDate := "2015-01-01"
	endDate := "2019-12-31"

	// Parse the start and end dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		fmt.Println("Error parsing end date:", err)
		return
	}

	// Loop through each day in the date range
	for current := start; current.Before(end) || current.Equal(end); current = current.AddDate(0, 0, 1) {
		dateStr := current.Format("2006-01-02")

		// Fetch data from the API for the current date
		resp, err := http.Get(fmt.Sprintf("https://archive-api.open-meteo.com/v1/archive?latitude=26&longitude=50.55&start_date=%s&end_date=%s&hourly=temperature_2m,relative_humidity_2m,cloud_cover,wind_speed_10m,direct_normal_irradiance&daily=sunrise,sunset,daylight_duration,sunshine_duration,rain_sum&timezone=auto", dateStr, dateStr))
		if err != nil {
			fmt.Println("Error fetching data for date", dateStr, ":", err)
			continue
		}
		defer resp.Body.Close()

		var data structure.APIResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Error decoding JSON for date", dateStr, ":", err)
			continue
		}

		// Initialize variables
		startIndex := -1
		endIndex := -1
		results := make([]map[string]interface{}, 0)

		// Find daylight period
		for i, irradiance := range data.Hourly.DirectNormalIrradiance {
			if irradiance > 0 && startIndex == -1 {
				startIndex = i
			} else if irradiance == 0 && startIndex != -1 {
				endIndex = i - 1
				break
			}
		}

		// Calculate averages and min/max for the specified parameters during daylight hours
		if startIndex != -1 && endIndex != -1 {
			// Get temperature data for the period
			tempSlice := data.Hourly.Temperature2m[startIndex:endIndex+1]
			minTemp := findMin(tempSlice)
			maxTemp := findMax(tempSlice)
			avgTemp := calculateAverage(tempSlice)

			// Calculate other averages
			avgHumidity := calculateAverage(data.Hourly.RelativeHumidity2m[startIndex:endIndex+1])
			avgCloudCover := calculateAverage(data.Hourly.CloudCover[startIndex:endIndex+1])
			avgWindSpeed := calculateAverage(data.Hourly.WindSpeed10m[startIndex:endIndex+1])
			avgIrradiance := calculateAverage(data.Hourly.DirectNormalIrradiance[startIndex:endIndex+1])

			// Format sunrise and sunset times (extract only time part)
			sunrise := formatTimeOnly(data.Daily.Sunrise[0])
			sunset := formatTimeOnly(data.Daily.Sunset[0])

			// Create result with new column names
			result := map[string]interface{}{
				"date":                      dateStr,
				"sunrise_time":              sunrise,
				"sunset_time":               sunset,
				"sunshine_duration_seconds": data.Daily.SunshineDuration[0],
				"daylight_duration_seconds": data.Daily.DaylightDuration[0],
				"min_temperature_C":         minTemp,
				"avg_temperature_C":         avgTemp,
				"max_temperature_C":         maxTemp,
				"avg_solar_irradiance_wm2":  avgIrradiance,
				"avg_relative_humidity_percent": avgHumidity,
				"avg_cloud_cover_percent":    avgCloudCover,
				"avg_wind_speed_kmh":        avgWindSpeed,
				"rainfall_mm":    data.Daily.RainSum[0],
			}
			results = append(results, result)
		}

		if len(results) > 0 {
			saveToDatabase(results)
		}
	}
}

func findMin(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return math.Round(min*100) / 100
}

func findMax(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return math.Round(max*100) / 100
}

func formatTimeOnly(datetime string) string {
	t, err := time.Parse("2006-01-02T15:04", datetime)
	if err != nil {
		return ""
	}
	return t.Format("15:04")
}

func saveToDatabase(results []map[string]interface{}) {
	for _, result := range results {
		stmt, err := db.Database.Prepare(`
			INSERT INTO weather_daily (
				date, sunrise_time, sunset_time, sunshine_duration_seconds,
				daylight_duration_seconds, min_temperature_C, avg_temperature_C,
				max_temperature_C, avg_solar_irradiance_wm2, avg_relative_humidity_percent,
				avg_cloud_cover_percent, avg_wind_speed_kmh, rainfall_mm
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`)
		if err != nil {
			log.Printf("Error preparing statement: %v", err)
			continue
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			result["date"],
			result["sunrise_time"],
			result["sunset_time"],
			result["sunshine_duration_seconds"],
			result["daylight_duration_seconds"],
			result["min_temperature_C"],
			result["avg_temperature_C"],
			result["max_temperature_C"],
			result["avg_solar_irradiance_wm2"],
			result["avg_relative_humidity_percent"],
			result["avg_cloud_cover_percent"],
			result["avg_wind_speed_kmh"],
			result["rainfall_mm"],
		)
		if err != nil {
			log.Printf("Error inserting data for date %v: %v", result["date"], err)
			continue
		}
		fmt.Printf("Successfully inserted weather data for date: %v\n", result["date"])
	}
}

func calculateAverage(data []float64) float64 {
	sum := 0.0
	for _, value := range data {
		sum += value
	}
	average := sum / float64(len(data))
	return math.Round(average*100) / 100
}

func InsertMonthlyWeatherData() {
	// Clear the table before inserting new data
	_, err := db.Database.Exec("DELETE FROM weather_monthly")
	if err != nil {
		log.Fatalf("Error clearing monthly weather table: %v", err)
	}
    // Query to aggregate daily data into monthly data
    query := `
        INSERT INTO weather_monthly (
            year,
            month,
            avg_sunshine_duration_seconds,
            avg_daylight_duration_seconds,
            min_temperature_C,
            avg_temperature_C,
            max_temperature_C,
            avg_solar_irradiance_wm2,
            avg_relative_humidity_percent,
            avg_cloud_cover_percent,
            avg_wind_speed_kmh,
            total_rainfall_mm
        )
        SELECT 
            CAST(strftime('%Y', date) AS INTEGER) as year,
            CAST(strftime('%m', date) AS INTEGER) as month,
            ROUND(AVG(sunshine_duration_seconds)) as avg_sunshine_duration_seconds,
            ROUND(AVG(daylight_duration_seconds)) as avg_daylight_duration_seconds,
            MIN(min_temperature_C) as min_temperature_C,
            AVG(avg_temperature_C) as avg_temperature_C,
            MAX(max_temperature_C) as max_temperature_C,
            AVG(avg_solar_irradiance_wm2) as avg_solar_irradiance_wm2,
            AVG(avg_relative_humidity_percent) as avg_relative_humidity_percent,
            AVG(avg_cloud_cover_percent) as avg_cloud_cover_percent,
            AVG(avg_wind_speed_kmh) as avg_wind_speed_kmh,
            SUM(rainfall_mm) as total_rainfall_mm
        FROM weather_daily
        GROUP BY year, month
        ORDER BY year, month;
    `

    // Execute the query
    result, err := db.Database.Exec(query)
    if err != nil {
        log.Printf("Error aggregating monthly weather data: %v", err)
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected: %v", err)
        return
    }

    log.Printf("Successfully inserted %d monthly records", rowsAffected)
}





