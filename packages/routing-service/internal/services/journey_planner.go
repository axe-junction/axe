package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/axe-junction/axe-server/internal/models"
	"github.com/google/uuid"
)

// JourneyPlannerService provides advanced journey planning capabilities
type JourneyPlannerService struct {
	stationRepo        models.StationRepo
	lineRepo           models.LigneRepo
	osrmRepo           models.OSRMRepo
	maxWalkingDistance float64
	maxTransfers       int
	walkingSpeed       float64 // m/s
	waitingTime        time.Duration
	transferPenalty    time.Duration
	realTimeDataSource RealTimeDataSource
	serviceAlertSource ServiceAlertSource
}

// RealTimeDataSource interface for real-time data integration
type RealTimeDataSource interface {
	GetRealTimeInfo(ctx context.Context, tripID uuid.UUID) (*models.RealTimeInfo, error)
}

// ServiceAlertSource interface for service alerts
type ServiceAlertSource interface {
	GetServiceAlerts(ctx context.Context, routeIDs []uuid.UUID) ([]models.ServiceAlert, error)
}

// NewJourneyPlannerService creates a new journey planner service
func NewJourneyPlannerService(
	stationRepo models.StationRepo,
	lineRepo models.LigneRepo,
	osrmRepo models.OSRMRepo,
) *JourneyPlannerService {
	return &JourneyPlannerService{
		stationRepo:        stationRepo,
		lineRepo:           lineRepo,
		osrmRepo:           osrmRepo,
		maxWalkingDistance: 1000, // 1km
		maxTransfers:       3,
		walkingSpeed:       1.4,             // 1.4 m/s (5 km/h)
		waitingTime:        5 * time.Minute, // average waiting time
		transferPenalty:    3 * time.Minute, // penalty for each transfer
		realTimeDataSource: &MockRealTimeDataSource{},
		serviceAlertSource: &MockServiceAlertSource{},
	}
}

// PlanJourney plans journeys based on the request
func (s *JourneyPlannerService) PlanJourney(ctx context.Context, req models.JourneyRequest) (*models.JourneyResponse, error) {
	startTime := time.Now()

	log.Printf("🗺️ Planning journey from %s to %s", req.From.Name, req.To.Name)

	// Set defaults
	if req.MaxResults == 0 {
		req.MaxResults = 5
	}
	if req.MaxWalkingDistance == 0 {
		req.MaxWalkingDistance = s.maxWalkingDistance
	}
	if req.MaxTransfers == 0 {
		req.MaxTransfers = s.maxTransfers
	}
	if req.DepartureTime == nil {
		now := time.Now()
		req.DepartureTime = &now
	}

	// Find journeys
	journeys, err := s.findJourneys(ctx, req)
	if err != nil {
		return &models.JourneyResponse{
			Request: req,
			Error:   err.Error(),
			Metadata: models.JourneyMetadata{
				PlanningTime: time.Since(startTime),
			},
		}, nil
	}

	// Rank and filter journeys
	rankedJourneys := s.rankJourneys(journeys, req.Preferences)
	if len(rankedJourneys) > req.MaxResults {
		rankedJourneys = rankedJourneys[:req.MaxResults]
	}

	// Add real-time data if requested
	if req.IncludeRealTime {
		s.enrichWithRealTimeData(ctx, rankedJourneys)
	}

	// Get service alerts
	serviceAlerts := s.getServiceAlerts(ctx, rankedJourneys)

	return &models.JourneyResponse{
		Request:  req,
		Journeys: rankedJourneys,
		Metadata: models.JourneyMetadata{
			PlanningTime:    time.Since(startTime),
			DataSources:     []string{"static_gtfs", "osrm"},
			RealTimeUsed:    req.IncludeRealTime,
			TimetableDate:   time.Now(),
			ServiceAlerts:   serviceAlerts,
			RegionsSearched: []string{"default"},
		},
	}, nil
}

