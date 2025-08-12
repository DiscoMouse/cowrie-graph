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

type Config struct {
	DatabaseDSN string `json:"database_dsn"`
}

type TopPassword struct {
	Password string `json:"password"`
	Count    int    `json:"count"`
}

func getTopPasswords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT password, COUNT(*) as count 
			FROM auth 
			GROUP BY password 
			ORDER BY count DESC 
			LIMIT 10;
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		passwords := []TopPassword{}
		for rows.Next() {
			var p TopPassword
			if err := rows.Scan(&p.Password, &p.Count); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			passwords = append(passwords, p)
		}
		c.JSON(http.StatusOK, passwords)
	}
}

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

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

	router := gin.Default()

	// --- THIS IS THE CORRECTED LINE ---
	// Serve the index.html file for the root URL path
	router.StaticFile("/", "./static/index.html")

	// API Routes
	api := router.Group("/api/v1")
	{
		api.GET("/top-passwords", getTopPasswords(db))
	}

	log.Println("Starting Gin server on :8080")
	router.Run(":8080")
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
