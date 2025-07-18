package main

import (
	"fmt"
	"log"
	"net/http"

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
		&models.Ligne{},
		&models.Stop{},
	)
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		panic(err)
	}

	// Enable UUID extension for PostgreSQL
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	// Seed database only if tables are empty
	var count int64
	db.Model(&models.Station{}).Count(&count)
	if count == 0 {
		if err := database.Seed(db); err != nil {
			log.Printf("Failed to seed database: %v", err)
		}
	} else {
		log.Printf("Database already has %d stations, skipping seeding", count)
		// Let's also check stops
		var stopCount int64
		db.Model(&models.Stop{}).Count(&stopCount)
		log.Printf("Database has %d stops", stopCount)

		// If we have stations but no stops, something went wrong - clean and reseed
		if stopCount == 0 {
			log.Println("⚠️ Found stations but no stops! Cleaning and reseeding...")
			db.Exec("TRUNCATE TABLE stations, lignes, stops RESTART IDENTITY CASCADE;")
			if err := database.Seed(db); err != nil {
				log.Printf("Failed to reseed database: %v", err)
			}
		}
	}

	userRepo := repo.NewUserRepository(db)

	oauthHandler := handlers.NewOAuthHandler(cfg, userRepo)
	ligneRepo := repo.NewLigneRepo(db)
	stationRepo := repo.NewStationRepo(db)
	routingService := services.NewRoutingService(stationRepo, ligneRepo)
	routingHandler := handlers.NewRoutingHandler(routingService)

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

	stationService := services.NewStationService(stationRepo)
	r.GET("/api/stations", func(c *gin.Context) {
		stations, err := stationService.GetAllStations()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stations)
	})

	// Debug endpoint to check database state
	r.GET("/debug/db", func(c *gin.Context) {
		var stations []models.Station
		var routes []models.Ligne
		var stops []models.Stop

		db.Find(&stations)
		db.Preload("Stops").Preload("Stations").Find(&routes)
		db.Preload("Route").Preload("Station").Find(&stops)

		c.JSON(http.StatusOK, gin.H{
			"stations": stations,
			"routes":   routes,
			"stops":    stops,
		})
	})

	r.GET("/debug/route-test/:startStation/:endStation", func(c *gin.Context) {
		startStationName := c.Param("startStation")
		endStationName := c.Param("endStation")

		var startStation, endStation models.Station
		if err := db.Where("name = ?", startStationName).First(&startStation).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Start station not found"})
			return
		}
		if err := db.Where("name = ?", endStationName).First(&endStation).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "End station not found"})
			return
		}

		var routes []models.Ligne
		query := `
			SELECT DISTINCT l.* FROM lignes l
			JOIN stops s1 ON l.id = s1.route_id AND s1.station_id = ?
			JOIN stops s2 ON l.id = s2.route_id AND s2.station_id = ?
			WHERE s1.sequence < s2.sequence
		`
		db.Raw(query, startStation.ID, endStation.ID).Scan(&routes)

		c.JSON(http.StatusOK, gin.H{
			"startStation": startStation,
			"endStation":   endStation,
			"routes":       routes,
			"query":        query,
		})
	})

	r.POST("/debug/reseed", func(c *gin.Context) {
		log.Println("🔄 Manual reseed requested...")
		db.Exec("DELETE FROM stops")
		db.Exec("DELETE FROM lignes")
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

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Irtiqaa Academy API is running",
		})
	})

	log.Printf("Server starting on port %d", cfg.PORT)
	r.Run(":" + fmt.Sprintf("%d", cfg.PORT))
}
