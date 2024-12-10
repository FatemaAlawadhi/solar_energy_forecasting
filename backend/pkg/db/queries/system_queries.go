package queries

import (
	"log"
	"backend/pkg/db"
)

type LocationData struct {
	Name             string
	InstalledCapacity float64
	NumberOfPanels   int
}

func GetLocationData() ([]LocationData, error) {
	query := `
		SELECT name, installed_capacity_kw, number_of_panels
		FROM locations
		WHERE name IN ('Awali', 'Refinery', 'UOB', 'Total System')
	`

	rows, err := db.Database.Query(query)
	if err != nil {
		log.Printf("Error querying location data: %v", err)
		return nil, err
	}
	defer rows.Close()

	var locations []LocationData
	for rows.Next() {
		var location LocationData
		if err := rows.Scan(&location.Name, &location.InstalledCapacity, &location.NumberOfPanels); err != nil {
			log.Printf("Error scanning location data: %v", err)
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}
