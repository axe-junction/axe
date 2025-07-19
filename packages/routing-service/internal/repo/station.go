package repo

import (
	"math"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StationRepo struct {
	db *gorm.DB
}

func NewStationRepo(db *gorm.DB) *StationRepo {
	return &StationRepo{db: db}
}

func (r *StationRepo) GetAll() ([]models.Station, error) {
	var stations []models.Station
	err := r.db.Find(&stations).Error
	return stations, err
}

func (r *StationRepo) GetByID(id string) ([]models.Station, error) {
	var station models.Station
	err := r.db.Where("id = ?", id).First(&station).Error
	if err != nil {
		return nil, err
	}
	return []models.Station{station}, nil
}

func (r *StationRepo) GetNearby(lat, lng float64, radius ...int) ([]models.Station, error) {
	searchRadius := 1000
	if len(radius) > 0 {
		searchRadius = radius[0]
	}

	const kmPerDegree = 111.0
	latDelta := float64(searchRadius) / (kmPerDegree * 1000)
	lngDelta := float64(searchRadius) / (kmPerDegree * 1000 * math.Cos(lat*math.Pi/180))

	var stations []models.Station
	err := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Find(&stations).Error

	return stations, err
}

// Enhanced methods for safety filtering and user-contributed stops
func (r *StationRepo) GetNearbyWithSafetyFilter(lat, lng float64, minSafetyRating float64, radius ...int) ([]models.Station, error) {
	// Use default radius of 1000m if not specified
	searchRadius := 1000
	if len(radius) > 0 {
		searchRadius = radius[0]
	}

	// Calculate approximate bounding box
	const kmPerDegree = 111.0
	latDelta := float64(searchRadius) / (kmPerDegree * 1000)
	lngDelta := float64(searchRadius) / (kmPerDegree * 1000 * math.Cos(lat*math.Pi/180))

	var stations []models.Station
	err := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Where("safety_rating >= ?", minSafetyRating).
		Find(&stations).Error

	return stations, err
}

func (r *StationRepo) UpdateSafetyRating(stationID uuid.UUID, rating float64) error {
	return r.db.Model(&models.Station{}).
		Where("id = ?", stationID).
		Updates(map[string]interface{}{
			"safety_rating": rating,
			"updated_at":    time.Now(),
		}).Error
}

func (r *StationRepo) CreateUserContributedStop(stop models.UserContributedStop) error {
	// Convert user-contributed stop to station format
	station := models.Station{
		ID:            stop.ID,
		Name:          stop.Name,
		Latitude:      stop.Latitude,
		Longitude:     stop.Longitude,
		Type:          "private_bus",
		SafetyRating:  3.0, // Default rating
		ReviewCount:   0,
		IsVerified:    false,
		ContributorID: &stop.ContributorID,
		DemandScore:   stop.DemandScore,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	return r.db.Create(&station).Error
}

func (r *StationRepo) GetUserContributedStops(lat, lng float64, radius int) ([]models.UserContributedStop, error) {
	const kmPerDegree = 111.0
	latDelta := float64(radius) / (kmPerDegree * 1000)
	lngDelta := float64(radius) / (kmPerDegree * 1000 * math.Cos(lat*math.Pi/180))

	var stations []models.Station
	err := r.db.Where("latitude BETWEEN ? AND ? AND longitude BETWEEN ? AND ?",
		lat-latDelta, lat+latDelta, lng-lngDelta, lng+lngDelta).
		Where("contributor_id IS NOT NULL").
		Find(&stations).Error

	// Convert to UserContributedStop format
	var stops []models.UserContributedStop
	for _, station := range stations {
		stop := models.UserContributedStop{
			ID:            station.ID,
			Name:          station.Name,
			Latitude:      station.Latitude,
			Longitude:     station.Longitude,
			ContributorID: *station.ContributorID,
			DemandScore:   station.DemandScore,
			Status:        "approved",
			CreatedAt:     station.CreatedAt,
			UpdatedAt:     station.UpdatedAt,
		}
		stops = append(stops, stop)
	}

	return stops, err
}

func (r *StationRepo) PromoteUserStopToOfficial(stopID uuid.UUID) error {
	return r.db.Model(&models.Station{}).
		Where("id = ?", stopID).
		Updates(map[string]interface{}{
			"is_verified":    true,
			"contributor_id": nil, // Remove contributor reference when promoted
			"updated_at":     time.Now(),
		}).Error
}
func (r *StationRepo) Update(station *models.Station) error {
	return r.db.Save(station).Error
}
func (r *StationRepo) Create(station *models.Station) error {
	return r.db.Create(station).Error
}
