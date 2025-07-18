package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

type JourneyPlannerHandler struct {
	service *services.JourneyPlannerService
}

func NewJourneyPlannerHandler(service *services.JourneyPlannerService) *JourneyPlannerHandler {
	return &JourneyPlannerHandler{service: service}
}

func (h *JourneyPlannerHandler) PlanJourney(c *gin.Context) {
	// Parse query parameters
	fromLatStr := c.Query("fromLat")
	fromLngStr := c.Query("fromLng")
	toLatStr := c.Query("toLat")
	toLngStr := c.Query("toLng")

	// Validate required parameters
	if fromLatStr == "" || fromLngStr == "" || toLatStr == "" || toLngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing required parameters: fromLat, fromLng, toLat, toLng",
		})
		return
	}

	// Parse coordinates
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

	// Optional parameters
	maxResults := 5
	if maxResultsStr := c.Query("maxResults"); maxResultsStr != "" {
		if parsed, err := strconv.Atoi(maxResultsStr); err == nil && parsed > 0 {
			maxResults = parsed
		}
	}

	maxWalkingDistance := 1000.0
	if maxWalkingStr := c.Query("maxWalking"); maxWalkingStr != "" {
		if parsed, err := strconv.ParseFloat(maxWalkingStr, 64); err == nil && parsed > 0 {
			maxWalkingDistance = parsed
		}
	}

	maxTransfers := 3
	if maxTransfersStr := c.Query("maxTransfers"); maxTransfersStr != "" {
		if parsed, err := strconv.Atoi(maxTransfersStr); err == nil && parsed >= 0 {
			maxTransfers = parsed
		}
	}

	// Parse departure time
	departureTime := time.Now()
	if departureTimeStr := c.Query("departureTime"); departureTimeStr != "" {
		if parsed, err := time.Parse("2006-01-02T15:04:05Z", departureTimeStr); err == nil {
			departureTime = parsed
		}
	}

	// Parse preferences
	preferences := models.JourneyPreferences{
		PreferFastest:        c.Query("preferFastest") == "true",
		PreferLeastWalking:   c.Query("preferLeastWalking") == "true",
		PreferLeastTransfers: c.Query("preferLeastTransfers") == "true",
		PreferCheapest:       c.Query("preferCheapest") == "true",
		PreferGreenest:       c.Query("preferGreenest") == "true",
	}

	// Create journey request
	request := models.JourneyRequest{
		From: models.Location{
			Name:      c.Query("fromName"),
			Latitude:  fromLat,
			Longitude: fromLng,
			Type:      "address",
		},
		To: models.Location{
			Name:      c.Query("toName"),
			Latitude:  toLat,
			Longitude: toLng,
			Type:      "address",
		},
		DepartureTime:        &departureTime,
		TimeType:             "departure",
		MaxResults:           maxResults,
		MaxWalkingDistance:   maxWalkingDistance,
		MaxTransfers:         maxTransfers,
		WheelchairAccessible: c.Query("wheelchairAccessible") == "true",
		Preferences:          preferences,
		IncludeRealTime:      c.Query("includeRealTime") == "true",
	}

	// Set default names if not provided
	if request.From.Name == "" {
		request.From.Name = "Origin"
	}
	if request.To.Name == "" {
		request.To.Name = "Destination"
	}

	// Plan journey
	response, err := h.service.PlanJourney(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to plan journey"})
		return
	}

	// Return successful response
	c.JSON(http.StatusOK, response)
}

func (h *JourneyPlannerHandler) GetJourneyModes(c *gin.Context) {
	modes := []models.TransportMode{
		models.ModeWalking,
		models.ModeBicycle,
		models.ModePublicTransport,
		models.ModeTram,
		models.ModeBus,
		models.ModeMetro,
		models.ModeTrain,
	}

	c.JSON(http.StatusOK, gin.H{
		"modes": modes,
	})
}

func (h *JourneyPlannerHandler) GetJourneyPreferences(c *gin.Context) {
	preferences := gin.H{
		"preferences": []string{
			"prefer_fastest",
			"prefer_least_walking",
			"prefer_least_transfers",
			"prefer_cheapest",
			"prefer_greenest",
		},
		"parameters": gin.H{
			"max_results":           "Maximum number of journey alternatives (default: 5)",
			"max_walking_distance":  "Maximum walking distance in meters (default: 1000)",
			"max_transfers":         "Maximum number of transfers (default: 3)",
			"departure_time":        "Departure time in ISO format (default: now)",
			"wheelchair_accessible": "Require wheelchair accessible routes (default: false)",
			"include_real_time":     "Include real-time data (default: false)",
		},
	}

	c.JSON(http.StatusOK, preferences)
}
