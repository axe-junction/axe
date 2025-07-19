package handlers

import (
	"net/http"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAllLinesHandler(lineService services.LineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lines, err := lineService.GetAllLines()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lines)
	}
}
