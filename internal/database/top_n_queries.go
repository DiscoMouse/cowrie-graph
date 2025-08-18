package database

func (s *Store) GetTopCountries() ([]TopCountry, error) {
	query := `
		SELECT ii.country_code, COUNT(s.id) as count FROM sessions s
		JOIN ip_intelligence ii ON s.ip = ii.ip
		WHERE ii.country_code IS NOT NULL AND ii.country_code != ''
		GROUP BY ii.country_code ORDER BY count DESC LIMIT 20;
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
		GROUP BY ii.city ORDER BY count DESC LIMIT 20;
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
		GROUP BY ii.organization ORDER BY count DESC LIMIT 20;
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

func (s *Store) GetTopPasswords() ([]TopPassword, error) {
	query := `SELECT password, COUNT(*) as count FROM auth GROUP BY password ORDER BY count DESC LIMIT 20;`
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
	query := `SELECT username, COUNT(*) as count FROM auth GROUP BY username ORDER BY count DESC LIMIT 20;`
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
	query := `SELECT ip, COUNT(*) as count FROM sessions GROUP BY ip ORDER BY count DESC LIMIT 20;`
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
		GROUP BY c.version ORDER BY count DESC LIMIT 20;
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
