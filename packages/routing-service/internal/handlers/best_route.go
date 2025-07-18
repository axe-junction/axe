package handlers

import (
	"net/http"
	"strconv"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	service *services.RoutingService
}

func NewRoutingHandler(service *services.RoutingService) *RoutingHandler {
	return &RoutingHandler{service: service}
}

func (h *RoutingHandler) GetBestRoute(c *gin.Context) {
	fromLatStr := c.Query("fromLat")
	fromLngStr := c.Query("fromLng")
	toLatStr := c.Query("toLat")
	toLngStr := c.Query("toLng")

	// Validate that all required parameters are provided
	if fromLatStr == "" || fromLngStr == "" || toLatStr == "" || toLngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters. Please provide fromLat, fromLng, toLat, and toLng",
		})
		return
	}

	fromLat, err := strconv.ParseFloat(fromLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromLat: must be a valid number"})
		return
	}
	fromLng, err := strconv.ParseFloat(fromLngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid fromLng: must be a valid number"})
		return
	}
	toLat, err := strconv.ParseFloat(toLatStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid toLat: must be a valid number"})
		return
	}
	toLng, err := strconv.ParseFloat(toLngStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid toLng: must be a valid number"})
		return
	}

	// Validate coordinate ranges
	if fromLat < -90 || fromLat > 90 || toLat < -90 || toLat > 90 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "latitude must be between -90 and 90"})
		return
	}
	if fromLng < -180 || fromLng > 180 || toLng < -180 || toLng > 180 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "longitude must be between -180 and 180"})
		return
	}

	routeResult, err := h.service.GetBestRoute(c.Request.Context(), fromLat, fromLng, toLat, toLng)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routeResult)
}
