package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/axe-junction/axe-server/vtc-service/internal/services"
)

type EstimateHandler struct {
	estimateService *services.EstimateService
}

func NewEstimateHandler(estimateService *services.EstimateService) *EstimateHandler {
	return &EstimateHandler{
		estimateService: estimateService,
	}
}

func (h *EstimateHandler) EstimateByNames(w http.ResponseWriter, r *http.Request) {
	var req services.EstimateNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	estimates, err := h.estimateService.GetEstimateByNames(req.OriginName, req.DestinationName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(estimates); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
