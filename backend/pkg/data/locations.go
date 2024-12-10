package data

import (
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"backend/pkg/db"
    "log"
)

func InitializeLocations() error {
    // Clear the table before inserting new data
	_, err := db.Database.Exec("DELETE FROM locations")
	if err != nil {
		log.Printf("Error clearing locations table: %v", err)
		return err
	}

    // location data including capacity and panels
    locations := []struct {
        id                  int
        name               string
        installedCapacity float64
        numberOfPanels    int
    }{
        {1, "Awali", 1590, 6625},      
        {2, "Refinery", 2892, 12050},  
        {3, "UOB", 518.4, 2160},        
        {4, "Total System", 5000, 20835}, 
    }

    // Insert locations with their details
    for _, loc := range locations {
        _, err := db.Database.Exec(`
            INSERT INTO locations (
                id, 
                name, 
                installed_capacity_kw, 
                number_of_panels
            ) VALUES (?, ?, ?, ?)
            ON CONFLICT (name) DO UPDATE SET
                installed_capacity_kw = excluded.installed_capacity_kw,
                number_of_panels = excluded.number_of_panels,
                last_updated = CURRENT_TIMESTAMP;`,
            loc.id,
            loc.name,
            loc.installedCapacity,
            loc.numberOfPanels,
        )
        if err != nil {
            return fmt.Errorf("error inserting location %s: %v", loc.name, err)
        }
    }

    fmt.Println("Successfully initialized locations table")
    return nil
}
