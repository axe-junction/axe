package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

type RoutingHandler struct {
	routingService        *services.RoutingService
	journeyPlannerService *services.JourneyPlannerService
}

func NewRoutingHandler(rs *services.RoutingService, jps *services.JourneyPlannerService) *RoutingHandler {
	return &RoutingHandler{
		routingService:        rs,
		journeyPlannerService: jps,
	}
}
func (h *RoutingHandler) GetBestRoute(c *gin.Context) {
	fromLatStr := c.Query("fromLat")
	fromLngStr := c.Query("fromLng")
	toLatStr := c.Query("toLat")
	toLngStr := c.Query("toLng")

	if fromLatStr != "" && fromLngStr != "" && toLatStr != "" && toLngStr != "" {
		fromLat, err1 := strconv.ParseFloat(fromLatStr, 64)
		fromLng, err2 := strconv.ParseFloat(fromLngStr, 64)
		toLat, err3 := strconv.ParseFloat(toLatStr, 64)
		toLng, err4 := strconv.ParseFloat(toLngStr, 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinates"})
			return
		}

		// Use the journey planner service for better routing
		if h.journeyPlannerService != nil {
			from := models.Location{
				Name:      "Start",
				Latitude:  fromLat,
				Longitude: fromLng,
				Type:      "address",
			}
			to := models.Location{
				Name:      "Destination",
				Latitude:  toLat,
				Longitude: toLng,
				Type:      "address",
			}

			now := time.Now()
			req := models.JourneyRequest{
				From:          from,
				To:            to,
				DepartureTime: &now,
				MaxResults:    3,
				Preferences: models.JourneyPreferences{
					PreferFastest: true,
				},
			}

			response, err := h.journeyPlannerService.PlanJourney(c.Request.Context(), req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, response)
			return
		}

		// Fallback to old routing service
		result, err := h.routingService.GetBestRoute(c.Request.Context(), fromLat, fromLng, toLat, toLng)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
		return
	}

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

	// Use the journey planner service for better routing
	if h.journeyPlannerService != nil {
		from := models.Location{
			Name:      "Start",
			Latitude:  req.FromLat,
			Longitude: req.FromLng,
			Type:      "address",
		}
		to := models.Location{
			Name:      "Destination",
			Latitude:  req.ToLat,
			Longitude: req.ToLng,
			Type:      "address",
		}

		now := time.Now()
		journeyReq := models.JourneyRequest{
			From:          from,
			To:            to,
			DepartureTime: &now,
			MaxResults:    3,
			Preferences: models.JourneyPreferences{
				PreferFastest: true,
			},
		}

		response, err := h.journeyPlannerService.PlanJourney(c.Request.Context(), journeyReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, response)
		return
	}

	fmt.Printf("JSON request: %v\n", req)
	result, err := h.routingService.GetBestRoute(c.Request.Context(), req.FromLat, req.FromLng, req.ToLat, req.ToLng)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
