package main

import (
	"database/sql"
	_ "embed" // <-- ADD THIS BLANK IMPORT
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/DiscoMouse/cowrie-graph/internal/config"
	"github.com/DiscoMouse/cowrie-graph/internal/database"
	"github.com/DiscoMouse/cowrie-graph/internal/handlers"
	"github.com/gin-gonic/gin"
)

//go:embed embed/config.example.json
var exampleConfigFile []byte

//go:embed embed/README.txt
var readmeFile []byte

func main() {
	// First-run setup logic
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		log.Println("Configuration file not found. Generating setup files...")
		if err := os.WriteFile("config.example.json", exampleConfigFile, 0644); err != nil {
			log.Fatalf("FATAL: Failed to write config.example.json: %v", err)
		}
		if err := os.WriteFile("README.txt", readmeFile, 0644); err != nil {
			log.Fatalf("FATAL: Failed to write README.txt: %v", err)
		}
		log.Println("Setup files created. Please create 'config.json' from the example and run again.")
		return // Exit gracefully
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		// This now handles a specific error from LoadConfig if the config is invalid
		if _, ok := err.(*json.SyntaxError); ok {
			log.Fatalf("FATAL: config.json is malformed: %v", err)
		}
		log.Println(err)
		return
	}

	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v. Please check your config.json.", err)
	}
	log.Println("Successfully connected to the database!")

	store, err := database.NewStore(db, cfg.GeoDBPath)
	if err != nil {
		log.Fatalf("Failed to create database store: %v", err)
	}
	log.Println("GeoIP database loaded successfully.")

	if err := store.CreateIntelligenceTable(); err != nil {
		log.Fatalf("Failed to run database migration: %v", err)
	}
	log.Println("Database migration successful.")

	apiHandler := handlers.NewAPIHandler(store)

	router := gin.Default()
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusMovedPermanently, "/dashboard") })
	router.GET("/dashboard", func(c *gin.Context) { c.File("./static/dashboard.html") })
	charts := router.Group("/charts")
	{
		charts.GET("/attacks-by-day", func(c *gin.Context) { c.File("./static/charts/attacks-by-day.html") })
		charts.GET("/attacks-by-month", func(c *gin.Context) { c.File("./static/charts/attacks-by-month.html") })
		charts.GET("/top-10s", func(c *gin.Context) { c.File("./static/charts/top-20s.html") })
		charts.GET("/top-geo", func(c *gin.Context) { c.File("./static/charts/top-geo.html") })
		charts.GET("/bar-race", func(c *gin.Context) { c.File("./static/charts/bar-race.html") })
		charts.GET("/country-race", func(c *gin.Context) { c.File("./static/charts/country-race.html") })
	}
	api := router.Group("/api/v1")
	{
		api.GET("/top-passwords", apiHandler.GetTopPasswords)
		api.GET("/attacks-by-day", apiHandler.GetAttacksByDay)
		api.GET("/top-usernames", apiHandler.GetTopUsernames)
		api.GET("/top-ips", apiHandler.GetTopIPs)
		api.GET("/top-clients", apiHandler.GetTopClients)
		api.GET("/attacks-by-month", apiHandler.GetAttacksByMonth)
		api.GET("/top-countries", apiHandler.GetTopCountries)
		api.GET("/top-cities", apiHandler.GetTopCities)
		api.GET("/top-orgs", apiHandler.GetTopOrgs)
		api.GET("/bar-race-data", apiHandler.GetBarRaceData)
		api.GET("/country-race-data", apiHandler.GetCountryBarRaceData)
	}

	log.Println("Starting Gin server on :8080")
	router.Run(":8080")
}
