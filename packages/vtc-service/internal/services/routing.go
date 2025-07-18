package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/axe-junction/axe-server/vtc-service/internal/config"
)

type OSRMResponse struct {
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
	} `json:"routes"`
	Code string `json:"code"`
}

type RoutingService struct {
	client *http.Client
}

func NewRoutingService(cfg *config.Config) *RoutingService {
	return &RoutingService{
		client: &http.Client{
			Timeout: cfg.Server.RequestTimeout,
		},
	}
}

func (r *RoutingService) GetRoute(origin, destination Coordinates) (*OSRMResponse, error) {
	osrmURL := fmt.Sprintf("https://router.project-osrm.org/route/v1/driving/%.6f,%.6f;%.6f,%.6f?overview=false",
		origin.Lng, origin.Lat, destination.Lng, destination.Lat)

	resp, err := r.client.Get(osrmURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make routing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("routing service returned status %d", resp.StatusCode)
	}

	var osrmResponse OSRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&osrmResponse); err != nil {
		return nil, fmt.Errorf("failed to decode routing response: %w", err)
	}

	if len(osrmResponse.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	return &osrmResponse, nil
}
