package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oschwald/maxminddb-golang"
)

// Store holds the database connection and the GeoIP reader.
type Store struct {
	DB    *sql.DB
	GeoDB *maxminddb.Reader
}

// NewStore creates a new Store and opens the GeoIP database.
func NewStore(db *sql.DB, geoDBPath string) (*Store, error) {
	geoDB, err := maxminddb.Open(geoDBPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open geoip database: %w", err)
	}
	return &Store{
		DB:    db,
		GeoDB: geoDB,
	}, nil
}

// CreateIntelligenceTable ensures the ip_intelligence table exists.
func (s *Store) CreateIntelligenceTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS ip_intelligence (
		ip VARCHAR(61) NOT NULL,
		country_code CHAR(2),
		city VARCHAR(100),
		latitude DECIMAL(10, 8),
		longitude DECIMAL(11, 8),
		asn INT,
		organization VARCHAR(255),
		is_tor BOOLEAN DEFAULT FALSE,
		last_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		PRIMARY KEY (ip)
	);`
	_, err := s.DB.Exec(query)
	return err
}
