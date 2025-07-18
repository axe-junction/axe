package models

import (
	"context"

	"github.com/google/uuid"
)

type Station struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name      string
	Latitude  float64
	Longitude float64
	Type      string // "bus", "tram", "train", "metro"
}
type Stations []Station
type Line struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name   string
	Type   string
	Agency string // "ETUSA", "SNTF", etc.
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

// New table for transfers
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

type StationRepo interface {
	GetAll() ([]Station, error)
	GetByID(id string) ([]Station, error)
	GetNearby(lat, lng float64, radius ...int) ([]Station, error)
}

type LineRepo interface {
	GetAll() (Lines, error)
	GetByID(id uuid.UUID) (Line, error)
	GetByType(typee string) (Lines, error)
	GetRouteBetweenStations(ctx context.Context, startID, endID uuid.UUID) (Lines, error)
	GetStopsForLine(lineID uuid.UUID) ([]Stop, error)
}

type OSRMRepo interface {
	GetRouteBetween(ctx context.Context, fromLng, fromLat, toLng, toLat float64, profile string) (*OSRMResponse, error)
}

type RouteSegment struct {
	From        Station     `json:"from"`
	To          Station     `json:"to"`
	Mode        string      `json:"mode"`                   // "walking", "bus", "tram", "train", "metro"
	PublicRoute *Line       `json:"public_route,omitempty"` // Only for public transport
	Distance    float64     `json:"distance"`               // in meters
	Duration    float64     `json:"duration"`               // in seconds
	Geometry    [][]float64 `json:"geometry,omitempty"`     // Optional: coordinates for the path
}

type RouteResult struct {
	From          Station        `json:"from"`
	To            Station        `json:"to"`
	TotalDistance float64        `json:"total_distance"` // in meters
	TotalDuration float64        `json:"total_duration"` // in seconds
	TotalTime     string         `json:"total_time"`     // human-readable, e.g., "1h 30m"
	Segments      []RouteSegment `json:"segments"`
	Route         Lines          `json:"route,omitempty"` // For backward compatibility
	Summary       RouteSummary   `json:"summary"`
}

type RouteSummary struct {
	WalkingDistance  float64  `json:"walking_distance"`
	PublicDistance   float64  `json:"public_distance"`
	WalkingDuration  float64  `json:"walking_duration"`
	PublicDuration   float64  `json:"public_duration"`
	TransfersCount   int      `json:"transfers_count"`
	PublicRouteNames []string `json:"public_route_names"`
}
