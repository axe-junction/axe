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
	// New repositories for enhanced functionality
	reviewRepo              models.ReviewRepo
	paymentRepo             models.PaymentRepo
	userContributedStopRepo models.UserContributedStopRepo
	routeCache              models.RouteCache
	paymentTxRepo           models.PaymentTransactionRepo
	// Fare calculator service
	fareCalculator          *FareCalculatorService
}

// For backward compatibility, keep the old constructor
func NewRoutingService(
	stationRepo models.StationRepo,
	lineRepo models.LineRepo,
	osrmRepo repo.OSRMRepo,
) *RoutingService {
	rs := &RoutingService{
		stationRepo:    stationRepo,
		lineRepo:       lineRepo,
		osrmRepo:       osrmRepo,
		fareCalculator: NewFareCalculatorService(),
	}

	// Initialize graph on startup
	go rs.BuildTransportGraph(context.Background())

	return rs
}

// Enhanced constructor with all repositories
func NewEnhancedRoutingService(
	stationRepo models.StationRepo,
	lineRepo models.LineRepo,
	osrmRepo repo.OSRMRepo,
	reviewRepo models.ReviewRepo,
	paymentRepo models.PaymentRepo,
	userContributedStopRepo models.UserContributedStopRepo,
	routeCache models.RouteCache,
	paymentTxRepo models.PaymentTransactionRepo,
) *RoutingService {
	rs := &RoutingService{
		stationRepo:             stationRepo,
		lineRepo:                lineRepo,
		osrmRepo:                osrmRepo,
		reviewRepo:              reviewRepo,
		paymentRepo:             paymentRepo,
		userContributedStopRepo: userContributedStopRepo,
		routeCache:              routeCache,
		paymentTxRepo:           paymentTxRepo,
		fareCalculator:          NewFareCalculatorService(),
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

// RouteMode represents different routing preferences
type RouteMode string

const (
	RouteModeEconomy  RouteMode = "economy"  // Cheapest routes
	RouteModeFastest  RouteMode = "fastest"  // Fastest routes
	RouteModeSafest   RouteMode = "safest"   // Highest safety rating
	RouteModeBalanced RouteMode = "balanced" // Balanced optimization
	RouteModeComfort  RouteMode = "comfort"  // Least walking, fewer transfers
)

// RouteWeight represents weights for different route criteria
type RouteWeight struct {
	Duration     float64 // Weight for travel time
	Fare         float64 // Weight for cost
	Safety       float64 // Weight for safety rating
	Walking      float64 // Weight for walking distance
	Transfers    float64 // Weight for number of transfers
	Comfort      float64 // Weight for comfort (less crowded)
}

// GetWeightsForMode returns routing weights based on selected mode
func GetWeightsForMode(mode RouteMode) RouteWeight {
	switch mode {
	case RouteModeEconomy:
		return RouteWeight{
			Duration:  0.2,
			Fare:      0.6, // Heavy weight on cost
			Safety:    0.1,
			Walking:   0.05,
			Transfers: 0.05,
			Comfort:   0.0,
		}
	case RouteModeFastest:
		return RouteWeight{
			Duration:  0.7, // Heavy weight on speed
			Fare:      0.1,
			Safety:    0.1,
			Walking:   0.05,
			Transfers: 0.05,
			Comfort:   0.0,
		}
	case RouteModeSafest:
		return RouteWeight{
			Duration:  0.1,
			Fare:      0.1,
			Safety:    0.7, // Heavy weight on safety
			Walking:   0.05,
			Transfers: 0.05,
			Comfort:   0.0,
		}
	case RouteModeComfort:
		return RouteWeight{
			Duration:  0.2,
			Fare:      0.2,
			Safety:    0.2,
			Walking:   0.1, // Minimize walking
			Transfers: 0.2, // Minimize transfers
			Comfort:   0.1,
		}
	default: // RouteModeBalanced
		return RouteWeight{
			Duration:  0.3,
			Fare:      0.3,
			Safety:    0.2,
			Walking:   0.1,
			Transfers: 0.1,
			Comfort:   0.0,
		}
	}
}

// Enhanced GetBestRoute with OSRM-based routing and transport modes
func (s *RoutingService) GetBestRoute(ctx context.Context, fromLat, fromLng, toLat, toLng float64, options models.RouteOptions) (*models.RouteResult, error) {
	// Determine the best transport mode based on distance and user preferences
	transportModes := s.determineTransportModes(fromLat, fromLng, toLat, toLng, options)
	
	var bestRoute *models.RouteResult
	bestScore := -1.0
	
	// Try different transport modes and get OSRM routes
	for _, mode := range transportModes {
		route, err := s.getOSRMRoute(ctx, fromLat, fromLng, toLat, toLng, mode, options)
		if err != nil {
			log.Printf("Failed to get %s route: %v", mode, err)
			continue
		}
		
		// Calculate fare for this route
		if s.fareCalculator != nil {
			fareInfo, err := s.fareCalculator.CalculateRouteFare(route.Segments)
			if err == nil {
				route.TotalFare = fareInfo.TotalFare
				route.Currency = fareInfo.Currency
				route.Summary.TotalFare = fareInfo.TotalFare
			}
		}
		
		// Check if route meets user constraints
		if !s.routeMeetsConstraints(route, options) {
			continue
		}
		
		// Score the route based on mode preferences
		routeMode := s.determineRouteMode(options)
		weights := GetWeightsForMode(routeMode)
		score := s.calculateRouteScore(route, weights, options)
		
		if score > bestScore {
			bestScore = score
			bestRoute = route
		}
	}
	
	if bestRoute == nil {
		return nil, errors.New("no suitable route found")
	}
	
	return bestRoute, nil
}

// getOSRMRoute gets a route from OSRM for a specific transport mode
func (s *RoutingService) getOSRMRoute(ctx context.Context, fromLat, fromLng, toLat, toLng float64, transportMode string, options models.RouteOptions) (*models.RouteResult, error) {
	// Map transport mode to OSRM profile
	profile := s.getOSRMProfile(transportMode)
	
	// Get route from OSRM
	osrmResp, err := s.osrmRepo.GetRouteBetween(ctx, fromLng, fromLat, toLng, toLat, profile)
	if err != nil {
		return nil, fmt.Errorf("OSRM error for %s: %w", transportMode, err)
	}
	
	if len(osrmResp.Routes) == 0 {
		return nil, fmt.Errorf("no %s route found", transportMode)
	}
	
	osrmRoute := osrmResp.Routes[0]
	
	// Convert OSRM response to our route format
	route := &models.RouteResult{
		From: models.Station{
			Latitude:  fromLat,
			Longitude: fromLng,
			Name:      "Origin",
		},
		To: models.Station{
			Latitude:  toLat,
			Longitude: toLng,
			Name:      "Destination",  
		},
		TotalDistance: osrmRoute.Distance,
		TotalDuration: osrmRoute.Duration,
		TotalTime:     formatDuration(osrmRoute.Duration),
		Segments:      s.convertOSRMToSegments(osrmRoute, transportMode),
		Summary: models.RouteSummary{
			WalkingDistance: s.getWalkingDistance(osrmRoute, transportMode),
			PublicDistance:  s.getPublicDistance(osrmRoute, transportMode),
			WalkingDuration: s.getWalkingDuration(osrmRoute, transportMode),
			PublicDuration:  s.getPublicDuration(osrmRoute, transportMode),
			TransfersCount:  s.getTransferCount(transportMode),
			PublicRouteNames: s.getRouteNames(transportMode),
		},
		OverallSafetyRating: s.calculateSafetyRating(transportMode),
		PaymentMethods:      s.getPaymentMethods(transportMode),
		PaymentRequired:     s.isPaymentRequired(transportMode),
	}
	
	return route, nil
}

// determineTransportModes determines which transport modes to try based on distance and preferences
func (s *RoutingService) determineTransportModes(fromLat, fromLng, toLat, toLng float64, options models.RouteOptions) []string {
	distance := s.calculateDistanceKm(fromLat, fromLng, toLat, toLng)
	modes := []string{}
	
	// Walking for short distances
	if distance <= 2.0 {
		modes = append(modes, "walking")
	}
	
	// Driving/taxi for any distance
	modes = append(modes, "driving")
	
	// Public transport for medium to long distances
	if distance > 1.0 {
		modes = append(modes, "bus")
		if distance > 3.0 {
			modes = append(modes, "train")
		}
	}
	
	// Cycling for short to medium distances
	if distance <= 10.0 {
		modes = append(modes, "cycling")
	}
	
	// Filter based on user preferences
	if options.Mode == "economy" {
		// Prioritize cheaper options
		return []string{"walking", "bus", "cycling", "driving"}
	} else if options.Mode == "fastest" {
		// Prioritize faster options
		return []string{"driving", "train", "bus", "cycling", "walking"}
	} else if options.Mode == "comfort" {
		// Prioritize comfortable options
		return []string{"driving", "train", "bus"}
	}
	
	return modes
}

// getOSRMProfile maps transport mode to OSRM routing profile
func (s *RoutingService) getOSRMProfile(transportMode string) string {
	switch transportMode {
	case "walking":
		return "foot"
	case "cycling":
		return "bike"
	case "driving", "taxi":
		return "car"
	case "bus", "train":
		return "car" // Use car profile for public transport approximation
	default:
		return "car"
	}
}

// convertOSRMToSegments converts OSRM route to our segment format
func (s *RoutingService) convertOSRMToSegments(osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry struct {
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Legs []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Summary  string  `json:"summary"`
		Steps    []struct {
			Distance    float64 `json:"distance"`
			Duration    float64 `json:"duration"`
			Instruction string  `json:"instruction"`
			Name        string  `json:"name"`
			Maneuver    struct {
				Type     string    `json:"type"`
				Modifier string    `json:"modifier"`
				Location []float64 `json:"location"`
			} `json:"maneuver"`
		} `json:"steps"`
	} `json:"legs"`
}, transportMode string) []models.RouteSegment {
	
	segments := []models.RouteSegment{}
	
	if len(osrmRoute.Legs) == 0 {
		// Single segment route
		segment := models.RouteSegment{
			Mode:           transportMode,
			Distance:       osrmRoute.Distance,
			Duration:       osrmRoute.Duration,
			Geometry:       osrmRoute.Geometry.Coordinates,
			Fare:           s.calculateModeFare(transportMode, osrmRoute.Distance),
			PaymentMethods: s.getPaymentMethods(transportMode),
			SafetyRating:   s.calculateSafetyRating(transportMode),
		}
		segments = append(segments, segment)
	} else {
		// Multi-leg route
		for _, leg := range osrmRoute.Legs {
			segment := models.RouteSegment{
				Mode:           transportMode,
				Distance:       leg.Distance,
				Duration:       leg.Duration,
				Fare:           s.calculateModeFare(transportMode, leg.Distance),
				PaymentMethods: s.getPaymentMethods(transportMode),
				SafetyRating:   s.calculateSafetyRating(transportMode),
			}
			segments = append(segments, segment)
		}
	}
	
	return segments
}

// Helper functions for route analysis
func (s *RoutingService) getWalkingDistance(osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry struct {
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Legs []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Summary  string  `json:"summary"`
		Steps    []struct {
			Distance    float64 `json:"distance"`
			Duration    float64 `json:"duration"`
			Instruction string  `json:"instruction"`
			Name        string  `json:"name"`
			Maneuver    struct {
				Type     string    `json:"type"`
				Modifier string    `json:"modifier"`
				Location []float64 `json:"location"`
			} `json:"maneuver"`
		} `json:"steps"`
	} `json:"legs"`
}, transportMode string) float64 {
	if transportMode == "walking" {
		return osrmRoute.Distance
	}
	return 0.0
}

