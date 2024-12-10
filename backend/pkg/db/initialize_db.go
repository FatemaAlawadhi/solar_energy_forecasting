package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Database *sql.DB

func InitializeDb() {
	var err error
	Database, err = sql.Open("sqlite3", "../../pkg/db/app.db")
	if err != nil {
		log.Fatalf("Error initializing new database: %v", err)
	}

	if err = Database.Ping(); err != nil {
		log.Fatalf("New database is not reachable: %v", err)
	}

	createTables := `

CREATE TABLE IF NOT EXISTS weather_daily (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date DATE NOT NULL UNIQUE,
    sunrise_time TIME NOT NULL,
    sunset_time TIME NOT NULL,
    sunshine_duration_seconds INTEGER,
    daylight_duration_seconds INTEGER,
    min_temperature_C DECIMAL(10, 2),
    avg_temperature_C DECIMAL(10, 2),
    max_temperature_C DECIMAL(10, 2),
    avg_solar_irradiance_wm2 DECIMAL(10, 2),
    avg_relative_humidity_percent DECIMAL(10, 2),
    avg_cloud_cover_percent DECIMAL(10, 2),
    avg_wind_speed_kmh DECIMAL(10, 2),
    rainfall_mm DECIMAL(10, 2)
);


CREATE TABLE IF NOT EXISTS weather_monthly (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year INT NOT NULL,
    month INT NOT NULL CHECK (month >= 1 AND month <= 12),
    avg_sunshine_duration_seconds INTEGER,
    avg_daylight_duration_seconds INTEGER,
    min_temperature_C DECIMAL(10, 2),
    avg_temperature_C DECIMAL(10, 2),
    max_temperature_C DECIMAL(10, 2),
    avg_solar_irradiance_wm2 DECIMAL(10, 2),
    avg_relative_humidity_percent DECIMAL(10, 2),
    avg_cloud_cover_percent DECIMAL(10, 2),
    avg_wind_speed_kmh DECIMAL(10, 2),
    total_rainfall_mm DECIMAL(10, 2),
    UNIQUE(year, month)
);

CREATE TABLE IF NOT EXISTS locations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(50) UNIQUE NOT NULL CHECK (name IN ('Awali', 'Refinery', 'UOB', 'Total System')),
    installed_capacity_kw DECIMAL(10, 2),
    number_of_panels INTEGER,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS monthly_generation (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year INT NOT NULL,
    month INT NOT NULL CHECK (month >= 1 AND month <= 12),
    location_id INTEGER NOT NULL,
    actual_kwh DECIMAL(10, 2),
    theoretical_kwh DECIMAL(10, 2),
    predicted_kwh DECIMAL(10, 2),
    FOREIGN KEY (location_id) REFERENCES locations(id),
    UNIQUE(year, month, location_id)
);

CREATE TABLE IF NOT EXISTS monthly_performance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year INT NOT NULL,
    month INT NOT NULL CHECK (month >= 1 AND month <= 12),
    location_id INTEGER NOT NULL,
    performance_ratio DECIMAL(10, 2),
    capacity_factor DECIMAL(10, 2),
    output_per_pv DECIMAL(10, 2),
    FOREIGN KEY (location_id) REFERENCES locations(id),
    UNIQUE(year, month, location_id)
);

CREATE TABLE IF NOT EXISTS yearly_performance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    year INT NOT NULL,
    location_id INTEGER NOT NULL,
    performance_ratio DECIMAL(10, 2),
    capacity_factor DECIMAL(10, 2),
    output_per_pv DECIMAL(10, 2),
    FOREIGN KEY (location_id) REFERENCES locations(id),
    UNIQUE(year, location_id)
);

CREATE TABLE IF NOT EXISTS overall_performance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_year INT NOT NULL,
    end_year INT NOT NULL,
    location_id INTEGER NOT NULL,
    performance_ratio DECIMAL(10, 2),
    capacity_factor DECIMAL(10, 2),
    output_per_pv DECIMAL(10, 2),
    FOREIGN KEY (location_id) REFERENCES locations(id),
    UNIQUE(start_year, end_year, location_id)
);

CREATE TABLE IF NOT EXISTS feature_importance (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    feature_name TEXT NOT NULL,
    importance_value DECIMAL(10, 4),
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(feature_name)
);`

	_, err = Database.Exec(createTables)
	if err != nil {
		log.Fatalf("Error creating tables in new database: %v", err)
	}

}


