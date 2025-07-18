package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/axe-junction/axe-server/vtc-service/internal/config"
	"github.com/axe-junction/axe-server/vtc-service/internal/handlers"
	"github.com/axe-junction/axe-server/vtc-service/internal/services"
)

func main() {
	cfg := config.NewConfig()

	estimateService := services.NewEstimateService(cfg)
	router := handlers.SetupRouter(estimateService)

	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	fmt.Printf("🌍 VTC Service running at http://localhost%s\n", cfg.Server.Port)
	fmt.Println("  POST /estimate - Get price estimates by place names")
	fmt.Println("  GET  /health - Health check")

	log.Fatal(server.ListenAndServe())
}
