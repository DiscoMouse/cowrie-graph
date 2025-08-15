package config

import (
	"encoding/json"
	"os"
)

// Config struct to hold our configuration
type Config struct {
	DatabaseDSN string `json:"database_dsn"`
	GeoDBPath   string `json:"geo_db_path"` // <-- ADD THIS LINE
}

// LoadConfig reads the configuration from config.json
func LoadConfig() (Config, error) {
	var config Config
	file, err := os.Open("config.json")
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
