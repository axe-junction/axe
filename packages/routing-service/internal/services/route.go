package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/axe-junction/axe-server/internal/repo"
	"github.com/google/uuid"
)

func sortStopsBySequence(stops []models.Stop) {
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].Sequence < stops[j].Sequence
	})
}

type RoutingService struct {
	stationRepo models.StationRepo
	lineRepo    models.LineRepo
	osrmRepo    repo.OSRMRepo
	graph       *TransportGraph
	graphMutex  sync.RWMutex
}

func NewRoutingService(stationRepo models.StationRepo, lineRepo models.LineRepo, osrmRepo repo.OSRMRepo) *RoutingService {
	rs := &RoutingService{
		stationRepo: stationRepo,
		lineRepo:    lineRepo,
		osrmRepo:    osrmRepo,
	}

	// Initialize graph on startup
	go rs.BuildTransportGraph(context.Background())

	return rs
}
func (s *RoutingService) estimateWalkingDuration(distance float64) float64 {
	walkingSpeed := 5.0 / 3.6 // 5 km/h in m/s
	return distance / walkingSpeed
}

// BuildTransportGraph precomputes network for efficient routing
func (s *RoutingService) BuildTransportGraph(ctx context.Context) {
	s.graphMutex.Lock()
	defer s.graphMutex.Unlock()

	log.Println("Building transport graph...")
	start := time.Now()

	graph := NewTransportGraph()

	// 1. Add all stations
	stations, _ := s.stationRepo.GetAll()
	for _, station := range stations {
		graph.AddStation(station)
	}

	// 2. Add line connections
	lines, _ := s.lineRepo.GetAll()
	for _, line := range lines {
		stops, _ := s.lineRepo.GetStopsForLine(line.ID)
		sortStopsBySequence(stops)

		for i := 0; i < len(stops)-1; i++ {
			from := stops[i].StationID
			to := stops[i+1].StationID

			// Get actual travel time from OSRM
			duration, distance, err := s.getPublicTransportTime(
				graph.Stations[from].Latitude,
				graph.Stations[from].Longitude,
				graph.Stations[to].Latitude,
				graph.Stations[to].Longitude,
				line.Type,
			)

			if err != nil {
				// Fallback to estimation
				duration = s.estimatePublicTransportDuration(
					graph.Stations[from],
					graph.Stations[to],
					line.Type,
				)
				distance = s.calculateDistance(
					graph.Stations[from],
					graph.Stations[to],
				)
			}

			graph.AddConnection(from, to, line.ID, duration, distance, line.Type)
		}
	}

	// 3. Add transfers between nearby stations
	for _, from := range stations {
		for _, to := range stations {
			if from.ID == to.ID {
				continue
			}

			distance := s.calculateDistance(from, to)
			if distance > 500 { // Only consider transfers < 500m
				continue
			}

			duration := s.estimateWalkingDuration(distance)
			graph.AddTransfer(from.ID, to.ID, duration, distance)
		}
	}

	s.graph = graph
	log.Printf("Transport graph built in %v with %d nodes and %d edges",
		time.Since(start), len(stations), graph.EdgeCount())
}
func (s *RoutingService) estimatePublicTransportDuration(from, to models.Station, transportType string) float64 {
	distance := s.calculateDistance(from, to)
	averageSpeed := 30.0 / 3.6 // 30 km/h in m/s

	// Adjust based on transport type
	switch transportType {
	case "tram":
		averageSpeed = 25.0 / 3.6
	case "train":
		averageSpeed = 50.0 / 3.6
	case "metro":
		averageSpeed = 40.0 / 3.6
	}

	// Add fixed time for stops
	stopsFactor := 1.0 + (distance/5000)*0.1 // 10% extra time per 5km
	return (distance / averageSpeed) * stopsFactor
}

