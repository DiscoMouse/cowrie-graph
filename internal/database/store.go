package database

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/oschwald/maxminddb-golang"
)

// --- NEW DATA STRUCTS ---
type TopCountry struct {
	CountryCode string `json:"country_code"`
	Count       int    `json:"count"`
}
type TopCity struct {
	City  string `json:"city"`
	Count int    `json:"count"`
}
type TopOrg struct {
	Organization string `json:"organization"`
	Count        int    `json:"count"`
}

// (Other structs remain the same)
type LocationStat struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	City        string  `json:"city"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Count       int     `json:"count"`
}
type TopPassword struct {
	Password string `json:"password"`
	Count    int    `json:"count"`
}
type TopUsername struct {
	Username string `json:"username"`
	Count    int    `json:"count"`
}
type TopIP struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
}
type TopClient struct {
	Version string `json:"version"`
	Count   int    `json:"count"`
}
type TimeSeriesStat struct {
	Date      string `json:"date"`
	Successes int    `json:"successes"`
	Failures  int    `json:"failures"`
}
type GeoData struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
	City struct {
		Names map[string]string `maxminddb:"names"`
	} `maxminddb:"city"`
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
	} `maxminddb:"location"`
}
type IPIntelligence struct {
	IP           string
	CountryCode  sql.NullString
	City         sql.NullString
	Latitude     sql.NullFloat64
	Longitude    sql.NullFloat64
	ASN          sql.NullInt64
	Organization sql.NullString
	IsTor        bool
}

// Store
type Store struct {
	DB    *sql.DB
	GeoDB *maxminddb.Reader
}

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

// --- NEW DATABASE METHODS ---
func (s *Store) GetTopCountries() ([]TopCountry, error) {
	query := `
		SELECT ii.country_code, COUNT(s.id) as count FROM sessions s
		JOIN ip_intelligence ii ON s.ip = ii.ip
		WHERE ii.country_code IS NOT NULL AND ii.country_code != ''
		GROUP BY ii.country_code ORDER BY count DESC LIMIT 10;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopCountry
	for rows.Next() {
		var item TopCountry
		if err := rows.Scan(&item.CountryCode, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) GetTopCities() ([]TopCity, error) {
	query := `
		SELECT ii.city, COUNT(s.id) as count FROM sessions s
		JOIN ip_intelligence ii ON s.ip = ii.ip
		WHERE ii.city IS NOT NULL AND ii.city != ''
		GROUP BY ii.city ORDER BY count DESC LIMIT 10;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopCity
	for rows.Next() {
		var item TopCity
		if err := rows.Scan(&item.City, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *Store) GetTopOrgs() ([]TopOrg, error) {
	query := `
		SELECT ii.organization, COUNT(s.id) as count FROM sessions s
		JOIN ip_intelligence ii ON s.ip = ii.ip
		WHERE ii.organization IS NOT NULL AND ii.organization != ''
		GROUP BY ii.organization ORDER BY count DESC LIMIT 10;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TopOrg
	for rows.Next() {
		var item TopOrg
		if err := rows.Scan(&item.Organization, &item.Count); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// (Other methods remain the same)
func (s *Store) GetIPCounts() (map[string]int, error) {
	query := `SELECT ip, COUNT(*) as count FROM sessions GROUP BY ip;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ipCounts := make(map[string]int)
	for rows.Next() {
		var ip string
		var count int
		if err := rows.Scan(&ip, &count); err != nil {
			return nil, err
		}
		ipCounts[ip] = count
	}
	return ipCounts, nil
}

func (s *Store) GetOrEnrichIP(ipString string) (*IPIntelligence, error) {
	intel := &IPIntelligence{IP: ipString}
	query := "SELECT country_code, city, latitude, longitude, organization, is_tor FROM ip_intelligence WHERE ip = ?"
	err := s.DB.QueryRow(query, ipString).Scan(&intel.CountryCode, &intel.City, &intel.Latitude, &intel.Longitude, &intel.Organization, &intel.IsTor)
	if err == nil {
		return intel, nil
	}
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("error querying ip_intelligence cache: %w", err)
	}

	geoData, err := s.getGeoDataForIP(ipString)
	if err != nil {
		if strings.Contains(err.Error(), "no record found") {
			geoData = nil
		} else {
			return nil, fmt.Errorf("error looking up geoip data: %w", err)
		}
	}

	var countryCode, city, org sql.NullString
	var latitude, longitude sql.NullFloat64
	var asn sql.NullInt64
	if geoData != nil {
		if geoData.Country.ISOCode != "" {
			countryCode = sql.NullString{String: geoData.Country.ISOCode, Valid: true}
		}
		if cityName, ok := geoData.City.Names["en"]; ok && cityName != "" {
			city = sql.NullString{String: cityName, Valid: true}
		}
		if geoData.Location.Latitude != 0 {
			latitude = sql.NullFloat64{Float64: geoData.Location.Latitude, Valid: true}
		}
		if geoData.Location.Longitude != 0 {
			longitude = sql.NullFloat64{Float64: geoData.Location.Longitude, Valid: true}
		}
	}

	insertQuery := `
		INSERT INTO ip_intelligence (ip, country_code, city, latitude, longitude, asn, organization)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE ip=ip;
	`
	_, err = s.DB.Exec(insertQuery, ipString, countryCode, city, latitude, longitude, asn, org)
	if err != nil {
		return nil, fmt.Errorf("error inserting into ip_intelligence cache: %w", err)
	}

	return s.GetOrEnrichIP(ipString)
}

func (s *Store) getGeoDataForIP(ipString string) (*GeoData, error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address format")
	}
	var geoData GeoData
	err := s.GeoDB.Lookup(ip, &geoData)
	if err != nil {
		return nil, err
	}
	return &geoData, nil
}

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

func (s *Store) GetTopPasswords() ([]TopPassword, error) {
	query := `SELECT password, COUNT(*) as count FROM auth GROUP BY password ORDER BY count DESC LIMIT 10;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passwords []TopPassword
	for rows.Next() {
		var p TopPassword
		if err := rows.Scan(&p.Password, &p.Count); err != nil {
			return nil, err
		}
		passwords = append(passwords, p)
	}
	return passwords, nil
}
func (s *Store) GetTopUsernames() ([]TopUsername, error) {
	query := `SELECT username, COUNT(*) as count FROM auth GROUP BY username ORDER BY count DESC LIMIT 10;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usernames []TopUsername
	for rows.Next() {
		var u TopUsername
		if err := rows.Scan(&u.Username, &u.Count); err != nil {
			return nil, err
		}
		usernames = append(usernames, u)
	}
	return usernames, nil
}
func (s *Store) GetTopIPs() ([]TopIP, error) {
	query := `SELECT ip, COUNT(*) as count FROM sessions GROUP BY ip ORDER BY count DESC LIMIT 10;`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ips []TopIP
	for rows.Next() {
		var i TopIP
		if err := rows.Scan(&i.IP, &i.Count); err != nil {
			return nil, err
		}
		ips = append(ips, i)
	}
	return ips, nil
}
func (s *Store) GetTopClients() ([]TopClient, error) {
	query := `
		SELECT c.version, COUNT(s.id) as count FROM sessions s
		JOIN clients c ON s.client = c.id
		GROUP BY c.version ORDER BY count DESC LIMIT 10;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []TopClient
	for rows.Next() {
		var cl TopClient
		if err := rows.Scan(&cl.Version, &cl.Count); err != nil {
			return nil, err
		}
		clients = append(clients, cl)
	}
	return clients, nil
}
func (s *Store) GetAttacksByDay() ([]TimeSeriesStat, error) {
	query := `
		SELECT DATE(timestamp) AS attack_date, SUM(success), COUNT(*) - SUM(success)
		FROM auth GROUP BY attack_date ORDER BY attack_date ASC;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []TimeSeriesStat
	for rows.Next() {
		var stat TimeSeriesStat
		var t time.Time
		if err := rows.Scan(&t, &stat.Successes, &stat.Failures); err != nil {
			return nil, err
		}
		stat.Date = t.Format("2006-01-02")
		stats = append(stats, stat)
	}
	return stats, nil
}
func (s *Store) GetAttacksByMonth() ([]TimeSeriesStat, error) {
	query := `
		SELECT DATE_FORMAT(timestamp, '%Y-%m-01') AS attack_month, SUM(success), COUNT(*) - SUM(success)
		FROM auth GROUP BY attack_month ORDER BY attack_month ASC;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []TimeSeriesStat
	for rows.Next() {
		var stat TimeSeriesStat
		var t time.Time
		if err := rows.Scan(&t, &stat.Successes, &stat.Failures); err != nil {
			return nil, err
		}
		stat.Date = t.Format("2006-01")
		stats = append(stats, stat)
	}
	return stats, nil
}
