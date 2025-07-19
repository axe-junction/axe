package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Station struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name          string
	Latitude      float64
	Longitude     float64
	Type          string     // "bus", "tram", "train", "metro", "private_bus"
	SafetyRating  float64    `json:"safety_rating"`            // Average safety rating (1-5)
	ReviewCount   int        `json:"review_count"`             // Number of reviews
	IsVerified    bool       `json:"is_verified"`              // Official vs user-contributed
	ContributorID *uuid.UUID `json:"contributor_id,omitempty"` // User who added this station
	DemandScore   float64    `json:"demand_score"`             // For user-contributed stops
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
type Stations []Station
// Line represents a transport line (bus, tram, metro, train, etc.)
type Line struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Name           string    `gorm:"column:name" json:"name"`
	Type           string    `gorm:"column:type" json:"type"`
	Agency         string    `gorm:"column:agency" json:"agency"`                         // "ETUSA", "SNTF", etc.
	BaseFare       float64   `gorm:"column:base_fare" json:"base_fare"`                   // Base fare in local currency
	FarePerKM      float64   `gorm:"column:fare_per_km" json:"fare_per_km"`              // Additional fare per km
	PaymentMethods []string  `gorm:"type:text[];column:payment_methods" json:"payment_methods"` // ["cash", "card", "mobile", "subscription"]
	SafetyRating   float64   `gorm:"column:safety_rating" json:"safety_rating"`          // Average safety rating (1-5)
	ReviewCount    int       `gorm:"column:review_count" json:"review_count"`            // Number of reviews
	IsActive       bool      `gorm:"column:is_active" json:"is_active"`                  // Whether line is currently operating
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`
	
	Stations []Station `gorm:"many2many:line_stations;" json:"stations,omitempty"`
	Stops    []Stop    `gorm:"foreignKey:LineID" json:"stops,omitempty"`
}
type Lines []Line



type Stop struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	LineID    uuid.UUID
	StationID uuid.UUID
	Sequence  int
	Line      Line    `gorm:"foreignKey:LineID"`
	Station   Station `gorm:"foreignKey:StationID"`
}

type Transfer struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	FromStationID   uuid.UUID
	ToStationID     uuid.UUID
	WalkingTime     float64 // seconds
	WalkingDistance float64 // meters
	FromStation     Station `gorm:"foreignKey:FromStationID"`
	ToStation       Station `gorm:"foreignKey:ToStationID"`
}

type OSRMResponse struct {
	Code   string `json:"code"`
	Routes []struct {
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
	} `json:"routes"`
}

// Enhanced repository interfaces

type StationRepo interface {
	GetAll() ([]Station, error)
	GetByID(id string) ([]Station, error)
	GetNearby(lat, lng float64, radius ...int) ([]Station, error)
	GetNearbyWithSafetyFilter(lat, lng float64, minSafetyRating float64, radius ...int) ([]Station, error)
	UpdateSafetyRating(stationID uuid.UUID, rating float64) error
	CreateUserContributedStop(stop UserContributedStop) error
	GetUserContributedStops(lat, lng float64, radius int) ([]UserContributedStop, error)
	PromoteUserStopToOfficial(stopID uuid.UUID) error
}

type LineRepo interface {
	GetAll() (Lines, error)
	GetByID(id uuid.UUID) (Line, error)
	GetByType(typee string) (Lines, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lines, error)
	GetStopsForLine(lineID uuid.UUID) ([]Stop, error)
	GetLinesWithSafetyFilter(minSafetyRating float64) (Lines, error)
	UpdateLineSafetyRating(lineID uuid.UUID, rating float64) error
	GetLinesByPaymentMethod(paymentMethods []string) (Lines, error)
}

type ReviewRepo interface {
	CreateReview(review Review) error
	GetReviewsForStation(stationID uuid.UUID) ([]Review, error)
	GetReviewsForLine(lineID uuid.UUID) ([]Review, error)
	GetAverageRating(entityID uuid.UUID, entityType string) (float64, int, error)
	GetSafetyRating(entityID uuid.UUID, entityType string) (float64, int, error)
}

type PaymentRepo interface {
	CreatePayment(payment Payment) error
	GetPaymentsByUser(userID uuid.UUID) ([]Payment, error)
	GetPaymentByID(paymentID uuid.UUID) (Payment, error)
	UpdatePaymentStatus(paymentID uuid.UUID, status string) error
	CalculateRouteFare(segments []RouteSegment) (float64, error)
}

type UserContributedStopRepo interface {
	CreateStop(stop UserContributedStop) error
	GetStops(lat, lng float64, radius int) ([]UserContributedStop, error)
	AddDemandRequest(request StopDemandRequest) error
	GetDemandScore(stopID uuid.UUID) (float64, error)
	UpdateDemandScore(stopID uuid.UUID, score float64) error
	GetHighDemandStops(threshold float64) ([]UserContributedStop, error)
	PromoteToOfficial(stopID uuid.UUID) error
}

type OSRMRepo interface {
	GetRouteBetween(ctx context.Context, fromLng, fromLat, toLng, toLat float64, profile string) (*OSRMResponse, error)
	GetMultiModalRoute(ctx context.Context, fromLng, fromLat, toLng, toLat float64, transportModes []string) (*OSRMResponse, error)
	GetOptimizedRoute(ctx context.Context, coordinates [][]float64, profile string) (*OSRMResponse, error)
	GetMatrixDurations(ctx context.Context, sources, destinations [][]float64, profile string) ([][]float64, error)
}

type RoutingService interface {
	GetBestRoute(ctx context.Context, fromLat, fromLng, toLat, toLng float64, options RouteOptions) (*RouteResult, error)
	GetMultipleRoutes(ctx context.Context, fromLat, fromLng, toLat, toLng float64, options RouteOptions) ([]*RouteResult, error)
	CalculateRouteFare(segments []RouteSegment) (float64, error)
	ProcessRoutePayment(userID uuid.UUID, routeResult *RouteResult, paymentMethod string) (*Payment, error)
	AddUserContributedStop(stop UserContributedStop) error
	RequestUserStop(stopID uuid.UUID, userID uuid.UUID, fromLat, fromLng, toLat, toLng float64) error
	GetNearbyUserStops(lat, lng float64, radius int) ([]UserContributedStop, error)
	BuildTransportGraph(ctx context.Context) error
}

type Review struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID       uuid.UUID  `json:"user_id"`
	StationID    *uuid.UUID `json:"station_id,omitempty"`
	LineID       *uuid.UUID `json:"line_id,omitempty"`
	Rating       float64    `json:"rating" binding:"required,min=1,max=5"`
	SafetyRating float64    `json:"safety_rating" binding:"required,min=1,max=5"`
	Comment      string     `json:"comment"`
	ReviewType   string     `json:"review_type"` // "station", "line", "route"
	IsVerified   bool       `json:"is_verified"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Payment struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID        uuid.UUID `json:"user_id"`
	RouteID       uuid.UUID `json:"route_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency" default:"DZD"`
	PaymentMethod string    `json:"payment_method"` // "cash", "card", "mobile", "subscription"
	Status        string    `json:"status"`         // "pending", "completed", "failed", "refunded"
	TransactionID string    `json:"transaction_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type UserContributedStop struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string    `json:"name" binding:"required"`
	Latitude       float64   `json:"latitude" binding:"required"`
	Longitude      float64   `json:"longitude" binding:"required"`
	ContributorID  uuid.UUID `json:"contributor_id"`
	RequestCount   int       `json:"request_count"` // How many users requested this stop
	DemandScore    float64   `json:"demand_score"`  // Calculated demand score
	Status         string    `json:"status"`        // "pending", "approved", "rejected", "promoted"
	Description    string    `json:"description"`
	NearbyLandmark string    `json:"nearby_landmark"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Demand tracking for user-contributed stops
