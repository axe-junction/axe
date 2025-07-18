package handlers

import (
	"net/http"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	routingService *services.RoutingService
}

func NewRoutingHandler(rs *services.RoutingService) *RoutingHandler {
	return &RoutingHandler{routingService: rs}
}

func (h *RoutingHandler) GetBestRoute(c *gin.Context) {
	var req struct {
		FromLat float64 `json:"from_lat" binding:"required"`
		FromLng float64 `json:"from_lng" binding:"required"`
		ToLat   float64 `json:"to_lat" binding:"required"`
		ToLng   float64 `json:"to_lng" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.routingService.GetBestRoute(c.Request.Context(),
		req.FromLat, req.FromLng, req.ToLat, req.ToLng)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

