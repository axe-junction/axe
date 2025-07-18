package handlers

import (
	"net/http"

	"github.com/axe-junction/axe-server/vtc-service/internal/services"
)

func SetupRouter(estimateService *services.EstimateService) http.Handler {
	mux := http.NewServeMux()

	estimateHandler := NewEstimateHandler(estimateService)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("/estimate", estimateHandler.EstimateByNames)

	return mux
}
