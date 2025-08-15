package handlers

import (
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

// These are our new, clean handlers.
// They just call a method on the store and return the result.

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
