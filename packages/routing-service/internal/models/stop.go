package models

import "github.com/google/uuid"

type LigneStop struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	RouteID   uuid.UUID `gorm:"column:route_id" json:"route_id"`
	StationID uuid.UUID `gorm:"column:station_id" json:"station_id"`
	Sequence  int       `gorm:"column:sequence" json:"sequence"`
	Route     Ligne     `json:"route" gorm:"foreignKey:RouteID"`
	Station   Station   `json:"station" gorm:"foreignKey:StationID"`
}
type LigneStops []LigneStop

// TableName sets the table name for LigneStop
func (LigneStop) TableName() string {
	return "stops"
}
