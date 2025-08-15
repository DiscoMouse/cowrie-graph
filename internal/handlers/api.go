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

// --- NEW: Handler for the world map data ---
func (h *APIHandler) GetAttacksByLocation(c *gin.Context) {
	// 1. Get all unique IPs and their counts
	ipCounts, err := h.Store.GetIPCounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve IP counts"})
		return
	}

	var locations []database.LocationStat

	// 2. For each IP, get its enriched data from our cache
	for ip, count := range ipCounts {
		intel, err := h.Store.GetOrEnrichIP(ip)
		if err != nil {
			// Log the error but continue; don't let one bad IP stop the whole process
			log.Printf("Could not enrich IP %s: %v", ip, err)
			continue
		}

		// 3. If it has coordinates, add it to our result list
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

// (Other handlers remain the same)
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
