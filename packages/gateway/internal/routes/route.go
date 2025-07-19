package routes

import (
	"github.com/axe-junction/axe-server/gateway/internal/config"
	"github.com/axe-junction/axe-server/gateway/internal/handlers"
	"github.com/axe-junction/axe-server/gateway/internal/middleware"
	"github.com/axe-junction/axe-server/gateway/internal/repository"
	"github.com/axe-junction/axe-server/gateway/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	cfg *config.Config,
	userRepo *repository.UserRepository,
	grpcClients *services.GRPCClients,
) *gin.Engine {
	// Set Gin mode
	if cfg.Server.Host != "localhost" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Session middleware
	store := cookie.NewStore([]byte(cfg.OAuth.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
	})
	router.Use(sessions.Sessions("gateway_session", store))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(cfg, userRepo)
	apiHandler := handlers.NewAPIHandler(grpcClients)
	authMiddleware := middleware.NewAuthMiddleware(userRepo)

	// Public routes
	public := router.Group("/")
	{
		public.GET("/health", apiHandler.HealthCheck)
		public.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "AXE Server API Gateway",
				"version": "1.0.0",
				"status":  "running",
			})
		})
	}

	// Authentication routes
	auth := router.Group("/auth")
	{
		auth.GET("/google", authHandler.GoogleLogin)
		auth.GET("/google/callback", authHandler.GoogleCallback)
		auth.POST("/logout", authHandler.Logout)
	}

	// Public API routes (no authentication required)
	publicAPI := router.Group("/api/public")
	{
		// Health check
		publicAPI.GET("/health", apiHandler.HealthCheck)

		// Public routing (with optional auth for better results)
		publicAPI.GET("/route", authMiddleware.OptionalAuth(), apiHandler.GetBestRouteByQuery)
		publicAPI.POST("/route", authMiddleware.OptionalAuth(), apiHandler.GetBestRoute)

		// Public VTC prices
		publicAPI.POST("/vtc/prices", authMiddleware.OptionalAuth(), apiHandler.GetVTCPrices)
	}

	// Protected API routes (authentication required)
	protectedAPI := router.Group("/api/v1")
	protectedAPI.Use(authMiddleware.RequireAuth())
	{
		// User profile routes
		user := protectedAPI.Group("/user")
		{
			user.GET("/profile", authHandler.GetProfile)
			user.PUT("/profile", authHandler.UpdateProfile)
		}

		// Enhanced routing services (for authenticated users)
		routing := protectedAPI.Group("/routing")
		{
			routing.POST("/route", apiHandler.GetBestRoute)
			routing.GET("/route", apiHandler.GetBestRouteByQuery)
		}

		// VTC services (for authenticated users)
		vtc := protectedAPI.Group("/vtc")
		{
			vtc.POST("/prices", apiHandler.GetVTCPrices)
		}
	}

	// Admin routes (require admin role)
	admin := router.Group("/api/admin")
	admin.Use(authMiddleware.RequireAuth())
	admin.Use(authMiddleware.RequireRole("admin"))
	{
		admin.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin users endpoint - not implemented yet"})
		})

		admin.GET("/stats", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin stats endpoint - not implemented yet"})
		})
	}

	return router
}
