package services

import (
	"fmt"
	"strings"

	"github.com/axe-junction/axe-server/vtc-service/internal/config"
)

type EstimateNameRequest struct {
	OriginName      string `json:"originName"`
	DestinationName string `json:"destinationName"`
}

type EstimateService struct {
	geocoding *GeocodingService
	routing   *RoutingService
	pricing   *PricingService
}

func NewEstimateService(cfg *config.Config) *EstimateService {
	return &EstimateService{
		geocoding: NewGeocodingService(cfg),
		routing:   NewRoutingService(cfg),
		pricing:   NewPricingService(cfg),
	}
}

func (e *EstimateService) GetEstimateByNames(originName, destinationName string) (*EstimateResponse, error) {
	if strings.TrimSpace(originName) == "" || strings.TrimSpace(destinationName) == "" {
		return nil, fmt.Errorf("missing origin or destination")
	}

	origin, err := e.geocoding.Geocode(originName)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode origin '%s': %w", originName, err)
	}

	destination, err := e.geocoding.Geocode(destinationName)
	if err != nil {
		return nil, fmt.Errorf("failed to geocode destination '%s': %w", destinationName, err)
	}

	route, err := e.routing.GetRoute(origin, destination)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate route: %w", err)
	}

	estimates := e.pricing.CalculateEstimates(
		route.Routes[0].Distance,
		route.Routes[0].Duration,
	)

	return &estimates, nil
}
