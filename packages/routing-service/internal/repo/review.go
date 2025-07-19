package repo

import (
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReviewRepo struct {
	db *gorm.DB
}

func NewReviewRepo(db *gorm.DB) *ReviewRepo {
	return &ReviewRepo{db: db}
}

func (r *ReviewRepo) CreateReview(review models.Review) error {
	review.ID = uuid.New()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()
	return r.db.Create(&review).Error
}

func (r *ReviewRepo) GetReviewsForStation(stationID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("station_id = ?", stationID).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepo) GetReviewsForLine(lineID uuid.UUID) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.Where("line_id = ?", lineID).Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepo) GetAverageRating(entityID uuid.UUID, entityType string) (float64, int, error) {
	var result struct {
		AvgRating float64
		Count     int64
	}

	query := r.db.Model(&models.Review{}).
		Select("AVG(rating) as avg_rating, COUNT(*) as count").
		Where("review_type = ?", entityType)

	switch entityType {
	case "station":
		query = query.Where("station_id = ?", entityID)
	case "line":
		query = query.Where("line_id = ?", entityID)
	}

	err := query.Scan(&result).Error
	return result.AvgRating, int(result.Count), err
}

func (r *ReviewRepo) GetSafetyRating(entityID uuid.UUID, entityType string) (float64, int, error) {
	var result struct {
		AvgSafety float64
		Count     int64
	}

	query := r.db.Model(&models.Review{}).
		Select("AVG(safety_rating) as avg_safety, COUNT(*) as count").
		Where("review_type = ?", entityType)

	switch entityType {
	case "station":
		query = query.Where("station_id = ?", entityID)
	case "line":
		query = query.Where("line_id = ?", entityID)
	}

	err := query.Scan(&result).Error
	return result.AvgSafety, int(result.Count), err
}
