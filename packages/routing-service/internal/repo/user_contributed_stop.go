package repo

import (
	"math"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserContributedStopRepo struct {
	db *gorm.DB
}

func NewUserContributedStopRepo(db *gorm.DB) *UserContributedStopRepo {
	return &UserContributedStopRepo{db: db}
}

func (r *UserContributedStopRepo) CreateStop(stop models.UserContributedStop) error {
	stop.ID = uuid.New()
	stop.CreatedAt = time.Now()
	stop.UpdatedAt = time.Now()
	stop.Status = "pending"
	stop.RequestCount = 1
	stop.DemandScore = 1.0
	return r.db.Create(&stop).Error
}

func (r *UserContributedStopRepo) GetStops(lat, lng float64, radius int) ([]models.UserContributedStop, error) {
	// Calculate approximate bounding box
	const kmPerDegree = 111.0
	latDelta := float64(radius) / (kmPerDegree * 1000)
	lngDelta := float64(radius) / (kmPerDegree * 1000 * math.Cos(lat*math.Pi/180))

	var stops []models.UserContributedStop
	err := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Where("status IN (?)", []string{"pending", "approved"}).
		Find(&stops).Error

	return stops, err
}

func (r *UserContributedStopRepo) AddDemandRequest(request models.StopDemandRequest) error {
	request.ID = uuid.New()
	request.RequestedAt = time.Now()
	return r.db.Create(&request).Error
}

func (r *UserContributedStopRepo) GetDemandScore(stopID uuid.UUID) (float64, error) {
	var stop models.UserContributedStop
	err := r.db.Where("id = ?", stopID).First(&stop).Error
	if err != nil {
		return 0, err
	}
	return stop.DemandScore, nil
}

func (r *UserContributedStopRepo) UpdateDemandScore(stopID uuid.UUID, score float64) error {
	return r.db.Model(&models.UserContributedStop{}).
		Where("id = ?", stopID).
		Updates(map[string]interface{}{
			"demand_score": score,
			"updated_at":   time.Now(),
		}).Error
}

func (r *UserContributedStopRepo) GetHighDemandStops(threshold float64) ([]models.UserContributedStop, error) {
	var stops []models.UserContributedStop
	err := r.db.Where("demand_score >= ? AND status = ?", threshold, "pending").
		Order("demand_score DESC").
		Find(&stops).Error
	return stops, err
}

func (r *UserContributedStopRepo) PromoteToOfficial(stopID uuid.UUID) error {
	// This would involve:
	// 1. Creating an official station record
	// 2. Updating the user-contributed stop status to "promoted"
	// 3. Potentially notifying users who requested this stop

	return r.db.Model(&models.UserContributedStop{}).
		Where("id = ?", stopID).
		Updates(map[string]interface{}{
			"status":     "promoted",
			"updated_at": time.Now(),
		}).Error
}
