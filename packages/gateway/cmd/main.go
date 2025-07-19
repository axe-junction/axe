package main

import (
	"log"
	"net/http"

	"github.com/axe-junction/axe-server/gateway/internal/config"
	"github.com/axe-junction/axe-server/gateway/internal/database"
	"github.com/axe-junction/axe-server/gateway/internal/repository"
	"github.com/axe-junction/axe-server/gateway/internal/routes"
	"github.com/axe-junction/axe-server/gateway/internal/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Validate required configuration
	if cfg.OAuth.GoogleClientID == "" || cfg.OAuth.GoogleClientSecret == "" {
		log.Fatal("Google OAuth credentials are required. Please set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables.")
	}

	// Connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("✅ Database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize gRPC clients
	grpcClients, err := services.NewGRPCClients(cfg)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to initialize some gRPC clients: %v", err)
	} else {
		log.Println("✅ gRPC clients initialized successfully")
	}

	// Setup routes
	router := routes.SetupRoutes(cfg, userRepo, grpcClients)

	// Server configuration
	server := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: router,
	}

	log.Printf("🚀 AXE Server API Gateway starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("📚 API Documentation:")
	log.Printf("   GET  /                        - Welcome message")
	log.Printf("   GET  /health                  - Health check")
	log.Printf("   GET  /auth/google             - Google OAuth login")
	log.Printf("   GET  /auth/google/callback    - Google OAuth callback")
	log.Printf("   POST /auth/logout             - Logout")
	log.Printf("   GET  /api/public/route        - Public route planning (GET)")
	log.Printf("   POST /api/public/route        - Public route planning (POST)")
	log.Printf("   POST /api/public/vtc/prices   - Public VTC price estimates")
	log.Printf("   GET  /api/v1/user/profile     - Get user profile (authenticated)")
	log.Printf("   PUT  /api/v1/user/profile     - Update user profile (authenticated)")
	log.Printf("   POST /api/v1/routing/route    - Enhanced route planning (authenticated)")
	log.Printf("   POST /api/v1/vtc/prices       - VTC price estimates (authenticated)")
	log.Printf("   GET  /api/admin/*             - Admin endpoints (admin role required)")

	// Graceful shutdown would be implemented here in production
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}
