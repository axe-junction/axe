package repo

import (
	"math"

	"github.com/axe-junction/axe-server/internal/models"
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
		Find(&stations).Error

	return stations, err
}
