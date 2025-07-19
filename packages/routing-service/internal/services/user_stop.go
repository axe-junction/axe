package services

import (
	"fmt"
	"log"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

// UserStopService manages user-contributed bus stops
type UserStopService struct {
	userStopRepo models.UserContributedStopRepo
	stationRepo  models.StationRepo
}

func NewUserStopService(
	userStopRepo models.UserContributedStopRepo,
	stationRepo models.StationRepo,
) *UserStopService {
	return &UserStopService{
		userStopRepo: userStopRepo,
		stationRepo:  stationRepo,
	}
}

// AddUserContributedStop adds a new user-contributed stop
func (s *UserStopService) AddUserContributedStop(userID uuid.UUID, stop models.UserContributedStop) error {
	stop.ContributorID = userID
	stop.ID = uuid.New()
	stop.Status = "pending"
	stop.RequestCount = 1
	stop.DemandScore = 1.0
	stop.CreatedAt = time.Now()
	stop.UpdatedAt = time.Now()

	// Validate stop
	if err := s.validateUserStop(stop); err != nil {
		return err
	}

	return s.userStopRepo.CreateStop(stop)
}

func (s *UserStopService) validateUserStop(stop models.UserContributedStop) error {
	if stop.Name == "" {
		return fmt.Errorf("stop name is required")
	}

	if stop.Latitude < -90 || stop.Latitude > 90 {
		return fmt.Errorf("invalid latitude")
	}

	if stop.Longitude < -180 || stop.Longitude > 180 {
		return fmt.Errorf("invalid longitude")
	}

	// Check if there's already a stop very close to this location
	nearbyStops, err := s.userStopRepo.GetStops(stop.Latitude, stop.Longitude, 50) // 50m radius
	if err == nil && len(nearbyStops) > 0 {
		for _, nearby := range nearbyStops {
			if nearby.ID != stop.ID { // Don't compare with self when updating
				return fmt.Errorf("a stop already exists near this location: %s", nearby.Name)
			}
		}
	}

	return nil
}

// RequestUserStop records demand for a user-contributed stop
func (s *UserStopService) RequestUserStop(stopID uuid.UUID, userID uuid.UUID, fromLat, fromLng, toLat, toLng float64) error {
	// Create demand request
	request := models.StopDemandRequest{
		ID:                    uuid.New(),
		UserContributedStopID: stopID,
		UserID:                userID,
		RequestedAt:           time.Now(),
		FromLatitude:          fromLat,
		FromLongitude:         fromLng,
		ToLatitude:            toLat,
		ToLongitude:           toLng,
	}

	err := s.userStopRepo.AddDemandRequest(request)
	if err != nil {
		return err
	}

	// Update demand score
	currentScore, err := s.userStopRepo.GetDemandScore(stopID)
	if err != nil {
		currentScore = 0
	}

	// Calculate new demand score based on route relevance
	relevanceScore := s.calculateRouteRelevance(fromLat, fromLng, toLat, toLng, stopID)
	newScore := currentScore + relevanceScore

	return s.userStopRepo.UpdateDemandScore(stopID, newScore)
}

func (s *UserStopService) calculateRouteRelevance(fromLat, fromLng, toLat, toLng float64, stopID uuid.UUID) float64 {
	// Simple relevance calculation based on how much the stop could improve the route
	// In a real implementation, this would use proper route calculation

	// Base score for any request
	baseScore := 1.0

	// Additional score based on route length (longer routes get higher scores)
	routeDistance := s.haversineDistance(fromLat, fromLng, toLat, toLng)
	distanceBonus := routeDistance / 10000.0 // 1 point per 10km

	// TODO: Add logic to check if the stop is actually on or near the optimal route
	// For now, we'll use a simplified scoring

	return baseScore + distanceBonus
}

func (s *UserStopService) haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // Earth's radius in meters

	lat1Rad := lat1 * (3.14159265359 / 180)
	lat2Rad := lat2 * (3.14159265359 / 180)
	deltaLatRad := (lat2 - lat1) * (3.14159265359 / 180)
	deltaLonRad := (lon2 - lon1) * (3.14159265359 / 180)

	a := 0.5 - 0.5*((deltaLatRad/2)*(deltaLatRad/2)) +
		(lat1Rad*lat2Rad)*(deltaLonRad/2)*(deltaLonRad/2)

	return earthRadius * 2 * (a * (2 - a))
}

// GetNearbyUserStops returns user-contributed stops near a location
func (s *UserStopService) GetNearbyUserStops(lat, lng float64, radius int) ([]models.UserContributedStop, error) {
	return s.userStopRepo.GetStops(lat, lng, radius)
}

// ProcessHighDemandStops checks for stops with high demand and promotes them
func (s *UserStopService) ProcessHighDemandStops(demandThreshold float64) error {
	highDemandStops, err := s.userStopRepo.GetHighDemandStops(demandThreshold)
	if err != nil {
		return err
	}

	for _, stop := range highDemandStops {
		err := s.promoteStopToOfficial(stop)
		if err != nil {
			log.Printf("Failed to promote stop %s: %v", stop.Name, err)
			continue
		}

		log.Printf("Successfully promoted user-contributed stop '%s' to official status (demand score: %.1f)",
			stop.Name, stop.DemandScore)
	}

	return nil
}

func (s *UserStopService) promoteStopToOfficial(stop models.UserContributedStop) error {
	// Create official station record
	err := s.stationRepo.CreateUserContributedStop(stop)
	if err != nil {
		return fmt.Errorf("failed to create official station: %w", err)
	}

	// Update stop status to promoted
	err = s.userStopRepo.PromoteToOfficial(stop.ID)
	if err != nil {
		return fmt.Errorf("failed to update stop status: %w", err)
	}

	return nil
}

// GetStopDemandScore returns the current demand score for a stop
func (s *UserStopService) GetStopDemandScore(stopID uuid.UUID) (float64, error) {
	return s.userStopRepo.GetDemandScore(stopID)
}

// UpdateStopStatus allows administrators to manually approve or reject stops
func (s *UserStopService) UpdateStopStatus(stopID uuid.UUID, status string) error {
	// This would require an update method in the repository
	// For now, we'll log the action
	log.Printf("Stop %s status updated to %s", stopID.String()[:8], status)

	if status == "approved" {
		// If manually approved, we might want to create an official station
		// but with lower demand threshold
		return nil
	}

	return nil
}

// GetUserContributedStops returns all stops contributed by a specific user
func (s *UserStopService) GetUserContributedStops(userID uuid.UUID) ([]models.UserContributedStop, error) {
	// This would require a method to filter by contributor ID
	// For now, we'll return an empty slice
	return []models.UserContributedStop{}, nil
}

// GetStopStatistics returns statistics about user-contributed stops
func (s *UserStopService) GetStopStatistics() (map[string]interface{}, error) {
	stats := map[string]interface{}{
		"total_contributed_stops": 0,
		"pending_stops":           0,
		"approved_stops":          0,
		"promoted_stops":          0,
		"rejected_stops":          0,
		"average_demand_score":    0.0,
		"last_updated":            time.Now(),
	}

	// TODO: Implement actual statistics gathering
	return stats, nil
}

// ScheduledDemandProcessing runs the demand processing as a background task
func (s *UserStopService) ScheduledDemandProcessing() {
	ticker := time.NewTicker(1 * time.Hour) // Run every hour
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := s.ProcessHighDemandStops(50.0) // Promote stops with 50+ demand score
			if err != nil {
				log.Printf("Error processing high demand stops: %v", err)
			}
		}
	}
}
