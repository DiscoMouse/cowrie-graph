package database

import "time"

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