// findJourneys finds possible journeys between origin and destination
func (s *JourneyPlannerService) findJourneys(ctx context.Context, req models.JourneyRequest) ([]models.Journey, error) {
	var journeys []models.Journey

	// Find nearby stations for origin and destination
	originStations, err := s.stationRepo.GetNearby(req.From.Latitude, req.From.Longitude)
	if err != nil {
		return nil, fmt.Errorf("failed to find origin stations: %w", err)
	}

	destStations, err := s.stationRepo.GetNearby(req.To.Latitude, req.To.Longitude)
	if err != nil {
		return nil, fmt.Errorf("failed to find destination stations: %w", err)
	}

	log.Printf("🔍 Found %d origin stations and %d destination stations", len(originStations), len(destStations))

	// Generate journeys for each combination
	var wg sync.WaitGroup
	journeyChan := make(chan models.Journey, 100)

	for _, originStation := range originStations {
		for _, destStation := range destStations {
			wg.Add(1)
			go func(from, to models.Station) {
				defer wg.Done()

				if journey, err := s.createJourney(ctx, req, from, to); err == nil {
					journeyChan <- journey
				}
			}(originStation, destStation)
		}
	}

	// Close channel when all goroutines are done
	go func() {
		wg.Wait()
		close(journeyChan)
	}()

	// Collect journeys
	for journey := range journeyChan {
		journeys = append(journeys, journey)
	}

	log.Printf("✅ Generated %d possible journeys", len(journeys))
	return journeys, nil
}

// createJourney creates a complete journey between two stations
func (s *JourneyPlannerService) createJourney(ctx context.Context, req models.JourneyRequest, fromStation, toStation models.Station) (models.Journey, error) {
	journeyID := uuid.New()

	// Find routes between stations
	routes, err := s.lineRepo.GetRouteBetweenStations(ctx, fromStation.ID, toStation.ID)
	if err != nil || len(routes) == 0 {
		return models.Journey{}, fmt.Errorf("no route found between stations")
	}

	// Create journey segments
	segments, err := s.createJourneySegments(ctx, req, fromStation, toStation, routes)
	if err != nil {
		return models.Journey{}, fmt.Errorf("failed to create journey segments: %w", err)
	}

	// Calculate journey metrics
	summary := s.calculateJourneySummary(segments)

	journey := models.Journey{
		ID:              journeyID,
		From:            req.From,
		To:              req.To,
		DepartureTime:   *req.DepartureTime,
		ArrivalTime:     req.DepartureTime.Add(summary.TotalDuration),
		Duration:        summary.TotalDuration,
		TotalDistance:   s.calculateTotalDistance(segments),
		WalkingDistance: s.calculateWalkingDistance(segments),
		PublicDistance:  s.calculatePublicDistance(segments),
		TransfersCount:  summary.TransfersCount,
		Segments:        segments,
		Summary:         summary,
		Score:           0, // Will be calculated during ranking
		CarbonFootprint: s.calculateCarbonFootprint(segments),
		RealTimeData:    false,
	}

	return journey, nil
}

// createJourneySegments creates all segments for a journey
func (s *JourneyPlannerService) createJourneySegments(ctx context.Context, req models.JourneyRequest, fromStation, toStation models.Station, routes models.Lignes) ([]models.JourneySegment, error) {
	var segments []models.JourneySegment
	currentTime := *req.DepartureTime

	// 1. Walking segment from origin to first station
	walkingSegment, err := s.createWalkingSegment(ctx, req.From, s.stationToLocation(fromStation), currentTime)
	if err != nil {
		log.Printf("⚠️ Warning: Could not create walking segment to first station: %v", err)
		walkingSegment = s.createBasicWalkingSegment(req.From, s.stationToLocation(fromStation), currentTime)
	}
	segments = append(segments, walkingSegment)
	currentTime = walkingSegment.ArrivalTime

	// 2. Public transport segments
	for _, route := range routes {
		if len(route.Stations) > 1 {
			for i := 0; i < len(route.Stations)-1; i++ {
				publicSegment := s.createPublicTransportSegment(
					route.Stations[i],
					route.Stations[i+1],
					route,
					currentTime,
				)
				segments = append(segments, publicSegment)
				currentTime = publicSegment.ArrivalTime
			}
		}
	}

	// 3. Walking segment from last station to destination
	walkingSegment2, err := s.createWalkingSegment(ctx, s.stationToLocation(toStation), req.To, currentTime)
	if err != nil {
		log.Printf("⚠️ Warning: Could not create walking segment from last station: %v", err)
		walkingSegment2 = s.createBasicWalkingSegment(s.stationToLocation(toStation), req.To, currentTime)
	}
	segments = append(segments, walkingSegment2)

	return segments, nil
}

