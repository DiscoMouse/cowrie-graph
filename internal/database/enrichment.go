package database

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
)

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