type StopDemandRequest struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserContributedStopID uuid.UUID `json:"user_contributed_stop_id"`
	UserID                uuid.UUID `json:"user_id"`
	RequestedAt           time.Time `json:"requested_at"`
	FromLatitude          float64   `json:"from_latitude"`
	FromLongitude         float64   `json:"from_longitude"`
	ToLatitude            float64   `json:"to_latitude"`
	ToLongitude           float64   `json:"to_longitude"`
}

type RouteOptions struct {
	Mode            string   `json:"mode"`              // "economy", "fastest", "safest", "balanced", "comfort"
	SafetyFilter    bool     `json:"safety_filter"`     // Filter out low-rated routes
	MinSafetyRating float64  `json:"min_safety_rating"` // Minimum safety rating (default: 3.0)
	PaymentMethods  []string `json:"payment_methods"`   // Preferred payment methods
	MaxWalkingTime  float64  `json:"max_walking_time"`  // Maximum walking time in seconds
	IncludePrivate  bool     `json:"include_private"`   // Include user-contributed stops
	MaxFare         float64  `json:"max_fare"`          // Maximum acceptable fare
}

type RouteSegment struct {
	From        Station     `json:"from"`
	To          Station     `json:"to"`
	Mode        string      `json:"mode"`                   // "walking", "bus", "tram", "train", "metro", "private_bus"
	PublicRoute *Line       `json:"public_route,omitempty"` // Only for public transport
	Distance    float64     `json:"distance"`               // in meters
	Duration    float64     `json:"duration"`               // in seconds
	Geometry    [][]float64 `json:"geometry,omitempty"`     // Optional: coordinates for the path
	// Payment and safety info
	Fare           float64  `json:"fare"`            // Segment fare
	PaymentMethods []string `json:"payment_methods"` // Available payment methods
	SafetyRating   float64  `json:"safety_rating"`   // Safety rating for this segment
	ReviewCount    int      `json:"review_count"`    // Number of reviews
}

