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

// TopUsername struct will hold the result of our top usernames query
type TopUsername struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}

// TopIP struct will hold the result of our top IPs query
type TopIP struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
}

// --- NEW ---
// TopClient struct will hold the result of our top clients query
type TopClient struct {
	Version string `json:"version"`
	Count   int    `json:"count"`
}

// DailyAttackStat struct will hold the result of our attacks-by-day query
type DailyAttackStat struct {
	Date      string `json:"date"`
	Successes int    `json:"successes"`
	Failures  int    `json:"failures"`
}

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

func getTopIPs(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `SELECT ip, COUNT(*) as count FROM sessions GROUP BY ip ORDER BY count DESC LIMIT 10;`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		ips := []TopIP{}
		for rows.Next() {
			var i TopIP
			if err := rows.Scan(&i.IP, &i.Count); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			ips = append(ips, i)
		}
		c.JSON(http.StatusOK, ips)
	}
}

// --- NEW ---
// getTopClients is our handler for the top SSH clients endpoint.
func getTopClients(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := `
			SELECT
				c.version,
				COUNT(s.id) as count
			FROM
				sessions s
			JOIN
				clients c ON s.client = c.id
			GROUP BY
				c.version
			ORDER BY
				count DESC
			LIMIT 10;
		`
		rows, err := db.Query(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		clients := []TopClient{}
		for rows.Next() {
			var cl TopClient
			if err := rows.Scan(&cl.Version, &cl.Count); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			clients = append(clients, cl)
		}
		c.JSON(http.StatusOK, clients)
	}
}

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
		api.GET("/top-usernames", getTopUsernames(db))
		api.GET("/top-ips", getTopIPs(db))
		// --- NEW ---
		api.GET("/top-clients", getTopClients(db))
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
