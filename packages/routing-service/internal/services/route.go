package services

import (
	"context"
	"errors"

	"github.com/axe-junction/axe-server/internal/models"
)

type RoutingService struct {
	stationRepo models.StationRepo
	ligneRepo   models.LigneRepo
}

func NewRoutingService(stationRepo models.StationRepo, ligneRepo models.LigneRepo) *RoutingService {
	return &RoutingService{
		stationRepo: stationRepo,
		ligneRepo:   ligneRepo,
	}
}

func (s *RoutingService) GetBestRoute(ctx context.Context, fromLat, fromLng, toLat, toLng float64) (*models.RouteResult, error) {
	// 1. Find nearby stations
	startStations, err := s.stationRepo.GetNearby(fromLat, fromLng)
	if err != nil || len(startStations) == 0 {
		return nil, errors.New("no start station found")
	}

	endStations, err := s.stationRepo.GetNearby(toLat, toLng)
	if err != nil || len(endStations) == 0 {
		return nil, errors.New("no end station found")
	}

	var bestRoute models.Lignes
	var found bool

	for _, fromStation := range startStations {
		for _, toStation := range endStations {
			route, err := s.ligneRepo.GetRouteBetweenStations(ctx, fromStation.ID, toStation.ID)
			if err == nil && len(route) > 0 {
				bestRoute = route
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		return nil, errors.New("no route found between the selected stations")
	}

	result := &models.RouteResult{
		From:  startStations[0],
		To:    endStations[0],
		Route: bestRoute,
	}

	return result, nil
}