func (s *RoutingService) getPublicDistance(osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry struct {
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Legs []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Summary  string  `json:"summary"`
		Steps    []struct {
			Distance    float64 `json:"distance"`
			Duration    float64 `json:"duration"`
			Instruction string  `json:"instruction"`
			Name        string  `json:"name"`
			Maneuver    struct {
				Type     string    `json:"type"`
				Modifier string    `json:"modifier"`
				Location []float64 `json:"location"`
			} `json:"maneuver"`
		} `json:"steps"`
	} `json:"legs"`
}, transportMode string) float64 {
	if transportMode == "bus" || transportMode == "train" {
		return osrmRoute.Distance
	}
	return 0.0
}

func (s *RoutingService) getWalkingDuration(osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry struct {
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Legs []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Summary  string  `json:"summary"`
		Steps    []struct {
			Distance    float64 `json:"distance"`
			Duration    float64 `json:"duration"`
			Instruction string  `json:"instruction"`
			Name        string  `json:"name"`
			Maneuver    struct {
				Type     string    `json:"type"`
				Modifier string    `json:"modifier"`
				Location []float64 `json:"location"`
			} `json:"maneuver"`
		} `json:"steps"`
	} `json:"legs"`
}, transportMode string) float64 {
	if transportMode == "walking" {
		return osrmRoute.Duration
	}
	return 0.0
}

