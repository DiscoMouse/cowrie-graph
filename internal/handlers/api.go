package handlers

import (
	"log"
	"net/http"

	"github.com/DiscoMouse/cowrie-graph/internal/database"
	"github.com/gin-gonic/gin"
)

// APIHandler holds the database store.
type APIHandler struct {
	Store *database.Store
}

// NewAPIHandler creates a new APIHandler.
func NewAPIHandler(store *database.Store) *APIHandler {
	return &APIHandler{Store: store}
}

// --- NEW HANDLER ---
func (h *APIHandler) GetCountryBarRaceData(c *gin.Context) {
	data, err := h.Store.GetCountryBarRaceData()
	if err != nil {
		log.Printf("ERROR: Failed to get country bar race data: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve country bar race data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

// (Other handlers remain the same)
func (h *APIHandler) GetBarRaceData(c *gin.Context) {
	data, err := h.Store.GetBarRaceData()
	if err != nil {
		log.Printf("ERROR: Failed to get bar race data: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bar race data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopCountries(c *gin.Context) {
	data, err := h.Store.GetTopCountries()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopCities(c *gin.Context) {
	data, err := h.Store.GetTopCities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopOrgs(c *gin.Context) {
	data, err := h.Store.GetTopOrgs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
func (h *APIHandler) GetAttacksByLocation(c *gin.Context) {
	ipCounts, err := h.Store.GetIPCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve IP counts"})
		return
	}
	var locations []database.LocationStat
	for ip, count := range ipCounts {
		intel, err := h.Store.GetOrEnrichIP(ip)
		if err != nil {
			log.Printf("Could not enrich IP %s: %v", ip, err)
			continue
		}
		if intel != nil && intel.Latitude.Valid && intel.Longitude.Valid {
			locations = append(locations, database.LocationStat{
				IP:          ip,
				CountryCode: intel.CountryCode.String,
				City:        intel.City.String,
				Latitude:    intel.Latitude.Float64,
				Longitude:   intel.Longitude.Float64,
				Count:       count,
			})
		}
	}
	c.JSON(http.StatusOK, locations)
}

func (h *APIHandler) GetTopPasswords(c *gin.Context) {
	data, err := h.Store.GetTopPasswords()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopUsernames(c *gin.Context) {
	data, err := h.Store.GetTopUsernames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopIPs(c *gin.Context) {
	data, err := h.Store.GetTopIPs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetTopClients(c *gin.Context) {
	data, err := h.Store.GetTopClients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetAttacksByDay(c *gin.Context) {
	data, err := h.Store.GetAttacksByDay()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetAttacksByMonth(c *gin.Context) {
	data, err := h.Store.GetAttacksByMonth()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
