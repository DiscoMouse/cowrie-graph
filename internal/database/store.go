package database

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// --- Data Structures ---
// These are the same structs we had in main.go
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

// --- Store ---
// Store holds the database connection pool.
type Store struct {
	DB *sql.DB
}

// NewStore creates a new Store.
func NewStore(db *sql.DB) *Store {
	return &Store{DB: db}
}

// --- Database Methods ---
// Each of our queries now becomes a method on the Store.

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
