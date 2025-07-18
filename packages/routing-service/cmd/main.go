package main

import (
	"fmt"
	"log"

	"github.com/axe-junction/axe-server/internal/config"
	"github.com/axe-junction/axe-server/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Get()
	db, err := database.Connect()
	if err != nil {
		panic(err)
	}

	_ = db

	r := gin.Default()

	publicAPI := r.Group("/auth")
	{
	}

	_ = publicAPI

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Irtiqaa Academy API is running",
		})
	})

	log.Printf("Server starting on port %d", cfg.PORT)
	r.Run(":" + fmt.Sprintf("%d", cfg.PORT))
}