// Distance calculation using Haversine formula
func (s *RoutingService) calculateDistance(from, to models.Station) float64 {
	const R = 6371000 // Earth radius in meters
	φ1 := from.Latitude * math.Pi / 180
	φ2 := to.Latitude * math.Pi / 180
	Δφ := (to.Latitude - from.Latitude) * math.Pi / 180
	Δλ := (to.Longitude - from.Longitude) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

// GetBestRoute using Dijkstra's algorithm with transfers
func (s *RoutingService) GetBestRoute(ctx context.Context, fromLat, fromLng, toLat, toLng float64) (*models.RouteResult, error) {
	s.graphMutex.RLock()
	defer s.graphMutex.RUnlock()

	if s.graph == nil {
		return nil, errors.New("transport graph not initialized")
	}

	// 1. Find nearest stations
	startStations, err := s.stationRepo.GetNearby(fromLat, fromLng, 1000) // 1km radius
	if err != nil {
		panic(err)
	}
	endStations, err := s.stationRepo.GetNearby(toLat, toLng, 1000)
	if err != nil {
		panic(err)
	}

	if len(startStations) == 0 || len(endStations) == 0 {
		return nil, errors.New("no stations found near start or end point")
	}

	// 2. Create temporary start/end nodes
	startNode := &Node{
		ID:        uuid.New(),
		StationID: uuid.Nil,
		Latitude:  fromLat,
		Longitude: fromLng,
	}

	endNode := &Node{
		ID:        uuid.New(),
		StationID: uuid.Nil,
		Latitude:  toLat,
		Longitude: toLng,
	}

	// 3. Connect to graph
	for _, station := range startStations {
		duration, distance, err := s.getWalkingTime(
			fromLng, fromLat,
			station.Longitude, station.Latitude,
		)
		if err != nil {
			continue
		}
		s.graph.AddConnection(startNode.ID, station.ID, uuid.Nil, duration, distance, "walking")
	}

	for _, station := range endStations {
		duration, distance, err := s.getWalkingTime(
			station.Longitude, station.Latitude,
			toLng, toLat,
		)
		if err != nil {
			continue
		}
		s.graph.AddConnection(station.ID, endNode.ID, uuid.Nil, duration, distance, "walking")
	}

	// 4. Find shortest path
	path := s.graph.ShortestPath(startNode.ID, endNode.ID)
	if path == nil {
		return nil, errors.New("no route found")
	}

	// 5. Convert to route segments
	return s.convertPathToRoute(path), nil
}

// Convert path to route result
func (s *RoutingService) convertPathToRoute(path []Edge) *models.RouteResult {
	segments := make([]models.RouteSegment, 0, len(path))
	totalDistance := 0.0
	totalDuration := 0.0

	for _, edge := range path {
		fromStation := s.graph.Stations[edge.From]
		toStation := s.graph.Stations[edge.To]

		segment := models.RouteSegment{
			From:     fromStation,
			To:       toStation,
			Distance: edge.Distance,
			Duration: edge.Duration,
			Mode:     edge.Mode,
		}

		if edge.LineID != uuid.Nil {
			// Fetch line details
			line, _ := s.lineRepo.GetByID(edge.LineID)
			segment.PublicRoute = &line
		}

		segments = append(segments, segment)
		totalDistance += edge.Distance
		totalDuration += edge.Duration
	}

	// Format duration for display
	totalTime := formatDuration(totalDuration)

	return &models.RouteResult{
		TotalDistance: totalDistance,
		TotalDuration: totalDuration,
		TotalTime:     totalTime,
		Segments:      segments,
	}
}

// Format duration as human-readable string
func formatDuration(seconds float64) string {
	duration := time.Duration(seconds * float64(time.Second))
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// OSRM Integration Helpers
func (s *RoutingService) getWalkingTime(fromLng, fromLat, toLng, toLat float64) (float64, float64, error) {
	osrmData, err := s.osrmRepo.GetRouteBetween(context.Background(), fromLng, fromLat, toLng, toLat, "walking")
	if err != nil || len(osrmData.Routes) == 0 {
		return 0, 0, err
	}
	return osrmData.Routes[0].Duration, osrmData.Routes[0].Distance, nil
}

func (s *RoutingService) getPublicTransportTime(fromLat, fromLng, toLat, toLng float64, transportType string) (float64, float64, error) {
	profile := "driving"
	switch transportType {
	case "tram", "train":
		profile = "train"
	}

	osrmData, err := s.osrmRepo.GetRouteBetween(context.Background(), fromLng, fromLat, toLng, toLat, profile)
	if err != nil || len(osrmData.Routes) == 0 {
		return 0, 0, err
	}
	return osrmData.Routes[0].Duration, osrmData.Routes[0].Distance, nil
}

// ... (Other helper methods)
