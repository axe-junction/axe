package models


type Stop struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RouteID   string  `json:"route_id"`
	StationID string  `json:"station_id"`
	Sequence int     `json:"sequence"`
}
type Stops []Stop