func (s *RoutingService) getPublicDuration(osrmRoute struct {
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
	Geometry struct {
		Coordinates [][]float64 `json:"coordinates"`
	} `json:"geometry"`
	Legs []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Summary  string  `json:"summary"`
		Steps    []struct {
			Distance    float64 `json:"distance"`
			Duration    float64 `json:"duration"`
			Instruction string  `json:"instruction"`
			Name        string  `json:"name"`
			Maneuver    struct {
				Type     string    `json:"type"`
				Modifier string    `json:"modifier"`
				Location []float64 `json:"location"`
			} `json:"maneuver"`
		} `json:"steps"`
	} `json:"legs"`
}, transportMode string) float64 {
	if transportMode == "bus" || transportMode == "train" {
		return osrmRoute.Duration
	}
	return 0.0
}

func (s *RoutingService) getTransferCount(transportMode string) int {
	// For now, assume single mode = no transfers
	return 0
}

func (s *RoutingService) getRouteNames(transportMode string) []string {
	switch transportMode {
	case "bus":
		return []string{"Bus Route"}
	case "train":
		return []string{"Train Line"}
	case "walking":
		return []string{"Walking"}
	case "driving":
		return []string{"Driving"}
	case "cycling":
		return []string{"Cycling"}
	default:
		return []string{transportMode}
	}
}

