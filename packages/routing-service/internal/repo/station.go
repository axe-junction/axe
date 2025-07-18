package repo

import (
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

// var _ models.StationRepo = (*StationRepo)(nil)