package models

import (
	"time"

	"github.com/google/uuid"
)

// Journey represents a complete journey with multiple segments
type Journey struct {
	ID                  uuid.UUID       `json:"id"`
	From                Location        `json:"from"`
	To                  Location        `json:"to"`
	DepartureTime       time.Time       `json:"departure_time"`
	ArrivalTime         time.Time       `json:"arrival_time"`
	Duration            time.Duration   `json:"duration"`
	TotalDistance       float64         `json:"total_distance"`
	WalkingDistance     float64         `json:"walking_distance"`
	PublicDistance      float64         `json:"public_distance"`
	TransfersCount      int             `json:"transfers_count"`
	Segments            []JourneySegment `json:"segments"`
	Summary             JourneySummary  `json:"summary"`
	Score               float64         `json:"score"`                // For ranking journeys
	Tags                []string        `json:"tags,omitempty"`       // e.g., ["fastest", "least_walking"]
	CarbonFootprint     float64         `json:"carbon_footprint"`     // grams CO2
	Cost                float64         `json:"cost,omitempty"`       // estimated cost
	RealTimeData        bool            `json:"real_time_data"`       // whether real-time data was used
	Warnings            []string        `json:"warnings,omitempty"`   // service alerts, delays
}

// JourneySegment represents a part of the journey
type JourneySegment struct {
	ID                uuid.UUID            `json:"id"`
	From              Location             `json:"from"`
	To                Location             `json:"to"`
	Mode              TransportMode        `json:"mode"`
	DepartureTime     time.Time            `json:"departure_time"`
	ArrivalTime       time.Time            `json:"arrival_time"`
	Duration          time.Duration        `json:"duration"`
	Distance          float64              `json:"distance"`
	Geometry          [][]float64          `json:"geometry"`
	Instructions      []Instruction        `json:"instructions,omitempty"`
	PublicTransport   *PublicTransportInfo `json:"public_transport,omitempty"`
	RealTimeInfo      *RealTimeInfo        `json:"real_time_info,omitempty"`
	Accessibility     AccessibilityInfo    `json:"accessibility"`
	CarbonFootprint   float64              `json:"carbon_footprint"`
	Cost              float64              `json:"cost,omitempty"`
}

// Location represents a geographic location
type Location struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Type        string    `json:"type"`        // "address", "station", "stop", "poi"
	Address     string    `json:"address,omitempty"`
	Platform    string    `json:"platform,omitempty"`
	StationID   uuid.UUID `json:"station_id,omitempty"`
}

// TransportMode represents different modes of transport
type TransportMode string

const (
	ModeWalking         TransportMode = "walking"
	ModeBicycle         TransportMode = "bicycle"
	ModePublicTransport TransportMode = "public_transport"
	ModeTram            TransportMode = "tram"
	ModeBus             TransportMode = "bus"
	ModeMetro           TransportMode = "metro"
	ModeTrain           TransportMode = "train"
	ModeRideshare       TransportMode = "rideshare"
	ModeTaxi            TransportMode = "taxi"
)

// PublicTransportInfo contains details about public transport segments
type PublicTransportInfo struct {
	RouteID       uuid.UUID `json:"route_id"`
	RouteName     string    `json:"route_name"`
	RouteShortName string   `json:"route_short_name"`
	RouteColor    string    `json:"route_color,omitempty"`
	RouteTextColor string   `json:"route_text_color,omitempty"`
	AgencyName    string    `json:"agency_name"`
	TripID        uuid.UUID `json:"trip_id,omitempty"`
	TripHeadsign  string    `json:"trip_headsign,omitempty"`
	StopsCount    int       `json:"stops_count"`
	IntermediateStops []Location `json:"intermediate_stops,omitempty"`
}

// RealTimeInfo contains real-time information
type RealTimeInfo struct {
	ExpectedDeparture time.Time `json:"expected_departure,omitempty"`
	ExpectedArrival   time.Time `json:"expected_arrival,omitempty"`
	Delay             time.Duration `json:"delay,omitempty"`
	IsRealTime        bool      `json:"is_real_time"`
	ServiceAlert      string    `json:"service_alert,omitempty"`
	Occupancy         string    `json:"occupancy,omitempty"` // "empty", "few_seats", "standing_room", "full"
}