type RouteResult struct {
	From          Station        `json:"from"`
	To            Station        `json:"to"`
	TotalDistance float64        `json:"total_distance"` // in meters
	TotalDuration float64        `json:"total_duration"` // in seconds
	TotalTime     string         `json:"total_time"`     // human-readable, e.g., "1h 30m"
	TotalFare     float64        `json:"total_fare"`     // Total route fare
	Currency      string         `json:"currency"`       // Currency (e.g., "DZD")
	Segments      []RouteSegment `json:"segments"`
	Route         Lines          `json:"route,omitempty"` // For backward compatibility
	Summary       RouteSummary   `json:"summary"`
	OverallSafetyRating float64  `json:"overall_safety_rating"`
	PaymentMethods      []string `json:"payment_methods"`
	PaymentRequired     bool     `json:"payment_required"`
}

type RouteSummary struct {
	WalkingDistance  float64  `json:"walking_distance"`
	PublicDistance   float64  `json:"public_distance"`
	WalkingDuration  float64  `json:"walking_duration"`
	PublicDuration   float64  `json:"public_duration"`
	TransfersCount   int      `json:"transfers_count"`
	PublicRouteNames []string `json:"public_route_names"`
	TotalFare        float64  `json:"total_fare"`
	AverageSafety    float64  `json:"average_safety"`
	RequiredPayments []string `json:"required_payments"`
}


type OSRMMatrixResponse struct {
	Durations [][]float64 `json:"durations"`
	Sources   []struct {
		Location []float64 `json:"location"`
		Name     string    `json:"name"`
	} `json:"sources"`
	Destinations []struct {
		Location []float64 `json:"location"`
		Name     string    `json:"name"`
	} `json:"destinations"`
}

type OSRMOptimizedResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Distance float64 `json:"distance"`
		Duration float64 `json:"duration"`
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"`
		} `json:"geometry"`
		Waypoints []struct {
			Location []float64 `json:"location"`
			Name     string    `json:"name"`
		} `json:"waypoints"`
	} `json:"routes"`
}

type CachedRoute struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	RouteHash string    `json:"route_hash" gorm:"unique"`
	FromLat   float64   `json:"from_lat"`
	FromLng   float64   `json:"from_lng"`
	ToLat     float64   `json:"to_lat"`
	ToLng     float64   `json:"to_lng"`
	RouteData []byte    `json:"route_data" gorm:"type:jsonb"` // Serialized RouteResult
	Options   []byte    `json:"options" gorm:"type:jsonb"`    // Serialized RouteOptions
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PaymentTransaction struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PaymentID       uuid.UUID  `json:"payment_id"`
	TransactionType string     `json:"transaction_type"` // "charge", "refund", "partial_refund"
	Amount          float64    `json:"amount"`
	Currency        string     `json:"currency"`
	ExternalRef     string     `json:"external_ref"` // Reference from payment provider
	Status          string     `json:"status"`       // "pending", "completed", "failed"
	FailureReason   string     `json:"failure_reason,omitempty"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
}

type RouteCache interface {
	StoreRoute(route CachedRoute) error
	GetRoute(routeHash string) (*CachedRoute, error)
	DeleteExpiredRoutes() error
	GenerateRouteHash(fromLat, fromLng, toLat, toLng float64, options RouteOptions) string
}

type PaymentTransactionRepo interface {
	CreateTransaction(transaction PaymentTransaction) error
	GetTransactionsByPayment(paymentID uuid.UUID) ([]PaymentTransaction, error)
	UpdateTransactionStatus(transactionID uuid.UUID, status, failureReason string) error
	GetPendingTransactions() ([]PaymentTransaction, error)
}

type FareInfo struct {
	TransportType string  `json:"transport_type"` // "bus", "metro", "train", "taxi"
	BaseFare      float64 `json:"base_fare"`      // Base fare in DZD
	Distance      float64 `json:"distance"`       // Distance in km
	DistanceFare  float64 `json:"distance_fare"`  // Additional fare per km
	TotalFare     float64 `json:"total_fare"`     // Total calculated fare
	Currency      string  `json:"currency"`       // "DZD"
}

type RouteFareInfo struct {
	RouteID       uuid.UUID  `json:"route_id"`
	TotalFare     float64    `json:"total_fare"`
	Currency      string     `json:"currency"`
	FareBreakdown []FareInfo `json:"fare_breakdown"`
	EstimatedCost string     `json:"estimated_cost"` // Human readable cost
}
