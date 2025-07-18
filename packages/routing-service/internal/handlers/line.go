package handlers

import (
	"net/http"

	"github.com/axe-junction/axe-server/internal/services"
	"github.com/gin-gonic/gin"
)

func GetAllLignesHandler(ligneService services.LineService) gin.HandlerFunc {
	return func(c *gin.Context) {
		lignes, err := ligneService.GetAllLignes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lignes)
	}
}
