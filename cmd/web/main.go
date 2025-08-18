package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/DiscoMouse/cowrie-graph/internal/config"
	"github.com/DiscoMouse/cowrie-graph/internal/database"
	"github.com/DiscoMouse/cowrie-graph/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := sql.Open("mysql", cfg.DatabaseDSN)
	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal("Failed to connect to database: ", err)
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
		charts.GET("/top-10s", func(c *gin.Context) { c.File("./static/charts/top-10s.html") })
		charts.GET("/world-map", func(c *gin.Context) { c.File("./static/charts/world-map.html") })
		charts.GET("/top-geo", func(c *gin.Context) { c.File("./static/charts/top-geo.html") })
		// --- ADD THIS ROUTE ---
		charts.GET("/bar-race", func(c *gin.Context) {
			c.File("./static/charts/bar-race.html")
		})
	}

	api := router.Group("/api/v1")
	{
		api.GET("/top-passwords", apiHandler.GetTopPasswords)
		api.GET("/attacks-by-day", apiHandler.GetAttacksByDay)
		api.GET("/top-usernames", apiHandler.GetTopUsernames)
		api.GET("/top-ips", apiHandler.GetTopIPs)
		api.GET("/top-clients", apiHandler.GetTopClients)
		api.GET("/attacks-by-month", apiHandler.GetAttacksByMonth)
		api.GET("/attacks-by-location", apiHandler.GetAttacksByLocation)
		api.GET("/top-countries", apiHandler.GetTopCountries)
		api.GET("/top-cities", apiHandler.GetTopCities)
		api.GET("/top-orgs", apiHandler.GetTopOrgs)
		// --- ADD THIS API ROUTE ---
		api.GET("/bar-race-data", apiHandler.GetBarRaceData)
	}

	log.Println("Starting Gin server on :8080")
	router.Run(":8080")
}
