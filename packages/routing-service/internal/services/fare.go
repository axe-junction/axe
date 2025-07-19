package services

import (
	"fmt"
	"math"

	"github.com/axe-junction/axe-server/internal/models"
)

// FareCalculatorService handles fare calculations for different transport types
type FareCalculatorService struct {
	// Fare rates in DZD (Algerian Dinar)
	busFare    float64
	metroFare  float64
	tramFare   float64
	taxiFare   float64
	walkingFee float64
}

// NewFareCalculatorService creates a new fare calculator with standard Algerian rates
func NewFareCalculatorService() *FareCalculatorService {
	return &FareCalculatorService{
		busFare:    30.0,  // 30 DZD base fare for bus
		metroFare:  50.0,  // 50 DZD base fare for metro
		tramFare:   40.0,  // 40 DZD base fare for tram
		taxiFare:   200.0, // 200 DZD base fare for taxi (per km)
		walkingFee: 0.0,   // Free walking
	}
}

// CalculateRouteFare calculates the total fare for a route
func (f *FareCalculatorService) CalculateRouteFare(segments []models.RouteSegment) (*models.RouteFareInfo, error) {
	fareBreakdown := []models.FareInfo{}
	totalFare := 0.0

	for _, segment := range segments {
		fareInfo := f.calculateSegmentFare(segment)
		fareBreakdown = append(fareBreakdown, fareInfo)
		totalFare += fareInfo.TotalFare
	}

	return &models.RouteFareInfo{
		TotalFare:     totalFare,
		Currency:      "DZD",
		FareBreakdown: fareBreakdown,
		EstimatedCost: f.formatCost(totalFare),
	}, nil
}

// calculateSegmentFare calculates fare for a single route segment
func (f *FareCalculatorService) calculateSegmentFare(segment models.RouteSegment) models.FareInfo {
	distance := segment.Distance / 1000.0 // Convert meters to km

	var baseFare, distanceFare, totalFare float64
	var transportType string

	switch segment.Mode {
	case "walking":
		transportType = "walking"
		baseFare = 0.0
		distanceFare = 0.0
		totalFare = 0.0
	case "bus", "public_bus":
		transportType = "bus"
		baseFare = f.busFare
		distanceFare = 0.0 // Flat fare for buses in Algeria
		totalFare = baseFare
	case "metro":
		transportType = "metro"
		baseFare = f.metroFare
		distanceFare = 0.0 // Flat fare for metro
		totalFare = baseFare
	case "tram":
		transportType = "tram"
		baseFare = f.tramFare
		distanceFare = 0.0 // Flat fare for tram
		totalFare = baseFare
	case "taxi", "private_bus":
		transportType = "taxi"
		baseFare = 50.0 // Base flag fare
		distanceFare = f.taxiFare
		totalFare = baseFare + (distance * distanceFare)
	default:
		// Default to bus fare for unknown transport types
		transportType = "bus"
		baseFare = f.busFare
		distanceFare = 0.0
		totalFare = baseFare
	}

	return models.FareInfo{
		TransportType: transportType,
		BaseFare:      baseFare,
		Distance:      distance,
		DistanceFare:  distanceFare,
		TotalFare:     math.Round(totalFare*100) / 100, // Round to 2 decimal places
		Currency:      "DZD",
	}
}

// formatCost formats the cost in a human-readable format
func (f *FareCalculatorService) formatCost(amount float64) string {
	if amount == 0 {
		return "Free"
	}
	return fmt.Sprintf("%.0f DZD", amount)
}

// GetTransportFareEstimate provides a quick fare estimate for a transport type and distance
func (f *FareCalculatorService) GetTransportFareEstimate(transportType string, distance float64) models.FareInfo {
	segment := models.RouteSegment{
		Mode:     transportType,
		Distance: distance * 1000, // Convert km to meters
	}
	return f.calculateSegmentFare(segment)
}

// CalculateMultiModalFare calculates fare for multiple transport options
func (f *FareCalculatorService) CalculateMultiModalFare(routes []*models.RouteResult) error {
	for _, route := range routes {
		fareInfo, err := f.CalculateRouteFare(route.Segments)
		if err != nil {
			return err
		}
		
		// Update route with fare information
		route.TotalFare = fareInfo.TotalFare
		route.Currency = fareInfo.Currency
		route.Summary.TotalFare = fareInfo.TotalFare
	}
	return nil
}
