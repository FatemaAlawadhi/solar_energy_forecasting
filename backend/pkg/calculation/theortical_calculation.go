package calculation

import (
	"backend/pkg/db"
	"fmt"
	"math"
	"time"
)

const (
	inverterEfficiency = 0.915 // 91.5%
)

type Location struct {
	ID                int
	Name              string
	InstalledCapacity float64
	NumberOfPV        int
}

func CalculateTheorticalOutput() error {
	// 1. Get all locations
	locations, err := getLocations()
	if err != nil {
		return fmt.Errorf("error getting locations: %v", err)
	}

	// 2. Get monthly weather data and calculate for each month
	query := `
		SELECT strftime('%Y', date) as year, 
			   strftime('%m', date) as month,
			   AVG(sunshine_duration_seconds) as avg_sunshine,
			   AVG(avg_solar_irradiance_wm2) as avg_irradiance,
			   COUNT(*) as days_in_month
		FROM weather_daily
		GROUP BY strftime('%Y', date), strftime('%m', date)
		ORDER BY year, month
	`
	
	rows, err := db.Database.Query(query)
	if err != nil {
		return fmt.Errorf("error querying weather data: %v", err)
	}
	defer rows.Close()

	tx, err := db.Database.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the update statement
	updateStmt, err := tx.Prepare(`
		INSERT INTO monthly_generation (year, month, location_id, theoretical_kwh)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(year, month, location_id) 
		DO UPDATE SET theoretical_kwh = excluded.theoretical_kwh
	`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	defer updateStmt.Close()

	for rows.Next() {
		var year, month int
		var avgSunshine, avgIrradiance float64
		var daysInMonth int

		if err := rows.Scan(&year, &month, &avgSunshine, &avgIrradiance, &daysInMonth); err != nil {
			return fmt.Errorf("error scanning row: %v", err)
		}

		for _, loc := range locations {
			dailyOutput := loc.InstalledCapacity * inverterEfficiency * (avgSunshine * avgIrradiance) / (1000 * 3600)
			monthlyOutput := math.Round(dailyOutput * float64(daysInMonth) * 100) / 100

			// Save to database
			_, err := updateStmt.Exec(year, month, loc.ID, monthlyOutput)
			if err != nil {
				return fmt.Errorf("error updating theoretical output for location %s: %v", loc.Name, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func getLocations() ([]Location, error) {
	query := `
		SELECT id, name, installed_capacity_kw, number_of_panels 
		FROM locations 
		ORDER BY id
	`
	
	rows, err := db.Database.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var loc Location
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.InstalledCapacity, &loc.NumberOfPV); err != nil {
			return nil, err
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

func CalculateMonthlyPerformance() {
	// Clear the monthly_performance table first
	_, err := db.Database.Exec("DELETE FROM monthly_performance")
	if err != nil {
		fmt.Printf("error clearing monthly_performance table: %v", err)
	}

	// Get all locations
	locations, err := getLocations()
	if err != nil {
		fmt.Printf("error getting locations: %v", err)
	}

	// Start a transaction for batch updates
	tx, err := db.Database.Begin()
	if err != nil {
		fmt.Printf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Prepare the update statement
	updateStmt, err := tx.Prepare(`
		INSERT INTO monthly_performance (
			year, month, location_id, 
			performance_ratio, capacity_factor, output_per_pv
		)
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Printf("error preparing statement: %v", err)
	}
	defer updateStmt.Close()

	// Query to get monthly data
	query := `
		SELECT 
			year, month, location_id,
			actual_kwh, theoretical_kwh
		FROM monthly_generation
		WHERE actual_kwh IS NOT NULL 
			AND theoretical_kwh IS NOT NULL
	`

	rows, err := db.Database.Query(query)
	if err != nil {
		fmt.Printf("error querying monthly generation: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var year, month, locationID int
		var actualKWH, theoreticalKWH float64

		if err := rows.Scan(&year, &month, &locationID, &actualKWH, &theoreticalKWH); err != nil {
			fmt.Printf("error scanning row: %v", err)
		}

		// Get location's installed capacity and number of panels
		var installedCapacity float64
		var numberOfPanels int
		for _, loc := range locations {
			if loc.ID == locationID {
				installedCapacity = loc.InstalledCapacity
				numberOfPanels = loc.NumberOfPV
				break
			}
		}

		// Get hours in month
		hoursInMonth := getHoursInMonth(year, month)

		// Calculate metrics
		performanceRatio := 0.0
		if theoreticalKWH > 0 {
			performanceRatio = math.Round((actualKWH/theoreticalKWH)*1000) / 1000  // 3 decimal places
		}

		capacityFactor := 0.0
		if installedCapacity > 0 {
			capacityFactor = math.Round((actualKWH/(installedCapacity*float64(hoursInMonth)))*1000) / 1000  // 3 decimal places
		}

		outputPerPV := 0.0
		if numberOfPanels > 0 {
			outputPerPV = math.Round((actualKWH/float64(numberOfPanels))*100) / 100  // 2 decimal places
		}

		_, err = updateStmt.Exec(
			year, month, locationID,
			performanceRatio, capacityFactor, outputPerPV,
		)
		if err != nil {
			fmt.Printf("error updating performance metrics: %v", err)
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("error committing transaction: %v", err)
	}

}

func getHoursInMonth(year, month int) int {
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	return lastDay.Day() * 24
}

func CalculateYearlyPerformance() {
	// Clear the yearly_performance table first
	_, err := db.Database.Exec("DELETE FROM yearly_performance")
	if err != nil {
		fmt.Printf("error clearing yearly_performance table: %v\n", err)
		return
	}

	locations, err := getLocations()
	if err != nil {
		fmt.Printf("error getting locations: %v\n", err)
		return
	}

	tx, err := db.Database.Begin()
	if err != nil {
		fmt.Printf("error starting transaction: %v\n", err)
		return
	}
	defer tx.Rollback()

	updateStmt, err := tx.Prepare(`
		INSERT INTO yearly_performance (
			year, location_id,
			performance_ratio, capacity_factor, output_per_pv
		)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Printf("error preparing statement: %v\n", err)
		return
	}
	defer updateStmt.Close()

	// Query to get yearly sums
	query := `
		SELECT 
			year, location_id,
			SUM(actual_kwh) as yearly_actual,
			SUM(theoretical_kwh) as yearly_theoretical
		FROM monthly_generation
		WHERE actual_kwh IS NOT NULL 
		AND theoretical_kwh IS NOT NULL
		GROUP BY year, location_id
	`

	rows, err := db.Database.Query(query)
	if err != nil {
		fmt.Printf("error querying yearly generation: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var year, locationID int
		var yearlyActual, yearlyTheoretical float64

		if err := rows.Scan(&year, &locationID, &yearlyActual, &yearlyTheoretical); err != nil {
			fmt.Printf("error scanning row: %v\n", err)
			continue
		}

		var installedCapacity float64
		var numberOfPanels int
		for _, loc := range locations {
			if loc.ID == locationID {
				installedCapacity = loc.InstalledCapacity
				numberOfPanels = loc.NumberOfPV
				break
			}
		}

		// Calculate hours in year (accounting for leap years)
		hoursInYear := 8760
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			hoursInYear = 8784
		}

		performanceRatio := 0.0
		if yearlyTheoretical > 0 {
			performanceRatio = math.Round((yearlyActual/yearlyTheoretical)*1000) / 1000
		}

		capacityFactor := 0.0
		if installedCapacity > 0 {
			capacityFactor = math.Round((yearlyActual/(installedCapacity*float64(hoursInYear)))*1000) / 1000
		}

		outputPerPV := 0.0
		if numberOfPanels > 0 {
			outputPerPV = math.Round((yearlyActual/float64(numberOfPanels))*100) / 100
		}

		_, err = updateStmt.Exec(
			year, locationID,
			performanceRatio, capacityFactor, outputPerPV,
		)
		if err != nil {
			fmt.Printf("error updating yearly performance metrics: %v\n", err)
			continue
		}
	}

	if err = tx.Commit(); err != nil {
		fmt.Printf("error committing transaction: %v\n", err)
		return
	}
}


func CalculateOverallPerformance() {
    // Clear the overall_performance table first
    _, err := db.Database.Exec("DELETE FROM overall_performance")
    if err != nil {
        fmt.Printf("error clearing overall_performance table: %v\n", err)
        return
    }

    locations, err := getLocations()
    if err != nil {
        fmt.Printf("error getting locations: %v\n", err)
        return
    }

    tx, err := db.Database.Begin()
    if err != nil {
        fmt.Printf("error starting transaction: %v\n", err)
        return
    }
    defer tx.Rollback()

    updateStmt, err := tx.Prepare(`
        INSERT INTO overall_performance (
            location_id,
            start_year, end_year,
            performance_ratio, capacity_factor, output_per_pv
        )
        VALUES (?, ?, ?, ?, ?, ?)
    `)
    if err != nil {
        fmt.Printf("error preparing statement: %v\n", err)
        return
    }
    defer updateStmt.Close()

    // First, get the min and max years from the data
    var startYear, endYear int
    err = db.Database.QueryRow(`
        SELECT MIN(year), MAX(year)
        FROM monthly_generation
        WHERE actual_kwh IS NOT NULL 
        AND theoretical_kwh IS NOT NULL
    `).Scan(&startYear, &endYear)
    if err != nil {
        fmt.Printf("error getting year range: %v\n", err)
        return
    }

    // Query to get overall sums for each location
    query := `
        SELECT 
            location_id,
            SUM(actual_kwh) as total_actual,
            SUM(theoretical_kwh) as total_theoretical
        FROM monthly_generation
        WHERE actual_kwh IS NOT NULL 
        AND theoretical_kwh IS NOT NULL
        GROUP BY location_id
    `

    rows, err := db.Database.Query(query)
    if err != nil {
        fmt.Printf("error querying overall generation: %v\n", err)
        return
    }
    defer rows.Close()

    for rows.Next() {
        var locationID int
        var totalActual, totalTheoretical float64

        if err := rows.Scan(&locationID, &totalActual, &totalTheoretical); err != nil {
            fmt.Printf("error scanning row: %v\n", err)
            continue
        }

        var installedCapacity float64
        var numberOfPanels int
        for _, loc := range locations {
            if loc.ID == locationID {
                installedCapacity = loc.InstalledCapacity
                numberOfPanels = loc.NumberOfPV
                break
            }
        }

        // Calculate total hours in the period
        totalHours := calculateTotalHours(startYear, endYear)

        performanceRatio := 0.0
        if totalTheoretical > 0 {
            performanceRatio = math.Round((totalActual/totalTheoretical)*1000) / 1000
        }

        capacityFactor := 0.0
        if installedCapacity > 0 {
            capacityFactor = math.Round((totalActual/(installedCapacity*float64(totalHours)))*1000) / 1000
        }

        outputPerPV := 0.0
        if numberOfPanels > 0 {
            outputPerPV = math.Round((totalActual/float64(numberOfPanels))*100) / 100
        }

        _, err = updateStmt.Exec(
            locationID,
            startYear, endYear,
            performanceRatio, capacityFactor, outputPerPV,
        )
        if err != nil {
            fmt.Printf("error updating overall performance metrics: %v\n", err)
            continue
        }
    }

    if err = tx.Commit(); err != nil {
        fmt.Printf("error committing transaction: %v\n", err)
        return
    }
}

// Helper function to calculate total hours between years
func calculateTotalHours(startYear, endYear int) int {
    totalHours := 0
    for year := startYear; year <= endYear; year++ {
        if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
            totalHours += 8784 // Leap year
        } else {
            totalHours += 8760 // Normal year
        }
    }
    return totalHours
}