// createWalkingSegment creates a walking segment using OSRM
func (s *JourneyPlannerService) createWalkingSegment(ctx context.Context, from, to models.Location, departureTime time.Time) (models.JourneySegment, error) {
	osrmData, err := s.osrmRepo.GetRouteBetween(ctx, from.Longitude, from.Latitude, to.Longitude, to.Latitude, "walking")
	if err != nil {
		return models.JourneySegment{}, err
	}

	var distance, duration float64
	var geometry [][]float64
	var instructions []models.Instruction

	if len(osrmData.Routes) > 0 {
		route := osrmData.Routes[0]
		distance = route.Distance
		duration = route.Duration
		geometry = route.Geometry.Coordinates

		// Convert OSRM steps to instructions (simplified for now)
		if len(route.Legs) > 0 {
			leg := route.Legs[0]
			instructions = append(instructions, models.Instruction{
				Text:       fmt.Sprintf("Walk %s", leg.Summary),
				Distance:   leg.Distance,
				Duration:   time.Duration(leg.Duration) * time.Second,
				Coordinate: geometry[0], // First coordinate
				Type:       "walk",
			})
		}
	}

	return models.JourneySegment{
		ID:            uuid.New(),
		From:          from,
		To:            to,
		Mode:          models.ModeWalking,
		DepartureTime: departureTime,
		ArrivalTime:   departureTime.Add(time.Duration(duration) * time.Second),
		Duration:      time.Duration(duration) * time.Second,
		Distance:      distance,
		Geometry:      geometry,
		Instructions:  instructions,
		Accessibility: models.AccessibilityInfo{
			WheelchairAccessible: true, // Walking is generally accessible
		},
		CarbonFootprint: 0, // Walking has no carbon footprint
	}, nil
}

// createBasicWalkingSegment creates a basic walking segment without OSRM
func (s *JourneyPlannerService) createBasicWalkingSegment(from, to models.Location, departureTime time.Time) models.JourneySegment {
	distance := s.calculateDistanceBetweenLocations(from, to)
	duration := time.Duration(distance/s.walkingSpeed) * time.Second

	return models.JourneySegment{
		ID:            uuid.New(),
		From:          from,
		To:            to,
		Mode:          models.ModeWalking,
		DepartureTime: departureTime,
		ArrivalTime:   departureTime.Add(duration),
		Duration:      duration,
		Distance:      distance,
		Geometry:      [][]float64{{from.Longitude, from.Latitude}, {to.Longitude, to.Latitude}},
		Instructions: []models.Instruction{{
			Text:     fmt.Sprintf("Walk to %s", to.Name),
			Distance: distance,
			Duration: duration,
			Type:     "walk",
		}},
		Accessibility: models.AccessibilityInfo{
			WheelchairAccessible: true,
		},
		CarbonFootprint: 0,
	}
}

// createPublicTransportSegment creates a public transport segment
func (s *JourneyPlannerService) createPublicTransportSegment(from, to models.Station, route models.Ligne, departureTime time.Time) models.JourneySegment {
	distance := s.calculateDistance(from, to)

	// Estimate travel time (30 km/h average speed)
	averageSpeed := 30.0 / 3.6 // 30 km/h in m/s
	travelTime := time.Duration(distance/averageSpeed) * time.Second

	// Add waiting time for the first segment
	adjustedDepartureTime := departureTime.Add(s.waitingTime)

	return models.JourneySegment{
		ID:            uuid.New(),
		From:          s.stationToLocation(from),
		To:            s.stationToLocation(to),
		Mode:          s.getTransportMode(route.Type),
		DepartureTime: adjustedDepartureTime,
		ArrivalTime:   adjustedDepartureTime.Add(travelTime),
		Duration:      travelTime,
		Distance:      distance,
		Geometry:      [][]float64{{from.Longitude, from.Latitude}, {to.Longitude, to.Latitude}},
		PublicTransport: &models.PublicTransportInfo{
			RouteID:        route.ID,
			RouteName:      route.Name,
			RouteShortName: route.Name,
			StopsCount:     2,
		},
		Accessibility: models.AccessibilityInfo{
			WheelchairAccessible: true,
			ElevatorAvailable:    true,
			AudioAnnouncements:   true,
			VisualAnnouncements:  true,
		},
		CarbonFootprint: s.calculatePublicTransportCarbonFootprint(distance),
	}
}

