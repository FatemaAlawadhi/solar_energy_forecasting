package queries 

import (
	"fmt"
	"backend/pkg/db"
	structure "backend/pkg/struct"
    "log"
    "database/sql"
)

// GetLastYearPowerGeneration returns the last power generation value for a specific location
func GetLastYearPowerGeneration(location string) string {
    var value sql.NullFloat64
    query := `
        WITH LastYear AS (
            SELECT year 
            FROM monthly_generation mg
            JOIN locations l ON mg.location_id = l.id
            WHERE l.name = ?
            ORDER BY year DESC 
            LIMIT 1
        )
        SELECT SUM(mg.actual_kwh)
        FROM monthly_generation mg
        WHERE mg.year = (SELECT year FROM LastYear)
        AND mg.location_id = (SELECT id FROM locations WHERE name = ?)
    `

    err := db.Database.QueryRow(query, location, location).Scan(&value)
    if err != nil {
        fmt.Printf("error getting last yearly power generation for %s: %v", location, err)
        return FormatPowerValue(0)
    }

    return FormatPowerValue(value.Float64)
}

// GetLastMonthPowerGeneration returns the last power generation value based on location
func GetLastMonthPowerGeneration(location string) string {
    var value sql.NullFloat64
    query := `
        SELECT actual_kwh
        FROM monthly_generation mg
        JOIN locations l ON mg.location_id = l.id
        WHERE l.name = ?
        ORDER BY year DESC, month DESC
        LIMIT 1
    `

    err := db.Database.QueryRow(query, location).Scan(&value)
    if err != nil {
        fmt.Printf("error getting last monthly power generation for %s: %v", location, err)
        return FormatPowerValue(0)
    }

    return FormatPowerValue(value.Float64)
}

// GetPowerGenerationForecast returns actual and predicted power generation values for a location
func GetPowerGenerationForecast(location string) []structure.ForecastResult {
    var results []structure.ForecastResult
    var query string
    var args []interface{}
    
    if location == "Total System" {
        query = `
            SELECT 
                mg.year,
                mg.month,
                SUM(mg.actual_kwh) as actual_kwh,
                SUM(mg.predicted_kwh) as predicted_kwh
            FROM monthly_generation mg
            JOIN locations l ON mg.location_id = l.id
            WHERE l.name != 'Total System'
            GROUP BY mg.year, mg.month
            ORDER BY mg.year DESC, mg.month DESC
        `
    } else {
        query = `
            SELECT 
                mg.year,
                mg.month,
                mg.actual_kwh,
                mg.predicted_kwh
            FROM monthly_generation mg
            JOIN locations l ON mg.location_id = l.id
            WHERE l.name = ?
            ORDER BY mg.year DESC, mg.month DESC
        `
        args = append(args, location)
    }

    rows, err := db.Database.Query(query, args...)
    if err != nil {
        fmt.Printf("error querying forecast data for %s: %v\n", location, err)
        return results
    }
    defer rows.Close()

    for rows.Next() {
        var year, month int
        var actual, predicted sql.NullFloat64
        
        if err := rows.Scan(&year, &month, &actual, &predicted); err != nil {
            fmt.Printf("error scanning forecast row for %s: %v\n", location, err)
            continue
        }

        result := structure.ForecastResult{
            Year:  year,
            Month: month,
            Actual: actual.Float64,
        }
        
        // Only set predicted if it's not NULL
        if predicted.Valid {
            result.Predicted = predicted.Float64
        }
        
        results = append(results, result)
    }
    
    return results
}

// FormatPowerValue converts kWh to the most appropriate unit (kWh, MWh, GWh, etc.) and returns as formatted string
func FormatPowerValue(valueInKWh float64) string {
    switch {
    case valueInKWh >= 1_000_000_000: // Billion kWh -> TWh
        return fmt.Sprintf("%.2f TWh", valueInKWh/1_000_000_000)
    case valueInKWh >= 1_000_000: // Million kWh -> GWh
        return fmt.Sprintf("%.2f GWh", valueInKWh/1_000_000)
    case valueInKWh >= 1_000: // Thousand kWh -> MWh
        return fmt.Sprintf("%.2f MWh", valueInKWh/1_000)
    default: // kWh
        return fmt.Sprintf("%.2f kWh", valueInKWh)
    }
}

func TotalPowerGeneration() (float64, float64, float64) {
    var totalUOB, totalRefinery, totalAwali float64

    // Query to get total for each location from monthly_generation
    query := `
        SELECT 
            l.name,
            COALESCE(SUM(mg.actual_kwh), 0) as total_generation
        FROM locations l
        LEFT JOIN monthly_generation mg ON l.id = mg.location_id
        WHERE l.name IN ('UOB', 'Refinery', 'Awali')
        GROUP BY l.name
    `

    rows, err := db.Database.Query(query)
    if err != nil {
        log.Printf("Error querying total power generation: %v", err)
        return 0, 0, 0
    }
    defer rows.Close()

    for rows.Next() {
        var location string
        var total float64
        if err := rows.Scan(&location, &total); err != nil {
            log.Printf("Error scanning total power generation row: %v", err)
            continue
        }

        switch location {
        case "UOB":
            totalUOB = total
        case "Refinery":
            totalRefinery = total
        case "Awali":
            totalAwali = total
        }
    }

    return totalUOB, totalRefinery, totalAwali
}