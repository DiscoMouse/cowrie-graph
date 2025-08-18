package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config struct to hold our configuration
type Config struct {
	DatabaseDSN string `json:"database_dsn"`
	GeoDBPath   string `json:"geo_db_path"`
}

// LoadConfig simply loads and parses config.json
func LoadConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return config, fmt.Errorf("failed to parse config.json: %w", err)
	}

	return config, nil
}