func (s *JourneyPlannerService) rankJourneys(journeys []models.Journey, preferences models.JourneyPreferences) []models.Journey {
	for i := range journeys {
		journeys[i].Score = s.calculateJourneyScore(journeys[i], preferences)
		journeys[i].Tags = s.generateJourneyTags(journeys[i], journeys)
	}

	sort.Slice(journeys, func(i, j int) bool {
		return journeys[i].Score > journeys[j].Score
	})

	return journeys
}

func (s *JourneyPlannerService) calculateJourneyScore(journey models.Journey, preferences models.JourneyPreferences) float64 {
	score := 1000.0

	timePenalty := journey.Duration.Minutes() * 0.5
	if preferences.PreferFastest {
		timePenalty *= 1.5
	}
	score -= timePenalty

	transferPenalty := float64(journey.TransfersCount) * 75.0
	if preferences.PreferLeastTransfers {
		transferPenalty *= 1.8
	}
	score -= transferPenalty

	walkingPenalty := journey.WalkingDistance / 8.0
	if preferences.PreferLeastWalking {
		walkingPenalty *= 1.6
	}
	score -= walkingPenalty

	// Carbon footprint penalty - reward eco-friendly journeys
	carbonPenalty := journey.CarbonFootprint / 50.0 // 1 point per 50g CO2
	if preferences.PreferGreenest {
		carbonPenalty *= 1.4 // Increase penalty weight
	}
	score -= carbonPenalty

	// Bonus for mixed transport modes (more interesting journeys)
	transportModes := make(map[models.TransportMode]bool)
	for _, segment := range journey.Segments {
		transportModes[segment.Mode] = true
	}
	if len(transportModes) > 2 {
		score += 25.0 // Bonus for variety
	}

	// Bonus for metro/tram usage (faster, more reliable)
	for _, segment := range journey.Segments {
		if segment.Mode == models.ModeMetro {
			score += 40.0
		} else if segment.Mode == models.ModeTram {
			score += 25.0
		}
	}

	// Penalty for excessive walking (>800m)
	if journey.WalkingDistance > 800 {
		score -= (journey.WalkingDistance - 800) / 5.0
	}

	// Bonus for reasonable journey duration (<60 mins)
	if journey.Duration.Minutes() < 60 {
		score += (60 - journey.Duration.Minutes()) * 0.5
	}

	// Apply preference multipliers
	if preferences.PreferFastest {
		score += 2000.0 / math.Max(journey.Duration.Minutes(), 1.0)
	}

	if preferences.PreferLeastWalking {
		score += (1200.0 - journey.WalkingDistance) / 10.0
	}

	if preferences.PreferLeastTransfers {
		score += (6.0 - float64(journey.TransfersCount)) * 35.0
	}

	if preferences.PreferGreenest {
		score += (1500.0 - journey.CarbonFootprint) / 15.0
	}

	return math.Max(0, score)
}

