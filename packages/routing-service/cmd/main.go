package main

import (
	"fmt"
	"log"

	"github.com/axe-junction/axe-server/internal/config"
	"github.com/axe-junction/axe-server/internal/database"
	"github.com/axe-junction/axe-server/internal/handlers"
	"github.com/axe-junction/axe-server/internal/middleware"
	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/repo"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
	}

	userRepo := repo.NewUserRepository(db)

	oauthHandler := handlers.NewOAuthHandler(cfg, userRepo)

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	store := cookie.NewStore([]byte(cfg.SESSION_SECRET))
	r.Use(sessions.Sessions("session", store))

	publicAPI := r.Group("/auth")
	{
		publicAPI.GET("/google", oauthHandler.GoogleLogin)
		publicAPI.GET("/google/callback", oauthHandler.GoogleCallback)
		publicAPI.POST("/logout", oauthHandler.Logout)
	}

	protectedAPI := r.Group("/api")
	protectedAPI.Use(middleware.RequireAuth())
	{
		protectedAPI.GET("/profile", oauthHandler.GetProfile)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Irtiqaa Academy API is running",
		})
	})

	log.Printf("Server starting on port %d", cfg.PORT)
	r.Run(":" + fmt.Sprintf("%d", cfg.PORT))
}
