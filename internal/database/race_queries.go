package database

func (s *Store) GetBarRaceData() ([]BarRaceDataPoint, error) {
	query := `
		SELECT
			DATE_FORMAT(starttime, '%Y-%m-%d %H:00') as hour,
			ip,
			COUNT(id) as count
		FROM sessions
		GROUP BY DATE_FORMAT(starttime, '%Y-%m-%d %H:00'), ip
		ORDER BY hour;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []BarRaceDataPoint
	for rows.Next() {
		var point BarRaceDataPoint
		if err := rows.Scan(&point.Hour, &point.IP, &point.Count); err != nil {
			return nil, err
		}
		data = append(data, point)
	}
	return data, nil
}

func (s *Store) GetCountryBarRaceData() ([]CountryBarRaceDataPoint, error) {
	query := `
		SELECT
			DATE_FORMAT(s.starttime, '%Y-%m-%d %H:00') as hour,
			ii.country_code,
			COUNT(s.id) as count
		FROM sessions s
		JOIN ip_intelligence ii ON s.ip = ii.ip
		WHERE ii.country_code IS NOT NULL AND ii.country_code != ''
		GROUP BY hour, ii.country_code
		ORDER BY hour;
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []CountryBarRaceDataPoint
	for rows.Next() {
		var point CountryBarRaceDataPoint
		if err := rows.Scan(&point.Hour, &point.CountryCode, &point.Count); err != nil {
			return nil, err
		}
		data = append(data, point)
	}
	return data, nil
}