// generateJourneyTags generates descriptive tags for journeys
func (s *JourneyPlannerService) generateJourneyTags(journey models.Journey, allJourneys []models.Journey) []string {
	var tags []string

	if len(allJourneys) == 0 {
		return tags
	}

	// Find min/max values for comparison
	minDuration := journey.Duration
	maxDuration := journey.Duration
	minWalking := journey.WalkingDistance
	maxWalking := journey.WalkingDistance
	minTransfers := journey.TransfersCount
	maxTransfers := journey.TransfersCount
	minCarbon := journey.CarbonFootprint
	maxCarbon := journey.CarbonFootprint

	for _, j := range allJourneys {
		if j.Duration < minDuration {
			minDuration = j.Duration
		}
		if j.Duration > maxDuration {
			maxDuration = j.Duration
		}
		if j.WalkingDistance < minWalking {
			minWalking = j.WalkingDistance
		}
		if j.WalkingDistance > maxWalking {
			maxWalking = j.WalkingDistance
		}
		if j.TransfersCount < minTransfers {
			minTransfers = j.TransfersCount
		}
		if j.TransfersCount > maxTransfers {
			maxTransfers = j.TransfersCount
		}
		if j.CarbonFootprint < minCarbon {
			minCarbon = j.CarbonFootprint
		}
		if j.CarbonFootprint > maxCarbon {
			maxCarbon = j.CarbonFootprint
		}
	}

	// Generate tags based on relative performance
	if journey.Duration == minDuration {
		tags = append(tags, "fastest")
	}
	if journey.WalkingDistance == minWalking {
		tags = append(tags, "least_walking")
	}
	if journey.TransfersCount == minTransfers {
		tags = append(tags, "least_transfers")
	}
	if journey.CarbonFootprint == minCarbon {
		tags = append(tags, "eco_friendly")
	}

	// Additional contextual tags
	if journey.WalkingDistance > 1000 {
		tags = append(tags, "long_walk")
	}
	if journey.Duration.Minutes() > 60 {
		tags = append(tags, "long_journey")
	}
	if journey.TransfersCount == 0 {
		tags = append(tags, "direct")
	}

	// Transport mode tags
	usesMetro := false
	usesTram := false
	usesBus := false

	for _, segment := range journey.Segments {
		switch segment.Mode {
		case models.ModeMetro:
			usesMetro = true
		case models.ModeTram:
			usesTram = true
		case models.ModeBus:
			usesBus = true
		}
	}

	if usesMetro {
		tags = append(tags, "metro")
	}
	if usesTram {
		tags = append(tags, "tram")
	}
	if usesBus {
		tags = append(tags, "bus")
	}

	return tags
}

// Helper methods

func (s *JourneyPlannerService) stationToLocation(station models.Station) models.Location {
	return models.Location{
		ID:        station.ID,
		Name:      station.Name,
		Latitude:  station.Latitude,
		Longitude: station.Longitude,
		Type:      "station",
		StationID: station.ID,
	}
}

func (s *JourneyPlannerService) getTransportMode(routeType string) models.TransportMode {
	switch routeType {
	case "Bus", "bus":
		return models.ModeBus
	case "Metro", "metro":
		return models.ModeMetro
	case "Tram", "tram":
		return models.ModeTram
	case "Train", "train":
		return models.ModeTrain
	default:
		return models.ModePublicTransport
	}
}

