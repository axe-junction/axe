package repo

import (
	"fmt"

	"github.com/axe-junction/axe-server/internal/models"
	"gorm.io/gorm"
)

type StationRepo struct {
	db *gorm.DB
}

func (stationRepo *StationRepo) GetAll() (models.Stations, error) {
	var stations models.Stations
	if err := stationRepo.db.Find(&stations).Error; err != nil {
		return nil, err
	}
	return stations, nil
}
func (stationRepo *StationRepo) GetByID(id string) (models.Stations, error) {
	var station models.Stations
	if err := stationRepo.db.Find(&station, "id = ?", id).Error; err != nil {
		return models.Stations{}, err
	}
	return station, nil
}

const earthRadius = 6371

func (stationRepo *StationRepo) GetNearby(latitude, longitude float64) (models.Stations, error) {
	var stations models.Stations

	query := fmt.Sprintf(`
		SELECT *, 
		(%[1]d * acos(
			cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
			sin(radians(?)) * sin(radians(latitude))
		)) AS distance
		FROM stations
		ORDER BY distance ASC
		LIMIT 5;
	`, earthRadius)

	if err := stationRepo.db.Raw(query, latitude, longitude, latitude).Scan(&stations).Error; err != nil {
		return nil, err
	}

	return stations, nil
}

// var _ models.StationRepo = (*StationRepo)(nil)

func NewStationRepo(db *gorm.DB) models.StationRepo {
	return &StationRepo{db: db}
}