func (s *RoutingService) calculateSafetyRating(transportMode string) float64 {
	switch transportMode {
	case "walking":
		return 3.5
	case "cycling":
		return 3.0
	case "bus":
		return 4.0
	case "train":
		return 4.5
	case "driving":
		return 3.8
	default:
		return 3.0
	}
}

func (s *RoutingService) getPaymentMethods(transportMode string) []string {
	switch transportMode {
	case "bus":
		return []string{"cash", "card", "mobile"}
	case "train":
		return []string{"card", "mobile", "ticket"}
	case "driving":
		return []string{"fuel", "toll"}
	case "walking", "cycling":
		return []string{}
	default:
		return []string{"cash"}
	}
}

func (s *RoutingService) isPaymentRequired(transportMode string) bool {
	return transportMode != "walking" && transportMode != "cycling"
}

func (s *RoutingService) calculateModeFare(transportMode string, distance float64) float64 {
	if s.fareCalculator != nil {
		fareInfo := s.fareCalculator.GetTransportFareEstimate(transportMode, distance/1000.0)
		return fareInfo.TotalFare
	}
	
	// Fallback fare calculation
	switch transportMode {
	case "bus":
		return 30.0
	case "train":
		return 50.0
	case "driving":
		return distance / 1000.0 * 20.0 // 20 DZD per km for fuel
	case "walking", "cycling":
		return 0.0
	default:
		return 30.0
	}
}

