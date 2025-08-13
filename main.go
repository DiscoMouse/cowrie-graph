package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Config struct to hold our configuration
type Config struct {
	DatabaseDSN string `json:"database_dsn"`
}

// TopPassword struct will hold the result of our top passwords query
type TopPassword struct {
	Password string `json:"password"`
	Count    int    `json:"count"`
}

// --- NEW ---
// TopUsername struct will hold the result of our top usernames query
type TopUsername struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}

// DailyAttackStat struct will hold the result of our attacks-by-day query
type DailyAttackStat struct {
	Date      string `json:"date"`
	Successes int    `json:"successes"`
	Failures  int    `json:"failures"`
}

// getTopPasswords is our handler function for the top passwords endpoint.
func getTopPasswords(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT password, COUNT(*) as count FROM auth GROUP BY password ORDER BY count DESC LIMIT 10;`
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

// --- NEW ---
// getTopUsernames is our new handler for the top usernames endpoint.
func getTopUsernames(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT username, COUNT(*) as count FROM auth GROUP BY username ORDER BY count DESC LIMIT 10;`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		usernames := []TopUsername{}
		for rows.Next() {
			var u TopUsername
			if err := rows.Scan(&u.Username, &u.Count); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			usernames = append(usernames, u)
		}
		c.JSON(http.StatusOK, usernames)
	}
}

// getAttacksByDay is our handler function for the attacks-by-day endpoint.
func getAttacksByDay(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT
				DATE(timestamp) AS attack_date,
				SUM(success) AS successful_logins,
				COUNT(*) - SUM(success) AS failed_logins
			FROM
				auth
			GROUP BY
				attack_date
			ORDER BY
				attack_date ASC;
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		stats := []DailyAttackStat{}
		for rows.Next() {
			var s DailyAttackStat
			var t time.Time
			if err := rows.Scan(&t, &s.Successes, &s.Failures); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			s.Date = t.Format("2006-01-02")
			stats = append(stats, s)
		}
		c.JSON(http.StatusOK, stats)
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
	router.StaticFile("/", "./static/index.html")

	// API Routes
	api := router.Group("/api/v1")
	{
		api.GET("/top-passwords", getTopPasswords(db))
		api.GET("/attacks-by-day", getAttacksByDay(db))
		// --- NEW ---
		api.GET("/top-usernames", getTopUsernames(db))
	}

	log.Println("Starting Gin server on :8080")
	router.Run(":8080")
}

// loadConfig loads the configuration from config.json
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
