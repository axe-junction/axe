package handlers

import (
	"net/http"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

type GetRouteRequest struct {
	FromLat float64 `json:"from_lat" binding:"required"`
	FromLng float64 `json:"from_lng" binding:"required"`
	ToLat   float64 `json:"to_lat" binding:"required"`
	ToLng   float64 `json:"to_lng" binding:"required"`
}

func GetRouteHandler(routeService services.RoutingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req GetRouteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp, err := routeService.GetBestRoute(
			c.Request.Context(),
			req.FromLat, req.FromLng,
			req.ToLat, req.ToLng,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