func (s *RoutingService) calculateDistanceKm(fromLat, fromLng, toLat, toLng float64) float64 {
	const R = 6371.0 // Earth radius in km
	φ1 := fromLat * math.Pi / 180
	φ2 := toLat * math.Pi / 180
	Δφ := (toLat - fromLat) * math.Pi / 180
	Δλ := (toLng - fromLng) * math.Pi / 180

	a := math.Sin(Δφ/2)*math.Sin(Δφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*
			math.Sin(Δλ/2)*math.Sin(Δλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *RoutingService) routeMeetsConstraints(route *models.RouteResult, options models.RouteOptions) bool {
	// Check fare constraint
	if options.MaxFare > 0 && route.TotalFare > options.MaxFare {
		return false
	}
	
	// Check walking time constraint
	if options.MaxWalkingTime > 0 && route.Summary.WalkingDuration > options.MaxWalkingTime {
		return false
	}
	
	// Check safety constraint
	if options.SafetyFilter && route.OverallSafetyRating < options.MinSafetyRating {
		return false
	}
	
	return true
}

// GetMultipleRoutes returns multiple route options
func (s *RoutingService) GetMultipleRoutes(ctx context.Context, fromLat, fromLng, toLat, toLng float64, options models.RouteOptions) ([]*models.RouteResult, error) {
	var routes []*models.RouteResult

	// Get the best route
	bestRoute, err := s.GetBestRoute(ctx, fromLat, fromLng, toLat, toLng, options)
	if err != nil {
		return nil, err
	}
	routes = append(routes, bestRoute)

	// Get alternative routes by modifying options
	alternativeOptions := options
	alternativeOptions.SafetyFilter = false // More permissive for alternatives
	if alternativeOptions.MaxWalkingTime > 0 {
		alternativeOptions.MaxWalkingTime *= 1.5 // Allow longer walking
	}

	altRoute, err := s.GetBestRoute(ctx, fromLat, fromLng, toLat, toLng, alternativeOptions)
	if err == nil && !routesEqual(bestRoute, altRoute) {
		routes = append(routes, altRoute)
	}

	// Try with private stops included if not already
	if !options.IncludePrivate {
		privateOptions := options
		privateOptions.IncludePrivate = true
		privateRoute, err := s.GetBestRoute(ctx, fromLat, fromLng, toLat, toLng, privateOptions)
		if err == nil && !routeContainsRoute(routes, privateRoute) {
			routes = append(routes, privateRoute)
		}
	}

	return routes, nil
}

// CalculateRouteFare calculates the total fare for a route
func (s *RoutingService) CalculateRouteFare(segments []models.RouteSegment) (float64, error) {
	return s.paymentRepo.CalculateRouteFare(segments)
}

// ProcessRoutePayment processes payment for a route
func (s *RoutingService) ProcessRoutePayment(userID uuid.UUID, routeResult *models.RouteResult, paymentMethod string) (*models.Payment, error) {
	payment := models.Payment{
		ID:            uuid.New(),
		UserID:        userID,
		RouteID:       uuid.New(), // This should be generated when route is calculated
		Amount:        routeResult.TotalFare,
		Currency:      routeResult.Currency,
		PaymentMethod: paymentMethod,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := s.paymentRepo.CreatePayment(payment)
	if err != nil {
		return nil, err
	}

	// Here you would integrate with actual payment processing
	// For now, we'll mark as completed
	err = s.paymentRepo.UpdatePaymentStatus(payment.ID, "completed")
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// AddUserContributedStop adds a user-contributed bus stop
func (s *RoutingService) AddUserContributedStop(stop models.UserContributedStop) error {
	stop.CreatedAt = time.Now()
	stop.UpdatedAt = time.Now()
	return s.userContributedStopRepo.CreateStop(stop)
}

// RequestUserStop records demand for a user-contributed stop
func (s *RoutingService) RequestUserStop(stopID uuid.UUID, userID uuid.UUID, fromLat, fromLng, toLat, toLng float64) error {
	request := models.StopDemandRequest{
		ID:                    uuid.New(),
		UserContributedStopID: stopID,
		UserID:                userID,
		RequestedAt:           time.Now(),
		FromLatitude:          fromLat,
		FromLongitude:         fromLng,
		ToLatitude:            toLat,
		ToLongitude:           toLng,
	}

	err := s.userContributedStopRepo.AddDemandRequest(request)
	if err != nil {
		return err
	}

	// Update demand score
	currentScore, err := s.userContributedStopRepo.GetDemandScore(stopID)
	if err != nil {
		currentScore = 0
	}

	newScore := currentScore + 1.0 // Simple increment, could be more sophisticated
	return s.userContributedStopRepo.UpdateDemandScore(stopID, newScore)
}

// GetNearbyUserStops gets user-contributed stops near a location
func (s *RoutingService) GetNearbyUserStops(lat, lng float64, radius int) ([]models.UserContributedStop, error) {
	return s.userContributedStopRepo.GetStops(lat, lng, radius)
}

// Enhanced convertPathToRoute with payment and safety info
func (s *RoutingService) convertPathToEnhancedRoute(path []Edge, options models.RouteOptions) *models.RouteResult {
	segments := make([]models.RouteSegment, 0, len(path))
	totalDistance := 0.0
	totalDuration := 0.0
	totalFare := 0.0
	safetyRatings := make([]float64, 0)
	paymentMethods := make(map[string]bool)

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

		// Add safety and payment info
		if edge.LineID != uuid.Nil {
			// Fetch line details
			line, err := s.lineRepo.GetByID(edge.LineID)
			if err == nil {
				segment.PublicRoute = &line
				
				// Use enhanced fare calculator
				fareInfo := s.fareCalculator.calculateSegmentFare(segment)
				segment.Fare = fareInfo.TotalFare
				
				segment.PaymentMethods = line.PaymentMethods
				segment.SafetyRating = line.SafetyRating
				segment.ReviewCount = line.ReviewCount

				// Track overall fare and payment methods
				totalFare += segment.Fare
				for _, method := range line.PaymentMethods {
					paymentMethods[method] = true
				}
			}
		} else {
			// Walking segment
			fareInfo := s.fareCalculator.calculateSegmentFare(segment)
			segment.Fare = fareInfo.TotalFare // Should be 0 for walking
			segment.PaymentMethods = []string{}
			segment.SafetyRating = 4.0 // Default walking safety rating
		}

		safetyRatings = append(safetyRatings, segment.SafetyRating)
		segments = append(segments, segment)
		totalDistance += edge.Distance
		totalDuration += edge.Duration
	}

	// Calculate overall safety rating
	overallSafetyRating := 0.0
	if len(safetyRatings) > 0 {
		sum := 0.0
		for _, rating := range safetyRatings {
			sum += rating
		}
		overallSafetyRating = sum / float64(len(safetyRatings))
	}

	// Convert payment methods map to slice
	availablePaymentMethods := make([]string, 0, len(paymentMethods))
	for method := range paymentMethods {
		availablePaymentMethods = append(availablePaymentMethods, method)
	}

	// Format duration for display
	totalTime := formatDuration(totalDuration)

	// Create enhanced summary
	summary := models.RouteSummary{
		TotalFare:        totalFare,
		AverageSafety:    overallSafetyRating,
		RequiredPayments: availablePaymentMethods,
	}

	return &models.RouteResult{
		TotalDistance:       totalDistance,
		TotalDuration:       totalDuration,
		TotalTime:           totalTime,
		TotalFare:           totalFare,
		Currency:            "DZD", // Default currency
		Segments:            segments,
		Summary:             summary,
		OverallSafetyRating: overallSafetyRating,
		PaymentMethods:      availablePaymentMethods,
		PaymentRequired:     totalFare > 0,
	}
}

// Helper to calculate fare for a segment
func (s *RoutingService) calculateSegmentFare(line models.Line, distance float64) float64 {
	// Simple fare calculation: base fare + distance-based fare
	return line.BaseFare + (line.FarePerKM * distance / 1000)
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

// Helper functions for route comparison
func routesEqual(route1, route2 *models.RouteResult) bool {
	if route1 == nil || route2 == nil {
		return route1 == route2
	}

	// Simple comparison based on segments
	if len(route1.Segments) != len(route2.Segments) {
		return false
	}

	for i, seg1 := range route1.Segments {
		seg2 := route2.Segments[i]
		if seg1.From.ID != seg2.From.ID || seg1.To.ID != seg2.To.ID {
			return false
		}
	}

	return true
}

func routeContainsRoute(routes []*models.RouteResult, route *models.RouteResult) bool {
	for _, r := range routes {
		if routesEqual(r, route) {
			return true
		}
	}
	return false
}

// OSRM Integration Helpers with enhanced profiles
func (s *RoutingService) getWalkingTime(fromLng, fromLat, toLng, toLat float64) (float64, float64, error) {
	osrmData, err := s.osrmRepo.GetRouteBetween(context.Background(), fromLng, fromLat, toLng, toLat, "walking")
	if err != nil || len(osrmData.Routes) == 0 {
		return 0, 0, err
	}
	return osrmData.Routes[0].Duration, osrmData.Routes[0].Distance, nil
}

func (s *RoutingService) getPublicTransportTime(fromLat, fromLng, toLat, toLng float64, transportType string) (float64, float64, error) {
	profile := s.getOSRMProfile(transportType)

	osrmData, err := s.osrmRepo.GetRouteBetween(context.Background(), fromLng, fromLat, toLng, toLat, profile)
	if err != nil || len(osrmData.Routes) == 0 {
		return 0, 0, err
	}
	return osrmData.Routes[0].Duration, osrmData.Routes[0].Distance, nil
}

// findRouteCandidates finds multiple route options using different algorithms
func (s *RoutingService) findRouteCandidates(ctx context.Context, fromLat, fromLng, toLat, toLng float64, options models.RouteOptions) ([]*models.RouteResult, error) {
	var candidates []*models.RouteResult

	// 1. Find nearest stations with safety filtering
	var startStations, endStations []models.Station
	var err error

	if options.SafetyFilter {
		startStations, err = s.stationRepo.GetNearbyWithSafetyFilter(fromLat, fromLng, options.MinSafetyRating, 1000)
		if err != nil {
			return nil, err
		}
		endStations, err = s.stationRepo.GetNearbyWithSafetyFilter(toLat, toLng, options.MinSafetyRating, 1000)
		if err != nil {
			return nil, err
		}
	} else {
		startStations, err = s.stationRepo.GetNearby(fromLat, fromLng, 1000)
		if err != nil {
			return nil, err
		}
		endStations, err = s.stationRepo.GetNearby(toLat, toLng, 1000)
		if err != nil {
			return nil, err
		}
	}

	// Include user-contributed stops if requested
	if options.IncludePrivate {
		userStartStops, err := s.userContributedStopRepo.GetStops(fromLat, fromLng, 1000)
		if err == nil {
			for _, userStop := range userStartStops {
				if userStop.Status == "approved" || userStop.DemandScore > 10.0 {
					station := models.Station{
						ID:        userStop.ID,
						Name:      userStop.Name,
						Latitude:  userStop.Latitude,
						Longitude: userStop.Longitude,
						Type:      "private_bus",
					}
					startStations = append(startStations, station)
				}
			}
		}

		userEndStops, err := s.userContributedStopRepo.GetStops(toLat, toLng, 1000)
		if err == nil {
			for _, userStop := range userEndStops {
				if userStop.Status == "approved" || userStop.DemandScore > 10.0 {
					station := models.Station{
						ID:        userStop.ID,
						Name:      userStop.Name,
						Latitude:  userStop.Latitude,
						Longitude: userStop.Longitude,
						Type:      "private_bus",
					}
					endStations = append(endStations, station)
				}
			}
		}
	}

	if len(startStations) == 0 || len(endStations) == 0 {
		return nil, errors.New("no stations found near start or end point")
	}

	// Try different combinations of start and end stations
	for _, startStation := range startStations[:min(3, len(startStations))] { // Limit to top 3
		for _, endStation := range endStations[:min(3, len(endStations))] {
			route := s.findSingleRoute(fromLat, fromLng, toLat, toLng, startStation, endStation, options)
			if route != nil {
				candidates = append(candidates, route)
			}
		}
	}

	// Remove duplicate routes
	candidates = s.removeDuplicateRoutes(candidates)

	return candidates, nil
}

// findSingleRoute finds a route between specific start and end stations
func (s *RoutingService) findSingleRoute(fromLat, fromLng, toLat, toLng float64, startStation, endStation models.Station, options models.RouteOptions) *models.RouteResult {
	// Check walking time limits
	walkToStart, _, err := s.getWalkingTime(fromLng, fromLat, startStation.Longitude, startStation.Latitude)
	if err != nil || (options.MaxWalkingTime > 0 && walkToStart > options.MaxWalkingTime) {
		return nil
	}

	walkFromEnd, _, err := s.getWalkingTime(endStation.Longitude, endStation.Latitude, toLng, toLat)
	if err != nil || (options.MaxWalkingTime > 0 && walkFromEnd > options.MaxWalkingTime) {
		return nil
	}

	// Find path through the transport network
	path := s.graph.ShortestPath(startStation.ID, endStation.ID)
	if path == nil {
		return nil
	}

	// Convert to route segments
	result := s.convertPathToEnhancedRoute(path, options)
	if result == nil {
		return nil
	}

	// Calculate total walking time and add to result
	result.Summary.WalkingDuration = walkToStart + walkFromEnd
	result.TotalDuration += walkToStart + walkFromEnd

	return result
}

// determineRouteMode determines the routing mode from options
func (s *RoutingService) determineRouteMode(options models.RouteOptions) RouteMode {
	// Use explicit mode if provided
	if options.Mode != "" {
		switch options.Mode {
		case "economy":
			return RouteModeEconomy
		case "fastest":
			return RouteModeFastest
		case "safest":
			return RouteModeSafest
		case "comfort":
			return RouteModeComfort
		case "balanced":
			return RouteModeBalanced
		}
	}

	// Infer mode from other options
	if options.MaxFare > 0 && options.MaxFare < 100 { // Very low fare limit
		return RouteModeEconomy
	}
	if options.SafetyFilter && options.MinSafetyRating >= 4.0 {
		return RouteModeSafest
	}
	if options.MaxWalkingTime > 0 && options.MaxWalkingTime < 300 { // Less than 5 minutes walking
		return RouteModeComfort
	}
	
	// Default to balanced
	return RouteModeBalanced
}

// rankRoutes scores and ranks routes based on weights and preferences
func (s *RoutingService) rankRoutes(candidates []*models.RouteResult, weights RouteWeight, options models.RouteOptions) *models.RouteResult {
	if len(candidates) == 0 {
		return nil
	}

	bestRoute := candidates[0]
	bestScore := s.calculateRouteScore(candidates[0], weights, options)

	for _, route := range candidates[1:] {
		score := s.calculateRouteScore(route, weights, options)
		if score > bestScore {
			bestScore = score
			bestRoute = route
		}
	}

	return bestRoute
}

// calculateRouteScore calculates a score for a route based on weights
func (s *RoutingService) calculateRouteScore(route *models.RouteResult, weights RouteWeight, options models.RouteOptions) float64 {
	// Normalize values for scoring (higher score = better)
	
	// Duration score (faster = better)
	durationScore := 1.0 / (1.0 + route.TotalDuration/3600.0) // Convert to hours and invert
	
	// Fare score (cheaper = better)
	fareScore := 1.0
	if route.TotalFare > 0 {
		fareScore = 1.0 / (1.0 + route.TotalFare/100.0) // Normalize by 100 DZD
	}
	
	// Safety score (higher rating = better)
	safetyScore := route.OverallSafetyRating / 5.0 // Normalize to 0-1
	
	// Walking score (less walking = better)
	walkingScore := 1.0 / (1.0 + route.Summary.WalkingDistance/1000.0) // Normalize by 1km
	
	// Transfer score (fewer transfers = better)
	transferScore := 1.0 / (1.0 + float64(route.Summary.TransfersCount))
	
	// Comfort score (placeholder - could consider crowding, vehicle quality, etc.)
	comfortScore := 0.8 // Default comfort score
	
	// Calculate weighted score
	totalScore := weights.Duration*durationScore +
		weights.Fare*fareScore +
		weights.Safety*safetyScore +
		weights.Walking*walkingScore +
		weights.Transfers*transferScore +
		weights.Comfort*comfortScore

	return totalScore
}

// GetFareEstimate provides a quick fare estimate for a transport type and distance
func (s *RoutingService) GetFareEstimate(transportType string, distance float64) models.FareInfo {
	if s.fareCalculator != nil {
		return s.fareCalculator.GetTransportFareEstimate(transportType, distance)
	}
	
	// Fallback basic calculation if fare calculator is not available
	return models.FareInfo{
		TransportType: transportType,
		Distance:      distance,
		BaseFare:      30.0, // Default bus fare
		TotalFare:     30.0,
		Currency:      "DZD",
	}
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (s *RoutingService) removeDuplicateRoutes(routes []*models.RouteResult) []*models.RouteResult {
	seen := make(map[string]bool)
	var unique []*models.RouteResult
	
	for _, route := range routes {
		// Create a simple hash of the route
		hash := fmt.Sprintf("%.2f-%.2f-%.0f", route.TotalDuration, route.TotalDistance, route.TotalFare)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, route)
		}
	}
	
	return unique
}