// AccessibilityInfo contains accessibility information
type AccessibilityInfo struct {
	WheelchairAccessible bool `json:"wheelchair_accessible"`
	StepsCount           int  `json:"steps_count,omitempty"`
	ElevatorAvailable    bool `json:"elevator_available"`
	AudioAnnouncements   bool `json:"audio_announcements"`
	VisualAnnouncements  bool `json:"visual_announcements"`
}

// Instruction represents a navigation instruction
type Instruction struct {
	Text        string    `json:"text"`
	Distance    float64   `json:"distance,omitempty"`
	Duration    time.Duration `json:"duration,omitempty"`
	Coordinate  []float64 `json:"coordinate,omitempty"`
	Type        string    `json:"type"` // "turn", "continue", "board", "alight"
	Modifier    string    `json:"modifier,omitempty"` // "left", "right", "straight"
}

// JourneySummary provides a high-level overview
type JourneySummary struct {
	TotalDuration       time.Duration `json:"total_duration"`
	WalkingDuration     time.Duration `json:"walking_duration"`
	PublicDuration      time.Duration `json:"public_duration"`
	WaitingDuration     time.Duration `json:"waiting_duration"`
	TransfersCount      int           `json:"transfers_count"`
	PublicRoutes        []string      `json:"public_routes"`
	DominantMode        TransportMode `json:"dominant_mode"`
	Efficiency          float64       `json:"efficiency"`    // 0-1 score
	Comfort             float64       `json:"comfort"`       // 0-1 score
	Environmental       float64       `json:"environmental"` // 0-1 score
}

// JourneyRequest represents a request for journey planning
type JourneyRequest struct {
	From              Location            `json:"from"`
	To                Location            `json:"to"`
	DepartureTime     *time.Time          `json:"departure_time,omitempty"`
	ArrivalTime       *time.Time          `json:"arrival_time,omitempty"`
	TimeType          string              `json:"time_type"` // "departure", "arrival"
	Modes             []TransportMode     `json:"modes,omitempty"`
	MaxWalkingDistance float64             `json:"max_walking_distance,omitempty"`
	MaxTransfers      int                 `json:"max_transfers,omitempty"`
	WheelchairAccessible bool             `json:"wheelchair_accessible,omitempty"`
	Preferences       JourneyPreferences  `json:"preferences,omitempty"`
	MaxResults        int                 `json:"max_results,omitempty"`
	IncludeRealTime   bool                `json:"include_real_time,omitempty"`
}

// JourneyPreferences allows users to specify their preferences
type JourneyPreferences struct {
	PreferFastest     bool    `json:"prefer_fastest"`
	PreferLeastWalking bool   `json:"prefer_least_walking"`
	PreferLeastTransfers bool `json:"prefer_least_transfers"`
	PreferCheapest    bool    `json:"prefer_cheapest"`
	PreferGreenest    bool    `json:"prefer_greenest"`
	WalkingSpeed      float64 `json:"walking_speed,omitempty"` // m/s
	MaxWalkingTime    time.Duration `json:"max_walking_time,omitempty"`
}

// JourneyResponse represents the response from journey planning
type JourneyResponse struct {
	Request     JourneyRequest `json:"request"`
	Journeys    []Journey      `json:"journeys"`
	Metadata    JourneyMetadata `json:"metadata"`
	Error       string         `json:"error,omitempty"`
	Warnings    []string       `json:"warnings,omitempty"`
}

// JourneyMetadata contains metadata about the journey planning response
type JourneyMetadata struct {
	PlanningTime      time.Duration `json:"planning_time"`
	DataSources       []string      `json:"data_sources"`
	RealTimeUsed      bool          `json:"real_time_used"`
	TimetableDate     time.Time     `json:"timetable_date"`
	ServiceAlerts     []ServiceAlert `json:"service_alerts,omitempty"`
	RegionsSearched   []string      `json:"regions_searched,omitempty"`
}

// ServiceAlert represents a service alert or disruption
type ServiceAlert struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Severity      string    `json:"severity"` // "info", "warning", "severe"
	Effect        string    `json:"effect"`   // "detour", "delay", "cancellation"
	StartTime     time.Time `json:"start_time,omitempty"`
	EndTime       time.Time `json:"end_time,omitempty"`
	AffectedRoutes []string `json:"affected_routes,omitempty"`
	AffectedStops []string  `json:"affected_stops,omitempty"`
}
