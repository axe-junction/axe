package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/axe-junction/axe-server/internal/config"
	"github.com/axe-junction/axe-server/internal/database"
	"github.com/axe-junction/axe-server/internal/handlers"
	"github.com/axe-junction/axe-server/internal/middleware"
	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/repo"
	"github.com/axe-junction/axe-server/internal/services"
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

	err = db.AutoMigrate(
		&models.User{},
		&models.Station{},
		&models.Line{},
		&models.Stop{},
		&models.Transfer{},
	)
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		panic(err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	var stationCount int64
	db.Model(&models.Station{}).Count(&stationCount)
	if stationCount == 0 {
		if err := database.Seed(db); err != nil {
			log.Printf("Failed to seed database: %v", err)
		}
	} else {
		log.Printf("Database already has %d stations, skipping seeding", stationCount)
	}

	userRepo := repo.NewUserRepository(db)
	lineRepo := repo.NewLineRepo(db)
	stationRepo := repo.NewStationRepo(db)
	osrmRepo := repo.NewOSRMRepo("https://osrm.walidbechar.dev")

	// Create a legacy wrapper for routing service
	legacyLineRepo := repo.NewLegacyLineRepoWrapper(db)
	routingService := services.NewRoutingService(stationRepo, legacyLineRepo, *osrmRepo)
	journeyPlannerService := services.NewJourneyPlannerService(stationRepo, lineRepo, osrmRepo)

	log.Println("Waiting for transport graph to initialize...")
	time.Sleep(10 * time.Second)

	oauthHandler := handlers.NewOAuthHandler(cfg, userRepo)
	routingHandler := handlers.NewRoutingHandler(routingService, journeyPlannerService)
	lineService := services.NewLineService(lineRepo)
	stationService := services.NewStationService(stationRepo)

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/test", func(c *gin.Context) {
		c.File("./static/test-routing.html")
	})

	routingAPI := r.Group("/routing")
	{
		routingAPI.GET("/best", routingHandler.GetBestRoute)
	}

	r.GET("/api/stations", func(c *gin.Context) {
		stations, err := stationService.GetAllStations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stations)
	})

	r.GET("/debug/db", func(c *gin.Context) {
		var stations []models.Station
		var lines []models.Line
		var stops []models.Stop
		var transfers []models.Transfer

		db.Find(&stations)
		db.Preload("Stops").Find(&lines)
		db.Find(&stops)
		db.Find(&transfers)

		c.JSON(http.StatusOK, gin.H{
			"stations":  stations,
			"lines":     lines,
			"stops":     stops,
			"transfers": transfers,
		})
	})

	r.POST("/debug/reseed", func(c *gin.Context) {
		log.Println("🔄 Manual reseed requested...")
		db.Exec("DELETE FROM transfers")
		db.Exec("DELETE FROM stops")
		db.Exec("DELETE FROM lines")
		db.Exec("DELETE FROM stations")

		if err := database.Seed(db); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Database reseeded successfully"})
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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Algiers Transit API is running",
		})
	})

	// Line API
	r.GET("/api/lines", func(c *gin.Context) {
		lines, err := lineService.GetAllLignes()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lines)
	})

	log.Printf("Server starting on port %d", cfg.PORT)
	r.Run(":" + fmt.Sprintf("%d", cfg.PORT))
}
