package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	"github.com/google/uuid"
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


	if err := database.MigrateEnhancedTables(db); err != nil {
		log.Printf("Failed to run enhanced migrations: %v", err)
	}

	var stationCount int64
	db.Model(&models.Station{}).Count(&stationCount)
	if stationCount == 0 {
		if err := database.Seed(db); err != nil {
			log.Printf("Failed to seed database: %v", err)
		}
		if err := database.SeedEnhancedData(db); err != nil {
			log.Printf("Failed to seed enhanced data: %v", err)
		}
	} else {
		log.Printf("Database already has %d stations, skipping seeding", stationCount)
		if err := database.SeedEnhancedData(db); err != nil {
			log.Printf("Failed to seed enhanced data: %v", err)
		}
	}

	userRepo := repo.NewUserRepository(db)
	lineRepo := repo.NewLineRepo(db)
	stationRepo := repo.NewStationRepo(db)
	osrmRepo := repo.NewOSRMRepo("https://osrm.walidbechar.dev")

	// Create enhanced repositories
	reviewRepo := repo.NewReviewRepo(db)
	paymentRepo := repo.NewPaymentRepo(db)
	userStopRepo := repo.NewUserContributedStopRepo(db)
	routeCacheRepo := repo.NewRouteCacheRepo(db)
	paymentTxRepo := repo.NewPaymentTransactionRepo(db)

	// Create routing service with unified line repository
	routingService := services.NewRoutingService(stationRepo, lineRepo, *osrmRepo)

	// Create enhanced services
	reviewService := services.NewReviewService(reviewRepo, stationRepo, lineRepo)
	paymentService := services.NewPaymentService(paymentRepo, paymentTxRepo, routeCacheRepo, routingService)
	userStopService := services.NewUserStopService(userStopRepo, stationRepo)

	log.Println("Waiting for transport graph to initialize...")
	time.Sleep(10 * time.Second)

	oauthHandler := handlers.NewOAuthHandler(cfg, userRepo)
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
		routingAPI.POST("/route", handlers.GetRouteHandler(*routingService))
		routingAPI.POST("/routes", handlers.GetMultipleRoutesHandler(*routingService))
	}

	// Enhanced API routes
	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.RequireAuth()) // All API routes require authentication
	{
		// Payment routes
		payments := apiV1.Group("/payments")
		{
			payments.POST("/process", handlers.ProcessRoutePaymentEnhancedHandler(paymentService))
			payments.GET("/:paymentId/status", handlers.GetPaymentStatusHandler(paymentService))
		}

		// Review routes
		reviews := apiV1.Group("/reviews")
		{
			reviews.POST("/", handlers.AddReviewHandler(reviewService))
			reviews.GET("/stations/:stationId", handlers.GetStationReviewsHandler(reviewService))
		}

		userStops := apiV1.Group("/user-stops")
		{
			userStops.POST("/", func(c *gin.Context) {
				var req handlers.AddUserStopRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				userID, exists := c.Get("user_id")
				if !exists {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
					return
				}

				stop := models.UserContributedStop{
					Name:           req.Name,
					Latitude:       req.Latitude,
					Longitude:      req.Longitude,
					Description:    req.Description,
					NearbyLandmark: req.NearbyLandmark,
				}

				err := userStopService.AddUserContributedStop(userID.(uuid.UUID), stop)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"message": "User-contributed stop added successfully",
					"stop":    stop,
				})
			})
			userStops.GET("/nearby", func(c *gin.Context) {
				lat := c.Query("lat")
				lng := c.Query("lng")
				radius := c.DefaultQuery("radius", "1000")

				if lat == "" || lng == "" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng parameters are required"})
					return
				}

				latVal, err := strconv.ParseFloat(lat, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
					return
				}
				lngVal, err := strconv.ParseFloat(lng, 64)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
					return
				}
				radiusVal, err := strconv.Atoi(radius)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
					return
				}

				stops, err := userStopService.GetNearbyUserStops(latVal, lngVal, radiusVal)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"stops": stops})
			})
			userStops.GET("/manage", handlers.ManageUserStopsHandler(userStopService))
		}

		routing := apiV1.Group("/routing")
		{
			routing.POST("/best", func(c *gin.Context) {
				var req models.RouteOptions
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				fromLat := c.Query("from_lat")
				fromLng := c.Query("from_lng")
				toLat := c.Query("to_lat")
				toLng := c.Query("to_lng")

				fromLatVal, _ := strconv.ParseFloat(fromLat, 64)
				fromLngVal, _ := strconv.ParseFloat(fromLng, 64)
				toLatVal, _ := strconv.ParseFloat(toLat, 64)
				toLngVal, _ := strconv.ParseFloat(toLng, 64)

				route, err := routingService.GetBestRoute(c.Request.Context(), fromLatVal, fromLngVal, toLatVal, toLngVal, req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, route)
			})
			routing.POST("/multiple", func(c *gin.Context) {
				var req models.RouteOptions
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				fromLat := c.Query("from_lat")
				fromLng := c.Query("from_lng")
				toLat := c.Query("to_lat")
				toLng := c.Query("to_lng")

				fromLatVal, _ := strconv.ParseFloat(fromLat, 64)
				fromLngVal, _ := strconv.ParseFloat(fromLng, 64)
				toLatVal, _ := strconv.ParseFloat(toLat, 64)
				toLngVal, _ := strconv.ParseFloat(toLng, 64)

				routes, err := routingService.GetMultipleRoutes(c.Request.Context(), fromLatVal, fromLngVal, toLatVal, toLngVal, req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{"routes": routes})
			})
		}
	}

	publicAPI := r.Group("/api/public")
	{
		publicAPI.GET("/stations", func(c *gin.Context) {
			stations, err := stationService.GetAllStations()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, stations)
		})
		publicAPI.GET("/fare-estimate", handlers.GetFareEstimateHandler(*routingService))
		publicAPI.POST("/route-with-mode", handlers.GetRouteWithModeHandler(*routingService))
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

	authAPI := r.Group("/auth")
	{
		authAPI.GET("/google", oauthHandler.GoogleLogin)
		authAPI.GET("/google/callback", oauthHandler.GoogleCallback)
		authAPI.POST("/logout", oauthHandler.Logout)
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
		lines, err := lineService.GetAllLines()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, lines)
	})

	backgroundService := services.NewBackgroundService(db, userStopService, paymentService)
	backgroundService.Start()

	// Graceful shutdown handling
	defer backgroundService.Stop()

	log.Printf("🚀 Enhanced routing service is running on :%s", cfg.PORT)
	r.Run(":" + fmt.Sprintf("%d", cfg.PORT))
}
