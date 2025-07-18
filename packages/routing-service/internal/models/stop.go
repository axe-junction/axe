package models

import "github.com/google/uuid"

type Stop struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RouteID   uuid.UUID `json:"route_id"`
	StationID uuid.UUID `json:"station_id"`
	Sequence  int       `json:"sequence"`
	Route     Ligne     `json:"route" gorm:"foreignKey:RouteID"`
	Station   Station   `json:"station" gorm:"foreignKey:StationID"`
}
type Stops []Stop
