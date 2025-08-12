package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// This struct matches the structure of our config.json file
type Config struct {
	DatabaseDSN string `json:"database_dsn"`
}

func loadConfig() (Config, error) {
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

func main() {
	// --- Load Configuration ---
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v. Did you create config.json?", err)
	}
	// --- End of Load Configuration ---

	// --- Database Connection ---
	db, err := sql.Open("mysql", config.DatabaseDSN)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	log.Println("Successfully connected to the database!")
	// --- End of Database Connection ---

	router := gin.Default()
	router.StaticFS("/", http.Dir("./static"))

	log.Println("Starting Gin server on :8080")
	router.Run(":8080")
}
