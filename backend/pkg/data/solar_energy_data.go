package data

import (
	"fmt"
	"strconv"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/xuri/excelize/v2"
	"backend/pkg/db"
	"log"
)

// ImportDataFromExcel reads data from an Excel file and inserts it into the database
func ImportEnergyData() {

	_, err := db.Database.Exec("DELETE FROM monthly_generation") 
	if err != nil {
		log.Fatalf("Error clearing solar energy table: %v", err)
	}

	f, err := excelize.OpenFile("../../pkg/db/BapcoSolarEnergy.xlsx")
	if err != nil {
		fmt.Println("Error opening Excel file:", err)
		return
	}
	defer f.Close()

	// Read UOB data
	if err := importUOBData(f); err != nil {
		fmt.Println("Error importing UOB data:", err)
		return
	}

	// Read Refinery data
	if err := importRefineryData(f); err != nil {
		fmt.Println("Error importing Refinery data:", err)
		return
	}

	// Read Awali data
	if err := importAwaliData(f); err != nil {
		fmt.Println("Error importing Awali data:", err)
		return
	}

	// Calculate and insert total system data
	if err := calculateTotalSystem(); err != nil {
		fmt.Println("Error calculating total system:", err)
		return
	}

	fmt.Println("Successfully imported all energy data")
}

func importUOBData(f *excelize.File) error {
	rows, err := f.GetRows("UOB")
	if err != nil {
		return err
	}

	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}

		year, _ := strconv.Atoi(row[0])
		month, _ := strconv.Atoi(row[1])
		totalUOB, _ := strconv.ParseFloat(row[2], 64)

		_, err := db.Database.Exec(`
			INSERT INTO monthly_generation (year, month, location_id, actual_kwh)
			VALUES (?, ?, 3, ?)
			ON CONFLICT (year, month, location_id) 
			DO UPDATE SET actual_kwh = excluded.actual_kwh;`,
			year, month, totalUOB)
		if err != nil {
			return err
		}
	}
	return nil
}

func importRefineryData(f *excelize.File) error {
	rows, err := f.GetRows("Refinery")
	if err != nil {
		return err
	}

	for _, row := range rows[1:] {
		if len(row) < 7 {
			continue
		}

		year, _ := strconv.Atoi(row[0])
		month, _ := strconv.Atoi(row[1])
		totalRefinery, _ := strconv.ParseFloat(row[6], 64) 

		_, err := db.Database.Exec(`
			INSERT INTO monthly_generation (year, month, location_id, actual_kwh)
			VALUES (?, ?, 2, ?)
			ON CONFLICT (year, month, location_id) 
			DO UPDATE SET actual_kwh = excluded.actual_kwh;`,
			year, month, totalRefinery)
		if err != nil {
			return err
		}
	}
	return nil
}

func importAwaliData(f *excelize.File) error {
	rows, err := f.GetRows("Awali")
	if err != nil {
		return err
	}

	for _, row := range rows[1:] {
		if len(row) < 13 {
			continue
		}

		year, _ := strconv.Atoi(row[0])
		month, _ := strconv.Atoi(row[1])
		totalAwali, _ := strconv.ParseFloat(row[12], 64) 

		_, err := db.Database.Exec(`
			INSERT INTO monthly_generation (year, month, location_id, actual_kwh)
			VALUES (?, ?, 1, ?)
			ON CONFLICT (year, month, location_id) 
			DO UPDATE SET actual_kwh = excluded.actual_kwh;`,
			year, month, totalAwali)
		if err != nil {
			return err
		}
	}
	return nil
}

func calculateTotalSystem() error {
	// Calculate and insert total system data
	query := `
		INSERT INTO monthly_generation (year, month, location_id, actual_kwh)
		SELECT 
			m1.year,
			m1.month,
			4 as location_id,
			SUM(m1.actual_kwh) as actual_kwh
		FROM monthly_generation m1
		WHERE m1.location_id IN (1, 2, 3)  -- Awali, Refinery, UOB
		GROUP BY m1.year, m1.month
		ON CONFLICT (year, month, location_id) 
		DO UPDATE SET actual_kwh = excluded.actual_kwh;
	`

	_, err := db.Database.Exec(query)
	return err
}

