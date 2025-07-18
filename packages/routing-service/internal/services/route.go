package services

import (
	"context"
	"errors"
	"log"

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
	log.Printf("🔍 Searching for route from (%.4f, %.4f) to (%.4f, %.4f)", fromLat, fromLng, toLat, toLng)

	// 1. Find nearby stations
	startStations, err := s.stationRepo.GetNearby(fromLat, fromLng)
	if err != nil || len(startStations) == 0 {
		log.Printf("❌ No start stations found near (%.4f, %.4f)", fromLat, fromLng)
		return nil, errors.New("no start station found")
	}
	log.Printf("📍 Found %d start stations near (%.4f, %.4f)", len(startStations), fromLat, fromLng)
	for i, station := range startStations {
		log.Printf("   %d. %s (ID: %s)", i+1, station.Name, station.ID)
	}

	endStations, err := s.stationRepo.GetNearby(toLat, toLng)
	if err != nil || len(endStations) == 0 {
		log.Printf("❌ No end stations found near (%.4f, %.4f)", toLat, toLng)
		return nil, errors.New("no end station found")
	}
	log.Printf("📍 Found %d end stations near (%.4f, %.4f)", len(endStations), toLat, toLng)
	for i, station := range endStations {
		log.Printf("   %d. %s (ID: %s)", i+1, station.Name, station.ID)
	}

	var bestRoute models.Lignes
	var found bool

	for _, fromStation := range startStations {
		for _, toStation := range endStations {
			log.Printf("🔗 Checking route between %s and %s", fromStation.Name, toStation.Name)
			route, err := s.ligneRepo.GetRouteBetweenStations(ctx, fromStation.ID, toStation.ID)
			if err != nil {
				log.Printf("   ❌ Error finding route: %v", err)
				continue
			}
			if len(route) > 0 {
				log.Printf("   ✅ Found %d routes", len(route))
				bestRoute = route
				found = true
				break
			} else {
				log.Printf("   ⚠️ No routes found")
			}
		}
		if found {
			break
		}
	}

	if !found {
		log.Printf("❌ No route found between any combination of stations")
		return nil, errors.New("no route found between the selected stations")
	}

	log.Printf("✅ Route found successfully")
	result := &models.RouteResult{
		From:  startStations[0],
		To:    endStations[0],
		Route: bestRoute,
	}

	return result, nil
}
