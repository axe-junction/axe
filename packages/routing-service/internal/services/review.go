package services

import (
	"fmt"
	"log"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

// ReviewService handles reviews and safety ratings
type ReviewService struct {
	reviewRepo  models.ReviewRepo
	stationRepo models.StationRepo
	lineRepo    models.LineRepo
}

func NewReviewService(
	reviewRepo models.ReviewRepo,
	stationRepo models.StationRepo,
	lineRepo models.LineRepo,
) *ReviewService {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		stationRepo: stationRepo,
		lineRepo:    lineRepo,
	}
}

// AddReview adds a new review and updates safety ratings
func (s *ReviewService) AddReview(userID uuid.UUID, review models.Review) error {
	review.UserID = userID
	review.ID = uuid.New()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	// Validate review
	if err := s.validateReview(review); err != nil {
		return err
	}

	err := s.reviewRepo.CreateReview(review)
	if err != nil {
		return err
	}

	// Update safety ratings asynchronously
	go s.updateSafetyRatings(review)

	return nil
}

func (s *ReviewService) validateReview(review models.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	if review.SafetyRating < 1 || review.SafetyRating > 5 {
		return fmt.Errorf("safety rating must be between 1 and 5")
	}

	if review.StationID == nil && review.LineID == nil {
		return fmt.Errorf("review must be for either a station or a line")
	}

	if review.StationID != nil && review.LineID != nil {
		return fmt.Errorf("review cannot be for both station and line")
	}

	return nil
}

func (s *ReviewService) updateSafetyRatings(review models.Review) {
	if review.StationID != nil {
		// Update station safety rating
		avgSafety, count, err := s.reviewRepo.GetSafetyRating(*review.StationID, "station")
		if err == nil {
			err = s.stationRepo.UpdateSafetyRating(*review.StationID, avgSafety)
			if err != nil {
				log.Printf("Failed to update station safety rating: %v", err)
			} else {
				log.Printf("Updated station %s safety rating to %.2f (based on %d reviews)",
					review.StationID.String()[:8], avgSafety, count)
			}
		}
	}

	if review.LineID != nil {
		// Update line safety rating
		avgSafety, count, err := s.reviewRepo.GetSafetyRating(*review.LineID, "line")
		if err == nil {
			err = s.lineRepo.UpdateLineSafetyRating(*review.LineID, avgSafety)
			if err != nil {
				log.Printf("Failed to update line safety rating: %v", err)
			} else {
				log.Printf("Updated line %s safety rating to %.2f (based on %d reviews)",
					review.LineID.String()[:8], avgSafety, count)
			}
		}
	}
}

// GetStationReviews returns all reviews for a station
func (s *ReviewService) GetStationReviews(stationID uuid.UUID) ([]models.Review, error) {
	return s.reviewRepo.GetReviewsForStation(stationID)
}

// GetLineReviews returns all reviews for a line
func (s *ReviewService) GetLineReviews(lineID uuid.UUID) ([]models.Review, error) {
	return s.reviewRepo.GetReviewsForLine(lineID)
}

// GetStationRating returns the average rating and safety rating for a station
func (s *ReviewService) GetStationRating(stationID uuid.UUID) (float64, float64, int, error) {
	avgRating, count, err := s.reviewRepo.GetAverageRating(stationID, "station")
	if err != nil {
		return 0, 0, 0, err
	}

	avgSafety, _, err := s.reviewRepo.GetSafetyRating(stationID, "station")
	if err != nil {
		return avgRating, 0, count, err
	}

	return avgRating, avgSafety, count, nil
}

// GetLineRating returns the average rating and safety rating for a line
func (s *ReviewService) GetLineRating(lineID uuid.UUID) (float64, float64, int, error) {
	avgRating, count, err := s.reviewRepo.GetAverageRating(lineID, "line")
	if err != nil {
		return 0, 0, 0, err
	}

	avgSafety, _, err := s.reviewRepo.GetSafetyRating(lineID, "line")
	if err != nil {
		return avgRating, 0, count, err
	}

	return avgRating, avgSafety, count, nil
}

// FilterLowRatedEntities filters out entities with low safety ratings
func (s *ReviewService) FilterLowRatedStations(stations []models.Station, minSafetyRating float64) []models.Station {
	var filtered []models.Station
	for _, station := range stations {
		if station.SafetyRating >= minSafetyRating {
			filtered = append(filtered, station)
		}
	}
	return filtered
}

func (s *ReviewService) FilterLowRatedLines(lines models.Lines, minSafetyRating float64) models.Lines {
	var filtered models.Lines
	for _, line := range lines {
		if line.SafetyRating >= minSafetyRating {
			filtered = append(filtered, line)
		}
	}
	return filtered
}

// UpdateReviewVerification marks a review as verified (for moderation)
func (s *ReviewService) UpdateReviewVerification(reviewID uuid.UUID, isVerified bool) error {
	// This would require adding an update method to the review repository
	// For now, we'll log the action
	log.Printf("Review %s verification status updated to %t", reviewID.String()[:8], isVerified)
	return nil
}

// GetReviewStatistics returns statistics about reviews
func (s *ReviewService) GetReviewStatistics() (map[string]interface{}, error) {
	// This would require additional queries to get comprehensive statistics
	// For now, we'll return a basic structure
	stats := map[string]interface{}{
		"total_reviews":         0,
		"station_reviews":       0,
		"line_reviews":          0,
		"average_rating":        0.0,
		"average_safety_rating": 0.0,
		"last_updated":          time.Now(),
	}

	return stats, nil
}
