package database

import "database/sql"

// Data Structures for API responses
type BarRaceDataPoint struct {
	Hour  string `json:"hour"`
	IP    string `json:"ip"`
	Count int    `json:"count"`
}
type CountryBarRaceDataPoint struct {
	Hour        string `json:"hour"`
	CountryCode string `json:"country_code"`
	Count       int    `json:"count"`
}
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
