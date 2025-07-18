package services

import "github.com/axe-junction/axe-server/internal/models"

type StationService struct {
	repo models.StationRepo
}

func NewStationService(repo models.StationRepo) *StationService {
	return &StationService{repo: repo}
}

func (s *StationService) GetAllStations() (models.Stations, error) {
	return s.repo.GetAll()
}

func (s *StationService) GetStationByID(id string) (models.Stations, error) {
	return s.repo.GetByID(id)
}

func (s *StationService) GetNearbyStations(lat, lng float64) (models.Stations, error) {
	return s.repo.GetNearby(lat, lng, 6371000)
}