func (s *JourneyPlannerService) calculateDistanceBetweenLocations(from, to models.Location) float64 {
	return s.calculateDistanceCoords(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
}

func (s *JourneyPlannerService) calculateDistance(from, to models.Station) float64 {
	return s.calculateDistanceCoords(from.Latitude, from.Longitude, to.Latitude, to.Longitude)
}

func (s *JourneyPlannerService) calculateDistanceCoords(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371000 // meters

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLatRad := (lat2 - lat1) * math.Pi / 180
	deltaLonRad := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLatRad/2)*math.Sin(deltaLatRad/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLonRad/2)*math.Sin(deltaLonRad/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func (s *JourneyPlannerService) calculateJourneySummary(segments []models.JourneySegment) models.JourneySummary {
	var totalDuration, walkingDuration, publicDuration, waitingDuration time.Duration
	var transfersCount int
	var publicRoutes []string
	var lastPublicRoute string

	for _, segment := range segments {
		totalDuration += segment.Duration

		switch segment.Mode {
		case models.ModeWalking:
			walkingDuration += segment.Duration
		default:
			publicDuration += segment.Duration

			if segment.PublicTransport != nil {
				routeName := segment.PublicTransport.RouteName
				if routeName != lastPublicRoute {
					if lastPublicRoute != "" {
						transfersCount++
					}
					publicRoutes = append(publicRoutes, routeName)
					lastPublicRoute = routeName
				}
			}
		}
	}

	// Remove duplicates from public routes
	publicRoutes = s.removeDuplicateStrings(publicRoutes)

	// Estimate waiting time
	waitingDuration = time.Duration(len(publicRoutes)) * s.waitingTime

	return models.JourneySummary{
		TotalDuration:   totalDuration,
		WalkingDuration: walkingDuration,
		PublicDuration:  publicDuration,
		WaitingDuration: waitingDuration,
		TransfersCount:  transfersCount,
		PublicRoutes:    publicRoutes,
		DominantMode:    s.getDominantMode(segments),
		Efficiency:      s.calculateEfficiency(segments),
		Comfort:         s.calculateComfort(segments),
		Environmental:   s.calculateEnvironmentalScore(segments),
	}
}

func (s *JourneyPlannerService) getDominantMode(segments []models.JourneySegment) models.TransportMode {
	modeDurations := make(map[models.TransportMode]time.Duration)

	for _, segment := range segments {
		modeDurations[segment.Mode] += segment.Duration
	}

	var dominantMode models.TransportMode
	var maxDuration time.Duration

	for mode, duration := range modeDurations {
		if duration > maxDuration {
			maxDuration = duration
			dominantMode = mode
		}
	}

	return dominantMode
}

func (s *JourneyPlannerService) calculateEfficiency(segments []models.JourneySegment) float64 {
	if len(segments) == 0 {
		return 0
	}

	// Calculate efficiency based on speed and directness
	totalDistance := s.calculateTotalDistance(segments)
	totalTime := s.calculateTotalDuration(segments)

	if totalTime.Seconds() == 0 {
		return 0
	}

	averageSpeed := totalDistance / totalTime.Seconds() // m/s

	// Normalize to 0-1 scale (assuming 20 m/s is excellent)
	efficiency := math.Min(averageSpeed/20.0, 1.0)

	return efficiency
}

func (s *JourneyPlannerService) calculateComfort(segments []models.JourneySegment) float64 {
	if len(segments) == 0 {
		return 0
	}

	comfort := 1.0

	// Reduce comfort for each transfer
	transfersCount := 0
	var lastMode models.TransportMode
	for _, segment := range segments {
		if segment.Mode != lastMode && lastMode != "" {
			transfersCount++
		}
		lastMode = segment.Mode
	}

	comfort -= float64(transfersCount) * 0.1

	// Reduce comfort for excessive walking
	walkingDistance := s.calculateWalkingDistance(segments)
	if walkingDistance > 500 {
		comfort -= (walkingDistance - 500) / 1000.0
	}

	return math.Max(0, comfort)
}

func (s *JourneyPlannerService) calculateEnvironmentalScore(segments []models.JourneySegment) float64 {
	totalCarbonFootprint := 0.0

	for _, segment := range segments {
		totalCarbonFootprint += segment.CarbonFootprint
	}

	// Lower carbon footprint = higher environmental score
	// Normalize to 0-1 scale (assuming 1000g CO2 is the threshold)
	score := math.Max(0, 1.0-(totalCarbonFootprint/1000.0))

	return score
}

func (s *JourneyPlannerService) calculateTotalDistance(segments []models.JourneySegment) float64 {
	total := 0.0
	for _, segment := range segments {
		total += segment.Distance
	}
	return total
}

func (s *JourneyPlannerService) calculateWalkingDistance(segments []models.JourneySegment) float64 {
	total := 0.0
	for _, segment := range segments {
		if segment.Mode == models.ModeWalking {
			total += segment.Distance
		}
	}
	return total
}

func (s *JourneyPlannerService) calculatePublicDistance(segments []models.JourneySegment) float64 {
	total := 0.0
	for _, segment := range segments {
		if segment.Mode != models.ModeWalking {
			total += segment.Distance
		}
	}
	return total
}

func (s *JourneyPlannerService) calculateTotalDuration(segments []models.JourneySegment) time.Duration {
	total := time.Duration(0)
	for _, segment := range segments {
		total += segment.Duration
	}
	return total
}

func (s *JourneyPlannerService) calculateCarbonFootprint(segments []models.JourneySegment) float64 {
	total := 0.0
	for _, segment := range segments {
		total += segment.CarbonFootprint
	}
	return total
}

func (s *JourneyPlannerService) calculatePublicTransportCarbonFootprint(distance float64) float64 {
	// Approximate carbon footprint for public transport (grams CO2 per meter)
	// Bus: ~0.08g CO2/m, Metro: ~0.03g CO2/m
	return distance * 0.05 // Average between bus and metro
}

func (s *JourneyPlannerService) removeDuplicateStrings(slice []string) []string {
	keys := make(map[string]bool)
	var result []string

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}

func (s *JourneyPlannerService) enrichWithRealTimeData(ctx context.Context, journeys []models.Journey) {
	for i := range journeys {
		for j := range journeys[i].Segments {
			if journeys[i].Segments[j].PublicTransport != nil {
				// Enhanced mock real-time data with realistic scenarios
				routeName := journeys[i].Segments[j].PublicTransport.RouteName

				var delay time.Duration
				var occupancy string

				// Simulate different conditions for different route types
				switch {
				case routeName == "Metro Line 1":
					delay = time.Duration(0) // Metro is usually punctual
					occupancy = "few_seats"
				case routeName == "Bus Line 23":
					delay = time.Duration(7) * time.Minute // Bus has some delay
					occupancy = "standing_room"
				case routeName == "Tram Line T1":
					delay = time.Duration(2) * time.Minute // Tram has minor delay
					occupancy = "few_seats"
				default:
					delay = time.Duration(0)
					occupancy = "few_seats"
				}

				journeys[i].Segments[j].RealTimeInfo = &models.RealTimeInfo{
					IsRealTime:        true,
					Delay:             delay,
					Occupancy:         occupancy,
					ExpectedDeparture: journeys[i].Segments[j].DepartureTime.Add(delay),
					ExpectedArrival:   journeys[i].Segments[j].ArrivalTime.Add(delay),
				}

				// Update segment times with delay
				journeys[i].Segments[j].DepartureTime = journeys[i].Segments[j].DepartureTime.Add(delay)
				journeys[i].Segments[j].ArrivalTime = journeys[i].Segments[j].ArrivalTime.Add(delay)
			}
		}
		journeys[i].RealTimeData = true

		// Recalculate journey summary with real-time data
		journeys[i].Summary = s.calculateJourneySummary(journeys[i].Segments)
	}
}

func (s *JourneyPlannerService) getServiceAlerts(ctx context.Context, journeys []models.Journey) []models.ServiceAlert {
	// Enhanced mock service alerts with realistic scenarios
	alerts := []models.ServiceAlert{
		{
			ID:             uuid.New(),
			Title:          "Metro Line 1 - Normal Service",
			Description:    "All metro services running on schedule",
			Severity:       "info",
			Effect:         "none",
			AffectedRoutes: []string{"Metro Line 1"},
		},
		{
			ID:             uuid.New(),
			Title:          "Bus Line 23 - Minor Delays",
			Description:    "Bus services experiencing 5-10 minute delays due to traffic",
			Severity:       "warning",
			Effect:         "delay",
			AffectedRoutes: []string{"Bus Line 23"},
		},
		{
			ID:             uuid.New(),
			Title:          "Tram Line T1 - Weekend Service",
			Description:    "Enhanced weekend service with additional departures",
			Severity:       "info",
			Effect:         "none",
			AffectedRoutes: []string{"Tram Line T1"},
		},
	}

	// Filter alerts based on routes used in journeys
	var relevantAlerts []models.ServiceAlert
	usedRoutes := make(map[string]bool)

	for _, journey := range journeys {
		for _, segment := range journey.Segments {
			if segment.PublicTransport != nil {
				usedRoutes[segment.PublicTransport.RouteName] = true
			}
		}
	}

	for _, alert := range alerts {
		for _, route := range alert.AffectedRoutes {
			if usedRoutes[route] {
				relevantAlerts = append(relevantAlerts, alert)
				break
			}
		}
	}

	return relevantAlerts
}

// Mock implementations for real-time data and service alerts
type MockRealTimeDataSource struct{}

func (m *MockRealTimeDataSource) GetRealTimeInfo(ctx context.Context, tripID uuid.UUID) (*models.RealTimeInfo, error) {
	return &models.RealTimeInfo{
		IsRealTime: true,
		Delay:      time.Duration(0),
		Occupancy:  "few_seats",
	}, nil
}

type MockServiceAlertSource struct{}

func (m *MockServiceAlertSource) GetServiceAlerts(ctx context.Context, routeIDs []uuid.UUID) ([]models.ServiceAlert, error) {
	return []models.ServiceAlert{}, nil
}
