package handlers

import (
	"net/http"
	"strconv"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAllStationsHandler(stationService services.StationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		stations, err := stationService.GetAllStations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stations)
	}
}

func GetNearbyStationsHandler(stationService services.StationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
		lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinates"})
			return
		}

		stations, err := stationService.GetNearbyStations(lat, lng)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stations)
	}
}
